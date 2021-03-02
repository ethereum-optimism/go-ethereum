package rollup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
)

// OVMContext represents the blocknumber and timestamp
// that exist during L2 execution
type OVMContext struct {
	blockNumber uint64
	timestamp   uint64
}

// SyncService implements the verifier functionality as well as the reorg
// protection for the sequencer.
type SyncService struct {
	ctx                       context.Context
	cancel                    context.CancelFunc
	verifier                  bool
	db                        ethdb.Database
	scope                     event.SubscriptionScope
	txFeed                    event.Feed
	txLock                    sync.Mutex
	enable                    bool
	eth1ChainId               uint64
	bc                        *core.BlockChain
	txpool                    *core.TxPool
	client                    RollupClient
	syncing                   atomic.Value
	OVMContext                OVMContext
	confirmationDepth         uint64
	pollInterval              time.Duration
	timestampRefreshThreshold time.Duration
}

// NewSyncService returns an initialized sync service
func NewSyncService(ctx context.Context, cfg Config, txpool *core.TxPool, bc *core.BlockChain, db ethdb.Database) (*SyncService, error) {
	if bc == nil {
		return nil, errors.New("Must pass BlockChain to SyncService")
	}

	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // satisfy govet

	if cfg.IsVerifier {
		log.Info("Running in verifier mode")
	} else {
		log.Info("Running in sequencer mode")
	}

	pollInterval := cfg.PollInterval
	if pollInterval == 0 {
		log.Info("Sanitizing poll interval to 15 seconds")
		pollInterval = time.Second * 15
	}
	timestampRefreshThreshold := cfg.TimestampRefreshThreshold
	if timestampRefreshThreshold == 0 {
		log.Info("Sanitizing timestamp refresh threshold to 15 minutes")
		timestampRefreshThreshold = time.Minute * 15
	}

	// Layer 2 chainid
	chainID := bc.Config().ChainID
	if chainID == nil {
		return nil, errors.New("Must configure with chain id")
	}
	// Initialize the rollup client
	client := NewClient(cfg.RollupClientHttp, chainID)
	log.Info("Configured rollup client", "url", cfg.RollupClientHttp, "chain-id", chainID.Uint64(), "ctc-deploy-height", cfg.CanonicalTransactionChainDeployHeight)
	service := SyncService{
		ctx:                       ctx,
		cancel:                    cancel,
		verifier:                  cfg.IsVerifier,
		enable:                    cfg.Eth1SyncServiceEnable,
		confirmationDepth:         cfg.Eth1ConfirmationDepth,
		syncing:                   atomic.Value{},
		bc:                        bc,
		txpool:                    txpool,
		eth1ChainId:               cfg.Eth1ChainId,
		client:                    client,
		db:                        db,
		pollInterval:              pollInterval,
		timestampRefreshThreshold: timestampRefreshThreshold,
	}

	// Initial sync service setup if it is enabled. This code depends on
	// a remote server that indexes the layer one contracts. Place this
	// code behind this if statement so that this can run without the
	// requirement of the remote server being up.
	if service.enable {
		// Ensure that the rollup client can connect to a remote server
		// before starting.
		err := service.ensureClient()
		if err != nil {
			return nil, fmt.Errorf("Rollup client unable to connect: %w", err)
		}

		// Ensure that the remote is still not syncing
		for {
			status, err := service.client.SyncStatus()
			if err != nil {
				log.Error("Cannot get sync status")
				continue
			}
			if !status.Syncing {
				break
			}
			log.Info("Still syncing", "index", status.CurrentTransactionIndex, "tip", status.HighestKnownTransactionIndex)
			time.Sleep(10 * time.Second)
		}

		// Initialize the latest L1 data here to make sure that
		// it happens before the RPC endpoints open up
		// Only do it if the sync service is enabled so that this
		// can be ran without needing to have a configured client.
		err = service.initializeLatestL1(cfg.CanonicalTransactionChainDeployHeight)
		if err != nil {
			return nil, fmt.Errorf("Cannot initialize latest L1 data: %w", err)
		}

		bn := service.GetLatestL1BlockNumber()
		ts := service.GetLatestL1Timestamp()
		log.Info("Initialized Latest L1 Info", "blocknumber", bn, "timestamp", ts)

		var i, q string
		index := service.GetLatestIndex()
		queueIndex := service.GetLatestEnqueueIndex()
		if index == nil {
			i = "<nil>"
		} else {
			i = strconv.FormatUint(*index, 10)
		}
		if queueIndex == nil {
			q = "<nil>"
		} else {
			q = strconv.FormatUint(*queueIndex, 10)
		}
		log.Info("Initialized Eth Context", "index", i, "queue-index", q)

		// The sequencer needs to sync to the tip at start up
		// By setting the sync status to true, it will prevent RPC calls.
		// Be sure this is set to false later.
		if !service.verifier {
			service.setSyncStatus(true)
		}
	}

	return &service, nil
}

