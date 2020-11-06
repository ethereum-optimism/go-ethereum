package rollup

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/addressmanager"
	ctc "github.com/ethereum/go-ethereum/contracts/canonicaltransactionchain"
	mgr "github.com/ethereum/go-ethereum/contracts/executionmanager"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Interface used for communicating with Ethereum 1 nodes
type EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	NetworkID(context.Context) (*big.Int, error)
	SyncProgress(context.Context) (*ethereum.SyncProgress, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
}

// Interface used for receiving events from the canonical transaction chain
type CTCEventFilterer interface {
	ParseTransactionEnqueued(types.Log) (*ctc.OVMCanonicalTransactionChainTransactionEnqueued, error)
	ParseQueueBatchAppended(types.Log) (*ctc.OVMCanonicalTransactionChainQueueBatchAppended, error)
	ParseSequencerBatchAppended(types.Log) (*ctc.OVMCanonicalTransactionChainSequencerBatchAppended, error)
}

type CTCCaller interface {
	GetNextQueueIndex(*bind.CallOpts) (*big.Int, error)
	GetNumPendingQueueElements(*bind.CallOpts) (*big.Int, error)
	GetQueueElement(*bind.CallOpts, *big.Int) (ctc.Lib_OVMCodecQueueElement, error)
	GetTotalElements(*bind.CallOpts) (*big.Int, error)
}

type ExecutionManagerCaller interface {
	GetMaxTransactionGasLimit(opts *bind.CallOpts) (*big.Int, error)
}

type RollupTxsByIndex []*RollupTransaction

func (l RollupTxsByIndex) Len() int           { return len(l) }
func (l RollupTxsByIndex) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l RollupTxsByIndex) Less(i, j int) bool { return l[i].index < l[j].index }

// Consider adding a processed bool for sanity check
type RollupTransaction struct {
	tx          *types.Transaction
	timestamp   time.Time
	blockHeight uint64
	index       uint64
	executed    bool
}

// Move this to its own file
func NewTransactionCache() *TransactionCache {
	return &TransactionCache{
		m: new(sync.Map),
	}
}

type TransactionCache struct {
	m *sync.Map
}

func (t *TransactionCache) Store(index uint64, rtx *RollupTransaction) {
	t.m.Store(index, rtx)
}

func (t *TransactionCache) Delete(index uint64) {
	t.m.Delete(index)
}

func (t *TransactionCache) Load(index uint64) (*RollupTransaction, bool) {
	result, ok := t.m.Load(index)
	if !ok {
		return nil, false
	}
	rtx, ok := result.(*RollupTransaction)
	if !ok {
		log.Error("Incorrect type in transaction cache", "type", fmt.Sprintf("%T", rtx))
		return nil, false
	}
	return rtx, true
}

func (t *TransactionCache) Range(f func(uint64, *RollupTransaction)) {
	t.m.Range(func(key interface{}, value interface{}) bool {
		rtx, ok := value.(*RollupTransaction)
		if !ok {
			log.Error("Unexpected value in transaction cache", "type", fmt.Sprintf("%T", value))
			return true
		}
		index, ok := key.(uint64)
		if !ok {
			log.Error("Unexpected key type in transaction cache", "type", fmt.Sprintf("%T", key))
			return true

		}
		f(index, rtx)
		return true
	})
}

// These variables represent the event signatures
var (
	transactionEnqueuedEventSignature    = crypto.Keccak256([]byte("TransactionEnqueued(address,address,uint256,bytes,uint256,uint256)"))
	queueBatchAppendedEventSignature     = crypto.Keccak256([]byte("QueueBatchAppended(uint256,uint256,uint256)"))
	sequencerBatchAppendedEventSignature = crypto.Keccak256([]byte("SequencerBatchAppended(uint256,uint256,uint256)"))
)

// Eth1Data represents the last processed ethereum 1 data.
type Eth1Data struct {
	BlockHeight uint64
	BlockHash   common.Hash
}

// SyncService implements the verifier functionality as well as the reorg
// protection for the sequencer.
type SyncService struct {
	ctx                              context.Context
	cancel                           context.CancelFunc
	verifier                         bool
	processingLock                   sync.RWMutex
	txLock                           sync.Mutex
	db                               ethdb.Database
	enable                           bool
	ctcFilterer                      CTCEventFilterer
	ctcCaller                        CTCCaller
	mgrCaller                        ExecutionManagerCaller
	txCache                          *TransactionCache
	ethclient                        EthereumClient
	ethrpcclient                     *ethclient.Client
	logClient                        bind.ContractFilterer
	eth1HTTPEndpoint                 string
	eth1ChainId                      uint64
	eth1NetworkId                    uint64
	txpool                           *core.TxPool
	bc                               *core.BlockChain
	clearTransactionsTicker          *time.Ticker
	clearTransactionsAfter           uint64
	heads                            chan *types.Header
	doneProcessing                   chan uint64
	signer                           types.Signer
	key                              ecdsa.PrivateKey
	address                          common.Address
	gasLimit                         uint64
	syncing                          bool
	Eth1Data                         Eth1Data
	HeaderCache                      [2048]*types.Header
	sequencerIngestTicker            *time.Ticker
	ctcDeployHeight                  *big.Int
	AddressResolverAddress           common.Address
	CanonicalTransactionChainAddress common.Address
	SequencerDecompressionAddress    common.Address
	StateCommitmentChainAddress      common.Address
	ExecutionManagerAddress          common.Address
}

