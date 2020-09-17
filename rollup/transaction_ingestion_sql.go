package rollup

import (
	"database/sql"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
)

type QueueOrigin int16

const (
	// Possible `queue_origin` values
	L1ToL2Queue QueueOrigin = 0
	SafetyQueue QueueOrigin = 1
	Sequencer   QueueOrigin = 2
)

const (
	MOST_RECENT_UNQUEUED_ROLLUP_TX = `
SELECT l1_tx_hash, l1_tx_log_index, queue_origin
FROM unqueued_rollup_tx
WHERE queue_origin IN (?)
ORDER BY block_number ASC, l1_tx_index ASC, l1_tx_log_index ASC
LIMIT 1`

	MAX_GETH_SUBMISSION_QUEUE_INDEX = `
SELECT MAX(queue_index) as queue_index
FROM geth_submission_queue`

	INSERT_GETH_SUBMISSION_QUEUE = `
INSERT INTO geth_submission_queue(l1_tx_hash, queue_index)
VALUES ($1, $2)`

	UPDATE_GETH_SUBMISSION_QUEUE_INDEX = `
UPDATE l1_rollup_tx
SET geth_submission_queue_index = $1, index_within_submission = 0
WHERE l1_tx_hash = $2 AND l1_tx_log_index = $3`

	GET_NEXT_QUEUED_GETH_SUBMISSION = `
SELECT
geth_submission_queue_index, target, calldata, block_timestamp, block_number, l1_tx_hash, l1_tx_index,
l1_tx_log_index, queue_origin, sender, l1_message_sender, gas_limit, nonce, signature
FROM next_queued_geth_submission`

	UPDATE_GETH_SUBMISSION_STATUS = `
UPDATE geth_submission_queue
SET status = $1
WHERE queue_index = $2`
)

type QueuedTransaction struct {
	GethSubmissionQueueIndex uint32 `db:"geth_submission_queue_index"`
	QueueIndex               int64  `db:"queue_index"` // ?
	Target                   string `db:"target"`
	Calldata                 string `db:"calldata"`
	BlockTimestamp           uint32 `db:"block_timestamp"`
	BlockNumber              uint32 `db:"block_number"`
	L1TxHash                 string `db:"l1_tx_hash"`
	L1TxIndex                uint16 `db:"l1_tx_index"`
	L1TxLogIndex             uint16 `db:"l1_tx_log_index"`
	QueueOrigin              uint8  `db:"queue_origin"`
	Sender                   string `db:"sender"`
	L1MessageSender          string `db:"l1_message_sender"`
	GasLimit                 uint64 `db:"gas_limit"`
	Nonce                    uint64 `db:"nonce"`
	Signature                string `db:"signature"`
}

type GethSubmissionQueueEntry struct {
	Id         int64  `db:"id"`
	L1TxHash   string `db:"l1_tx_hash"`
	QueueIndex int64  `db:"queue_index"`
	Status     string `db:"status"`
	Created    string `db:"created"`
}

type L1RollupTx struct {
	Id                       int64     `db:"id"`
	Sender                   string    `db:"sender"`
	L1MessageSender          string    `db:"l1_message_sender"`
	Target                   string    `db:"target"`
	Calldata                 string    `db:"calldata"`
	QueueOrigin              int16     `db:"queue_origin"`
	Nonce                    string    `db:"nonce"`     // numeric type
	GasLimit                 string    `db:"gas_limit"` // numeric type
	Signature                string    `db:"signature"`
	GethSubmissionQueueIndex int64     `db:"geth_submission_queue_index"`
	IndexWithinSubmission    int32     `db:"index_within_submission"`
	L1TxHash                 string    `db:"l1_tx_hash"`
	L1TxIndex                int32     `db:"l1_tx_index"`
	L1TxLogIndex             int32     `db:"l1_tx_log_index"`
	Created                  time.Time `db:"created"`
}

type UnqueuedRollupTxEntry struct {
	L1TxHash     string `db:"l1_tx_hash"`
	L1TxLogIndex int32  `db:"l1_tx_log_index"`
	QueueOrigin  int16  `db:"queue_origin"`
}

func GetMostRecentQueuedTransaction(db *sqlx.DB) (*types.Transaction, error) {
	tx := QueuedTransaction{}
	err := db.Get(&tx, MAX_GETH_SUBMISSION_QUEUE_INDEX)
	if err != nil {
		return nil, err
	}

	// Deserialize from database serialization to Transaction serialization
	nonce := uint64(tx.Nonce) // how to best handle numeric type?
	amount, gasPrice := big.NewInt(0), big.NewInt(0)
	gasLimit := uint64(tx.GasLimit)
	data, err := hexutil.Decode(tx.Calldata)
	if err != nil {
		return nil, err
	}

	to := common.HexToAddress(tx.Target)
	l1From := common.HexToAddress(tx.L1MessageSender)
	l1Txid, err := hexutil.DecodeUint64(tx.L1TxHash)
	if err != nil {
		return nil, err
	}
	l1TxHash := hexutil.Uint64(l1Txid)

	// TODO: sighash type needs to be in the database
	sighash := types.SighashEIP155
	queueOrigin := types.QueueOrigin(tx.QueueOrigin)

	txn := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data, &l1From, &l1TxHash, queueOrigin, sighash)

	return txn, nil
}

func UpdateSentSubmissionStatus(db *sqlx.DB, index uint32) error {
	_, err := db.Exec(UPDATE_GETH_SUBMISSION_STATUS, "SENT", index)
	if err != nil {
		return err
	}

	return nil
}

// GetMostRecentUnqueuedRollupTx queries for the most recent item in the
// unqueued_rollup_tx view.
func GetMostRecentUnqueuedRollupTx(db *sqlx.DB, origins []QueueOrigin) (*UnqueuedRollupTxEntry, error) {
	query, args, err := sqlx.In(MOST_RECENT_UNQUEUED_ROLLUP_TX, origins)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	item := UnqueuedRollupTxEntry{}
	err = db.QueryRowx(query, args...).StructScan(&item)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func GetMaxGethSubmissionQueueIndex(db *sqlx.DB) (int64, error) {
	var index int64
	err := db.Get(&index, MAX_GETH_SUBMISSION_QUEUE_INDEX)
	if err != nil {
		return 0, err
	}
	return index, nil
}

// InsertGethSubmissionQueueEntry inserts the L1TxHash and QueueIndex
func InsertGethSubmissionQueueEntry(db *sqlx.DB, hash string, index uint64) (sql.Result, error) {
	return db.Exec(INSERT_GETH_SUBMISSION_QUEUE, hash, index)
}

// TODO(mark): this should parse the sql.Result into a different struct
func UpdateGethSubmissionQueueIndex(db *sqlx.DB, index int64, hash string, logIndex uint64) (sql.Result, error) {
	return db.Exec(UPDATE_GETH_SUBMISSION_QUEUE_INDEX, index, hash, logIndex)
}
