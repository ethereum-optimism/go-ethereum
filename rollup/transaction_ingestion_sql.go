package rollup

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
)

const (
	SQLMaxGethSubmissionQueueIndex = `
SELECT MAX(queue_index) as queue_index
FROM geth_submission_queue`

	SQLGetNextQueuedGethSubmission = `
SELECT
geth_submission_queue_index, target, calldata, block_timestamp, block_number, l1_tx_hash, l1_tx_index,
l1_tx_log_index, queue_origin, sender, l1_message_sender, gas_limit, nonce, signature
FROM next_queued_geth_submission ORDER BY index_within_submission ASC`

	SQLUpdateGethSubmissionStatus = `
UPDATE geth_submission_queue
SET status = $1
WHERE queue_index = $2`
)

type QueuedTransaction struct {
	Id                       uint64         `db:"id"`
	GethSubmissionQueueIndex uint32         `db:"geth_submission_queue_index"`
	QueueIndex               uint64         `db:"queue_index"`
	Target                   string         `db:"target"`
	Calldata                 string         `db:"calldata"`
	BlockTimestamp           uint32         `db:"block_timestamp"`
	BlockNumber              uint32         `db:"block_number"`
	L1TxHash                 string         `db:"l1_tx_hash"`
	L1TxIndex                uint16         `db:"l1_tx_index"`
	L1TxLogIndex             uint16         `db:"l1_tx_log_index"`
	QueueOrigin              uint16         `db:"queue_origin"`
	Sender                   string         `db:"sender"`
	L1MessageSender          string         `db:"l1_message_sender"`
	GasLimit                 uint64         `db:"gas_limit"`
	Nonce                    sql.NullInt64  `db:"nonce"`
	Signature                sql.NullString `db:"signature"`
}

func GetMostRecentQueuedTransactions(db *sqlx.DB) ([]*types.Transaction, uint32, []uint32, error) {
	txs := []QueuedTransaction{}
	err := db.Select(&txs, SQLGetNextQueuedGethSubmission)
	if err != nil {
		return nil, 0, nil, err
	}

	if len(txs) == 0 {
		return []*types.Transaction{}, 0, nil, nil
	}

	log.Debug("Ingesting L1 to L2 Transactions", "count", len(txs))

	transactions := make([]*types.Transaction, len(txs))
	submissionIndices := make([]uint32, len(txs))
	timestamps := make([]uint32, len(txs))

	for i, tx := range txs {
		submissionIndices[i] = tx.GethSubmissionQueueIndex
		timestamps[i] = tx.BlockTimestamp

		nonce := tx.Nonce.Int64
		if !tx.Nonce.Valid {
			nonce = 0
		}

		amount, gasPrice := big.NewInt(0), big.NewInt(0)
		gasLimit := tx.GasLimit
		data, err := hexutil.Decode(tx.Calldata)
		if err != nil {
			return nil, 0, nil, err
		}

		to := common.HexToAddress(tx.Target)
		l1From := common.HexToAddress(tx.L1MessageSender)
		l1TxId := hexutil.Uint64(tx.Id)

		// TODO: sighash type needs to be in the database
		sighash := types.SighashEIP155
		queueOrigin := types.QueueOrigin(tx.QueueOrigin)

		txn := types.NewTransaction(uint64(nonce), to, amount, gasLimit, gasPrice, data, &l1From, &l1TxId, queueOrigin, sighash)

		transactions[i] = txn
	}

	// All of the submission indices must be the same
	// The case in which len(txs) == 0 was checked above.
	submissionIndex := submissionIndices[0]
	for i := 1; i < len(submissionIndices); i++ {
		if submissionIndices[i] != submissionIndex {
			return nil, 0, nil, fmt.Errorf("Submission index mismatch %d and %d", submissionIndex, submissionIndices[i])
		}
	}

	return transactions, submissionIndex, timestamps, nil
}

func UpdateSentSubmissionStatus(db *sqlx.DB, status string, index uint32) error {
	_, err := db.Exec(SQLUpdateGethSubmissionStatus, status, index)
	if err != nil {
		return err
	}

	return nil
}