// NewSyncService returns an initialized sync service
func NewSyncService(ctx context.Context, cfg Config, txpool *core.TxPool, bc *core.BlockChain, db ethdb.Database) (*SyncService, error) {
	if txpool == nil {
		return nil, errors.New("Must pass TxPool to SyncService")
	}
	if bc == nil {
		return nil, errors.New("Must pass BlockChain to SyncService")
	}

	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // satisfy govet

	if cfg.TxIngestionSignerKey == nil {
		cfg.TxIngestionSignerKey, _ = crypto.GenerateKey()
	}
	address := crypto.PubkeyToAddress(cfg.TxIngestionSignerKey.PublicKey)

	// Layer 2 chainid to use for signing
	chainID := bc.Config().ChainID
	service := SyncService{
		ctx:                              ctx,
		cancel:                           cancel,
		verifier:                         cfg.IsVerifier,
		enable:                           cfg.Eth1SyncServiceEnable,
		heads:                            make(chan *types.Header),
		doneProcessing:                   make(chan uint64),
		eth1HTTPEndpoint:                 cfg.Eth1HTTPEndpoint,
		AddressResolverAddress:           cfg.AddressResolverAddress,
		CanonicalTransactionChainAddress: cfg.CanonicalTransactionChainAddress,
		SequencerDecompressionAddress:    cfg.SequencerDecompressionAddress,
		signer:                           types.NewOVMSigner(chainID),
		key:                              *cfg.TxIngestionSignerKey,
		address:                          address,
		gasLimit:                         cfg.GasLimit,
		txpool:                           txpool,
		bc:                               bc,
		eth1ChainId:                      cfg.Eth1ChainId,
		eth1NetworkId:                    cfg.Eth1NetworkId,
		ctcDeployHeight:                  cfg.CanonicalTransactionChainDeployHeight,
		db:                               db,
		clearTransactionsAfter:           (5760 * 15), // 15 days worth of blocks
		clearTransactionsTicker:          time.NewTicker(time.Hour),
		sequencerIngestTicker:            time.NewTicker(15 * time.Second),
		txCache:                          NewTransactionCache(),
		HeaderCache:                      [2048]*types.Header{},
	}

	// Always initialize syncing to true to start, the sequencer can toggle off
	// syncing while the verifier is always syncing
	service.setSyncStatus(true)

	return &service, nil
}

// Start initializes the service, connecting to Ethereum1 and starting the
// subservices required for the operation of the SyncService.
// txs through syncservice go to mempool.locals
// txs through rpc go to mempool.remote
func (s *SyncService) Start() error {
	if !s.enable {
		return nil
	}

	log.Info("Initializing Sync Service", "endpoint", s.eth1HTTPEndpoint, "chainid", s.eth1ChainId, "networkid", s.eth1NetworkId, "address resolver", s.AddressResolverAddress)
	log.Info("Watching topics", "transaction-enqueued", hexutil.Encode(transactionEnqueuedEventSignature), "queue-batch-appened", hexutil.Encode(queueBatchAppendedEventSignature), "sequencer-batch-appended", hexutil.Encode(sequencerBatchAppendedEventSignature))

	blockHeight := rawdb.ReadHeadEth1HeaderHeight(s.db)
	if blockHeight == 0 {
		if s.ctcDeployHeight == nil {
			return errors.New("Must configure with canonical transaction chain deploy height")
		}
		cfgHeight := s.ctcDeployHeight.Uint64()
		// Do not underflow in the case where the ctc contract is in the genesis state
		if cfgHeight == 0 {
			blockHeight = cfgHeight
		} else {
			blockHeight = cfgHeight - 1
		}
	}
	blockHash := rawdb.ReadHeadEth1HeaderHash(s.db)
	eth1Data := Eth1Data{
		BlockHeight: blockHeight,
		BlockHash:   blockHash,
	}
	s.Eth1Data = eth1Data

	_, client, err := s.dialEth1Node()
	if err != nil {
		return fmt.Errorf("Cannot dial eth1 nodes: %w", err)
	}
	s.ethrpcclient = client
	s.ethclient = client
	s.logClient = client

	err = s.verifyNetwork()
	if err != nil {
		return fmt.Errorf("Wrong network: %w", err)
	}
	// Resolve addresses and set them globally
	err = s.resolveAddresses()
	if err != nil {
		return fmt.Errorf("Error resolving addresses: %w", err)
	}
	// Bind to the contracts
	err = s.bindContracts()
	if err != nil {
		return fmt.Errorf("Error binding to contracts: %w", err)
	}
	// Check the sync status of the eth1 node
	err = s.checkSyncStatus()
	if err != nil {
		return fmt.Errorf("Bad sync status: %w", err)
	}

	go s.LogDoneProcessing()
	// Catch up to the tip of the eth1 chain
	err = s.processHistoricalLogs()
	if err != nil {
		return fmt.Errorf("Cannot process historical logs: %w", err)
	}

	gasLimit, err := s.mgrCaller.GetMaxTransactionGasLimit(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(s.Eth1Data.BlockHeight),
		Context:     s.ctx,
	})
	if err != nil {
		return fmt.Errorf("Cannot fetch gas limit: %w", err)
	}
	s.gasLimit = gasLimit.Uint64()
	log.Info("Setting max transaction gas limit", "gas limit", s.gasLimit)

	go s.Loop()
	go s.pollHead()
	go s.ClearTransactionLoop()

	if !s.verifier {
		go s.sequencerIngestQueue()
	}

	return nil
}