func (s *SyncService) ensureClient() error {
	_, err := s.client.GetLatestEthContext()
	if err != nil {
		return fmt.Errorf("Cannot connect to data service: %w", err)
	}
	return nil
}

// Start initializes the service, connecting to Ethereum1 and starting the
// subservices required for the operation of the SyncService.
// txs through syncservice go to mempool.locals
// txs through rpc go to mempool.remote
func (s *SyncService) Start() error {
	if !s.enable {
		return nil
	}
	log.Info("Initializing Sync Service", "eth1-chainid", s.eth1ChainId)

	// When a sequencer, be sure to sync to the tip of the ctc before allowing
	// user transactions.
	if !s.verifier {
		err := s.syncTransactionsToTip()
		if err != nil {
			return fmt.Errorf("Cannot sync transactions to the tip: %w", err)
		}
		// TODO: This should also sync the enqueue'd transactions that have not
		// been synced yet
		s.setSyncStatus(false)
	}

	if s.verifier {
		go s.VerifierLoop()
	} else {
		go s.SequencerLoop()
	}
	return nil
}

// initializeLatestL1 sets the initial values of the `L1BlockNumber`
// and `L1Timestamp` to the deploy height of the Canonical Transaction
// chain if the chain is empty, otherwise set it from the last
// transaction processed. This must complete before transactions
// are accepted via RPC when running as a sequencer.
func (s *SyncService) initializeLatestL1(ctcDeployHeight *big.Int) error {
	index := s.GetLatestIndex()
	if index == nil {
		if ctcDeployHeight == nil {
			return errors.New("Must configure with canonical transaction chain deploy height")
		}
		context, err := s.client.GetEthContext(ctcDeployHeight.Uint64())
		if err != nil {
			return fmt.Errorf("Cannot fetch ctc deploy block at height %d: %w", ctcDeployHeight.Uint64(), err)
		}
		s.SetLatestL1Timestamp(context.Timestamp)
		s.SetLatestL1BlockNumber(context.BlockNumber)
	} else {
		log.Info("Found latest index", "index", *index)
		block := s.bc.GetBlockByNumber(*index - 1)
		txs := block.Transactions()
		if len(txs) != 1 {
			log.Error("Unexpected number of transactions in block: %d", len(txs))
		}
		tx := txs[0]
		s.SetLatestL1Timestamp(tx.L1Timestamp())
		s.SetLatestL1BlockNumber(tx.L1BlockNumber().Uint64())
	}
	// Only the sequencer cares about latest queue index
	if !s.verifier {
		queueIndex := s.GetLatestEnqueueIndex()
		if queueIndex == nil {
			enqueue, err := s.client.GetLastConfirmedEnqueue()
			if err != nil {
				return fmt.Errorf("Cannot fetch last confirmed queue tx: %w", err)
			}
			// There are no enqueues yet
			if enqueue == nil {
				return nil
			}
			queueIndex = enqueue.GetMeta().QueueIndex
		}
		s.SetLatestEnqueueIndex(queueIndex)
	}
	return nil
}

// setSyncStatus sets the `syncing` field as well as prevents
// any transactions from coming in via RPC.
// `syncing` should never be set directly outside of this function.
func (s *SyncService) setSyncStatus(status bool) {
	log.Info("Setting sync status", "status", status)
	s.syncing.Store(status)
}

