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

// TODO: look at instantiation of TxPool
// is NewTxPool(&config) pattern used?
// could add various config options to the
// TxIngestion struct
// Add:
//   - db
//   - logger

type TxIngestion struct {
	loopTicker    *time.Ticker
	db            *sqlx.DB
	signer        types.Signer
	bc            *core.BlockChain
	currentState  *state.StateDB // Current state in the blockchain head
	currentMaxGas uint64         // Current gas limit for transaction caps
}

// What is the correct thing to pass in here?
// Processor? Might need to pull some logic out of the
// TxPool
// TxPool.addTx calls TxPool.verifyTx
func NewTxIngestion(cfg Config, chaincfg *params.ChainConfig, chain *core.BlockChain) *TxIngestion {
	interval := cfg.TxIngestionPollInterval * time.Second

	txIngestion := TxIngestion{
		loopTicker: time.NewTicker(interval),
		signer:     types.NewEIP155Signer(chaincfg.ChainID), // should be NewOVMSigner
		bc:         chain,
	}

	if cfg.IsTxIngestionEnabled() {
		conn := txIngestion.makeConn(&cfg)
		db, err := sqlx.Connect("postgres", conn)
		if err != nil {
			log.Error("Cannot connect to postgres", "msg", err.Error())
			return nil
		}

		txIngestion.db = db

		head := txIngestion.bc.CurrentBlock().Header()
		statedb, err := txIngestion.bc.StateAt(head.Root)
		if err != nil {
			log.Error("Cannot get statedb", "msg", err.Error())
			return nil
		}

		txIngestion.currentState = statedb
		txIngestion.currentMaxGas = head.GasLimit // this needs to be updated

		go txIngestion.loop()
	}

	return &txIngestion
}

func (t *TxIngestion) loop() {
	for range t.loopTicker.C {
		tx, err := GetMostRecentQueuedTransaction(t.db)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		fmt.Println(tx)

		err = t.validateTx(tx)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// There is another codepath:
		// core.ApplyTransaction

		header := types.Header{}
		txs := []*types.Transaction{}
		block := types.NewBlock(&header, txs, []*types.Header{}, []*types.Receipt{})

		// Potentially fetch state like this instead
		// of using t.currentState
		// statedb, err := state.New(parent.Root, bc.stateCache)

		processor, vmConfig := t.bc.Processor(), t.bc.GetVMConfig()
		if vmConfig == nil {
			log.Error("No blockchain.VMConfig")
			continue
		}
		// usedGas is used in bc.verifyer.Verify, which checks the output of
		// processor.Process against the block. I don't think we need to compute
		// the receipts root for the block?
		receipts, logs, _, err := processor.Process(block, t.currentState, *vmConfig)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		status, err := t.bc.WriteBlockWithState(block, receipts, logs, t.currentState, false)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if status == core.NonStatTy {
			log.Error("Unable to write block to disk")
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
	// Ensure the
	// transaction
	// doesn't exceed
	// the current block
	// limit gas.
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
	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true, true)
	if err != nil {
		return err
	}

	if tx.Gas() < intrGas {
		return core.ErrIntrinsicGas
	}

	return nil
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
