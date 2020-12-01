package diffdb

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/mattn/go-sqlite3"
	"math/big"
)

type Key struct {
	Key     common.Hash
	Mutated bool
}

type Diff map[common.Address][]Key

type DiffDb struct {
	db *sql.DB
}

var insertStatement = `
INSERT INTO diffs
    (block, address, key, mutated)
    VALUES
    ($1, $2, $3, $4)
`
var createStmt = `
CREATE TABLE IF NOT EXISTS diffs (
    block INTEGER,
    address STRING,
    key STRING,
    mutated BOOL
)
`
var selectStmt = `
SELECT * from diffs WHERE block = $1
`

/// Inserts a new row to the sqlite with the provided diff data.
func (diff *DiffDb) SetDiffKey(block *big.Int, address common.Address, key common.Hash, mutated bool) error {
	_, err := diff.db.Exec(insertStatement, block.Uint64(), address, key, mutated)
	return err
}

/// Gets all the rows for the matching block and converts them to a Diff map.
func (diff *DiffDb) GetDiff(blockNum *big.Int) (Diff, error) {
	// make the query
	rows, err := diff.db.Query(selectStmt, blockNum.Uint64())
	if err != nil {
		return nil, err
	}

	// initialize our data
	res := make(Diff)
	var block uint64
	var address common.Address
	var key common.Hash
	var mutated bool
	for rows.Next() {
		// deserialize the line
		err = rows.Scan(&block, &address, &key, &mutated)
		if err != nil {
			return nil, err
		}
		// add the data to the map
		res[address] = append(res[address], Key{key, mutated})
	}

	return res, rows.Err()
}

func NewDiffDb(path string) (*DiffDb, error) {
	// get a handle
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// create the table if it does not exist
	_, err = db.Exec(createStmt)
	if err != nil {
		return nil, err
	}

	// retturn
	return &DiffDb{db: db}, nil
}