func (s *SyncService) getCommonAncestor(index *big.Int, list *[]*types.Header) (uint64, error) {
	header, err := s.ethclient.HeaderByNumber(s.ctx, index)
	if err != nil {
		return 0, fmt.Errorf(":%w", err)
	}
	number := header.Number.Uint64()
	// Do not allow for reorgs past the deployment height
	// of the contracts.
	if number == s.ctcDeployHeight.Uint64() {
		return number, nil
	}
	cached := s.HeaderCache[number%2048]
	if cached != nil && bytes.Equal(header.Hash().Bytes(), cached.Hash().Bytes()) {
		return number, nil
	}
	prevNumber := new(big.Int).SetUint64(number - 1)
	*list = append(*list, header)
	return s.getCommonAncestor(prevNumber, list)
}

func (s *SyncService) pollHead() {
	headTicker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-headTicker.C:
			head, err := s.ethclient.HeaderByNumber(s.ctx, nil)
			if err != nil {
				log.Error("Cannot fetch tip")
				continue
			}
			// The tip is the same
			if bytes.Equal(head.Hash().Bytes(), s.Eth1Data.BlockHash.Bytes()) {
				continue
			}
			process := new([]*types.Header)
			index, err := s.getCommonAncestor(head.Number, process)
			if err != nil {
				log.Error("Cannot get common ancestor", "message", err.Error())
				continue
			}
			log.Debug("get common ancestor", "index", index, "count", len(*process))
			blocks := (*process)[:]
			for i := len(blocks) - 1; i >= 0; i-- {
				block := blocks[i]
				s.heads <- block
			}
		case <-s.ctx.Done():
			break
		}
	}
}

// resolveAddresses will resolve the addresses from the address resolver on
// layer one.
func (s *SyncService) resolveAddresses() error {
	if s.ethrpcclient == nil {
		return errors.New("Must initialize eth rpc client first")
	}
	resolver, err := addressmanager.NewLibAddressManager(s.AddressResolverAddress, s.ethrpcclient)
	if err != nil {
		return fmt.Errorf("Cannot create new address manager: %w", err)
	}
	// TODO(mark): using the correct block height is a consensus critical thing.
	// Be sure to use the correct height by setting BlockNumber in the context
	opts := bind.CallOpts{Context: s.ctx}

	s.CanonicalTransactionChainAddress, err = resolver.GetAddress(&opts, "OVM_CanonicalTransactionChain")
	if err != nil {
		return fmt.Errorf("Cannot resolve canonical transaction chain: %w", err)
	}
	s.SequencerDecompressionAddress, err = resolver.GetAddress(&opts, "OVM_SequencerDecompression")
	if err != nil {
		return fmt.Errorf("Cannot resolve sequencer decompression: %w", err)
	}
	s.StateCommitmentChainAddress, err = resolver.GetAddress(&opts, "OVM_StateCommitmentChain")
	if err != nil {
		return fmt.Errorf("Cannot resolve state commitment chain: %w", err)
	}
	s.ExecutionManagerAddress, err = resolver.GetAddress(&opts, "OVM_ExecutionManager")
	if err != nil {
		return fmt.Errorf("Cannot resolve execution manager: %w", err)
	}
	return nil
}

// bindContracts will create bindings for the contracts on layer one
func (s *SyncService) bindContracts() error {
	if s.ethrpcclient == nil {
		return errors.New("Must initialize eth rpc client first")
	}

	var err error
	log.Info("Binding to OVM_CanonicalTransactionChain", "address", s.CanonicalTransactionChainAddress)
	s.ctcFilterer, err = ctc.NewOVMCanonicalTransactionChainFilterer(s.CanonicalTransactionChainAddress, s.ethrpcclient)
	if err != nil {
		return fmt.Errorf("Cannot initialize ctc filterer: %w", err)
	}
	s.ctcCaller, err = ctc.NewOVMCanonicalTransactionChainCaller(s.CanonicalTransactionChainAddress, s.ethrpcclient)
	if err != nil {
		return fmt.Errorf("Cannot initialize ctc caller: %w", err)
	}
	log.Info("Binding to OVM_ExecutionManager", "address", s.ExecutionManagerAddress)
	s.mgrCaller, err = mgr.NewOVMExecutionManagerCaller(s.ExecutionManagerAddress, s.ethrpcclient)
	if err != nil {
		return fmt.Errorf("Cannot initialize execution manager caller: %w", err)
	}
	return nil
}

// setSyncStatus sets the `syncing` field as well as manages the
// lock around adding "remote" transactions to the mempool. The
// remote transactions correspond to transactions from RPC, like
// `eth_sendRawTransaction`. `syncing` should never be set directly
// outside of this function.
func (s *SyncService) setSyncStatus(status bool) {
	log.Info("Setting sync status", "status", status)
	if status {
		s.txpool.LockAddRemote()
	} else {
		s.txpool.UnlockAddRemote()
	}
	s.syncing = status
}

// GetSigningKey returns the public key that is used for signing by the
// syncservice.
func (s *SyncService) GetSigningKey() ecdsa.PublicKey {
	return s.key.PublicKey
}

