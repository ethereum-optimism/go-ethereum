package rollup

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ctc "github.com/ethereum/go-ethereum/contracts/canonicaltransactionchain"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// --rollup.ethereumrpc
// do you need websocket connection for eth subscribe?

type EthereumClient interface {
	ChainID() *big.Int
	NetworkID() *big.Int
	SyncProgress() *ethereum.SyncProgress
	HeaderByNumber() *types.Header
	TransactionByHash() *types.Transaction
}

// Consider adding a processed bool for sanity check
type RollupTransaction struct {
	tx          *types.Transaction
	timestamp   time.Time
	blockHeight uint64
}

// Implement the Sort interface for []types.Log by Index
type ByIndex []types.Log

func (l ByIndex) Len() int           { return len(l) }
func (l ByIndex) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ByIndex) Less(i, j int) bool { return l[i].Index < l[i].Index }

type ctcBatchContext struct {
	NumSequencedTransactions       *big.Int
	NumSubsequentQueueTransactions *big.Int
	Timestamp                      *big.Int
	BlockNumber                    *big.Int
}

type chainElement struct {
	IsSequenced bool
	Timestamp   *big.Int // Origin Sequencer
	BlockNumber *big.Int // Origin Sequencer
	TxData      []byte   // Origin Sequencer
}

type appendSequencerBatchCallData struct {
	ShouldStartAtBatch    *big.Int
	TotalElementsToAppend *big.Int
	Contexts              []ctcBatchContext
	ChainElements         []chainElement
}

