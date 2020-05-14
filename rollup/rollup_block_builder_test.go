package rollup

import (
    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/rawdb"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/params"
    "math/big"
    "testing"
    "time"
)

var (
    timeoutDuration = time.Millisecond * 100

    testTxPoolConfig  core.TxPoolConfig
    cliqueChainConfig *params.ChainConfig

    // Test accounts
    testBankKey, _  = crypto.GenerateKey()
    testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
    testBankFunds   = big.NewInt(1000000000000000000)

    testUserKey, _  = crypto.GenerateKey()
    testUserAddress = crypto.PubkeyToAddress(testUserKey.PublicKey)
)

func init() {
    cliqueChainConfig = params.AllCliqueProtocolChanges
    cliqueChainConfig.Clique = &params.CliqueConfig{
        Period: 10,
        Epoch:  30000,
    }
}

type TestBlockStore struct {
    blocks map[uint64]*types.Block
}

func newTestBlockStore(blocks []*types.Block) *TestBlockStore {
    store := &TestBlockStore{blocks: make(map[uint64]*types.Block, len(blocks))}
    for _, block := range blocks {
        store.blocks[block.NumberU64()] = block
    }

    return store
}

func (t *TestBlockStore) GetBlockByNumber(number uint64) *types.Block {
    if block, found := t.blocks[number]; found {
        return block
    }
    return nil
}

type TestBlockSubmitter struct {
    submittedBlocks []*RollupBlock
    submitCh chan *RollupBlock
}

func newTestBlockSubmitter(submittedBlocks []*RollupBlock, submitCh chan *RollupBlock) *TestBlockSubmitter {
    return &TestBlockSubmitter{
        submittedBlocks: submittedBlocks,
        submitCh: submitCh,
    }
}

func (t *TestBlockSubmitter) submit(block *RollupBlock) error {
    t.submittedBlocks = append(t.submittedBlocks, block)
    t.submitCh <-block
    return nil
}

func createBlocks(number int, startIndex int, withTx bool) types.Blocks {
    blocks := make(types.Blocks, number)
    for i := 0; i < number; i++ {
        header := &types.Header{Number: big.NewInt(int64(i + startIndex))}
        txs := make(types.Transactions, 0)
        if withTx {
            tx, _ := types.SignTx(types.NewTransaction(uint64(i), testUserAddress, big.NewInt(1), params.TxGas, big.NewInt(0), nil), types.HomesteadSigner{}, testBankKey)
            txs = append(txs, tx)
        }
        block := types.NewBlock(header, txs, make([]*types.Header,0), make([]*types.Receipt, 0))
        blocks[i] = block
    }
    return blocks
}

func assertTransitionFromBlock(t *testing.T, transition *Transition, block *types.Block) {
    if transition.postState != block.Root() {
        t.Fatal("expecting transition postState to equal block root", "postState", transition.postState, "block.Hash()", block.Root())
    }
    if transition.transaction.Hash() != block.Transactions()[0].Hash() {
        t.Fatal("expecting transition tx hash to equal block tx hash", "transition tx", transition.transaction.Hash(), "block tx", block.Transactions()[0].Hash())
    }
}

func newTestBlockBuilder(blockStore *TestBlockStore, blockSubmitter *TestBlockSubmitter, lastProcessedBlock uint64, maxBlockTime time.Duration, maxBlockGas uint64, maxBlockTransactions int) (*RollupBlockBuilder, error) {
    db := rawdb.NewMemoryDatabase()
    if err := db.Put(LastProcessedDBKey, SerializeBlockNumber(lastProcessedBlock)); err != nil {
        return nil, err
    }

    return NewRollupBlockBuilder(db, blockStore, blockSubmitter, maxBlockTime, maxBlockGas, maxBlockTransactions)
}

