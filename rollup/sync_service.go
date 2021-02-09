package rollup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TODO: rename this to OVMContext
// LatestL1ToL2 represents the latest blocknumber and timestamp.
type LatestL1ToL2 struct {
	blockNumber uint64
	timestamp   uint64
	queueIndex  *uint64
	index       *uint64
}

// TODO: update timestamp/blocknumber logic

// SyncService implements the verifier functionality as well as the reorg
// protection for the sequencer.
type SyncService struct {
	ctx               context.Context
	cancel            context.CancelFunc
	verifier          bool
	db                ethdb.Database
	scope             event.SubscriptionScope
	txFeed            event.Feed
	enable            bool
	eth1ChainId       uint64
	bc                *core.BlockChain
	txpool            *core.TxPool
	client            RollupClient
	syncing           bool
	LatestL1ToL2      LatestL1ToL2
	confirmationDepth uint64
	pollInterval      time.Duration
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

	// Layer 2 chainid
	chainID := bc.Config().ChainID
	// Initialize the rollup client
	client := NewClient(cfg.RollupClientHttp, chainID)
	service := SyncService{
		ctx:               ctx,
		cancel:            cancel,
		verifier:          cfg.IsVerifier,
		enable:            cfg.Eth1SyncServiceEnable,
		confirmationDepth: cfg.Eth1ConfirmationDepth,
		bc:                bc,
		txpool:            txpool,
		eth1ChainId:       cfg.Eth1ChainId,
		client:            client,
		db:                db,
		pollInterval:      time.Second * 15,
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

		// Initialization logic
		block := service.bc.CurrentBlock()
		if block == nil {
			return nil, errors.New("Current block is nil")
		}
		if block != service.bc.Genesis() {
			// Roll back the chain to force some amount of resync
			depth := block.Number().Uint64() - cfg.InitialReorgDepth
			if depth > block.Number().Uint64() {
				return nil, fmt.Errorf("Overflow with initial reorg depth %d and tip %d", cfg.InitialReorgDepth, block.Number().Uint64())
			}
			err = service.reorganize(depth)
			if err != nil {
				return nil, fmt.Errorf("Cannot reorg with depth %d: %w", cfg.InitialReorgDepth, err)
			}
		}

		// Initialize the latest L1 data here to make sure that
		// it happens before the RPC endpoints open up
		// Only do it if the sync service is enabled so that this
		// can be ran without needing to have a configured client.
		err = service.initializeLatestL1(cfg.CanonicalTransactionChainDeployHeight)
		if err != nil {
			return nil, fmt.Errorf("Cannot initialize latest L1 data: %w", err)
		}
		log.Info("Initialized Latest L1 Info")
	}

	// The sequencer needs to sync to the tip at start up
	if !service.verifier {
		service.setSyncStatus(true)
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
		s.setSyncStatus(false)
	}

	go s.Loop()
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
		block := s.bc.GetBlockByNumber(*index + 1)
		txs := block.Transactions()
		if len(txs) != 1 {
			log.Error("Unexpected number of transactions in block: %d", len(txs))
		}
		tx := txs[0]
		s.SetLatestL1Timestamp(tx.L1Timestamp())
		s.SetLatestL1BlockNumber(tx.L1BlockNumber().Uint64())

	}
	// Set the last enqueue index for the sequencer
	// to properly begin syncing L1 to L2 transactions from
	enqueue, err := s.client.GetLatestEnqueue()
	if err != nil {
		return fmt.Errorf("Cannot get latest enqueue: %w", err)
	}
	if enqueue == nil {
		s.SetLatestEnqueueIndex(nil)
		return nil
	}
	for {
		// When the ctc index is not nil then the enqueue has
		// been included in the chain already
		meta := enqueue.GetMeta()
		if meta.Index != nil {
			s.SetLatestEnqueueIndex(meta.QueueIndex)
			return nil
		}
		next, err := s.client.GetEnqueue(*meta.QueueIndex - 1)
		if err != nil {
			log.Error("Cannot get enqueue", "index", *meta.Index)
			continue
		}
		enqueue = next
	}
}

// setSyncStatus sets the `syncing` field as well as prevents
// any transactions from coming in via RPC.
// `syncing` should never be set directly outside of this function.
func (s *SyncService) setSyncStatus(status bool) {
	log.Info("Setting sync status", "status", status)
	s.syncing = status
}