func (c *ctcBatchContext) Encode(w io.Writer) error {
	elements := [][]byte{
		common.LeftPadBytes(c.NumSequencedTransactions.Bytes(), 3),
		common.LeftPadBytes(c.NumSubsequentQueueTransactions.Bytes(), 3),
		common.LeftPadBytes(c.Timestamp.Bytes(), 5),
		common.LeftPadBytes(c.BlockNumber.Bytes(), 5),
	}
	for _, element := range elements {
		_, err := w.Write(element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ctcBatchContext) Len() int {
	return 3 + 3 + 5 + 5
}

func (c *ctcBatchContext) Decode(r io.ReaderAt) error {
	offset := int64(0)
	elements := [][]byte{
		make([]byte, 3),
		make([]byte, 3),
		make([]byte, 5),
		make([]byte, 5),
	}

	for i, element := range elements {
		sr := io.NewSectionReader(r, offset, int64(len(element)))
		off, err := sr.Read(element)
		if err != nil {
			return err
		}

		switch i {
		case 0:
			c.NumSequencedTransactions = new(big.Int).SetBytes(element)
		case 1:
			c.NumSubsequentQueueTransactions = new(big.Int).SetBytes(element)
		case 2:
			c.Timestamp = new(big.Int).SetBytes(element)
		case 3:
			c.BlockNumber = new(big.Int).SetBytes(element)
		}
		offset += int64(off)
	}

	return nil
}

func (c *appendSequencerBatchCallData) Encode(w io.Writer) error {
	contexts := new(bytes.Buffer)
	for _, context := range c.Contexts {
		buf := new(bytes.Buffer)
		err := context.Encode(buf)
		if err != nil {
			return err
		}
		contexts.Write(buf.Bytes())
	}
	transactions := new(bytes.Buffer)
	for _, el := range c.ChainElements {
		if !el.IsSequenced {
			continue
		}
		header := make([]byte, 0, 4)
		buf := bytes.NewBuffer(header)
		err := binary.Write(buf, binary.BigEndian, uint32(len(el.TxData)))
		if err != nil {
			return err
		}
		_ = buf.Next(1) // Move forward a byte
		transactions.Write(buf.Bytes())
		transactions.Write(el.TxData)
	}
	elements := [][]byte{
		common.LeftPadBytes(c.ShouldStartAtBatch.Bytes(), 5),
		common.LeftPadBytes(c.TotalElementsToAppend.Bytes(), 3),
		common.LeftPadBytes(new(big.Int).SetUint64(uint64(len(c.Contexts))).Bytes(), 3),
		contexts.Bytes(),
		transactions.Bytes(),
	}
	for _, element := range elements {
		_, err := w.Write(element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *appendSequencerBatchCallData) Decode(r io.ReaderAt) error {
	offset := int64(0)
	elements := [][]byte{
		make([]byte, 5),
		make([]byte, 3),
		make([]byte, 3), // context count
	}
	ctxCount := new(big.Int)
	for i, element := range elements {
		sr := io.NewSectionReader(r, offset, int64(len(element)))
		off, err := sr.Read(element)
		if err != nil {
			return err
		}

		switch i {
		case 0:
			a.ShouldStartAtBatch = new(big.Int).SetBytes(element)
		case 1:
			a.TotalElementsToAppend = new(big.Int).SetBytes(element)
		case 2:
			ctxCount.SetBytes(element)
		}
		offset += int64(off)
	}

	a.Contexts = make([]ctcBatchContext, ctxCount.Uint64())
	for i := uint64(0); i < ctxCount.Uint64(); i++ {
		batchCtx := ctcBatchContext{}
		sr := io.NewSectionReader(r, offset, offset+int64(batchCtx.Len()))
		err := batchCtx.Decode(sr)
		if err != nil {
			return fmt.Errorf("Cannot decode batch context: %w", err)
		}
		a.Contexts[i] = batchCtx
		offset += int64(batchCtx.Len())
	}

	txCount := uint64(0)
	for _, ctx := range a.Contexts {
		txCount += ctx.NumSequencedTransactions.Uint64()
		txCount += ctx.NumSubsequentQueueTransactions.Uint64()
	}

	if txCount != a.TotalElementsToAppend.Uint64() {
		return errors.New("Incorrect number of elements")
	}

	a.ChainElements = []chainElement{}
	for _, ctx := range a.Contexts {
		timestamp := ctx.Timestamp
		blockNumber := ctx.BlockNumber
		for i := uint64(0); i < ctx.NumSequencedTransactions.Uint64(); i++ {
			header := make([]byte, 3)
			sr := io.NewSectionReader(r, offset, offset+3)
			off, err := sr.Read(header)

			if err != nil {
				return fmt.Errorf("Cannot read tx header: %w", err)
			}
			offset += int64(off)

			var sizeHi uint16
			var sizeLo uint8
			hr := bytes.NewReader(header)
			err = binary.Read(hr, binary.BigEndian, &sizeHi)
			if err != nil {
				return fmt.Errorf("Cannot read tx header hi bits: %w", err)
			}
			err = binary.Read(hr, binary.BigEndian, &sizeLo)
			if err != nil {
				return fmt.Errorf("Cannot read tx header lo bits: %w", err)
			}

			size := (sizeHi << 8) | uint16(sizeLo)
			tx := make([]byte, size)
			tsr := io.NewSectionReader(r, offset, offset+int64(size))
			off, err = tsr.Read(tx)
			if err != nil {
				return fmt.Errorf("Cannot read tx: %w", err)
			}
			offset += int64(off)

			element := chainElement{
				IsSequenced: true,
				Timestamp:   timestamp,
				BlockNumber: blockNumber,
				TxData:      tx,
			}
			a.ChainElements = append(a.ChainElements, element)
		}

		for i := uint64(0); i < ctx.NumSubsequentQueueTransactions.Uint64(); i++ {
			element := chainElement{
				IsSequenced: false,
				Timestamp:   nil,
				BlockNumber: nil,
				TxData:      []byte{},
			}
			a.ChainElements = append(a.ChainElements, element)
		}
	}
	return nil
}

// TODO: double check these signatures
var (
	transactionEnqueuedEventSignature    = crypto.Keccak256([]byte("TransactionEnqueued(address,address,uint256,bytes,uint256,uint256)"))
	queueBatchAppendedEventSignature     = crypto.Keccak256([]byte("QueueBatchAppended(uint256,uint256,uint256)"))
	sequencerBatchAppendedEventSignature = crypto.Keccak256([]byte("SequencerBatchAppended(uint256,uint256,uint256)"))
)

// This needs to be indexed
type latestEth1Data struct {
	BlockHeight        uint64
	BlockHash          common.Hash
	LastRequestedBlock uint64
}

type SyncService struct {
	ctx    context.Context
	cancel context.CancelFunc

	processingLock sync.RWMutex

	db ethdb.Database

	// this needs to be an interface for testing purposes
	ctcFilterer *ctc.OVMCanonicalTransactionChainFilterer
	// TODO: investigate this one
	//ctcFilterer    *bind.ContractFilterer
	ctcABI         abi.ABI
	l1ToL2Enqueued sync.Map
	// turn ethclient into an interface for testing purposes
	ethclient    *ethclient.Client
	httpEndpoint string

	eth1ChainID   big.Int
	eth1NetworkID big.Int

	// might not need backend
	backend    bind.ContractBackend
	httpLogger bind.ContractFilterer

	txpool *core.TxPool
	bc     *core.BlockChain

	clearTransactionsTicker *time.Ticker
	clearTransactionsAfter  uint64

	heads            chan *types.Header
	headSubscription ethereum.Subscription

	signer  types.Signer
	key     ecdsa.PrivateKey
	address common.Address

	latestEth1Data                latestEth1Data
	StateCommitmentChainAddress   common.Address
	L1ToL2TransactionQueueAddress common.Address
}

// Testing strategy
// set up mock Subscription to send Headers over
// set up mock httpLogger
// set up mock ethclient
// send headers over them
// assert that state is correct

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
	parsed, err := abi.JSON(strings.NewReader(ctc.OVMCanonicalTransactionChainABI))
	if err != nil {
		return nil, err
	}

	blockHeight := rawdb.ReadHeadEth1HeightKey(db)
	if blockHeight == 0 {
		blockHeight = cfg.CanonicalTransactionChainDeployHeight.Uint64() - 1
	}
	blockHash := rawdb.ReadHeadEth1HeaderHash(db)

	eth1Data := latestEth1Data{
		BlockHeight: blockHeight,
		BlockHash:   blockHash,
	}

	chainID := bc.Config().ChainID

	service := SyncService{
		ctx:                           ctx,
		cancel:                        cancel,
		heads:                         make(chan *types.Header),
		httpEndpoint:                  cfg.httpEndpoint,
		StateCommitmentChainAddress:   cfg.StateCommitmentChainAddress,
		L1ToL2TransactionQueueAddress: cfg.L1ToL2TransactionQueueAddress,
		signer:                        types.NewOVMSigner(chainID),
		key:                           *cfg.TxIngestionSignerKey,
		address:                       address,
		txpool:                        txpool,
		bc:                            bc,
		ctcABI:                        parsed,
		latestEth1Data:                eth1Data,
		eth1ChainID:                   cfg.Eth1ChainID,
		eth1NetworkID:                 cfg.Eth1NetworkID,
		db:                            db,
		clearTransactionsAfter:        (5760 * 15), // 15 days worth of blocks
		clearTransactionsTicker:       time.NewTicker(time.Hour),
	}

	return &service, nil
}

func (s *SyncService) Start() error {
	// See if rpcClient should be used
	_, client, err := s.dialEth1Node()
	s.ethclient = client

	err = s.verifyNetwork()
	if err != nil {
		return err
	}

	ctcFilterer, err := ctc.NewOVMCanonicalTransactionChainFilterer(s.StateCommitmentChainAddress, client)
	if err != nil {
		return err
	}
	s.ctcFilterer = ctcFilterer

	err = s.checkSyncStatus()
	if err != nil {
		return err
	}
	err = s.processHistoricalLogs()
	if err != nil {
		return err
	}

	sub, err := client.SubscribeNewHead(s.ctx, s.heads)
	s.headSubscription = sub

	go s.Loop()
	go s.ClearTransactionLoop()

	return nil
}

func (s *SyncService) ClearTransactionLoop() {
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

			s.l1ToL2Enqueued.Range(func(key interface{}, value interface{}) bool {
				rtx, ok := value.(RollupTransaction)
				if !ok {
					log.Error("Unexpected value in transaction cache", "type", fmt.Sprintf("%T", value))
					return true
				}
				if rtx.blockHeight+s.clearTransactionsAfter > currentHeight {
					index, ok := key.(uint64)
					if !ok {
						log.Error("Unexpected key type in transaction cache", "type", fmt.Sprintf("%T", key))
						return true
					}
					log.Debug("Clearing transaction from transaction cache", "hash", rtx.tx.Hash(), "index", index)
					s.l1ToL2Enqueued.Delete(key)
				}
				return true
			})
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *SyncService) dialEth1Node() (*rpc.Client, *ethclient.Client, error) {
	connErrCh := make(chan error, 1)
	defer close(connErrCh)

	var rpcClient *rpc.Client
	var err error

	go func() {
		retries := 0
		for {
			rpcClient, err = rpc.Dial(s.httpEndpoint)
			if err != nil {
				log.Error("Error connecting to Eth1", "endpoint", s.httpEndpoint)
				if retries > 10 {
					connErrCh <- err
					return
				}
				retries++
				select {
				case <-s.ctx.Done():
					break
				case <-time.After(time.Second):
					continue
				}
			}
			connErrCh <- err // sending `nil`
		}
	}()

	select {
	case err = <-connErrCh:
		break
	case <-s.ctx.Done():
		return nil, nil, errors.New("Cancelled connection to Eth1")
	}

	if err != nil {
		return nil, nil, errors.New("Connection to Eth1 timed out")
	}

	client := ethclient.NewClient(rpcClient)
	return rpcClient, client, nil
}

func (s *SyncService) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}

	if s.headSubscription != nil {
		defer s.headSubscription.Unsubscribe()
	}

	// close other Loop

	return nil
}

func (s *SyncService) Loop() {
	for {
		select {
		case header := <-s.heads:
			blockHeight := header.Number.Uint64()
			s.latestEth1Data.LastRequestedBlock = blockHeight
			s.ProcessETHBlock(s.ctx, header)
			rawdb.WriteHeadEth1HeaderHash(s.db, header.Hash())
			rawdb.WriteHeadEth1HeightKey(s.db, blockHeight)
		case <-s.ctx.Done():
			break
		}
	}
}

func (s *SyncService) verifyNetwork() error {
	cid, err := s.ethclient.ChainID(s.ctx)
	if err != nil {
		return err
	}
	if cid.Cmp(&s.eth1ChainID) != 0 {
		return fmt.Errorf("Received incorrect chain id %d", cid.Uint64())
	}

	nid, err := s.ethclient.NetworkID(s.ctx)
	if err != nil {
		return err
	}
	if nid.Cmp(&s.eth1NetworkID) != 0 {
		return fmt.Errorf("Received incorrect network id %d", nid.Uint64())
	}
	return nil
}

func (s *SyncService) checkSyncStatus() error {
	for {
		syncProg, err := s.ethclient.SyncProgress(s.ctx)
		if err != nil {
			log.Error("Cannot fetch sync progress", "message", err.Error())
			return err
		}

		if syncProg == nil {
			return nil
		} else {
			log.Info("Ethereum node not fully synced", "current block", syncProg.CurrentBlock)
			time.Sleep(1 * time.Minute)
		}
	}
}

func (s *SyncService) processHistoricalLogs() error {
	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
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
			if tipHeight == s.latestEth1Data.BlockHeight {
				errCh <- nil
			}
			if tipHeight < s.latestEth1Data.BlockHeight {
				log.Error("Historical block processing tip is earlier than last processed block height")
				errCh <- fmt.Errorf("Eth1 chain not synced")
			}

			// Fetch the next header and process it
			header, err := s.ethclient.HeaderByNumber(s.ctx, new(big.Int).SetUint64(s.latestEth1Data.BlockHeight+1))
			headerHeight := header.Number.Uint64()
			headerHash := header.Hash()
			s.latestEth1Data.LastRequestedBlock = headerHeight

			err = s.ProcessETHBlock(s.ctx, header)
			if err != nil {
				log.Error("Cannot process block", "message", err.Error(), "height", headerHeight, "hash", headerHash.Hex())
				time.Sleep(1 * time.Second)
				continue
			}
			log.Info("Processed historical block", "height", headerHeight, "hash", headerHash.Hex())
		}
	}()

	select {
	case <-s.ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *SyncService) ProcessETHBlock(ctx context.Context, header *types.Header) error {
	// Roll back the chain when the header is old
	blockHeight := header.Number.Uint64()
	blockHash := header.Hash()

	if blockHeight <= s.latestEth1Data.BlockHeight {
		// this is a reorg on layer 1
		// need to delete queued transactions from old block
	}

	if blockHeight > s.latestEth1Data.BlockHeight+1 {
		// TODO: instead of returning error, call processHistoricalLogs
		return fmt.Errorf("Unexpected future block at height %d", blockHeight)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			s.StateCommitmentChainAddress,
		},
		BlockHash: &blockHash,
	}

	// TODO: need to be able to mock httpLogger...
	// TODO: make sure httpLogger is set correctly
	logs, err := s.httpLogger.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	// sort the logs by Index
	// TODO: test ByIndex, must be ascending
	sort.Sort(ByIndex(logs))
	for _, ethlog := range logs {
		if ethlog.BlockNumber != blockHeight {
			log.Warn("Unexpected block height from log", "got", ethlog.BlockNumber, "expected", blockHeight)
			continue
		}
		if !bytes.Equal(ethlog.Address.Bytes(), s.StateCommitmentChainAddress.Bytes()) {
			continue
		}
		// If an error happens here, then the latestEth1Data will not be updated
		if err := s.ProcessLog(ctx, ethlog); err != nil {
			return err
		}
	}

	// Set the last seen information
	s.latestEth1Data.BlockHeight = blockHeight
	s.latestEth1Data.BlockHash = blockHash

	return nil
}