// IsSyncing returns the syncing status of the syncservice.
func (s *SyncService) IsSyncing() bool {
	return s.syncing
}

// sequencerIngestQueue will ingest transactions from the queue. This
// is only for sequencer mode and will panic if called in verifier mode.
func (s *SyncService) sequencerIngestQueue() {
	if s.verifier {
		panic("Cannot run sequencer ingestion in verifier mode")
	}

	for {
		select {
		case <-s.sequencerIngestTicker.C:
			switch s.syncing {
			case false:
				// Get the tip
				tip, err := s.ethclient.HeaderByNumber(s.ctx, nil)
				if err != nil {
					log.Error("Sequencer ingest queue cannot get tip", "message", err.Error())
					continue
				}
				tipHeight := tip.Number.Uint64()
				// The transactions need to be played in order and there is no
				// guarantee of order when it comes to the txcache iteration, so
				// collect an array of pointers and then sort them by index.
				txs := []*RollupTransaction{}
				s.txCache.Range(func(index uint64, rtx *RollupTransaction) {
					// The transaction has not been executed
					// TODO(mark): possibly add sufficiently old logic
					if !rtx.executed {
						txs = append(txs, rtx)
					}
				})

				// Sort in ascending order
				sort.Sort(RollupTxsByIndex(txs))
				log.Info("Ingesting transactions from L1", "count", len(txs))
				for i := 0; i < len(txs); i++ {
					rtx := txs[i]
					log.Debug("Sequencer ingesting", "local-index", i, "rtx-index", rtx.index)
					// set the timestamp
					s.bc.SetCurrentTimestamp(rtx.timestamp.Unix())
					// The god key needs to sign L1ToL2 transactions
					err := s.maybeReorgAndApplyTx(rtx.index, rtx.tx, true)
					if err != nil {
						log.Error("Sequencer ingest queue transaction failed: %w", err)
					}
					rtx.executed = true
					s.txCache.Store(rtx.index, rtx)
				}
				log.Info("Sequencer Ingest Queue Status", "syncing", s.syncing, "tip-height", tipHeight)
			case true:
				opts := bind.CallOpts{Pending: false, Context: s.ctx}
				totalElements, err := s.ctcCaller.GetTotalElements(&opts)
				// Also check that the chain is synced to the tip
				tip := s.bc.CurrentBlock()
				isAtTip := tip.Number().Uint64() == totalElements.Uint64()

				pending, err := s.ctcCaller.GetNumPendingQueueElements(&opts)
				if pending.Uint64() == 0 && isAtTip {
					s.setSyncStatus(false)
					continue
				}
				// Get the next queue index
				index, err := s.ctcCaller.GetNextQueueIndex(&opts)
				if err != nil {
					log.Error("Cannot get next queue index", "message", err.Error())
					continue
				}
				log.Info("Sequencer Ingest Queue Status", "syncing", s.syncing, "at-tip", isAtTip, "local-tip-height", tip.Number().Uint64(), "next-queue-index", index, "pending-queue-elements", pending.Uint64())
			}
		case <-s.ctx.Done():
			return
		}
	}
}

// LogDoneProcessing reads from the doneProcessing channel
// as to prevent the entire application from stalling. This
// is used for testing, but we could log here too.
func (s *SyncService) LogDoneProcessing() {
	for {
		<-s.doneProcessing
	}
}

// ClearTransactionLoop will clear transactions from the transaction
// cache after they are considered finalized. It currently uses an estimation,
// this could be improved so that the guarantees are better around removing
// exactly when the transactions are finalized.
func (s *SyncService) ClearTransactionLoop() {
	log.Info("Starting transaction clearing loop")
	for {
		select {
		case <-s.clearTransactionsTicker.C:
			tip, err := s.ethclient.HeaderByNumber(s.ctx, nil)
			if err != nil {
				log.Error("Unable to fetch tip in clear transaction loop")
				continue
			}
			if tip.Number == nil {
				log.Error("Unable to fetch tip in clear transaction loop")
				continue
			}
			currentHeight := tip.Number.Uint64()

			count := 0
			s.txCache.Range(func(index uint64, rtx *RollupTransaction) {
				if rtx.blockHeight+s.clearTransactionsAfter > currentHeight {
					log.Debug("Clearing transaction from transaction cache", "hash", rtx.tx.Hash(), "index", index)
					s.txCache.Delete(index)
					count++
				}
			})
			log.Info("SyncService: cleared transactions from cache", "count", count)
		case <-s.ctx.Done():
			return
		}
	}
}

// dialEth1Node will connect to an eth1 node
func (s *SyncService) dialEth1Node() (*rpc.Client, *ethclient.Client, error) {
	rpcClient, err := rpc.Dial(s.eth1HTTPEndpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to connect to eth1 at %s: %w", s.eth1HTTPEndpoint, err)
	}

	client := ethclient.NewClient(rpcClient)
	return rpcClient, client, nil
}

// Stop will close the open channels and cancel the goroutines
// started by this service.
func (s *SyncService) Stop() error {
	defer close(s.heads)
	defer close(s.doneProcessing)

	if s.cancel != nil {
		defer s.cancel()
	}
	return nil
}