// IsSyncing returns the syncing status of the syncservice.
func (s *SyncService) IsSyncing() bool {
	return s.syncing
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

// Loop is the main processing loop for the sync service.
// If running as a sequencer, it will pull in any enqueue transactions
// and apply them. It pulls in as many sequential enqueue transactions
// at once and applies them sequentially as to create the most efficient
// batch contexts. Note that this function assumes that the historical
// state has already been synced.
func (s *SyncService) Loop() {
	log.Info("Starting Tip processing loop")
	for {
		log.Info("Start loop")

		// Only the sequencer needs to poll for enqueue transactions
		// and then can choose when to apply them. We choose to apply
		// transactions such that it makes for efficient batch submitting.
		// Place as many L1ToL2 transactions in the same context as possible
		// by executing them one after another.
		if !s.verifier {
			latest, err := s.client.GetLatestEnqueue()
			if err != nil {
				log.Error("Cannot get latest enqueue")
				continue
			}
			// This should never happen unless the backend is empty
			if latest == nil {
				log.Error("Latest enqueue is nil in Loop")
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

			log.Info("loop enqueue", "start", start, "end", end)
			for i := start; i <= end; i++ {
				enqueue, err := s.client.GetEnqueue(i)
				if err != nil {
					log.Error("Cannot get enqueue in loop", "index", i)
					continue
				}
				log.Info("loop enqueue apply tx", "index", i, "hash", enqueue.Hash().Hex())
				err = s.applyTransaction(enqueue)
				if err != nil {
					log.Error("Cannot apply transaction", "msg", err)
				}
				s.SetLatestEnqueueIndex(enqueue.GetMeta().QueueIndex)
				if enqueue.GetMeta().Index == nil {
					index := *s.GetLatestIndex() + 1
					s.SetLatestIndex(&index)
				} else {
					s.SetLatestIndex(enqueue.GetMeta().Index)
				}
			}
		}

		// Both the verifier and the sequencer poll for ctc transactions.
		// For the sequencer, ctc transactions are in the past while for
		// the verifier, ctc transactions are extending the chain.
		// The sequencer essentially runs a verifier to make sure that
		// it reflects the ultimate source of truth which is the L1 contracts.
		latest, err := s.client.GetLatestTransaction()
		if err != nil {
			log.Error("Cannot fetch transaction")
			continue
		}

		var start uint64
		if s.GetLatestIndex() == nil {
			start = 0
		} else {
			start = *s.GetLatestIndex() + 1
		}
		// TODO: end if broken here
		end := *latest.GetMeta().Index

		log.Info("loop transaction", "start", start, "end", end)
		for i := start; i < end; i++ {
			tx, err := s.client.GetTransaction(i)
			if err != nil {
				log.Error("Cannot get tx in loop", "index", i)
				continue
			}
			log.Info("loop apply tx", "index", i, "hash", tx.Hash().Hex())
			err = s.applyTransaction(tx)
			if err != nil {
				log.Error("Cannot apply transaction", "msg", err)
			}
			s.SetLatestIndex(&i)
		}

		time.Sleep(s.pollInterval)
		log.Info("End loop")
	}
}

// This function must sync all the way to the tip
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
		tipHeight := latest.GetMeta().Index
		block := s.bc.CurrentBlock()
		start := block.Number().Uint64()
		if start != 0 {
			start = start - 1
		}

		log.Info("Syncing transactions to tip", "start", start, "end", *tipHeight)
		for i := start; i <= *tipHeight; i++ {
			tx, err := s.client.GetTransaction(i)
			if err != nil {
				return fmt.Errorf("Cannot get transaction: %w", err)
			}
			// The transaction does not yet exist in the ctc
			if tx == nil {
				log.Trace("Transaction in ctc does not yet exist", "index", i)
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
			if types.QueueOrigin(tx.QueueOrigin().Uint64()) == types.QueueOriginSequencer {
				queueIndex := tx.GetMeta().QueueIndex
				log.Info("queue index", "index", queueIndex)
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

// TODO: move to db
func (s *SyncService) GetLatestEnqueueIndex() *uint64 {
	if s.LatestL1ToL2.queueIndex == nil {
		return nil
	}
	cpy := *s.LatestL1ToL2.queueIndex
	return &cpy
}

// TODO: move to db
func (s *SyncService) SetLatestEnqueueIndex(index *uint64) {
	if index == nil {
		s.LatestL1ToL2.queueIndex = nil
	} else {
		if s.LatestL1ToL2.queueIndex == nil {
			s.LatestL1ToL2.queueIndex = index
		} else {
			*s.LatestL1ToL2.queueIndex = *index
		}
	}
}

// TODO: move to db
func (s *SyncService) SetLatestIndex(index *uint64) {
	if index == nil {
		s.LatestL1ToL2.index = nil
	} else {
		if s.LatestL1ToL2.index == nil {
			s.LatestL1ToL2.index = index
		} else {
			*s.LatestL1ToL2.index = *index
		}
	}
}

// TODO: move to db
func (s *SyncService) GetLatestIndex() *uint64 {
	if s.LatestL1ToL2.index == nil {
		return nil
	}
	cpy := *s.LatestL1ToL2.index
	return &cpy
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
	return nil
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (s *SyncService) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return s.scope.Track(s.txFeed.Subscribe(ch))
}

// TODO: This function needs to be rethought, its no longer
// easily possible to say "start syncing from L1 height x"
func (s *SyncService) SetL1Head(number uint64) error {
	return nil
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

	// TODO: temp debug
	log.Info("maybeApplyTransaction", "index", *index)

	// The transaction has yet to be played, so it is safe to apply
	if block == nil {
		return s.applyTransaction(tx)
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

// Adds the transaction to the mempool so that downstream services
// can apply it to the state. This should directly play against
// the state eventually, skipping the mempool.
func (s *SyncService) applyTransaction(tx *types.Transaction) error {
	err := s.txpool.ValidateTx(tx)
	// The sequencer needs to prevent transactions that fail the mempool
	// checks. The verifier needs to play the transactions no matter what
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
