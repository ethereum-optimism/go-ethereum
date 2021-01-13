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
	"sync/atomic"
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
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const headerCacheSize = 2048

// Interface used for communicating with Ethereum 1 nodes
type EthereumClient interface {
	ChainID(context.Context) (*big.Int, error)
	NetworkID(context.Context) (*big.Int, error)
	SyncProgress(context.Context) (*ethereum.SyncProgress, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
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

// RollupTransaction represents a transaction parsed from L1
type RollupTransaction struct {
	tx          *types.Transaction
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
	transactionEnqueuedEventSignature      = crypto.Keccak256([]byte("TransactionEnqueued(address,address,uint256,bytes,uint256,uint256)"))
	queueBatchAppendedEventSignature       = crypto.Keccak256([]byte("QueueBatchAppended(uint256,uint256,uint256)"))
	sequencerBatchAppendedEventSignature   = crypto.Keccak256([]byte("SequencerBatchAppended(uint256,uint256,uint256)"))
	transactionBatchAppendedEventSignature = crypto.Keccak256([]byte("TransactionBatchAppended(uint256,bytes32,uint256,uint256,bytes)"))
)

// Eth1Data represents the last processed ethereum 1 data
// The sync service updates this as it syncs blocks.
type Eth1Data struct {
	BlockHeight uint64
	BlockHash   common.Hash
}

// LatestL1ToL2 represents the latest blocknumber and timestamp.
// This must be separate than Eth1Data because it is only updated
// each time that there is a L1 to L2 transaction.
type LatestL1ToL2 struct {
	blockNumber uint64
	timestamp   uint64
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
	scope                            event.SubscriptionScope
	txFeed                           event.Feed
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
	LatestL1ToL2                     LatestL1ToL2
	confirmationDepth                uint64
	HeaderCache                      [headerCacheSize]*types.Header
	sequencerIngestTicker            *time.Ticker
	ctcDeployHeight                  *big.Int
	AddressResolverAddress           common.Address
	CanonicalTransactionChainAddress common.Address
	SequencerDecompressionAddress    common.Address
	StateCommitmentChainAddress      common.Address
	L1CrossDomainMessengerAddress    common.Address
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

	if cfg.IsVerifier {
		log.Info("Running in verifier mode")
	}

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
		L1CrossDomainMessengerAddress:    cfg.L1CrossDomainMessengerAddress,
		confirmationDepth:                cfg.Eth1ConfirmationDepth,
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
		clearTransactionsAfter:           (5760 * 18), // 18 days worth of blocks
		clearTransactionsTicker:          time.NewTicker(time.Hour),
		sequencerIngestTicker:            time.NewTicker(15 * time.Second),
		txCache:                          NewTransactionCache(),
		HeaderCache:                      [headerCacheSize]*types.Header{},
	}
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
	log.Info("Initializing Sync Service", "endpoint", s.eth1HTTPEndpoint, "eth1-chainid", s.eth1ChainId, "eth1-networkid", s.eth1NetworkId, "address-resolver", s.AddressResolverAddress, "tx-ingestion-address", s.address, "confirmation-depth", s.confirmationDepth)
	log.Info("Watching topics", "transaction-enqueued", hexutil.Encode(transactionEnqueuedEventSignature), "queue-batch-appened", hexutil.Encode(queueBatchAppendedEventSignature), "sequencer-batch-appended", hexutil.Encode(sequencerBatchAppendedEventSignature))

	// Always initialize syncing to true to start, the sequencer can toggle off
	// syncing while the verifier is always syncing
	s.setSyncStatus(true)

	blockHeight := rawdb.ReadHeadEth1HeaderHeight(s.db)
	blockHash := rawdb.ReadHeadEth1HeaderHash(s.db)
	if blockHeight == 0 {
		if s.ctcDeployHeight == nil {
			return errors.New("Must configure with canonical transaction chain deploy height")
		}
		blockHeight = s.ctcDeployHeight.Uint64()
		// TODO: need to fetch the correct blockHash in this case
	}
	log.Info("Starting Eth1 sync heights", "height", blockHeight, "hash", blockHash.Hex())
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

	// Set the initial values of the `LatestL1BlockNumber`
	// and `LatestL1Timestamp`
	err = s.initializeLatestL1()
	if err != nil {
		return fmt.Errorf("Cannot set latest L1: %w", err)
	}

	go s.LogDoneProcessing()
	// Catch up to the tip of the eth1 chain
	err = s.processHistoricalLogs()
	if err != nil {
		return fmt.Errorf("Cannot process historical logs: %w", err)
	}

	gasLimit, err := s.mgrCaller.GetMaxTransactionGasLimit(&bind.CallOpts{
		Context: s.ctx,
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

// initializeLatestL1 sets the initial values of the `L1BlockNumber`
// and `L1Timestamp` to the deploy height of the Canonical Transaction
// chain if the chain is empty, otherwise set it from the last
// transaction processed.
func (s *SyncService) initializeLatestL1() error {
	block := s.bc.CurrentBlock()
	if block == nil {
		return errors.New("Current block is nil")
	}
	if block == s.bc.Genesis() {
		if s.ctcDeployHeight == nil {
			return errors.New("Must configure with canonical transaction chain deploy height")
		}
		var err error
		block, err = s.ethclient.BlockByNumber(s.ctx, s.ctcDeployHeight)
		if err != nil {
			return fmt.Errorf("Cannot fetch ctc deploy block at height %d", s.ctcDeployHeight)
		}
		s.SetLatestL1Timestamp(block.Time())
		s.SetLatestL1BlockNumber(block.Number().Uint64())
	} else {
		txs := block.Transactions()
		if len(txs) != 1 {
			log.Error("Unexpected number of transactions in block: %d", len(txs))
		}
		if len(txs) > 0 {
			tx := txs[0]
			s.SetLatestL1Timestamp(tx.L1Timestamp())
			s.SetLatestL1BlockNumber(tx.L1BlockNumber().Uint64())
		}
	}
	return nil
}

func (s *SyncService) getCommonAncestor(index *big.Int, list *[]*types.Header) (uint64, error) {
	header, err := s.ethclient.HeaderByNumber(s.ctx, index)
	if err != nil {
		return 0, fmt.Errorf("Cannot fetch header: %w", err)
	}
	number := header.Number.Uint64()
	// Do not allow for reorgs past the deployment height
	// of the contracts.
	if number == s.ctcDeployHeight.Uint64() {
		return number, nil
	}
	cached := s.HeaderCache[number%headerCacheSize]
	if cached != nil && bytes.Equal(header.Hash().Bytes(), cached.Hash().Bytes()) {
		return number, nil
	}
	prevNumber := new(big.Int).SetUint64(number - 1)
	*list = append(*list, header)
	return s.getCommonAncestor(prevNumber, list)
}

func (s *SyncService) pollHead() {
	headTicker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-headTicker.C:
			// Fetch the tip by passing `nil`. We want to consume the
			// blockNumber, but we need to cherry-pick in support for
			// `eth_blockNumber` from upstream geth. For now, just fetch
			// the tip and use the blockNumber from there.
			head, err := s.ethclient.HeaderByNumber(s.ctx, nil)
			if err != nil {
				log.Error("Cannot fetch tip", "height", "tip")
				continue
			}
			// We want to trail the tip by a confirmation number of blocks, so
			// subtract the confirmationDepth from the tip height and fetch the
			// block header that will be consumed.
			blockNumber := head.Number.Sub(head.Number, new(big.Int).SetUint64(s.confirmationDepth))
			head, err = s.ethclient.HeaderByNumber(s.ctx, blockNumber)
			if err != nil {
				log.Error("Cannot fetch tip", "height", blockNumber.Uint64())
				continue
			}
			// The tip is the same, do not ingest the block.
			if bytes.Equal(head.Hash().Bytes(), s.Eth1Data.BlockHash.Bytes()) {
				continue
			}
			// It is possible that multiple blocks have passed since this
			// function last polled, so recursively fetch them.
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

	// For now, only handle the case where syncing is true
	for {
		select {
		case <-s.sequencerIngestTicker.C:
			switch s.syncing {
			case true:
				opts := bind.CallOpts{Pending: false, Context: s.ctx}
				totalElements, err := s.ctcCaller.GetTotalElements(&opts)
				// Also check that the chain is synced to the tip
				tip := s.bc.CurrentBlock()
				isAtTip := tip.Number().Uint64() == totalElements.Uint64()

				pending, err := s.ctcCaller.GetNumPendingQueueElements(&opts)
				// For now always disable sync service
				if true {
					// TODO: Remove this
					// Set all txs found during sync to executed
					s.txCache.Range(func(index uint64, rtx *RollupTransaction) {
						rtx.executed = true
						s.txCache.Store(rtx.index, rtx)
					})
					s.setSyncStatus(false)
					continue
				}
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
				if rtx.executed && rtx.blockHeight+s.clearTransactionsAfter <= currentHeight {
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

	s.scope.Close()

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
			// Update the LatestL1ToL2 info if it is too old.
			// This does not impact the verifier, as all timestamps
			// are set by information that comes in the calldata.
			s.maybeUpdateLatestL1ToL2(header)

			blockHeight := header.Number.Uint64()
			eth1data, err := s.ProcessETHBlock(s.ctx, header)
			if err != nil {
				log.Error("Error processing eth block", "message", err.Error(), "height", blockHeight)
				s.doneProcessing <- blockHeight
				continue
			}
			s.Eth1Data = eth1data
			if !s.verifier {
				s.ApplyLogs(header)
			}
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
		return fmt.Errorf("Received incorrect chain id %d, expected %d", cid.Uint64(), s.eth1ChainId)
	}
	nid, err := s.ethclient.NetworkID(s.ctx)
	if err != nil {
		return fmt.Errorf("Cannot fetch network id: %w", err)
	}
	if nid.Uint64() != s.eth1NetworkId {
		return fmt.Errorf("Received incorrect network id %d, expected %d", nid.Uint64(), s.eth1NetworkId)
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
			tipHeight := tip.Number.Uint64() - headerCacheSize

			// Break when we are up to the header cache size or if the tip
			// number is less than the header cache size. This protects from
			// an undeflow.
			if tipHeight == s.Eth1Data.BlockHeight || tip.Number.Uint64() < headerCacheSize {
				log.Info("Done fetching historical logs", "height", tipHeight)
				errCh <- nil
			}

			if tipHeight < s.Eth1Data.BlockHeight {
				log.Error("Historical block processing tip is earlier than last processed block height")
				errCh <- fmt.Errorf("Eth1 chain not synced: height %d", tipHeight)
			}

			// The above checks prevent `fromBlock` from being
			// greater than the tip
			fromBlock := s.Eth1Data.BlockHeight + 1
			// Use the tip height as the max value
			toBlock := s.Eth1Data.BlockHeight + 1000
			if tipHeight < toBlock {
				toBlock = tipHeight
			}

			query := ethereum.FilterQuery{
				Addresses: []common.Address{
					s.CanonicalTransactionChainAddress,
				},
				FromBlock: new(big.Int).SetUint64(fromBlock),
				ToBlock:   new(big.Int).SetUint64(toBlock),
				Topics:    [][]common.Hash{},
			}

			logs, err := s.logClient.FilterLogs(s.ctx, query)
			if err != nil {
				log.Error("Cannot query logs", "message", err)
				continue
			}
			if len(logs) == 0 {
				height := s.Eth1Data.BlockHeight + 1000
				if tipHeight < height {
					height = tipHeight
				}

				header, err := s.ethclient.HeaderByNumber(s.ctx, new(big.Int).SetUint64(height))
				if err != nil {
					log.Debug("Problem fetching block", "messsage", err)
					continue
				}
				headerHeight := header.Number.Uint64()
				headerHash := header.Hash()

				eth1data, err := s.ProcessETHBlock(s.ctx, header)
				if err != nil {
					log.Error("Cannot process block", "message", err.Error(), "height", headerHeight, "hash", headerHash.Hex())
					continue
				}
				s.Eth1Data = eth1data

				log.Info("Processed historical block", "height", headerHeight, "hash", headerHash.Hex())
				s.doneProcessing <- headerHeight
			} else {
				sort.Sort(LogsByIndex(logs))
				for _, ethlog := range logs {
					// Prevent logs emitted from other contracts from being processed
					if !bytes.Equal(ethlog.Address.Bytes(), s.CanonicalTransactionChainAddress.Bytes()) {
						continue
					}
					if err := s.ProcessLog(s.ctx, ethlog); err != nil {
						log.Error("Cannot process historical log", "message", err)
						continue
					}

					rawdb.WriteHeadEth1HeaderHash(s.db, ethlog.BlockHash)
					rawdb.WriteHeadEth1HeaderHeight(s.db, ethlog.BlockNumber)

					s.Eth1Data = Eth1Data{
						BlockHash:   ethlog.BlockHash,
						BlockHeight: ethlog.BlockNumber,
					}
					log.Info("Processed historical block", "height", ethlog.BlockNumber, "hash", ethlog.BlockHash.Hex())
					s.doneProcessing <- ethlog.BlockNumber
				}
			}
			// Set the last processed header in the cache
			for {
				processed, err := s.ethclient.HeaderByNumber(s.ctx, new(big.Int).SetUint64(toBlock))
				if err == nil {
					s.HeaderCache[processed.Number.Uint64()%headerCacheSize] = processed
					break
				}
			}
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

	// Create a filter for all logs from the ctc at a specific block hash
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			s.CanonicalTransactionChainAddress,
		},
		// currently unsupported in hardhat
		// see: https://github.com/nomiclabs/hardhat/pull/948/
		//BlockHash: &blockHash,
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
	s.HeaderCache[blockHeight%headerCacheSize] = header

	return Eth1Data{
		BlockHash:   blockHash,
		BlockHeight: blockHeight,
	}, nil
}

// ApplyLogs will apply cached transactions from `enqueue` logs.
// This function should only be called in the case of sequencer.
func (s *SyncService) ApplyLogs(tip *types.Header) error {
	// Handle the enqueue'd transactions. This codepath is only useful
	// for the sequencer, as the verifier should only handle transactions
	// from sequencer batch append and queue batch append.
	// The transactions need to be played in order and there is no
	// guarantee of order when it comes to the txcache iteration, so
	// collect an array of pointers and then sort them by index.
	txs := []*RollupTransaction{}
	s.txCache.Range(func(index uint64, rtx *RollupTransaction) {
		// The transaction has not been executed. We know that it is
		// sufficiently old because we only ingest blocks that are
		// sufficiently old.
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
		// The reorg code is no longer used here, after it is better
		// tested we should move back to using it here and remove
		// the verifier check in maybeReorgAndApplyTx
		err := s.applyTransaction(rtx.tx)
		if err != nil {
			log.Error("Sequencer ingest queue transaction failed", "msg", err)
		}
		rtx.executed = true
		s.txCache.Store(rtx.index, rtx)
	}
	return nil
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

// Methods for safely accessing and storing the latest
// L1 blocknumber and timestamp
func (s *SyncService) GetLatestL1Timestamp() uint64 {
	return atomic.LoadUint64(&s.LatestL1ToL2.timestamp)
}

func (s *SyncService) GetLatestL1BlockNumber() uint64 {
	return atomic.LoadUint64(&s.LatestL1ToL2.blockNumber)
}

func (s *SyncService) SetLatestL1Timestamp(ts uint64) {
	atomic.StoreUint64(&s.LatestL1ToL2.timestamp, ts)
}

func (s *SyncService) SetLatestL1BlockNumber(bn uint64) {
	atomic.StoreUint64(&s.LatestL1ToL2.blockNumber, bn)
}

// maybeUpdateLatestL1ToL2 updates the latest L1 block information
// if the timestamp is greater than 5 minutes old. Based on the
// simple equation: now - 5 < prev, where now is the timestamp
// of the L1 tip and prev is the previous latest L1 timestamp.
func (s *SyncService) maybeUpdateLatestL1ToL2(tip *types.Header) {
	prev := time.Unix(int64(s.GetLatestL1Timestamp()), 0)
	now := time.Unix(int64(tip.Time), 0).Add(-5 * time.Minute)
	if now.Before(prev) {
		s.SetLatestL1Timestamp(tip.Time)
		s.SetLatestL1BlockNumber(tip.Number.Uint64())
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
	if bytes.Equal(topic, transactionBatchAppendedEventSignature) {
		return nil
	}

	log.Error("Unknown log topic", "topic", hexutil.Encode(topic), "tx-hash", ethlog.TxHash.Hex())
	return nil
}

func (s *SyncService) ProcessTransactionEnqueuedLog(ctx context.Context, ethlog types.Log) error {
	event, err := s.ctcFilterer.ParseTransactionEnqueued(ethlog)
	if err != nil {
		return fmt.Errorf("Cannot parse transaction enqueued event log: %w", err)
	}

	// Nonce is the queue index
	// Value and gasPrice are set to 0
	nonce := event.QueueIndex.Uint64()
	tx := types.NewTransaction(nonce, event.Target, big.NewInt(0), event.GasLimit.Uint64(), big.NewInt(0), event.Data, &event.L1TxOrigin, new(big.Int).SetUint64(ethlog.BlockNumber), types.QueueOriginL1ToL2, types.SighashEIP155)
	// Set the timestamp in the txmeta
	timestamp := uint64(event.Timestamp.Int64())
	tx.SetL1Timestamp(timestamp)

	rtx := RollupTransaction{tx: tx, blockHeight: ethlog.BlockNumber, executed: false, index: event.QueueIndex.Uint64()}
	// In the case of a reorg, the rtx at a certain index can be overwritten
	s.txCache.Store(event.QueueIndex.Uint64(), &rtx)
	log.Debug("Transaction enqueued", "queue-index", event.QueueIndex.Uint64(), "timestamp", timestamp, "l1-blocknumber", ethlog.BlockNumber, "to", event.Target.Hex())

	// Ensure monotonicity
	latest := s.GetLatestL1Timestamp()
	if timestamp < latest {
		log.Error("Timestamp unexpectedly early", "latest", latest, "new", timestamp)
	} else {
		// Set the timestamp and blocknumber so that transactions from
		// queue origin sequencer can access this information
		s.SetLatestL1Timestamp(timestamp)
		s.SetLatestL1BlockNumber(ethlog.BlockNumber)
	}

	return nil
}

// ProcessSequencerBatchAppendedLog processes the sequencerbatchappended log
// from the canonical transaction chain contract.
func (s *SyncService) ProcessSequencerBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	log.Debug("Processing sequencer batch appended")
	event, err := s.ctcFilterer.ParseSequencerBatchAppended(ethlog)
	if err != nil {
		return fmt.Errorf("Unable to parse sequencer batch appended log data: %w", err)
	}
	log.Debug("Sequencer Batch Appended Event Log", "startingQueueIndex", event.StartingQueueIndex.Uint64(), "numQueueElements", event.NumQueueElements.Uint64(), "totalElements", event.TotalElements.Uint64())

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
	// Keep track of the number of enqueued elements so that the queue index
	// can be calculated in the case of `element.IsSequenced` is false.
	enqueuedCount := uint64(0)
	for i, element := range cd.ChainElements {
		var tx *types.Transaction
		index := (event.TotalElements.Uint64() - uint64(len(cd.ChainElements))) + uint64(i)
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

			switch ctcTx.typ {
			case CTCTransactionTypeEIP155:
				eip155, ok := ctcTx.tx.(*CTCTxEIP155)
				if !ok {
					return fmt.Errorf("Unexpected type when parsing ctc tx eip155: %T", ctcTx.tx)
				}
				nonce, gasLimit := uint64(eip155.nonce), uint64(eip155.gasLimit)
				to := eip155.target

				gasPrice := new(big.Int).SetUint64(uint64(eip155.gasPrice))
				data := eip155.data
				l1BlockNumber := element.BlockNumber
				// Set the L1TxOrigin to `nil`
				if to == (common.Address{}) {
					tx = types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data, nil, l1BlockNumber, types.QueueOriginSequencer)
				} else {
					tx = types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data, nil, l1BlockNumber, types.QueueOriginSequencer, types.SighashEIP155)
				}
				tx.SetIndex(index)
				tx.SetL1Timestamp(element.Timestamp.Uint64())
				// `WithSignature` accepts:
				// r || s || v where v is normalized to 0 or 1
				tx, err = tx.WithSignature(s.signer, eip155.Signature[:])
				if err != nil {
					return fmt.Errorf("Cannot add signature to eip155 tx: %w", err)
				}
				from, err := s.signer.Sender(tx)
				if err != nil {
					from = common.Address{}
					log.Error("Unable to compute from", "signature", hexutil.Encode(eip155.Signature[:]))
				}
				t := "<nil>"
				if tx.To() != nil {
					t = tx.To().Hex()
				}
				log.Debug("Deserialized CTC EIP155 transaction", "index", index, "to", t, "gasPrice", tx.GasPrice().Uint64(), "gasLimit", tx.Gas(), "from", from.Hex(), "sig", hexutil.Encode(eip155.Signature[:]), "hash", tx.Hash().Hex())
			case CTCTransactionTypeEthSign:
				ethsign, ok := ctcTx.tx.(*CTCTxEthSign)
				if !ok {
					return fmt.Errorf("Unexpected type when parsing ctc tx eip155: %T", ctcTx.tx)
				}
				nonce, gasLimit := uint64(ethsign.nonce), uint64(ethsign.gasLimit)
				to := ethsign.target
				gasPrice := new(big.Int).SetUint64(uint64(ethsign.gasPrice))
				data := ethsign.data
				l1BlockNumber := element.BlockNumber
				if to == (common.Address{}) {
					tx = types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data, nil, l1BlockNumber, types.QueueOriginSequencer)
				} else {
					tx = types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, data, nil, l1BlockNumber, types.QueueOriginSequencer, types.SighashEthSign)
				}
				tx.SetIndex(index)
				tx.SetL1Timestamp(element.Timestamp.Uint64())
				// `WithSignature` accepts:
				// r || s || v where v is normalized to 0 or 1
				tx, err = tx.WithSignature(s.signer, ethsign.Signature[:])
				if err != nil {
					return fmt.Errorf("Cannot add signature to ethsign tx: %w", err)
				}
				log.Debug("Deserialized CTC EthSign transaction", "index", index, "to", tx.To().Hex(), "gasPrice", tx.GasPrice().Uint64(), "gasLimit", tx.Gas())
			default:
				// TODO(mark): still need to pass along this transaction and
				// execute it. The `to` should be the sequencer entrypoint,
				// the calldata should be the all the data, max gas limit,
				// gas price of 0.
				log.Info("Unknown tx type", "data", hexutil.Encode(element.TxData))
				continue
			}
		} else {
			// Queue transaction
			queueIndex := event.StartingQueueIndex.Uint64() + enqueuedCount
			enqueuedCount++
			rtx, ok := s.txCache.Load(queueIndex)
			if !ok {
				log.Error("Cannot find transaction in transaction cache", "queue-index", queueIndex)
				continue
			}
			tx = rtx.tx
			tx.SetIndex(index)
			rtx.executed = true
			s.txCache.Store(rtx.index, rtx)
		}

		log.Debug("Sequencer batch appended applying tx", "index", index)
		err = s.maybeReorgAndApplyTx(index, tx)
		if err != nil {
			return fmt.Errorf("Sequencer batch appended error with index %d: %w", index, err)
		}
		log.Info("Sequencer Batch appended success", "index", index)
	}
	return nil
}

// maybeReorg will check to see if the transaction at the index is different
// and then reorg the chain to `index-1` if it is.
func (s *SyncService) maybeReorg(index uint64, tx *types.Transaction) error {
	// Handle the special case of never reorging the genesis block and the off
	// by one case that exists between the CTC and geth 2 state.
	if index == 0 || index == 1 {
		return nil
	}
	// Check if there is already a transaction at the index
	if block := s.bc.GetBlockByNumber(index - 1); block != nil {
		// A transaction exists at the current index
		if count := len(block.Transactions()); count != 1 {
			// Don't return an error here to handle the case of the genesis
			// block not having a transaction, until the genesis tx is included
			log.Debug("Unexpected number of transactions in block", "count", count)
			return nil
		}
		prev := block.Transactions()[0]
		// The transaction hash is not the canonical identifier of a transaction
		// due to nonces technically needing to be incremented.
		// Do an equality check using `to`, `data`, `l1TxOrigin` and `gasLimit`
		if !isCtcTxEqual(tx, prev) {
			log.Info("Different tx detected", "index", index, "new", tx.Hash().Hex(), "previous", prev.Hash().Hex())
		} else {
			log.Info("Same tx detected", "index", index)
		}
	}
	return nil
}

// maybeReorgAndApplyTx will reorg based on the transaction found at the index
// and then maybe apply the transaction if it is the correct index.
func (s *SyncService) maybeReorgAndApplyTx(index uint64, tx *types.Transaction) error {
	err := s.maybeReorg(index, tx)
	if err != nil {
		return fmt.Errorf("Cannot reorganize before applying tx: %w", err)
	}
	// Only apply transactions in the case where it is the verifier.
	// This is here so that we can observe the behavior of `maybeReorg`
	// to make sure that it is operating as expected.
	if s.verifier {
		err = s.applyTransaction(tx)
		if err != nil {
			return fmt.Errorf("Cannot apply tx: %w", err)
		}
	}
	return nil
}

// ProcessQueueBatchAppendedLog handles the queue batch appended event that is
// emitted from the canonical transaction chain.
func (s *SyncService) ProcessQueueBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	log.Debug("Processing queue batch appended")
	event, err := s.ctcFilterer.ParseQueueBatchAppended(ethlog)
	if err != nil {
		return fmt.Errorf("Unable to parse queue batch appended log data: %w", err)
	}
	log.Debug("Queue Batch Appended Event Log", "startingQueueIndex", event.StartingQueueIndex.Uint64(), "numQueueElements", event.NumQueueElements.Uint64(), "totalElements", event.TotalElements.Uint64())

	// Disable queue batch appended for minnet
	if true {
		log.Debug("Queue batch append disabled")
		return nil
	}

	start := event.StartingQueueIndex.Uint64()
	end := start + event.NumQueueElements.Uint64()
	for i := start; i < end; i++ {
		rtx, ok := s.txCache.Load(i)
		if !ok {
			log.Error("Cannot find transaction in transaction cache", "index", i)
			continue
		}
		err = s.maybeReorgAndApplyTx(i, rtx.tx)
		if err != nil {
			log.Error("Error applying transaction", "message", err.Error())
			continue
		}
		rtx.executed = true
		s.txCache.Store(rtx.index, rtx)
	}
	return nil
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (s *SyncService) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return s.scope.Track(s.txFeed.Subscribe(ch))
}

// SetL1Head resets the Eth1Data
// and the LatestL1BlockNumber and LatestL1Timestamp
func (s *SyncService) SetL1Head(number uint64) error {
	header, err := s.ethclient.HeaderByNumber(s.ctx, new(big.Int).SetUint64(number))
	if err != nil {
		return fmt.Errorf("Cannot fetch block in SetL1Head: %w", err)
	}

	// Reset the header cache
	for i := 0; i < len(s.HeaderCache); i++ {
		s.HeaderCache[i] = nil
	}

	// Reset the last synced L1 heights
	rawdb.WriteHeadEth1HeaderHash(s.db, header.Hash())
	rawdb.WriteHeadEth1HeaderHeight(s.db, header.Number.Uint64())
	s.HeaderCache[number%headerCacheSize] = header

	s.Eth1Data = Eth1Data{
		BlockHeight: header.Number.Uint64(),
		BlockHash:   header.Hash(),
	}
	return nil
}

// Adds the transaction to the mempool so that downstream services
// can apply it to the state. This should directly play against
// the state eventually, skipping the mempool.
func (s *SyncService) applyTransaction(tx *types.Transaction) error {
	err := s.txpool.ValidateTx(tx)
	// The sequencer needs to prevent transactions that fail the mempool
	// checks. The verifier needs to play the transactions no matter what.
	if !s.verifier {
		qo := tx.QueueOrigin()
		if err != nil && qo.Uint64() == uint64(types.QueueOriginSequencer) {
			return fmt.Errorf("invalid transaction: %w", err)
		}
	}
	txs := types.Transactions{tx}
	s.txFeed.Send(core.NewTxsEvent{Txs: txs})
	return nil
}

func (s *SyncService) ApplyTransaction(tx *types.Transaction) error {
	return s.applyTransaction(tx)
}
