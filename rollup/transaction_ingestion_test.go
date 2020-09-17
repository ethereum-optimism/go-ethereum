package rollup

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var (
	key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr    = crypto.PubkeyToAddress(key.PublicKey)
	chainId = big.NewInt(1)
)

func TestApplyTransaction(t *testing.T) {
	cfg := Config{TxIngestionPollInterval: 100}
	chainCfg := params.AllEthashProtocolChanges // wtf is this
	chainCfg.ChainID = chainId

	engine := ethash.NewFaker()
	db := rawdb.NewMemoryDatabase()
	_ = new(core.Genesis).MustCommit(db)

	chain, err := core.NewBlockChain(db, nil, chainCfg, engine, vm.Config{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	chaincfg := params.ChainConfig{ChainID: chainId}

	txPool := core.NewTxPool(core.TxPoolConfig{}, &chaincfg, chain)
	txIngestion := NewTxIngestion(cfg, &chaincfg, chain, txPool)

	signer := types.NewOVMSigner(chainId)
	tx, err := types.SignTx(types.NewTransaction(0, addr, new(big.Int), 21000, new(big.Int), []byte{}, &addr, nil, types.QueueOriginL1ToL2, types.SighashEIP155), signer, key)

	err = txIngestion.applyTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}

	got := txPool.Get(tx.Hash())
	if got == nil {
		t.Fatal("Transaction not found in pool")
	}
}
