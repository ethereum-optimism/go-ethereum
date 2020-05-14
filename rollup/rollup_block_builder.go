package rollup

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"sync"
	"time"
)

type BlockStore interface {
	GetBlockByNumber(number uint64) *types.Block
}

const (
	MinTxBytes                 = uint64(100)
	MinTxGas                   = MinTxBytes*params.TxDataNonZeroGasEIP2028 + params.SstoreSetGas
	RollupBlockGasBuffer       = uint64(1_000_000)
)

var (
	logger = log.New(RollupBlockBuilder{})
	ErrTransactionLimitReached = errors.New("transaction limit reached")
	ErrMoreThanOneTxInBlock    = errors.New("block contains more than one transaction")
	LastProcessedDBKey         = []byte("lastProcessedRollupBlock")
)

type Transition struct {
	transaction *types.Transaction
	postState   common.Hash
}

func newTransition(tx *types.Transaction, postState common.Hash) *Transition {
	return &Transition{
		transaction: tx,
		postState:   postState,
	}
}

type RollupBlock struct {
	transitions []*Transition
}

func newRollupBlock(defaultSize int) *RollupBlock {
	return &RollupBlock{transitions: make([]*Transition, 0, defaultSize)}
}

func (r *RollupBlock) addBlock(block *types.Block) {
	r.transitions = append(r.transitions, newTransition(block.Transactions()[0], block.Root()))
}

type BuildingBlock struct {
	firstBlockNumber uint64
	lastBlockNumber  uint64
	gasUsed          uint64

	rollupBlock *RollupBlock
}

func newBuildingBlock(defaultTxCapacity int) *BuildingBlock {
	return &BuildingBlock{
		firstBlockNumber: 0,
		lastBlockNumber:  0,
		gasUsed:          RollupBlockGasBuffer,
		rollupBlock:      newRollupBlock(defaultTxCapacity),
	}
}

func (b *BuildingBlock) addBlock(block *types.Block, maxBlockGas uint64, maxBlockTransactions int) error {
	if maxBlockTransactions < len(b.rollupBlock.transitions)+1 {
		return ErrTransactionLimitReached
	}
	blockGasCost := GetBlockRollupGasUsage(block)
	if maxBlockGas < b.gasUsed + blockGasCost {
		return core.ErrGasLimitReached
	}

	b.rollupBlock.addBlock(block)
	b.gasUsed += blockGasCost
	if b.firstBlockNumber == 0 {
		b.firstBlockNumber = block.NumberU64()
	}
	b.lastBlockNumber = block.NumberU64()

	return nil
}

type RollupBlockBuilder struct {
	db            ethdb.Database
	blockProvider        BlockStore
	rollupBlockSubmitter RollupBlockSubmitter
	pendingMu     sync.RWMutex

	newBlockCh    chan *types.Block

	maxRollupBlockTime         time.Duration
	maxRollupBlockGas          uint64
	maxRollupBlockTransactions int

	lastProcessedBlockNumber uint64
	blockInProgress          *BuildingBlock
}

func NewRollupBlockBuilder(db ethdb.Database, blockStore interface{}, rollupBlockSubmitter interface{}, maxBlockTime time.Duration, maxBlockGas uint64, maxBlockTransactions int) (*RollupBlockBuilder, error) {
	lastBlock, err := fetchLastProcessed(db)
	if err != nil {
		return nil, err
	}

	builder := &RollupBlockBuilder{
		db:            db,
		blockProvider: blockStore.(BlockStore),
		rollupBlockSubmitter: rollupBlockSubmitter.(RollupBlockSubmitter),
		newBlockCh:    make(chan *types.Block, 10_000),

		maxRollupBlockTime:         maxBlockTime,
		maxRollupBlockGas:          maxBlockGas,
		maxRollupBlockTransactions: maxBlockTransactions,

		lastProcessedBlockNumber: lastBlock,
		blockInProgress:          newBuildingBlock(maxBlockTransactions),
	}

	go builder.buildLoop(maxBlockTime)

	return builder, nil
}

