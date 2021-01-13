package eth

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rollup"
)

func TestGasLimit(t *testing.T) {
	backend := &EthAPIBackend{
		extRPCEnabled:    false,
		eth:              nil,
		gpo:              nil,
		verifier:         false,
		DisableTransfers: false,
		GasLimit:         0,
		UsingOVM:         true,
	}

	nonce := uint64(0)
	to := common.HexToAddress("0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c")
	value := big.NewInt(0)
	gasPrice := big.NewInt(0)
	data := []byte{}
	qo := types.QueueOriginSequencer
	sighash := types.SighashEIP155

	// Set the gas limit to 1 so that the transaction will not be
	// able to be added.
	gasLimit := uint64(1)
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data, nil, nil, qo, sighash)

	err := backend.SendTx(context.Background(), tx)
	if err == nil {
		t.Fatal("Transaction with too large of gas limit accepted")
	}
	if err.Error() != fmt.Sprintf("Transaction gasLimit (%d) is greater than max gasLimit (%d)", gasLimit, backend.GasLimit) {
		t.Fatalf("Unexpected error type: %s", err)
	}
}

func TestSetHead(t *testing.T) {
	s, bc, db, err := newTestSyncService()
	if err != nil {
		t.Fatal(err)
	}

	backend := &EthAPIBackend{
		eth: &Ethereum{
			syncService: s,
			blockchain:  bc,
		},
		UsingOVM: true,
	}

	// Insert a chain of 3
	blocks := makeBlockChain(bc.CurrentBlock(), 3, ethash.NewFaker(), *db, 10)
	_, err = bc.InsertChain(blocks)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the current block is height 3
	current := backend.CurrentBlock()
	if current.NumberU64() != 3 {
		t.Fatal("Current doesn't match")
	}

	// Set the head to 1 and then make sure the tip is on 1
	backend.SetHead(1)
	current = backend.CurrentBlock()
	if current.NumberU64() != 1 {
		t.Fatal("Current doesn't match")
	}
}

func newTestSyncService() (*rollup.SyncService, *core.BlockChain, *ethdb.Database, error) {
	chainCfg := params.AllEthashProtocolChanges
	chainID := big.NewInt(420)
	chainCfg.ChainID = chainID

	engine := ethash.NewFaker()
	db := rawdb.NewMemoryDatabase()
	_ = new(core.Genesis).MustCommit(db)
	chain, err := core.NewBlockChain(db, nil, chainCfg, engine, vm.Config{}, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Cannot initialize blockchain: %w", err)
	}
	chaincfg := params.ChainConfig{ChainID: chainID}

	txPool := core.NewTxPool(core.TxPoolConfig{PriceLimit: 0}, &chaincfg, chain)
	cfg := rollup.Config{
		IsVerifier: false,
	}

	service, err := rollup.NewSyncService(context.Background(), cfg, txPool, chain, db)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Cannot initialize syncservice: %w", err)
	}

	return service, chain, &db, nil
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *types.Block, n int, engine consensus.Engine, db ethdb.Database, seed int) []*types.Block {
	blocks, _ := core.GenerateChain(params.TestChainConfig, parent, engine, db, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}
