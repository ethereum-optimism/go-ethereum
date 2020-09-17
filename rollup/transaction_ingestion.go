package rollup

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/jmoiron/sqlx"
)

type TxIngestion struct {
	loopTicker *time.Ticker
	db         *sqlx.DB
	signer     types.Signer
	txpool     *core.TxPool
}

// Should this have a safety check on cfg.TxIngestionPollInterval?
func NewTxIngestion(cfg Config, chaincfg *params.ChainConfig, txpool *core.TxPool) *TxIngestion {
	txIngestion := TxIngestion{
		signer:     types.NewOVMSigner(chaincfg.ChainID),
		txpool:     txpool,
		loopTicker: time.NewTicker(cfg.TxIngestionPollInterval),
	}

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
	return t.txpool.AddLocal(tx)
}

// It would be nice to block on a channel that will
// unblock when the transaction finishes computing the
// state transition.
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
