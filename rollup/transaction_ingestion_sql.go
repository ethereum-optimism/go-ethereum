package rollup

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
)

const (
	SQLMaxGethSubmissionQueueIndex = `
SELECT MAX(queue_index) as queue_index
FROM geth_submission_queue`
)

type QueuedTransaction struct {
	GethSubmissionQueueIndex uint32 `db:"geth_submission_queue_index"`
	QueueIndex               int64  `db:"queue_index"`
	Target                   string `db:"target"`
	Calldata                 string `db:"calldata"`
	BlockTimestamp           uint32 `db:"block_timestamp"`
	BlockNumber              uint32 `db:"block_number"`
	L1TxHash                 string `db:"l1_tx_hash"`
	L1TxIndex                uint16 `db:"l1_tx_index"`
	L1TxLogIndex             uint16 `db:"l1_tx_log_index"`
	QueueOrigin              uint16 `db:"queue_origin"`
	Sender                   string `db:"sender"`
	L1MessageSender          string `db:"l1_message_sender"`
	GasLimit                 uint64 `db:"gas_limit"`
	Nonce                    uint64 `db:"nonce"`
	Signature                string `db:"signature"`
}

func GetMostRecentQueuedTransaction(db *sqlx.DB) (*types.Transaction, error) {
	tx := QueuedTransaction{}
	err := db.Get(&tx, SQLMaxGethSubmissionQueueIndex)
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
