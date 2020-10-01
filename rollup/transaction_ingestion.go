package rollup

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TODO: docker-compose, docker entrypoint
// need to inject env variable in for dbhost
// name of postgres service

// --txingestion.enable true
// --txingestion.dbhost
// --txingestion.dbport
// --txingestion.dbuser test
// --txingestion.dbpassword test

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
		// TODO: If geth is closing gracefully, need to get event
		// and pass through to cancel this context
		ctx := context.Background()

		db, err := txIngestion.DBConnectWithRetry(ctx, &cfg)
		if err != nil {
			log.Debug("Could not connect to database")
			return nil
		}

		txIngestion.db = db

		go txIngestion.loop()
	}

	return &txIngestion
}

func (t *TxIngestion) DBConnectWithRetry(ctx context.Context, cfg *Config) (*sqlx.DB, error) {
	connErrCh := make(chan error, 1)
	defer close(connErrCh)

	var db *sqlx.DB
	var err error

	go func() {
		try := 0
		for {
			try++
			db, err = t.DBConnect(cfg)
			if err != nil {
				log.Error("Cannot connect to postgres", "msg", err.Error(), "try", try)
				select {
				case <-ctx.Done():
					break
				case <-time.After(time.Second):
					continue
				}
			}
			break
		}
		connErrCh <- err
	}()

	select {
	case err = <-connErrCh:
		break
	case <-time.After(time.Minute * 3):
		return nil, errors.New("db connection timed out")
	case <-ctx.Done():
		return nil, errors.New("db connection cancelled")
	}

	return db, err
}

func (t *TxIngestion) DBConnect(cfg *Config) (*sqlx.DB, error) {
	conn := DbConnectionString(cfg)
	log.Info("Connecting to postgres", "host", cfg.TxIngestionDBHost, "port", cfg.TxIngestionDBPort)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (t *TxIngestion) applyTransaction(tx *types.Transaction) error {
	return t.txpool.AddLocal(tx)
}

// It would be nice to block on a channel that will
// unblock when the transaction finishes computing the
// state transition.
func (t *TxIngestion) loop() {
	for range t.loopTicker.C {
		// TODO: check to make sure the mempool isn't too full
		// if too full, continue

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

// DbConnectionString resolves Postgres config params to a connection string
func DbConnectionString(cfg *Config) string {
	if len(cfg.TxIngestionDBUser) > 0 && len(cfg.TxIngestionDBPassword) > 0 {
		return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.TxIngestionDBUser, cfg.TxIngestionDBPassword, cfg.TxIngestionDBHost, cfg.TxIngestionDBPort, cfg.TxIngestionDBName)
	}
	if len(cfg.TxIngestionDBUser) > 0 && len(cfg.TxIngestionDBPassword) == 0 {
		return fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable",
			cfg.TxIngestionDBUser, cfg.TxIngestionDBHost, cfg.TxIngestionDBPort, cfg.TxIngestionDBName)
	}
	return fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", cfg.TxIngestionDBHost, cfg.TxIngestionDBPort, cfg.TxIngestionDBName)
}