func (s *SyncService) ProcessLog(ctx context.Context, ethlog types.Log) error {
	// defer catchPanic()
	s.processingLock.RLock()
	defer s.processingLock.RUnlock()

	if len(ethlog.Topics) == 0 {
		return fmt.Errorf("No logs for block %d", ethlog.BlockNumber)
	}
	topic := ethlog.Topics[0].Bytes()

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
		// add context to error
		return err
	}

	// Nonce is set by god key at execution time
	// Value and gasPrice are set to 0
	// nil is the txid (unused)
	tx := types.NewTransaction(uint64(0), event.Target, big.NewInt(0), event.GasLimit.Uint64(), big.NewInt(0), event.Data, &event.L1TxOrigin, nil, types.QueueOriginL1ToL2, types.SighashEIP155)

	// Timestamp is used to update the blockchains clocktime
	timestamp := time.Unix(event.Timestamp.Int64(), 0)
	rtx := RollupTransaction{tx: tx, timestamp: timestamp, blockHeight: ethlog.BlockNumber}
	s.l1ToL2Enqueued.Store(event.QueueIndex.Uint64(), rtx)

	return nil
}

func (s *SyncService) ProcessSequencerBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	event, err := s.ctcFilterer.ParseSequencerBatchAppended(ethlog)
	if err != nil {
		return err
	}

	tx, pending, err := s.ethclient.TransactionByHash(ctx, ethlog.TxHash)
	if err == ethereum.NotFound {
		log.Error("Transaction not found", "hash", ethlog.TxHash.Hex())
		return err
	}
	if err != nil {
		log.Error("Cannot fetch transaction", "hash", ethlog.TxHash.Hex())
		return err
	}
	if pending {
		log.Error("Transaction unexpectedly in mempool", "hash", ethlog.TxHash.Hex())
		return err
	}

	cd := appendSequencerBatchCallData{}
	err = cd.Decode(bytes.NewReader(tx.Data()))
	if err != nil {
		return err
	}

	// event.StartingQueueIndex
	// event.NumQueueElements

	// event.TotalElements
	for i, element := range cd.ChainElements {
		var tx *types.Transaction
		index := event.TotalElements.Uint64() - (uint64(i) + event.NumQueueElements.Uint64())
		// Sequencer transaction
		if element.IsSequenced {
			nonce := uint64(0)
			to := common.Address{} // sequencerDecompressionAddress
			gasLimit := uint64(0)  // max gas limit
			l1TxOrigin := common.Address{}
			tx = types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, big.NewInt(0), element.TxData, &l1TxOrigin, nil, types.QueueOriginSequencer, types.SighashEIP155)
			// Check if there is already a transaction at the index
			if block := s.bc.GetBlockByNumber(index); block != nil {
				// Check if the transaction is different
				// this check will not work, need to check the items in the tx
				// instead.
				if included := block.Transaction(tx.Hash()); included == nil {
					previous := block.Transactions()[0]
					log.Info("Different transaction detected, reorganizing", "new", tx.Hash().Hex(), "previous", previous.Hash().Hex())
					err := s.bc.SetHead(index - 1)
					if err != nil {
						log.Error("Cannot reorganize")
						return err
					}
				}
			}
		} else {
			// Queue transaction
			result, ok := s.l1ToL2Enqueued.Load(index)
			if !ok {
				log.Error("Cannot find transaction in transaction cache", "index", index)
				continue
			}
			rtx, ok := result.(RollupTransaction)
			if !ok {
				log.Error("Incorrect type in transaction cache")
				continue
			}
			tx = rtx.tx
		}

		s.applyTransaction(tx)
	}

	return nil
}

