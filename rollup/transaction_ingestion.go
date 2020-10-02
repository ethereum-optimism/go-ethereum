package rollup

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TxIngestion struct {
	loopTicker *time.Ticker
	db         *sqlx.DB
	signer     types.Signer
	key        *ecdsa.PrivateKey
	txpool     *core.TxPool
}

// Should this have a safety check on cfg.TxIngestionPollInterval?
func NewTxIngestion(cfg Config, chaincfg *params.ChainConfig, txpool *core.TxPool) *TxIngestion {
	if cfg.TxIngestionSignerKey == nil {
		cfg.TxIngestionSignerKey, _ = crypto.GenerateKey()
	}

	txIngestion := TxIngestion{
		signer:     types.NewOVMSigner(chaincfg.ChainID),
		txpool:     txpool,
		loopTicker: time.NewTicker(cfg.TxIngestionPollInterval),
		key:        cfg.TxIngestionSignerKey,
	}

	if cfg.IsTxIngestionEnabled() {
		log.Info("Transaction ingestion connecting to database", "host", cfg.TxIngestionDBHost, "port", cfg.TxIngestionDBPort)

		// TODO: If geth is closing gracefully, need to get event
		// and pass through to cancel this context
		ctx := context.Background()
		db, err := txIngestion.DBConnectWithRetry(ctx, &cfg)
		if err != nil {
			panic("Timed out connecting to database")
		}
		log.Info("TxIngestion connected to database")
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
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (t *TxIngestion) applyTransaction(tx *types.Transaction) error {
	return t.txpool.AddLocal(tx)
}

func (t *TxIngestion) loop() {
	if t.key == nil {
		panic("Transaction ingestion requires a private key")
	}

	hex := hexutil.Encode(crypto.FromECDSAPub(&t.key.PublicKey))
	address := crypto.PubkeyToAddress(t.key.PublicKey)
	log.Info("Starting transaction ingestion", "key", hex, "address", address.Hex())

	for range t.loopTicker.C {
		txs, index, err := GetMostRecentQueuedTransactions(t.db)
		if err != nil {
			log.Error("Error getting most recently queued transactions: " + err.Error())
			continue
		}

		for i, tx := range txs {
			log.Debug("Transaction Ingestion", "hash", tx.Hash().Hex(), "submission index", index, "element", i)

			nonce := t.txpool.Nonce(address)
			tx.SetNonce(nonce)
			tx, err := types.SignTx(tx, t.signer, t.key)

			if err != nil {
				log.Error("Cannot sign transaction", "hash", tx.Hash().Hex(), "message", err.Error())
				continue
			}

			err = t.applyTransaction(tx)
			if err != nil {
				log.Error("Cannot apply transaction", "hash", tx.Hash().Hex(), "message", err.Error())
				continue
			}
		}

		err = UpdateSentSubmissionStatus(t.db, "Sent", index)
		if err != nil {
			// TODO(mark): this should probably panic to prevent playing the
			// same transaction twice, prefer safety over liveliness
			log.Error("Cannot update submission status", "message", err.Error())
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
