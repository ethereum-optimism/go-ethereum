// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	// FreezerHeaderTable indicates the name of the freezer header table.
	FreezerHeaderTable = "headers"

	// FreezerHashTable indicates the name of the freezer canonical hash table.
	FreezerHashTable = "hashes"

	// FreezerBodiesTable indicates the name of the freezer block body table.
	FreezerBodiesTable = "bodies"

	// FreezerReceiptTable indicates the name of the freezer receipts table.
	FreezerReceiptTable = "receipts"

	// FreezerDifficultyTable indicates the name of the freezer total difficulty table.
	FreezerDifficultyTable = "diffs"

	// ancient append Postgres statements
	appendAncientHeaderPgStr   = "INSERT INTO ancient_headers (block_number, header) VALUES ($1, $2) ON CONFLICT (block_number) DO UPDATE SET header = $2"
	appendAncientHashPgStr     = "INSERT INTO ancient_hashes (block_number, hash) VALUES ($1, $2) ON CONFLICT (block_number) DO UPDATE SET hash = $2"
	appendAncientBodyPgStr     = "INSERT INTO ancient_bodies (block_number, body) VALUES ($1, $2) ON CONFLICT (block_number) DO UPDATE SET body = $2"
	appendAncientReceiptsPgStr = "INSERT INTO ancient_receipts (block_number, receipts) VALUES ($1, $2) ON CONFLICT (block_number) DO UPDATE SET receipts = $2"
	appendAncientTDPgStr       = "INSERT INTO ancient_tds (block_number, td) VALUES ($1, $2) ON CONFLICT (block_number) DO UPDATE SET td = $2"

	// ancient truncate Postgres statements
	truncateAncientHeaderPgStr   = "DELETE FROM ancient_headers WHERE block_number > $1"
	truncateAncientHashPgStr     = "DELETE FROM ancient_hashes WHERE block_number > $1"
	truncateAncientBodiesPgStr   = "DELETE FROM ancient_bodies WHERE block_number > $1"
	truncateAncientReceiptsPgStr = "DELETE FROM ancient_receipts WHERE block_number > $1"
	truncateAncientTDPgStr       = "DELETE FROM ancient_tds WHERE block_number > $1"

	// ancient size Postgres statement
	ancientSizePgStr = "SELECT pg_total_relation_size($1)"

	// ancients Postgres statement
	ancientsPgStr = "SELECT block_number FROM ancient_headers ORDER BY block_number DESC LIMIT 1"

	// ancient has Postgres statements
	hasAncientHeaderPgStr   = "SELECT exists(SELECT 1 FROM ancient_headers WHERE block_number = $1)"
	hasAncientHashPgStr     = "SELECT exists(SELECT 1 FROM ancient_hashes WHERE block_number = $1)"
	hasAncientBodyPgStr     = "SELECT exists(SELECT 1 FROM ancient_bodies WHERE block_number = $1)"
	hasAncientReceiptsPgStr = "SELECT exists(SELECT 1 FROM ancient_receipts WHERE block_number = $1)"
	hasAncientTDPgStr       = "SELECT exists(SELECT 1 FROM ancient_tds WHERE block_number = $1)"

	// ancient get Postgres statements
	getAncientHeaderPgStr   = "SELECT header FROM ancient_headers WHERE block_number = $1"
	getAncientHashPgStr     = "SELECT hash FROM ancient_hashes WHERE block_number = $1"
	getAncientBodyPgStr     = "SELECT body FROM ancient_bodies WHERE block_number = $1"
	getAncientReceiptsPgStr = "SELECT receipts FROM ancient_receipts WHERE block_number = $1"
	getAncientTDPgStr       = "SELECT td FROM ancient_tds WHERE block_number = $1"
)

// HasAncient satisfies the ethdb.AncientReader interface
// HasAncient returns an indicator whether the specified data exists in the ancient store
func (d *Database) HasAncient(kind string, number uint64) (bool, error) {
	var pgStr string
	switch kind {
	case FreezerHeaderTable:
		pgStr = hasAncientHeaderPgStr
	case FreezerHashTable:
		pgStr = hasAncientHashPgStr
	case FreezerBodiesTable:
		pgStr = hasAncientBodyPgStr
	case FreezerReceiptTable:
		pgStr = hasAncientReceiptsPgStr
	case FreezerDifficultyTable:
		pgStr = hasAncientTDPgStr
	default:
		return false, fmt.Errorf("unexpected ancient kind: %s", kind)
	}
	has := new(bool)
	return *has, d.db.Get(has, pgStr, number)
}