func (s *SyncService) ProcessQueueBatchAppendedLog(ctx context.Context, ethlog types.Log) error {
	event, err := s.ctcFilterer.ParseQueueBatchAppended(ethlog)
	if err != nil {
		return err
	}

	start := event.StartingQueueIndex.Uint64()
	end := start + event.NumQueueElements.Uint64()

	for i := start; i < end; i++ {
		result, ok := s.l1ToL2Enqueued.Load(i)
		if !ok {
			// This is a very bad error, how to recover?
			log.Error("Unknown transaction")
			continue
		}
		rtx, ok := result.(RollupTransaction)
		if !ok {
			log.Error("Incorrect type in map")
			continue
		}

		tx := rtx.tx
		nonce := s.txpool.Nonce(s.address)
		tx.SetNonce(nonce)
		tx, err := types.SignTx(tx, s.signer, &s.key)
		if err != nil {
			log.Error("Error signing transaction")
			continue
		}

		s.bc.SetCurrentTimestamp(rtx.timestamp.Unix())

		// the reorg logic lives inside applyTransaction?
		err = s.applyTransaction(tx)
		if err != nil {
			log.Error("Error applying transaction")
			continue
		}
	}
	return nil
}

func (s *SyncService) applyTransaction(tx *types.Transaction) error {
	return s.txpool.AddLocal(tx)
}