// IsSyncing returns the syncing status of the syncservice.
// Returns false if not yet set.
func (s *SyncService) IsSyncing() bool {
	value := s.syncing.Load()
	val, ok := value.(bool)
	if !ok {
		return false
	}
	return val
}

// Stop will close the open channels and cancel the goroutines
// started by this service.
func (s *SyncService) Stop() error {
	s.scope.Close()

	if s.cancel != nil {
		defer s.cancel()
	}
	return nil
}

func (s *SyncService) VerifierLoop() {
	log.Info("Starting Verifier Loop", "poll-interval", s.pollInterval, "timestamp-refresh-threshold", s.timestampRefreshThreshold)
	for {
		// The verifier polls for ctc transactions.
		// the ctc transactions are extending the chain.
		latest, err := s.client.GetLatestTransaction()
		if err != nil {
			log.Error("Cannot fetch transaction")
			continue
		}

		if latest == nil {
			time.Sleep(s.pollInterval)
			continue
		}

		var start uint64
		if s.GetLatestIndex() == nil {
			start = 0
		} else {
			start = *s.GetLatestIndex() + 1
		}
		end := *latest.GetMeta().Index
		log.Info("Polling transactions", "start", start, "end", end)
		for i := start; i <= end; i++ {
			tx, err := s.client.GetTransaction(i)
			if err != nil {
				log.Error("Cannot get tx in loop", "index", i)
				continue
			}
			log.Debug("Applying transaction", "index", i)
			err = s.maybeApplyTransaction(tx)
			if err != nil {
				log.Error("Cannot apply transaction", "msg", err)
			}
			s.SetLatestIndex(&i)
		}
		time.Sleep(s.pollInterval)
	}
}

func (s *SyncService) SequencerLoop() {
	log.Info("Starting Sequencer Loop", "poll-interval", s.pollInterval, "timestamp-refresh-threshold", s.timestampRefreshThreshold)
	for {
		// Only the sequencer needs to poll for enqueue transactions
		// and then can choose when to apply them. We choose to apply
		// transactions such that it makes for efficient batch submitting.
		// Place as many L1ToL2 transactions in the same context as possible
		// by executing them one after another.
		// TODO: break this routine out into a function so that lock
		// management is more simple. For now, be sure to unlock before
		// each outer continue
		s.txLock.Lock()
		latest, err := s.client.GetLatestEnqueue()
		if err != nil {
			log.Error("Cannot get latest enqueue")
			s.txLock.Unlock()
			time.Sleep(s.pollInterval)
			continue
		}
		// This should never happen unless the backend is empty
		if latest == nil {
			log.Debug("No enqueue transactions found")
			s.txLock.Unlock()
			time.Sleep(s.pollInterval)
			continue
		}
		// Compare the remote latest queue index to the local latest
		// queue index. If the remote latest queue index is greater
		// than the local latest queue index, be sure to ingest more
		// enqueued transactions
		var start uint64
		if s.GetLatestEnqueueIndex() == nil {
			start = 0
		} else {
			start = *s.GetLatestEnqueueIndex() + 1
		}
		end := *latest.GetMeta().QueueIndex

		log.Info("Polling enqueued transactions", "start", start, "end", end)
		for i := start; i <= end; i++ {
			enqueue, err := s.client.GetEnqueue(i)
			if err != nil {
				log.Error("Cannot get enqueue in loop", "index", i, "message", err)
				continue
			}

			if enqueue == nil {
				log.Debug("No enqueue transaction found")
				break
			}

			// This should never happen
			if enqueue.L1BlockNumber() == nil {
				log.Error("No blocknumber for enqueue", "index", i, "timestamp", enqueue.L1Timestamp(), "blocknumber", enqueue.L1BlockNumber())
				continue
			}

			// Update the timestamp and blocknumber based on the enqueued
			// transactions
			if enqueue.L1Timestamp() > s.GetLatestL1Timestamp() {
				ts := enqueue.L1Timestamp()
				bn := enqueue.L1BlockNumber().Uint64()
				s.SetLatestL1Timestamp(ts)
				s.SetLatestL1BlockNumber(bn)
				log.Info("Updated Eth Context from enqueue", "index", i, "timestamp", ts, "blocknumber", bn)
			}

			log.Debug("Applying enqueue transaction", "index", i)
			err = s.applyTransaction(enqueue)
			if err != nil {
				log.Error("Cannot apply transaction", "msg", err)
			}

			s.SetLatestEnqueueIndex(enqueue.GetMeta().QueueIndex)
			if enqueue.GetMeta().Index == nil {
				latest := s.GetLatestIndex()
				index := uint64(0)
				if latest != nil {
					index = *latest + 1
				}
				s.SetLatestIndex(&index)
			} else {
				s.SetLatestIndex(enqueue.GetMeta().Index)
			}
		}
		s.txLock.Unlock()

		// Update the execution context's timestamp and blocknumber
		// over time. This is only necessary for the sequencer.
		context, err := s.client.GetLatestEthContext()
		if err != nil {
			log.Error("Cannot get latest eth context", "msg", err)
			continue
		}
		current := time.Unix(int64(s.GetLatestL1Timestamp()), 0)
		next := time.Unix(int64(context.Timestamp), 0)
		if next.Sub(current) > s.timestampRefreshThreshold {
			log.Info("Updating Eth Context", "timetamp", context.Timestamp, "blocknumber", context.BlockNumber)
			s.SetLatestL1BlockNumber(context.BlockNumber)
			s.SetLatestL1Timestamp(context.Timestamp)
		}
		time.Sleep(s.pollInterval)
	}
}