func (s *SyncService) Loop() {
	log.Info("Starting Tip processing loop")
	for {
		select {
		case header := <-s.heads:
			if header == nil {
				continue
			}
			blockHeight := header.Number.Uint64()
			eth1data, err := s.ProcessETHBlock(s.ctx, header)
			if err != nil {
				log.Error("Error processing eth block", "message", err.Error(), "height", blockHeight)
				s.doneProcessing <- blockHeight
				// TODO(mark): consider checking the error type here and calling
				// syncHistoricalBlocks in case the eth_subscribe starts to
				// return blocks in the future relative to what is known locally.
				continue
			}
			s.Eth1Data = eth1data
			s.doneProcessing <- blockHeight
		case <-s.ctx.Done():
			break
		}
	}
}

// verifyNetwork ensures that the remote eth1 node is the expected type of node
// based on the chainid and networkid. Log processing should not begin until
// after this check passes.
func (s *SyncService) verifyNetwork() error {
	cid, err := s.ethclient.ChainID(s.ctx)
	if err != nil {
		return fmt.Errorf("Cannot fetch chain id: %w", err)
	}
	if cid.Uint64() != s.eth1ChainId {
		return fmt.Errorf("Received incorrect chain id %d", cid.Uint64())
	}
	nid, err := s.ethclient.NetworkID(s.ctx)
	if err != nil {
		return fmt.Errorf("Cannot fetch network id: %w", err)
	}
	if nid.Uint64() != s.eth1NetworkId {
		return fmt.Errorf("Received incorrect network id %d", nid.Uint64())
	}
	return nil
}

// checkSyncStatus checks the syncing status of the remote eth1 node.
// Log processing should not begin until the remote node is fully synced.
func (s *SyncService) checkSyncStatus() error {
	for {
		syncProg, err := s.ethclient.SyncProgress(s.ctx)
		if err != nil {
			log.Error("Cannot fetch sync progress", "message", err.Error())
			return err
		}

		if syncProg == nil {
			return nil
		}
		log.Info("Ethereum node not fully synced", "current block", syncProg.CurrentBlock)
		time.Sleep(1 * time.Minute)
	}
}

// processHistoricalLogs will sync block by block of the eth1 chain, looking for
// events it can process.
func (s *SyncService) processHistoricalLogs() error {
	errCh := make(chan error)

	go func(c chan error) {
		log.Info("Processing historical logs")
		defer func() { close(c) }()
		for {
			// Get the tip of the chain
			tip, err := s.ethclient.HeaderByNumber(s.ctx, nil)
			if err != nil {
				log.Error("Problem fetching tip for historical log sync")
				time.Sleep(1 * time.Second)
				continue
			}
			// Check to see if the tip is the last processed block height
			tipHeight := tip.Number.Uint64()
			if tipHeight == s.Eth1Data.BlockHeight {
				log.Info("Done fetching historical logs", "height", tipHeight)
				errCh <- nil
			}
			if tipHeight < s.Eth1Data.BlockHeight {
				log.Error("Historical block processing tip is earlier than last processed block height")
				errCh <- fmt.Errorf("Eth1 chain not synced: height %d", tipHeight)
			}

			// Fetch the next header and process it
			header, err := s.ethclient.HeaderByNumber(s.ctx, new(big.Int).SetUint64(s.Eth1Data.BlockHeight+1))
			if err != nil {
				errCh <- fmt.Errorf("Cannot fetch header by number %d: %w", s.Eth1Data.BlockHeight+1, err)
			}
			if header.Number == nil {
				errCh <- fmt.Errorf("Header has nil number")
			}
			headerHeight := header.Number.Uint64()
			headerHash := header.Hash()

			eth1data, err := s.ProcessETHBlock(s.ctx, header)
			if err != nil {
				log.Error("Cannot process block", "message", err.Error(), "height", headerHeight, "hash", headerHash.Hex())
				time.Sleep(1 * time.Second)
				continue
			}
			s.Eth1Data = eth1data
			log.Info("Processed historical block", "height", headerHeight, "hash", headerHash.Hex())
			s.doneProcessing <- headerHeight
		}
	}(errCh)

	select {
	case <-s.ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

// ProcessETHBlock will process all of the logs for a single ethereum block.
func (s *SyncService) ProcessETHBlock(ctx context.Context, header *types.Header) (Eth1Data, error) {
	if header == nil {
		return s.Eth1Data, errors.New("Cannot process nil header")
	}
	s.processingLock.RLock()
	defer s.processingLock.RUnlock()

	blockHeight := header.Number.Uint64()
	blockHash := header.Hash()
	log.Debug("Processing block", "height", blockHeight, "hash", blockHash.Hex())
	// This indicates a reorg on layer 1. Need to delete transactions
	// from the cache that correspond to the block height.
	if blockHeight <= s.Eth1Data.BlockHeight {
		log.Info("Reorganize on eth1 detected, removing transactions from cache", "new height", blockHeight, "old height", s.Eth1Data.BlockHeight, "new hash", header.Hash().Hex())
		count := 0
		s.txCache.Range(func(index uint64, rtx *RollupTransaction) {
			if blockHeight < rtx.blockHeight {
				log.Debug("Clearing transaction from transaction cache", "hash", rtx.tx.Hash(), "index", index)
				s.txCache.Delete(index)
				count++
			}
		})
		log.Info("Reorganize cleared transactions from cache", "count", count)
	}

	// This should never happen and means that historical logs need to be
	// processed.
	if blockHeight > s.Eth1Data.BlockHeight+1 {
		return s.Eth1Data, fmt.Errorf("Unexpected future block at height %d", blockHeight)
	}

	// Create a filter for all logs from the ctc at a specific block hash
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			s.CanonicalTransactionChainAddress,
		},
		// currently unsupported in hardhat
		// see: https://github.com/nomiclabs/hardhat/pull/948/
		//BlockHash: &blockHash
		FromBlock: header.Number,
		ToBlock:   header.Number,
		Topics:    [][]common.Hash{},
	}

	logs, err := s.logClient.FilterLogs(ctx, query)
	if err != nil {
		return s.Eth1Data, fmt.Errorf("Cannot query for logs at block %s: %w", blockHash.Hex(), err)
	}

	// sort the logs by Index
	// TODO: test ByIndex, must be ascending
	sort.Sort(LogsByIndex(logs))
	for _, ethlog := range logs {
		if ethlog.BlockNumber != blockHeight {
			log.Warn("Unexpected block height from log", "got", ethlog.BlockNumber, "expected", blockHeight)
			continue
		}
		// Prevent logs emitted from other contracts from being processed
		if !bytes.Equal(ethlog.Address.Bytes(), s.CanonicalTransactionChainAddress.Bytes()) {
			continue
		}
		if err := s.ProcessLog(ctx, ethlog); err != nil {
			// TODO: reorg out the applied transactions and remove the
			// transactions that were added to the cache so that none are
			// replayed. The same Eth1Data is returned, so it will not be
			// updated. In `processHistoricalLogs`, this will result in the same
			// block being queried. In the `Loop`, the next block should arrive
			// via a notification. Think about good solutions for this.
			return s.Eth1Data, fmt.Errorf("Cannot process log at height %d: %w", blockHeight, err)
		}
	}

	// Write to the database for term persistence
	rawdb.WriteHeadEth1HeaderHash(s.db, header.Hash())
	rawdb.WriteHeadEth1HeaderHeight(s.db, blockHeight)
	s.HeaderCache[blockHeight%2018] = header

	return Eth1Data{
		BlockHash:   blockHash,
		BlockHeight: blockHeight,
	}, nil
}