func (b *RollupBlockBuilder) NewBlock(block *types.Block) {
	b.newBlockCh <- block
}

func (b *RollupBlockBuilder) Stop() {
	close(b.newBlockCh)
}

func (b *RollupBlockBuilder) buildLoop(maxBlockTime time.Duration) {
	lastProcessed := b.lastProcessedBlockNumber

	if err := b.sync(); err != nil {
		panic(fmt.Errorf("error syncing: %+v", err))
	}

	timer := time.NewTimer(0)
	<-timer.C // discard the initial tick

	for {
		select {
		case block, ok := <-b.newBlockCh:
			if !ok {
				timer.Stop()
				logger.Info("Closing rollup block builder new block channel. If not shutting down, this is an error")
				return
			}

			built, err := b.handleNewBlock(block)
			if err != nil {
				panic(fmt.Errorf("error handling new block. Error: %v. Block: %+v", err, block))
			}
			if timer != nil && built {
				timer.Reset(b.maxRollupBlockTime)
			}
		case <-timer.C:
			if lastProcessed != b.lastProcessedBlockNumber && b.blockInProgress.firstBlockNumber != 0 {
				if _, err := b.buildRollupBlock(true); err != nil {
					panic(fmt.Errorf("error buidling block: %v", err))
				}
			}

			lastProcessed = b.lastProcessedBlockNumber
			timer.Reset(maxBlockTime)
		}
	}
}

func (b *RollupBlockBuilder) handleNewBlock(block *types.Block) (bool, error) {
	logger.Debug("handling new block in rollup block builder", "block", block)
	if block.NumberU64() <= b.lastProcessedBlockNumber {
		logger.Debug("handling old block -- ignoring", "block", block)
		return false, nil
	}

	if txCount := len(block.Transactions()); txCount > 1 {
		// should never happen
		logger.Error("received block with more than one transaction", "block", block)
		return false, ErrMoreThanOneTxInBlock
	} else if txCount == 0 {
		logger.Debug("handling empty block -- ignoring", "block", block)
		b.lastProcessedBlockNumber = block.NumberU64()
		return false, nil
	}

	switch err := b.addBlock(block); err {
	case core.ErrGasLimitReached, ErrTransactionLimitReached:
		if _, e := b.buildRollupBlock(false); e != nil {
			logger.Error("unable to build rollup block", "error", e, "rollup block", b.blockInProgress)
			return false, e
		}
		if addErr := b.addBlock(block); addErr != nil {
			// TODO: Retry and whatnot instead of instant panic
			logger.Error("unable to build rollup block", "error", addErr, "rollup block", b.blockInProgress)
			return false, addErr
		}
	default:
		if err != nil {
			logger.Error("unrecognized error adding to rollup block in progress", "error", err, "rollup block", b.blockInProgress)
			return false, err
		} else {
			logger.Debug("successfully added block to rollup block in progress", "number", block.NumberU64())
		}
	}

	built, err := b.tryBuildRollupBlock()
	if err != nil {
		logger.Error("error building block", "error", err, "block", block)
		return false, err
	}

	return built, nil
}

func (b *RollupBlockBuilder) sync() error {
	logger.Info("syncing blocks in rollup block builder", "starting block", b.lastProcessedBlockNumber)

	for {
		blockNum := b.lastProcessedBlockNumber + uint64(1)
		block := b.blockProvider.GetBlockByNumber(blockNum)
		logger.Info("got block number", "number", blockNum, "block", block)
		if block == nil {
			logger.Info("done syncing blocks in rollup block builder", "number", b.lastProcessedBlockNumber)
			return nil
		}
		if _, err := b.handleNewBlock(block); err != nil {
			logger.Error("Error handling new block", "error", err)
			return err
		} else {
			logger.Debug("successfully synced block", "number", blockNum, "last processed", b.lastProcessedBlockNumber)
		}
	}
}