// This function must sync all the way to the tip
// TODO: it should then sync all of the enqueue transactions
func (s *SyncService) syncTransactionsToTip() error {
	// Then set up a while loop that only breaks when the latest
	// transaction does not change through two runs of the loop.
	// The latest transaction can change during the timeframe of
	// all of the transactions being sync'd.
	for {
		// This function must be sure to sync all the way to the tip.
		// First query the latest transaction
		latest, err := s.client.GetLatestTransaction()
		if err != nil {
			return fmt.Errorf("Cannot get latest transaction: %w", err)
		}
		if latest == nil {
			log.Info("No transactions to sync")
			return nil
		}
		tipHeight := latest.GetMeta().Index
		index := rawdb.ReadHeadIndex(s.db)
		start := uint64(0)
		if index != nil {
			start = *index + 1
		}

		log.Info("Syncing transactions to tip", "start", start, "end", *tipHeight)
		for i := start; i <= *tipHeight; i++ {
			tx, err := s.client.GetTransaction(i)
			if err != nil {
				return fmt.Errorf("Cannot get transaction: %w", err)
			}
			// The transaction does not yet exist in the ctc
			if tx == nil {
				log.Info("Transaction in ctc does not yet exist", "index", i)
				return nil
			}
			err = s.maybeApplyTransaction(tx)
			if err != nil {
				return fmt.Errorf("Cannot apply transaction: %w", err)
			}
			if err != nil {
				log.Error("Cannot ingest transaction", "index", i)
			}
			s.SetLatestIndex(tx.GetMeta().Index)
			if types.QueueOrigin(tx.QueueOrigin().Uint64()) == types.QueueOriginL1ToL2 {
				queueIndex := tx.GetMeta().QueueIndex
				s.SetLatestEnqueueIndex(queueIndex)
			}
		}
		// Be sure to check that no transactions came in while
		// the above loop was running
		post, err := s.client.GetLatestTransaction()
		if err != nil {
			return fmt.Errorf("Cannot get latest transaction: %w", err)
		}
		// These transactions should always have an index since they
		// are already in the ctc.
		if *latest.GetMeta().Index == *post.GetMeta().Index {
			log.Info("Done syncing transactions to tip")
			return nil
		}
	}
}

// Methods for safely accessing and storing the latest
// L1 blocknumber and timestamp. These are held in memory.
func (s *SyncService) GetLatestL1Timestamp() uint64 {
	return atomic.LoadUint64(&s.OVMContext.timestamp)
}

func (s *SyncService) GetLatestL1BlockNumber() uint64 {
	return atomic.LoadUint64(&s.OVMContext.blockNumber)
}