// GetLastProcessedEth1Data will read the last processed information from the
// database and return it in an Eth1Data struct.
func (s *SyncService) GetLastProcessedEth1Data() Eth1Data {
	hash := rawdb.ReadHeadEth1HeaderHash(s.db)
	height := rawdb.ReadHeadEth1HeaderHeight(s.db)

	return Eth1Data{
		BlockHash:   hash,
		BlockHeight: height,
	}
}

// ProcessLog will process a single log and handle it depending on its source.
// It assumes that the log came from the ctc contract, so be sure to filter out
// other logs before calling this method.
func (s *SyncService) ProcessLog(ctx context.Context, ethlog types.Log) error {
	// This should not happen, but don't trust service providers.
	if len(ethlog.Topics) == 0 {
		return fmt.Errorf("No topics for block %d", ethlog.BlockNumber)
	}
	topic := ethlog.Topics[0].Bytes()
	log.Debug("Processing log", "topic", hexutil.Encode(topic))

	if bytes.Equal(topic, transactionEnqueuedEventSignature) {
		return s.ProcessTransactionEnqueuedLog(ctx, ethlog)
	}
	if bytes.Equal(topic, sequencerBatchAppendedEventSignature) {
		return s.ProcessSequencerBatchAppendedLog(ctx, ethlog)
	}
	if bytes.Equal(topic, queueBatchAppendedEventSignature) {
		return s.ProcessQueueBatchAppendedLog(ctx, ethlog)
	}

	return fmt.Errorf("Unknown log topic %s", hexutil.Encode(topic))
}

func (s *SyncService) ProcessTransactionEnqueuedLog(ctx context.Context, ethlog types.Log) error {
	event, err := s.ctcFilterer.ParseTransactionEnqueued(ethlog)
	if err != nil {
		return fmt.Errorf("Cannot parse transaction enqueued event log: %w", err)
	}

	// Nonce is set by god key at execution time
	// Value and gasPrice are set to 0
	tx := types.NewTransaction(uint64(0), event.Target, big.NewInt(0), event.GasLimit.Uint64(), big.NewInt(0), event.Data, &event.L1TxOrigin, new(big.Int).SetUint64(ethlog.BlockNumber), types.QueueOriginL1ToL2, types.SighashEIP155)
	// Set the index on the transaction so that it can be sorted by index.
	tx.SetIndex(event.QueueIndex.Uint64())

	// Timestamp is used to update the blockchains clocktime
	timestamp := time.Unix(event.Timestamp.Int64(), 0)
	rtx := RollupTransaction{tx: tx, timestamp: timestamp, blockHeight: ethlog.BlockNumber, executed: false, index: event.QueueIndex.Uint64()}
	// In the case of a reorg, the rtx at a certain index can be overwritten
	s.txCache.Store(event.QueueIndex.Uint64(), &rtx)
	log.Debug("Transaction enqueued", "index", event.QueueIndex.Uint64(), "timestamp", timestamp, "l1-blocknumber", ethlog.BlockNumber, "to", event.Target.Hex())

	return nil
}

