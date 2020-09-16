package rollup

// TODO:
//    config.IsByzantium - always return true
//    config.IsEIP158 - always return true
//    ChainConfig.IsHomestead ?
//	  ChainConfig.IsIstanbul

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/jmoiron/sqlx"
)

const (
	// txSlotSize is used to calculate how many data slots a single
	// transaction takes up based on its size. The slots are used as DoS
	// protection, ensuring that validating a new transaction remains
	// a constant operation (in reality O(maxslots), where max slots are
	// 4 currently).
	txSlotSize = 32 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder
	// and more expensive to propagate; larger transactions also take
	// more resources to validate whether they fit into the pool or not.
	txMaxSize = 2 * txSlotSize // 64KB, don't bump without EIP-2464 support
)

type TxIngestion struct {
	loopTicker    *time.Ticker
	db            *sqlx.DB
	signer        types.Signer
	bc            *core.BlockChain
	txpool        *core.TxPool
	currentState  *state.StateDB // Current state in the blockchain head
	currentMaxGas uint64         // Current gas limit for transaction caps
	vmConfig      *vm.Config
}

// Should this have a safety check on cfg.TxIngestionPollInterval?
func NewTxIngestion(cfg Config, chaincfg *params.ChainConfig, chain *core.BlockChain, txpool *core.TxPool) *TxIngestion {
	vmConfig := chain.GetVMConfig()
	if vmConfig == nil {
		log.Error("no Blockchain.VMConfig")
		return nil
	}

	txIngestion := TxIngestion{
		loopTicker: time.NewTicker(cfg.TxIngestionPollInterval),
		signer:     types.NewOVMSigner(chaincfg.ChainID),
		bc:         chain,
		txpool:     txpool,
		vmConfig:   vmConfig,
	}

	head := txIngestion.bc.CurrentBlock().Header()
	statedb, err := txIngestion.bc.StateAt(head.Root)
	if err != nil {
		log.Error("Cannot get statedb", "msg", err.Error())
		return nil
	}

	txIngestion.currentState = statedb
	txIngestion.currentMaxGas = head.GasLimit // this needs to be updated

	if cfg.IsTxIngestionEnabled() {
		conn := txIngestion.makeConn(&cfg)
		db, err := sqlx.Connect("postgres", conn)
		if err != nil {
			log.Error("Cannot connect to postgres", "msg", err.Error())
			return nil
		}

		txIngestion.db = db

		go txIngestion.loop()
	}

	return &txIngestion
}

func (t *TxIngestion) applyTransaction(tx *types.Transaction) error {
	err := t.validateTx(tx)
	if err != nil {
		return err
	}

	return t.txpool.AddLocal(tx)

	/*
		header := types.Header{GasLimit: 210000000000}
		txs := []*types.Transaction{tx}
		// how do i set the gas limit on the block?
		block := types.NewBlock(&header, txs, []*types.Header{}, []*types.Receipt{})

		processor := t.bc.Processor()
		// usedGas is used in bc.verifyer.Verify, which checks the output of
		// processor.Process against the block. I don't think we need to compute
		// the receipts root for the block?
		receipts, logs, _, err := processor.Process(block, t.currentState, *t.vmConfig)
		if err != nil {
			return err
		}

		status, err := t.bc.WriteBlockWithState(block, receipts, logs, t.currentState, false)
		if err != nil {
			return err
		}

		if status == core.NonStatTy {
			return fmt.Errorf("Unable to write tx %x to disk", tx.Hash())
		}

		return nil
	*/
}

func (t *TxIngestion) loop() {
	for range t.loopTicker.C {
		tx, err := GetMostRecentQueuedTransaction(t.db)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// TODO: some sort of performance metering
		err = t.applyTransaction(tx)

		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
}

// validateTx checks whether a transaction is valid according to the
// consensus rules
func (t *TxIngestion) validateTx(tx *types.Transaction) error {
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > txMaxSize {
		return core.ErrOversizedData
	}
	// Transactions can't be negative. This may
	// never happen using RLP decoded
	// transactions but may occur if you
	// create a transaction using the
	// RPC.
	if tx.Value().Sign() < 0 {
		return core.ErrNegativeValue
	}
	// Ensure the transaction doesn't exceed
	// the current block limit gas.
	if t.currentMaxGas < tx.Gas() {
		return core.ErrGasLimit
	}
	// Make sure the transaction is signed properly
	from, err := types.Sender(t.signer, tx)
	if err != nil {
		return core.ErrInvalidSender
	}
	// Ensure the transaction adheres to nonce ordering
	if t.currentState.GetNonce(from) > tx.Nonce() {
		return core.ErrNonceTooLow
	}
	// Transactor should have enough funds
	// to cover the costs cost == V + GP * GL
	if t.currentState.GetBalance(from).Cmp(tx.Cost()) < 0 {
		return core.ErrInsufficientFunds
	}
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := t.IntrinsicGas()
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return core.ErrIntrinsicGas
	}
	return nil
}

func (t *TxIngestion) IntrinsicGas() (uint64, error) {
	return 0, nil
}

func (t *TxIngestion) Stop() {
	t.loopTicker.Stop()
}

func (t *TxIngestion) makeConn(cfg *Config) string {
	str := "host=%s port=%s user=%s password=%s dbname=%s"
	host := cfg.TxIngestionDBHost
	port := cfg.TxIngestionDBPort
	user := cfg.TxIngestionDBUser
	pass := cfg.TxIngestionDBPassword
	name := cfg.TxIngestionDBName

	return fmt.Sprintf(str, host, port, user, pass, name)
}
