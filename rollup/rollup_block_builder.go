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

var (
	ErrTransactionLimitReached = errors.New("transaction limit reached")
	ErrMoreThanOneTxInBlock    = errors.New("block contains more than one transaction")
	LastProcessedDBKey         = []byte("lastProcessedRollupBlock")
	MinTxBytes                 = uint64(100)
	MinTxGas                   = MinTxBytes*params.TxDataNonZeroGasEIP2028 + params.SstoreSetGas
	RollupBlockGasBuffer       = uint64(1_000_000)
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
	if maxBlockGas < b.gasUsed+block.GasUsed() {
		return core.ErrGasLimitReached
	}

	b.rollupBlock.addBlock(block)
	b.gasUsed += params.SstoreSetGas + uint64(len(block.Transactions()[0].Data()))*params.TxDataNonZeroGasEIP2028
	if b.firstBlockNumber == 0 {
		b.firstBlockNumber = block.NumberU64()
	}
	b.lastBlockNumber = block.NumberU64()

	return nil
}

type RollupBlockBuilder struct {
	db         ethdb.Database
	blockchain *core.BlockChain
	newBlockCh chan *types.Block
	pendingMu  sync.RWMutex

	maxRollupBlockTime         time.Duration
	maxRollupBlockGas          uint64
	maxRollupBlockTransactions int

	lastProcessedBlockNumber uint64
	blockInProgress          *BuildingBlock
}

func NewRollupBlockBuilder(db ethdb.Database, blockchain *core.BlockChain, maxBlockTime time.Duration, maxBlockGas uint64, maxBlockTransactions int) (*RollupBlockBuilder, error) {
	lastBlock, err := fetchLastProcessed(db)
	if err != nil {
		return nil, err
	}

	builder := &RollupBlockBuilder{
		db:         db,
		blockchain: blockchain,
		newBlockCh: make(chan *types.Block),

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

func (b *RollupBlockBuilder) buildLoop(maxBlockTime time.Duration) {
	lastProcessed := b.lastProcessedBlockNumber

	if err := b.sync(); err != nil {
		panic(fmt.Errorf("error syncing: %+v", err))
	}

	timer := time.NewTimer(0)
	<-timer.C // discard the initial tick

	for {
		select {
		case block := <-b.newBlockCh:
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
	log.Debug("handling new block in rollup block builder", "block", block)
	if block.NumberU64() <= b.lastProcessedBlockNumber {
		log.Debug("handling old block -- ignoring", "block", block)
		return false, nil
	}

	if txCount := len(block.Transactions()); txCount > 1 {
		// should never happen
		log.Error("received block with more than one transaction", "block", block)
		return false, ErrMoreThanOneTxInBlock
	} else if txCount == 0 {
		log.Debug("handling empty block -- ignoring", "block", block)
		b.lastProcessedBlockNumber = block.NumberU64()
		return false, nil
	}

	switch err := b.addBlock(block); err {
	case core.ErrGasLimitReached, ErrTransactionLimitReached:
		if _, e := b.buildRollupBlock(false); e != nil {
			log.Error("unable to build rollup block", "error", e, "rollup block", b.blockInProgress)
			return false, e
		}
		if addErr := b.addBlock(block); addErr != nil {
			// TODO: Retry and whatnot instead of instant panic
			log.Error("unable to build rollup block", "error", addErr, "rollup block", b.blockInProgress)
			return false, addErr
		}
	default:
		if err != nil {
			log.Error("unrecognized error adding to rollup block in progress", "error", err, "rollup block", b.blockInProgress)
			return false, err
		} else {
			log.Debug("successfully added block to rollup block in progress", "number", block.NumberU64())
		}
	}

	built, err := b.tryBuildRollupBlock()
	if err != nil {
		log.Error("error building block", "error", err, "block", block)
		return false, err
	}

	return built, nil
}

func (b *RollupBlockBuilder) sync() error {
	log.Info("syncing blocks in rollup block builder", "starting block", b.lastProcessedBlockNumber)

	for {
		blockNum := b.lastProcessedBlockNumber + uint64(1)
		block := b.blockchain.GetBlockByNumber(blockNum)
		log.Info("got block number", "number", blockNum, "block", block)
		if block == nil {
			log.Info("done syncing blocks in rollup block builder", "number", b.lastProcessedBlockNumber)
			return nil
		}
		if _, err := b.handleNewBlock(block); err != nil {
			log.Error("Error handling new block", "error", err)
			return err
		} else {
			log.Debug("successfully synced block", "number", blockNum, "last processed", b.lastProcessedBlockNumber)
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
		log.Debug("rollup block is not full, so not finalizing it", "txCount", txCount, "gasAfterOneMoreTx", gasAfterOneMoreTx)
		return false, nil
	}
	log.Debug("rollup block is full, finalizing it", "txCount", txCount, "gasAfterOneMoreTx", gasAfterOneMoreTx)

	return b.buildRollupBlock(false)
}

func (b *RollupBlockBuilder) buildRollupBlock(force bool) (bool, error) {
	var toSubmit *BuildingBlock
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()

	txCount := len(b.blockInProgress.rollupBlock.transitions)
	if force && txCount == 0 {
		log.Debug("rollup block is empty so not finalizing it, even though force = true")
		return false, nil
	}
	if txCount < b.maxRollupBlockTransactions && b.blockInProgress.gasUsed+MinTxGas <= b.maxRollupBlockGas {
		log.Debug("rollup block is not full, so not finalizing it")
		return false, nil
	}
	log.Debug("building rollup block")

	toSubmit = b.blockInProgress
	b.blockInProgress = newBuildingBlock(b.maxRollupBlockTransactions)
	if err := b.submit(toSubmit); err != nil {
		log.Error("error submitting rollup block", "lastBlockNumber", toSubmit.lastBlockNumber, "error", err)
		return false, err
	}
	log.Debug("successfully built rollup block", "lastBlockNumber", toSubmit.lastBlockNumber)

	return true, nil
}

func (b *RollupBlockBuilder) submit(block *BuildingBlock) error {
	// TODO: Submit to chain & get hash
	log.Debug("submitting rollup block", "block", block)
	if err := b.db.Put(LastProcessedDBKey, serializeBlockNumber(block.lastBlockNumber)); err != nil {
		log.Error("error saving last processed rollup block", "block", block)
		// TODO: Something here
	}
	log.Debug("rollup block submitted", "block", block)
	return nil
}

func fetchLastProcessed(db ethdb.Database) (uint64, error) {
	has, err := db.Has(LastProcessedDBKey)
	if err != nil {
		log.Error("received error checking if LastProcessedDBKey exists in DB", "error", err)
		return 0, err
	}
	if has {
		lastProcessedBytes, e := db.Get(LastProcessedDBKey)
		if e != nil {
			log.Error("error fetching LastProcessedDBKey from DB", "error", err)
			return 0, err
		}
		lastProcessedBlock := deserializeBlockNumber(lastProcessedBytes)
		log.Info("fetched last processed block from database", "number", lastProcessedBlock)
		return lastProcessedBlock, nil
	} else {
		log.Info("no last processed block found in the db -- returning 0")
		return 0, nil
	}
}

func serializeBlockNumber(blockNumber uint64) []byte {
	numberAsByteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsByteArray, blockNumber)
	return numberAsByteArray
}

func deserializeBlockNumber(blockNumber []byte) uint64 {
	return binary.LittleEndian.Uint64(blockNumber)
}