// ProcessSequencerBatchAppendedLog processes the sequencerbatchappended log
// from the canonical transaction chain contract.
func (s *SyncService) ProcessSequencerBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	log.Debug("Processing sequencer batch appended")
	// TODO(mark): temporary fix to disable sequencer batch append
	if true {
		return nil
	}

	event, err := s.ctcFilterer.ParseSequencerBatchAppended(ethlog)
	if err != nil {
		return fmt.Errorf("Unable to parse sequencer batch appended log data: %w", err)
	}

	tx, pending, err := s.ethclient.TransactionByHash(ctx, ethlog.TxHash)
	if err == ethereum.NotFound {
		return fmt.Errorf("Transaction %s not found: %w", ethlog.TxHash.Hex(), err)
	}
	if err != nil {
		return fmt.Errorf("Cannot fetch transaction %s: %w", ethlog.TxHash.Hex(), err)
	}
	if pending {
		return fmt.Errorf("Transaction %s unexpectedly in mempool", ethlog.TxHash.Hex())
	}

	calldata := tx.Data()
	if len(calldata) < 4 {
		return fmt.Errorf("Unexpected calldata %s", hexutil.Encode(calldata))
	}

	cd := appendSequencerBatchCallData{}
	// Remove the function selector
	err = cd.Decode(bytes.NewReader(calldata[4:]))
	if err != nil {
		return fmt.Errorf("Cannot decode sequencer batch appended calldata: %w", err)
	}

	log.Debug("Decoded chain elements", "count", len(cd.ChainElements))
	for i, element := range cd.ChainElements {
		var tx *types.Transaction
		// The number of queue elements must be the number *after* all of the
		// operations for this to work properly.
		index := event.TotalElements.Uint64() - (uint64(i) + event.NumQueueElements.Uint64())
		// Certain types of transactions require a signature from the god key.
		// Keep track of this so that the god key can sign after reorganizing,
		// to ensure that nonces are correct.
		godKeyShouldSign := false
		// Sequencer transaction
		if element.IsSequenced {
			// Different types of transactions can be included in the canonical
			// transaction chain. The first byte specifies what kind of
			// transaction it is. Parse the data emitted from the log and then
			// build the tx to be played against the evm based on the type
			ctcTx := CTCTransaction{}
			err = ctcTx.Decode(element.TxData)
			if err != nil {
				return fmt.Errorf("Cannot deserialize txdata at index %d: %w", index, err)
			}

			// TODO: QueueOriginSequencer transactions need to include the last
			// L1BlockNumber of a L1ToL2 transaction, not `nil`.
			switch ctcTx.typ {
			case CTCTransactionTypeEOA:
				// The god key needs to sign in this case
				godKeyShouldSign = true
				nonce := uint64(0)
				to := s.SequencerDecompressionAddress
				// TEMP: replacement of s.gasLimit, which is fetched from
				// the contracts, it breaks things
				gasLimit := uint64(8000000)

				tx = types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, big.NewInt(0), element.TxData, nil, nil, types.QueueOriginSequencer, types.SighashEIP155)
				tx.SetIndex(index)
				log.Debug("Deserialized CTC EOA transaction", "index", index, "to", tx.To().Hex(), "data", hexutil.Encode(element.TxData))
			case CTCTransactionTypeEIP155:
				// The signature is deserialized so the god key does not need to
				// sign in this case.
				eip155, ok := ctcTx.tx.(*CTCTxEIP155)
				if !ok {
					return fmt.Errorf("Unexpected type when parsing ctc tx eip155: %T", ctcTx.tx)
				}
				nonce, gasLimit := uint64(eip155.nonce), uint64(eip155.gasLimit)
				to, l1TxOrigin := eip155.target, common.Address{}
				gasPrice := new(big.Int).SetUint64(uint64(eip155.gasPrice))
				data := eip155.data
				tx = types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data, &l1TxOrigin, nil, types.QueueOriginSequencer, types.SighashEIP155)
				tx.SetIndex(index)
				// `WithSignature` accepts:
				// r || s || v where v is normalized to 0 or 1
				tx, err = tx.WithSignature(s.signer, eip155.Signature[:])
				if err != nil {
					return fmt.Errorf("Cannot add signature to eip155 tx: %w", err)
				}
				log.Debug("Deserialized CTC EIP155 transaction", "index", index, "to", tx.To().Hex(), "gasPrice", tx.GasPrice().Uint64(), "gasLimit", tx.Gas())
			default:
				// This should never happen
				return fmt.Errorf("Unknown tx type: %x", element.TxData)
			}
		} else {
			// Queue transaction
			// The god key needs to sign in this case
			godKeyShouldSign = true
			rtx, ok := s.txCache.Load(index)
			if !ok {
				log.Error("Cannot find transaction in transaction cache", "index", index)
				continue
			}
			s.bc.SetCurrentTimestamp(rtx.timestamp.Unix())
			log.Debug("Setting timestamp", "timestamp", rtx.timestamp.Unix())
			rtx.executed = true
			s.txCache.Store(rtx.index, rtx)
		}

		log.Debug("Sequencer batch appended applying tx", "index", index)
		err = s.maybeReorgAndApplyTx(index, tx, godKeyShouldSign)
		if err != nil {
			return fmt.Errorf("Sequencer batch appended error with index %d: %w", index, err)
		}
		log.Info("Sequencer Batch appended success", "index", index, "to", tx.To().Hex(), "god-key-used", godKeyShouldSign)
	}
	return nil
}