func (s *SyncService) SetLatestL1Timestamp(ts uint64) {
	atomic.StoreUint64(&s.OVMContext.timestamp, ts)
}

func (s *SyncService) SetLatestL1BlockNumber(bn uint64) {
	atomic.StoreUint64(&s.OVMContext.blockNumber, bn)
}

func (s *SyncService) GetLatestEnqueueIndex() *uint64 {
	return rawdb.ReadHeadQueueIndex(s.db)
}

func (s *SyncService) SetLatestEnqueueIndex(index *uint64) {
	if index != nil {
		rawdb.WriteHeadQueueIndex(s.db, *index)
	}
}

func (s *SyncService) SetLatestIndex(index *uint64) {
	if index != nil {
		rawdb.WriteHeadIndex(s.db, *index)
	}
}

func (s *SyncService) GetLatestIndex() *uint64 {
	return rawdb.ReadHeadIndex(s.db)
}

// reorganize will reorganize to directly to the index passed in.
// The caller must handle the offset relative to the ctc.
func (s *SyncService) reorganize(index uint64) error {
	if index == 0 {
		return nil
	}
	err := s.bc.SetHead(index)
	if err != nil {
		return fmt.Errorf("Cannot reorganize in syncservice: %w", err)
	}

	// TODO: make sure no off by one error here
	s.SetLatestIndex(&index)

	// When in sequencer mode, be sure to roll back the latest queue
	// index as well.
	if !s.verifier {
		enqueue, err := s.client.GetLastConfirmedEnqueue()
		if err != nil {
			return fmt.Errorf("cannot reorganize: %w", err)
		}
		s.SetLatestEnqueueIndex(enqueue.GetMeta().QueueIndex)
	}
	log.Info("Reorganizing", "height", index)
	return nil
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (s *SyncService) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return s.scope.Track(s.txFeed.Subscribe(ch))
}

// maybeApplyTransaction will potentially apply the transaction after first
// inspecting the local database. This is mean to prevent transactions from
// being replayed.
func (s *SyncService) maybeApplyTransaction(tx *types.Transaction) error {
	index := tx.GetMeta().Index
	if index == nil {
		return fmt.Errorf("nil index in maybeApplyTransaction")
	}
	// Handle off by one
	block := s.bc.GetBlockByNumber(*index + 1)

	// The transaction has yet to be played, so it is safe to apply
	if block == nil {
		err := s.applyTransaction(tx)
		if err != nil {
			return fmt.Errorf("Maybe apply transaction failed on index %d: %w", *index, err)
		}
		return nil
	}
	// There is already a transaction at that index, so check
	// for its equality.
	txs := block.Transactions()
	if len(txs) != 1 {
		log.Info("block", "txs", len(txs), "number", block.Number().Uint64())
		return fmt.Errorf("More than 1 transaction in block")
	}
	if isCtcTxEqual(tx, txs[0]) {
		log.Info("Matching transaction found", "index", *index)
	} else {
		log.Warn("Non matching transaction found", "index", *index)
	}
	return nil
}

// Lower level API used to apply a transaction, must only be used with
// transactions that came from L1.
func (s *SyncService) applyTransaction(tx *types.Transaction) error {
	txs := types.Transactions{tx}
	s.txFeed.Send(core.NewTxsEvent{Txs: txs})
	return nil
}

// Higher level API for applying transactions. Should only be called for
// queue origin sequencer transactions, as the contracts on L1 manage the same
// validity checks that are done here.
func (s *SyncService) ApplyTransaction(tx *types.Transaction) error {
	s.txLock.Lock()
	defer s.txLock.Unlock()
	if s.verifier {
		return errors.New("Verifier does not accept transactions out of band")
	}
	qo := tx.QueueOrigin()
	if qo == nil {
		return errors.New("invalid transaction with no queue origin")
	}
	if qo.Uint64() != uint64(types.QueueOriginSequencer) {
		return fmt.Errorf("invalid transaction with queue origin %d", qo.Uint64())
	}
	err := s.txpool.ValidateTx(tx)
	if err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}
	return s.applyTransaction(tx)
}
