package rollup

import (
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethdb"
    "sync"
    "time"
)

var (
    ErrTransactionLimitReached = errors.New("transaction limit reached")
    LastProcessedDBKey = []byte("lastProcessedRollupBlock")
    MinTxGas = uint64(21000)
)

type Transition struct {
    transaction *types.Transaction
    postState common.Hash
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
    lastBlockNumber uint64
    gasUsed uint64

    rollupBlock *RollupBlock
}

func newBuildingBlock(defaultTxCapacity int) *BuildingBlock {
    return &BuildingBlock{
        firstBlockNumber: 0,
        lastBlockNumber:  0,
        gasUsed:          0,
        rollupBlock:      newRollupBlock(defaultTxCapacity),
    }
}

func (b *BuildingBlock) addBlock(block *types.Block, maxBlockGas uint64, maxBlockTransactions int) error {
    if maxBlockTransactions < len(b.rollupBlock.transitions) + 1 {
        return ErrTransactionLimitReached
    }
    if maxBlockGas < b.gasUsed + block.GasUsed(){
        return core.ErrGasLimitReached
    }

    b.rollupBlock.addBlock(block)
    b.gasUsed += block.Transactions()[0].Gas()
    if b.firstBlockNumber == 0 {
        b.firstBlockNumber = block.NumberU64()
    }
    b.lastBlockNumber = block.NumberU64()

    return nil
}

type RollupBlockBuilder struct {
    db          ethdb.Database
    blockchain  *core.BlockChain
    newBlockCh  chan *types.Block
    pendingMu   sync.RWMutex

    maxRollupBlockTime   time.Duration
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
        db: db,
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
    b.newBlockCh <-block
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
                if err := b.buildRollupBlock(); err != nil {
                    panic(fmt.Errorf("error buidling block: %v", err))
                }
            }

            lastProcessed = b.lastProcessedBlockNumber
            timer.Reset(maxBlockTime)
        }
    }
}

func (b *RollupBlockBuilder) handleNewBlock(block *types.Block) (bool, error) {
    if block.NumberU64() <= b.lastProcessedBlockNumber {
        return false, nil
    }
    if txCount := len(block.Transactions()); txCount > 1 {
        panic(fmt.Errorf("received block with more than 1 tx: %+v", block))
    } else if txCount == 0 {
        b.lastProcessedBlockNumber = block.NumberU64()
        return false, nil
    }

    switch err := b.addBlock(block); err {
    case core.ErrGasLimitReached, ErrTransactionLimitReached:
        if e := b.buildRollupBlock(); e != nil {
            panic(fmt.Errorf("unable to build rollup block. Error: %v, rollup block: %+v", e, b.blockInProgress))
        }
        if addErr := b.addBlock(block); addErr != nil {
            // TODO: Retry and whatnot instead of instant panic
            panic(fmt.Errorf("unable to build rollup block. Error: %v, rollup block: %+v", addErr, b.blockInProgress))
        }
    default:
        if err != nil {
            panic(fmt.Errorf("unrecognized error adding to rollup block in progress. Error: %v, rollup block: %+v", err, b.blockInProgress))
        }
    }

    built, err := b.tryBuildRollupBlock()
    if err != nil {
        panic(fmt.Errorf("error buidling block: %v. Block: %+v ", err, block))
    }
    return built, nil
}

func (b *RollupBlockBuilder) sync() error {
    for {
        block := b.blockchain.GetBlockByNumber(b.lastProcessedBlockNumber + uint64(1))
        if block == nil {
            break
        }
        if _, err := b.handleNewBlock(block); err != nil {
            panic(fmt.Errorf("error handling block: %v. Block: %+v ", err, block))
        }
    }
    return nil
}

func (b *RollupBlockBuilder) addBlock(block *types.Block) error {
    b.pendingMu.Lock()
    defer b.pendingMu.Unlock()
    return b.blockInProgress.addBlock(block, b.maxRollupBlockGas, b.maxRollupBlockTransactions)
}

func (b *RollupBlockBuilder) tryBuildRollupBlock() (bool, error) {
    b.pendingMu.RLock()
    if len(b.blockInProgress.rollupBlock.transitions) < b.maxRollupBlockTransactions && b.blockInProgress.gasUsed + MinTxGas <= b.maxRollupBlockGas {
        return false, nil
    }
    b.pendingMu.RUnlock()

    return true, b.buildRollupBlock()
}

func (b *RollupBlockBuilder) buildRollupBlock() error {
    var toSubmit *BuildingBlock
    b.pendingMu.Lock()

    toSubmit = b.blockInProgress
    b.blockInProgress = newBuildingBlock(b.maxRollupBlockTransactions)
    b.pendingMu.Unlock()

    go b.submit(toSubmit)
    return nil
}

func (b *RollupBlockBuilder) submit(block *BuildingBlock) {
    // TODO: Submit to chain & get hash

    b.db.Put(LastProcessedDBKey, serializeBlockNumber(b.blockInProgress.lastBlockNumber))
}

func fetchLastProcessed(db ethdb.Database) (uint64, error) {
    has, err := db.Has(LastProcessedDBKey)
    if err != nil {
        panic(fmt.Errorf("error checking if last processed block number was set: %v", err))
    }
    if has {
        lastProcessedBytes, e := db.Get(LastProcessedDBKey)
        if e != nil {
            panic(fmt.Errorf("error fetching last processed block number: %v", err))
        }
        return deserializeBlockNumber(lastProcessedBytes), nil
    } else {
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