// maybeReorg will check to see if the transaction at the index is different
// and then reorg the chain to `index-1` if it is.
func (s *SyncService) maybeReorg(index uint64, tx *types.Transaction) error {
	// Handle the special case of never reorging the genesis block
	if index == 0 {
		return nil
	}
	// Check if there is already a transaction at the index
	if block := s.bc.GetBlockByNumber(index); block != nil {
		// A transaction exists at the current index
		// Sanity check that there is a transaction in the block
		if count := len(block.Transactions()); count != 1 {
			return fmt.Errorf("Unexpected number of transactions in a block %d", count)
		}
		prev := block.Transactions()[0]
		// The transaction hash is not the canonical identifier of a transaction
		// due to nonces coming from the god key. Do an equality check using
		// `to`, `data`, `l1TxOrigin` and `gasLimit`
		if !isCtcTxEqual(tx, prev) {
			log.Info("Different transaction detected, reorganizing", "new", tx.Hash().Hex(), "previous", prev.Hash().Hex())
			// Set the sync status to true. This will grab a lock around
			// the mempool such that transactions will no longer be able to come
			// via RPC.
			s.setSyncStatus(true)
			// Reorganize the chain
			err := s.bc.SetHead(index - 1)
			// TODO: need to iterate through the transactions in the txcache and
			// set `rtx.executed = false` for ones that have a blockheight where:
			// blockheight > index -1
			if err != nil {
				return fmt.Errorf("Cannot reorganize to %d: %w", index-1, err)
			}
		}
	}
	return nil
}

// maybeReorgAndApplyTx will reorg based on the transaction found at the index
// and then maybe sign the transaction if it needs to be signed by the god key
// and then maybe apply the transaction if it is the correct index.
func (s *SyncService) maybeReorgAndApplyTx(index uint64, tx *types.Transaction, godKeyShouldSign bool) error {
	err := s.maybeReorg(index, tx)
	if err != nil {
		return fmt.Errorf("Cannot reorganize before applying tx: %w", err)
	}
	if godKeyShouldSign {
		tx, err = s.signTransaction(tx)
		if err != nil {
			return fmt.Errorf("Cannot sign transaction with god key: %w", err)
		}
	}
	err = s.applyTransaction(tx)
	if err != nil {
		return fmt.Errorf("Cannot apply tx: %w", err)
	}
	return nil
}

// maybeApplyTransaction will look at the tips height and apply the transaction
// if the transaction is at the next index. This allows the codepath to work for
// the verifier as it syncs as well as the sequencer for reorgs.
// This is currently subject to race conditions because the block production
// goes through the miner. Cannot use.
func (s *SyncService) maybeApplyTransaction(index uint64, tx *types.Transaction) error {
	block := s.bc.CurrentBlock()
	// Special case for the transaction at index 0
	blockNumber := block.Number().Uint64()
	if blockNumber+1 == index || index == 0 {
		log.Debug("Can apply transaction", "index", index)
		return s.applyTransaction(tx)
	}
	log.Debug("Skipping application of transaction", "index", index, "hash", tx.Hash().Hex(), "to", tx.To().Hex(), "local-tip-height", blockNumber)
	return nil
}

// signTransaction gets the nonce for the god key, sets the nonce on the
// transaction and then signs the transaction with the god key. This is not safe
// if a reorg happens after this method, so be sure to only use this function
// after reorganizing.
func (s *SyncService) signTransaction(tx *types.Transaction) (*types.Transaction, error) {
	nonce := s.txpool.Nonce(s.address)
	tx.SetNonce(nonce)
	tx, err := types.SignTx(tx, s.signer, &s.key)
	if err != nil {
		return nil, fmt.Errorf("Transaction signing failed: %w", err)
	}
	return tx, nil
}

// ProcessQueueBatchAppendedLog handles the queue batch appended event that is
// emitted from the canonical transaction chain.
func (s *SyncService) ProcessQueueBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	log.Debug("Processing queue batch appended")
	event, err := s.ctcFilterer.ParseQueueBatchAppended(ethlog)
	if err != nil {
		return fmt.Errorf("Unable to parse queue batch appended log data: %w", err)
	}

	start := event.StartingQueueIndex.Uint64()
	end := start + event.NumQueueElements.Uint64()

	for i := start; i < end; i++ {
		rtx, ok := s.txCache.Load(i)
		if !ok {
			log.Error("Cannot find transaction in transaction cache", "index", i)
			continue
		}
		s.bc.SetCurrentTimestamp(rtx.timestamp.Unix())
		// The god key needs to sign in this case
		err = s.maybeReorgAndApplyTx(i, rtx.tx, true)
		if err != nil {
			log.Error("Error applying transaction", "message", err.Error())
			continue
		}
		// TODO: make sure that this mutates the item in the cache and not
		// a copy of the item here.
		rtx.executed = true
		s.txCache.Store(rtx.index, rtx)
	}
	return nil
}

// Adds the transaction to the mempool so that downstream services
// can apply it to the state. This should directly play against
// the state eventually, skipping the mempool.
func (s *SyncService) applyTransaction(tx *types.Transaction) error {
	err := s.txpool.AddLocal(tx)
	if err != nil {
		return fmt.Errorf("Cannot add tx to mempool: %w", err)
	}
	return nil
}