// Ancient satisfies the ethdb.AncientReader interface
// Ancient retrieves an ancient binary blob from the append-only immutable files
func (d *Database) Ancient(kind string, number uint64) ([]byte, error) {
	var pgStr string
	switch kind {
	case FreezerHeaderTable:
		pgStr = getAncientHeaderPgStr
	case FreezerHashTable:
		pgStr = getAncientHashPgStr
	case FreezerBodiesTable:
		pgStr = getAncientBodyPgStr
	case FreezerReceiptTable:
		pgStr = getAncientReceiptsPgStr
	case FreezerDifficultyTable:
		pgStr = getAncientTDPgStr
	default:
		return nil, fmt.Errorf("unexpected ancient kind: %s", kind)
	}
	data := new([]byte)
	return *data, d.db.Get(data, pgStr, number)
}

// Ancients satisfies the ethdb.AncientReader interface
// Ancients returns the ancient item numbers in the ancient store
func (d *Database) Ancients() (uint64, error) {
	num := new(uint64)
	if err := d.db.Get(num, ancientsPgStr); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return *num, nil
}

// AncientSize satisfies the ethdb.AncientReader interface
// AncientSize returns the ancient size of the specified category
func (d *Database) AncientSize(kind string) (uint64, error) {
	var tableName string
	switch kind {
	case FreezerHeaderTable:
		tableName = "eth.ancient_headers"
	case FreezerHashTable:
		tableName = "eth.ancient_hashes"
	case FreezerBodiesTable:
		tableName = "eth.ancient_bodies"
	case FreezerReceiptTable:
		tableName = "eth.ancient_receipts"
	case FreezerDifficultyTable:
		tableName = "eth.ancient_tds"
	default:
		return 0, fmt.Errorf("unexpected ancient kind: %s", kind)
	}
	size := new(uint64)
	return *size, d.db.Get(size, ancientSizePgStr, tableName)
}

// AppendAncient satisfies the ethdb.AncientWriter interface
// AppendAncient injects all binary blobs belong to block at the end of the append-only immutable table files
func (d *Database) AppendAncient(number uint64, hash, header, body, receipts, td []byte) error {
	// append in batch
	var err error
	if d.ancientTx == nil {
		d.ancientTx, err = d.db.Beginx()
		if err != nil {
			return err
		}
	}
	defer func() {
		if err != nil {
			if err := d.ancientTx.Rollback(); err != nil {
				logrus.Error(err)
				d.ancientTx = nil
			}
		}
	}()

	if _, err := d.ancientTx.Exec(appendAncientHashPgStr, number, hash); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(appendAncientHeaderPgStr, number, header); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(appendAncientBodyPgStr, number, body); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(appendAncientReceiptsPgStr, number, receipts); err != nil {
		return err
	}
	_, err = d.ancientTx.Exec(appendAncientTDPgStr, number, td)
	return err
}

// TruncateAncients satisfies the ethdb.AncientWriter interface
// TruncateAncients discards all but the first n ancient data from the ancient store
func (d *Database) TruncateAncients(n uint64) error {
	// truncate in batch
	var err error
	if d.ancientTx == nil {
		d.ancientTx, err = d.db.Beginx()
		if err != nil {
			return err
		}
	}
	defer func() {
		if err != nil {
			if err := d.ancientTx.Rollback(); err != nil {
				logrus.Error(err)
				d.ancientTx = nil
			}
		}
	}()
	if _, err := d.ancientTx.Exec(truncateAncientHeaderPgStr, n); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(truncateAncientHashPgStr, n); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(truncateAncientBodiesPgStr, n); err != nil {
		return err
	}
	if _, err := d.ancientTx.Exec(truncateAncientReceiptsPgStr, n); err != nil {
		return err
	}
	_, err = d.ancientTx.Exec(truncateAncientTDPgStr, n)
	return err
}

// Sync satisfies the ethdb.AncientWriter interface
// Sync flushes all in-memory ancient store data to disk
func (d *Database) Sync() error {
	if d.ancientTx == nil {
		return nil
	}
	if err := d.ancientTx.Commit(); err != nil {
		return err
	}
	d.ancientTx = nil
	return nil
}
