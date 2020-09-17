package rollup

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var err error

func init() {
	conn := "user=test password=test sslmode=disable dbname=rollup"

	db, err = sqlx.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
}

func insertMockBlock(db *sqlx.DB, blockNumber uint64) error {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	blockHash := common.Bytes2Hex(b)

	stmt := `
INSERT INTO l1_block (block_hash, parent_hash, block_number, block_timestamp, gas_limit, gas_used, processed)
VALUES ($1, $2, $3, 0, 0, 0, true)`

	_, err := db.Exec(stmt, blockHash, blockHash, blockNumber)
	return err
}

func insertMockL1Tx(db *sqlx.DB, hash string, blockNumber uint64, blockIndex uint64) error {
	addr := "0000000000000000000000000000000000000000"
	stmt := `
INSERT INTO l1_tx (
    block_number, tx_index, tx_hash, from_address, to_address,
	nonce, gas_limit, gas_price, calldata, signature
) VALUES ($1, $2, $3, $4, $5, $6, 0, 0, 00, 00)`

	_, err := db.Exec(stmt, blockNumber, blockIndex, hash, addr, addr, blockNumber)
	return err
}

func insertMockL1RollupTx(db *sqlx.DB, hash string) error {
	addr := "0000000000000000000000000000000000000000"
	stmt := `
INSERT INTO l1_rollup_tx (
	sender, l1_message_sender, target, calldata, queue_origin, nonce,
	gas_limit, signature, l1_tx_hash, l1_tx_index, l1_tx_log_index
) VALUES ($1, $2, $3, 00, 0, 0, 0, 0, $4, 0, 0)`

	_, err := db.Exec(stmt, addr, addr, addr, hash)
	return err
}

func insertUnqueuedRollupTxMock(db *sqlx.DB, txid string, height uint64) error {
	var err error
	err = insertMockBlock(db, height)
	err = insertMockL1Tx(db, txid, height, 0)
	err = insertMockL1RollupTx(db, txid)

	return err
}

func TestGetMostRecentUnqueuedRollupTx(t *testing.T) {
	t.Skip()

	origins := []types.QueueOrigin{types.QueueOriginL1ToL2}

	// Create a random txid
	txid := randomHex(32)

	// Insert mock data
	insertUnqueuedRollupTxMock(db, txid, 0)

	// Call the function being tested
	tx, err := GetMostRecentUnqueuedRollupTx(db, origins)
	if err != nil {
		t.Fatal(err)
	}

	if tx == nil {
		t.Fatal("transaction entry not found")
	}

	// Make sure the correct data is being returned
	if tx.L1TxHash != txid {
		t.Fatalf("txid doesn't match. Got %s, expected %s", tx.L1TxHash, txid)
	}

	// Clean up the database
	db.Exec("TRUNCATE l1_rollup_tx, l1_tx, l1_block CASCADE")
}

func TestGetMaxGethSubmissionQueueIndex(t *testing.T) {
	t.Skip()

	cases := []struct {
		txid  string
		index uint64
	}{
		{
			txid:  randomHex(32),
			index: 10,
		},
		{
			txid:  randomHex(32),
			index: 5,
		},
		{
			txid:  randomHex(32),
			index: 0,
		},
	}

	big, err := rand.Int(rand.Reader, big.NewInt(32))
	if err != nil {
		t.Fatal(err)
	}
	blockNumber := big.Uint64()

	// insert a mock block
	err = insertMockBlock(db, blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	// insert mock transactoins and corresponding submission queue entries
	for i, c := range cases {
		err := insertMockL1Tx(db, c.txid, blockNumber, uint64(i))
		if err != nil {
			t.Fatal(err)
		}
		_, err = InsertGethSubmissionQueueEntry(db, c.txid, c.index)
		if err != nil {
			t.Fatal(err)
		}
	}

	max, err := GetMaxGethSubmissionQueueIndex(db)
	if err != nil {
		t.Fatal(err)
	}

	// 10 is the max `index` value in the table test above
	if max != 10 {
		t.Fatalf("Unexpected max value: %d, expected %d", max, 10)
	}

	db.Exec("TRUNCATE geth_submission_queue, l1_tx, l1_block CASCADE")
}

func TestUpdateGethSubmissionQueueIndex(t *testing.T) {
	t.Skip()

	txid := randomHex(32)
	insertUnqueuedRollupTxMock(db, txid, 0)

	_, err = InsertGethSubmissionQueueEntry(db, txid, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err := UpdateGethSubmissionQueueIndex(db, 0, txid, 0)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: test to make sure that the submission queue
	// index is actually updated

	db.Exec("TRUNCATE geth_submission_queue, l1_tx, l1_block CASCADE")
}

func TestGetMostRecentQueuedTransaction(t *testing.T) {
	t.Skip()

	txid := randomHex(32)
	insertUnqueuedRollupTxMock(db, txid, 0)

	_, err = InsertGethSubmissionQueueEntry(db, txid, 0)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := GetMostRecentQueuedTransaction(db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", tx)

	db.Exec("TRUNCATE geth_submission_queue, l1_tx, l1_block CASCADE")
}

func TestUpdateSentSubmissionStatus(t *testing.T) {
	t.Skip()

	txid := randomHex(32)
	insertUnqueuedRollupTxMock(db, txid, 0)

	_, err := InsertGethSubmissionQueueEntry(db, txid, 0)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: finish
	err = UpdateSentSubmissionStatus(db, 0)
	if err != nil {
		t.Fatal(err)
	}

	queue := []GethSubmissionQueueEntry{}
	err = db.Select(queue, "SELECT * FROM geth_submission_queue")
	fmt.Printf("%#v\n", queue)

	db.Exec("TRUNCATE geth_submission_queue, l1_tx, l1_block CASCADE")
}

func randomHex(size uint) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	return common.ToHex(b)
}