func getSubmitChBlockStoreAndSubmitter() (chan *RollupBlock, *TestBlockStore, *TestBlockSubmitter) {
   submitCh := make(chan *RollupBlock, 10)
   return submitCh, newTestBlockStore(make([]*types.Block, 0)), newTestBlockSubmitter(make([]*RollupBlock, 0), submitCh)
}

/***************
 * Tests Start *
 ***************/

// Single block submission tests

func TestBlockSubmissionMaxTransactions(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 1)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(1, 1, true)
    blockBuilder.NewBlock(blocks[0])

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
        if len(blockSubmitter.submittedBlocks) > 1 {
            t.Fatal("Expected 1 block to have been submitted", "numSubmitted", len(blockSubmitter.submittedBlocks))
        }
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}

func TestBlockLessThanMaxTransactions(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 2)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(1, 1, true)
    blockBuilder.NewBlock(blocks[0])

    timeout := time.After(timeoutDuration)
    select {
    case <-blockSubmitCh:
        t.Fatalf("should not have submitted a block")
    case <-timeout:
        return
    }
}

func TestBlockSubmissionMaxGas(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()

    blocks := createBlocks(1, 1, true)
    gasLimit := GetBlockRollupGasUsage(blocks[0]) + RollupBlockGasBuffer

    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, gasLimit, 2)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blockBuilder.NewBlock(blocks[0])

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
        if len(blockSubmitter.submittedBlocks) > 1 {
            t.Fatal("Expected 1 block to have been submitted", "numSubmitted", len(blockSubmitter.submittedBlocks))
        }
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}

func TestBlockLessThanMaxGas(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()

    blocks := createBlocks(1, 1, true)
    gasLimit := GetBlockRollupGasUsage(blocks[0]) + RollupBlockGasBuffer + MinTxGas

    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, gasLimit, 2)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blockBuilder.NewBlock(blocks[0])

    timeout := time.After(timeoutDuration)
    select {
    case <-blockSubmitCh:
        t.Fatalf("should not have submitted a block")
    case <-timeout:
        return
    }
}


// Multiple block submission tests

func TestMultipleBlockSubmissionMaxTransactions(t *testing.T) {
  blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
  blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 1)
  if err != nil {
      t.Fatalf("unable to make test block builder, error: %v", err)
  }

  blocks := createBlocks(2, 1, true)
  blockBuilder.NewBlock(blocks[0])
  blockBuilder.NewBlock(blocks[1])

  timeout := time.After(timeoutDuration)
  select {
  case rollupBlock := <-blockSubmitCh:
      assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
      time.Sleep(time.Microsecond * 10)
      if len(blockSubmitter.submittedBlocks) != 2 {
          t.Fatal("Expected 2 block to have been submitted", "numSubmitted", len(blockSubmitter.submittedBlocks))
      }
      assertTransitionFromBlock(t, blockSubmitter.submittedBlocks[1].transitions[0], blocks[1])
      return
  case <-timeout:
      t.Fatalf("test timeout")
  }
}

func TestMultipleBlocksLessThanMaxTransactions(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 3)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(2, 1, true)
    blockBuilder.NewBlock(blocks[0])
    blockBuilder.NewBlock(blocks[1])

    timeout := time.After(timeoutDuration)
    select {
    case <-blockSubmitCh:
        t.Fatalf("should not have submitted a block")
    case <-timeout:
        return
    }
}

func TestMultipleBlockSubmissionMaxGas(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()

    blocks := createBlocks(2, 1, true)
    gasLimit := GetBlockRollupGasUsage(blocks[0]) + RollupBlockGasBuffer

    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, gasLimit, 3)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blockBuilder.NewBlock(blocks[0])
    blockBuilder.NewBlock(blocks[1])

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
        time.Sleep(time.Microsecond * 10)
        if len(blockSubmitter.submittedBlocks) != 2 {
            t.Fatal("Expected 2 block to have been submitted", "numSubmitted", len(blockSubmitter.submittedBlocks))
        }
        assertTransitionFromBlock(t, blockSubmitter.submittedBlocks[1].transitions[0], blocks[1])
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}