func (b *RollupBlockBuilder) addBlock(block *types.Block) error {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	if err := b.blockInProgress.addBlock(block, b.maxRollupBlockGas, b.maxRollupBlockTransactions); err != nil {
		return err
	}
	b.lastProcessedBlockNumber = block.NumberU64()
	return nil
}

func (b *RollupBlockBuilder) tryBuildRollupBlock() (bool, error) {
	txCount := len(b.blockInProgress.rollupBlock.transitions)
	gasAfterOneMoreTx := b.blockInProgress.gasUsed + MinTxGas
	if txCount < b.maxRollupBlockTransactions && gasAfterOneMoreTx <= b.maxRollupBlockGas {
		logger.Debug("rollup block is not full, so not finalizing it", "txCount", txCount, "gasAfterOneMoreTx", gasAfterOneMoreTx)
		return false, nil
	}
	logger.Debug("rollup block is full, finalizing it", "txCount", txCount, "gasAfterOneMoreTx", gasAfterOneMoreTx)

	return b.buildRollupBlock(false)
}

func (b *RollupBlockBuilder) buildRollupBlock(force bool) (bool, error) {
	var toSubmit *BuildingBlock
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()

	txCount := len(b.blockInProgress.rollupBlock.transitions)

	if force && txCount == 0 {
		logger.Debug("rollup block is empty so not finalizing it, even though force = true")
		return false, nil
	}
	if txCount < b.maxRollupBlockTransactions && b.blockInProgress.gasUsed+MinTxGas <= b.maxRollupBlockGas {
		logger.Debug("rollup block is not full, so not finalizing it")
		return false, nil
	}
	logger.Debug("building rollup block")

	toSubmit = b.blockInProgress
	b.blockInProgress = newBuildingBlock(b.maxRollupBlockTransactions)

	if err := b.submitBlock(toSubmit); err != nil {
		logger.Error("error submitting rollup block", "lastBlockNumber", toSubmit.lastBlockNumber, "error", err)
		return false, err
	}
	logger.Debug("successfully built rollup block", "lastBlockNumber", toSubmit.lastBlockNumber)

	return true, nil
}

func (b *RollupBlockBuilder) submitBlock(block *BuildingBlock) error {
	// TODO: Submit to chain & get hash
	logger.Debug("submitting rollup block", "block", block)

	if err := b.rollupBlockSubmitter.submit(block.rollupBlock); err != nil {
		return err
	}

	if err := b.db.Put(LastProcessedDBKey, SerializeBlockNumber(block.lastBlockNumber)); err != nil {
		logger.Error("error saving last processed rollup block", "block", block)
		// TODO: Something here
	}
	logger.Debug("rollup block submitted", "block", block)
	return nil
}

func fetchLastProcessed(db ethdb.Database) (uint64, error) {
	has, err := db.Has(LastProcessedDBKey)
	if err != nil {
		logger.Error("received error checking if LastProcessedDBKey exists in DB", "error", err)
		return 0, err
	}
	if has {
		lastProcessedBytes, e := db.Get(LastProcessedDBKey)
		if e != nil {
			logger.Error("error fetching LastProcessedDBKey from DB", "error", err)
			return 0, err
		}
		lastProcessedBlock := DeserializeBlockNumber(lastProcessedBytes)
		logger.Info("fetched last processed block from database", "number", lastProcessedBlock)
		return lastProcessedBlock, nil
	} else {
		logger.Info("no last processed block found in the db -- returning 0")
		return 0, nil
	}
}

func SerializeBlockNumber(blockNumber uint64) []byte {
	numberAsByteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsByteArray, blockNumber)
	return numberAsByteArray
}

func DeserializeBlockNumber(blockNumber []byte) uint64 {
	return binary.LittleEndian.Uint64(blockNumber)
}

func GetBlockRollupGasUsage(block *types.Block) uint64 {
	return params.SstoreSetGas + uint64(len(block.Transactions()[0].Data()))*params.TxDataNonZeroGasEIP2028
}