func TestMultipleBlocksLessThanMaxGas(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()

    blocks := createBlocks(2, 1, true)
    gasLimit := 2*(GetBlockRollupGasUsage(blocks[0]) + RollupBlockGasBuffer + MinTxGas)

    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, gasLimit, 3)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blockBuilder.NewBlock(blocks[0])
    blockBuilder.NewBlock(blocks[1])

    timeout := time.After(timeoutDuration)
    select {
    case <-blockSubmitCh:
        t.Fatalf("should not have submitted a block")
    case <-timeout:
        return
    }
}

// Empty block tests

func TestEmptyBlocksIgnored(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 1)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(2, 1, false)
    blockBuilder.NewBlock(blocks[0])
    blockBuilder.NewBlock(blocks[1])

    timeout := time.After(timeoutDuration)
    select {
    case <-blockSubmitCh:
        t.Fatalf("should not have submitted a block")
    case <-timeout:
        return
    }
}

func TestEmptyBlocksIgnoredWithNonEmpty(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()

    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Minute * 1, 1_000_000_000, 1)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    emptyBlocks := createBlocks(2, 1, false)

    blockBuilder.NewBlock(emptyBlocks[0])
    blockBuilder.NewBlock(emptyBlocks[1])

    nonEmpty :=createBlocks(1, 3, true)[0]
    blockBuilder.NewBlock(nonEmpty)

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], nonEmpty)
        if len(blockSubmitter.submittedBlocks) > 1 {
            t.Fatal("Expected 1 block to have been submitted", "numSubmitted", len(blockSubmitter.submittedBlocks))
        }
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}

// timer submission

func TestBlockSubmissionMaxTimeBetweenBlocks(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Microsecond * 1, 1_000_000_000, 10)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(2, 1, true)
    blockBuilder.NewBlock(blocks[0])
    blockBuilder.NewBlock(blocks[1])

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
        time.Sleep(time.Microsecond * 10)
        if len(blockSubmitter.submittedBlocks) != 2 && len(rollupBlock.transitions) != 2 {
            t.Fatal("Expected 2 transitions to have been submitted", "blocksSubmitted", len(blockSubmitter.submittedBlocks), "transitionsInFirst", len(rollupBlock.transitions))
        }
        var secondTransition *Transition
        switch true {
        case len(blockSubmitter.submittedBlocks) == 2:
            secondTransition = blockSubmitter.submittedBlocks[1].transitions[0]
        case len(rollupBlock.transitions) == 2:
            secondTransition = rollupBlock.transitions[1]
        }
        assertTransitionFromBlock(t, secondTransition, blocks[1])
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}


func TestBlockSubmissionMaxTimeBetweenBlocksReset(t *testing.T) {
    blockSubmitCh, blockStore, blockSubmitter := getSubmitChBlockStoreAndSubmitter()
    blockBuilder, err := newTestBlockBuilder(blockStore, blockSubmitter, 0, time.Microsecond * 1, 1_000_000_000, 10)
    if err != nil {
        t.Fatalf("unable to make test block builder, error: %v", err)
    }

    blocks := createBlocks(2, 1, true)
    blockBuilder.NewBlock(blocks[0])

    timeout := time.After(timeoutDuration)
    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[0])
        if len(blockSubmitter.submittedBlocks) != 1{
            t.Fatal("Expected 1 block to have been submitted", "blocksSubmitted", len(blockSubmitter.submittedBlocks))
        }
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }

    blockBuilder.NewBlock(blocks[1])

    select {
    case rollupBlock := <-blockSubmitCh:
        assertTransitionFromBlock(t, rollupBlock.transitions[0], blocks[1])
        if len(blockSubmitter.submittedBlocks) != 2{
            t.Fatal("Expected 2 blocks to have been submitted", "blocksSubmitted", len(blockSubmitter.submittedBlocks))
        }
        return
    case <-timeout:
        t.Fatalf("test timeout")
    }
}