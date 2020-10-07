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
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/jmoiron/sqlx"
)

var (
	unsupportedTableTypeErr = errors.New("postgres ethdb: unsupported table")
)

const (
	dbSizePgStr = "SELECT pg_database_size(current_database())"

	hasKVPgStr    = "SELECT exists(SELECT 1 FROM kvstore WHERE eth_key = $1)"
	getKVPgStr    = "SELECT eth_data FROM kvstore WHERE eth_key = $1"
	putKVPgStr    = "INSERT INTO kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3) ON CONFLICT (eth_key) DO NOTHING"
	deleteKVPgStr = "DELETE FROM kvstore WHERE eth_key = $1"

	hasHeaderPgStr = "SELECT exists(SELECT 1 FROM headers WHERE header_key = $1)"
	getHeaderPgStr = "SELECT header FROM headers WHERE header_key = $1"
	putHeaderPgStr = "INSERT INTO headers (header_key, header, height, hash) VALUES ($1, $2, $3, $4) ON CONFLICT (header_key) DO NOTHING"
	deleteHeaderPgStr = "DELETE FROM headers WHERE header_key = $1"

	hasHashPgStr = "SELECT exists(SELECT 1 FROM hashes WHERE hash_key = $1)"
	getHashPgStr = "SELECT hash FROM hashes WHERE hash_key = $1"
	putHashPgStr = "INSERT INTO hashes (hash_key, hash, header_fk) VALUES ($1, $2, $3) ON CONFLICT (hash_key) DO NOTHING"
	deleteHashPgStr = "DELETE FROM hashes WHERE hash_key = $1"

	hasBodyPgStr = "SELECT exists(SELECT 1 FROM bodies WHERE body_key = $1)"
	getBodyPgStr = "SELECT body FROM bodies WHERE body_key = $1"
	putBodyPgStr = "INSERT INTO bodies (body_key, body, header_fk) VALUES ($1, $2, $3) ON CONFLICT (body_key) DO NOTHING"
	deleteBodyPgStr = "DELETE FROM bodies WHERE body_key = $1"

	hasReceiptPgStr = "SELECT exists(SELECT 1 FROM receipts WHERE receipt_key = $1)"
	getReceiptPgStr = "SELECT receipts FROM receipts WHERE receipt_key = $1"
	putReceiptPgStr = "INSERT INTO receipts (receipt_key, receipts, header_fk) VALUES ($1, $2, $3) ON CONFLICT (receipt_key) DO NOTHING"
	deleteReceiptPgStr = "DELETE FROM receipts WHERE receipt_key = $1"

	hasTDPgStr = "SELECT exists(SELECT 1 FROM tds WHERE td_key = $1)"
	getTDPgStr = "SELECT td FROM tds WHERE td_key = $1"
	putTDPgStr = "INSERT INTO tds (td_key, td, header_fk) VALUES ($1, $2, $3) ON CONFLICT (td_key) DO NOTHING"
	deleteTDPgStr = "DELETE FROM tds WHERE td_key = $1"

	hasBloomBitsPgStr = "SELECT exists(SELECT 1 FROM bloom_bits WHERE bb_key = $1)"
	getBloomBitsPgStr = "SELECT bits FROM bloom_bits WHERE bb_key = $1"
	putBloomBitsPgStr = "INSERT INTO bloom_bits (bb_key, bits) VALUES ($1, $2) ON CONFLICT (bb_key) DO NOTHING"
	deleteBloomBitsPgStr = "DELETE FROM bloom_bits WHERE bb_key = $1"

	hasTxLookupPgStr = "SELECT exists(SELECT 1 FROM tx_lookups WHERE lookup_key = $1)"
	getTxLookupPgStr = "SELECT lookup FROM tx_lookups WHERE lookup_key = $1"
	putTxLookupPgStr = "INSERT INTO tx_lookups (lookup_key, lookup) VALUES ($1, $2) ON CONFLICT (lookup_key) DO NOTHING"
	deleteTxLookupPgStr = "DELETE FROM tx_lookups WHERE lookup_key = $1"

	hasPreimagePgStr = "SELECT exists(SELECT 1 FROM preimages WHERE preimage_key = $1)"
	getPreimagePgStr = "SELECT preimage FROM preimages WHERE preimage_key = $1"
	putPreimagePgStr = "INSERT INTO preimages (preimage_key, preimage) VALUES ($1, $2) ON CONFLICT (preimage_key) DO NOTHING"
	deletePreimagePgStr = "DELETE FROM preimages WHERE preimage_key = $1"

	hasNumberPgStr = "SELECT exists(SELECT 1 FROM numbers WHERE number_key = $1)"
	getNumberPgStr = "SELECT number FROM numbers WHERE number_key = $1"
	putNumberPgStr = "INSERT INTO numbers (number_key, number, header_fk) VALUES ($1, $2, $3) ON CONFLICT (number_key) DO NOTHING"
	deleteNumberPgStr = "DELETE FROM numbers WHERE number_key = $1"

	hasConfigPgStr = "SELECT exists(SELECT 1 FROM configs WHERE config_key = $1)"
	getConfigPgStr = "SELECT config FROM configs WHERE config_key = $1"
	putConfigPgStr = "INSERT INTO configs (config_key, config) VALUES ($1, $2) ON CONFLICT (config_key) DO NOTHING"
	deleteConfigPgStr = "DELETE FROM configs WHERE config_key = $1"

	hasBloomIndexPgStr = "SELECT exists(SELECT 1 FROM bloom_indexes WHERE bbi_key = $1)"
	getBloomIndexPgStr = "SELECT index FROM bloom_indexes WHERE bbi_key = $1"
	putBloomIndexPgStr = "INSERT INTO bloom_indexes (bbi_key, index) VALUES ($1, $2) ON CONFLICT (bbi_key) DO NOTHING"
	deleteBloomIndexPgStr = "DELETE FROM bloom_indexes WHERE bbi_key = $1"

	hasTxMetaPgStr = "SELECT exists(SELECT 1 FROM tx_meta WHERE meta_key = $1)"
	getTxMetaPgStr = "SELECT meta FROM tx_meta WHERE meta_key = $1"
	putTxMetaPgStr = "INSERT INTO tx_meta (meta_key, meta) VALUES ($1, $2) ON CONFLICT (meta_key) DO NOTHING"
	deleteTxMetaPgStr = "DELETE FROM tx_meta WHERE meta_key = $1"
)

// Database is the type that satisfies the ethdb.Database and ethdb.KeyValueStore interfaces for PG-IPFS Ethereum data using a direct Postgres connection
type Database struct {
	db        *sqlx.DB
	ancientTx *sqlx.Tx
}

// NewKeyValueStore returns a ethdb.KeyValueStore interface for PG-IPFS
func NewKeyValueStore(db *sqlx.DB) ethdb.KeyValueStore {
	return &Database{
		db: db,
	}
}

// NewDatabase returns a ethdb.Database interface for PG-IPFS
func NewDatabase(db *sqlx.DB) ethdb.Database {
	return &Database{
		db: db,
	}
}

// Has satisfies the ethdb.KeyValueReader interface
// Has retrieves if a key is present in the key-value data store
func (d *Database) Has(key []byte) (bool, error) {
	table, err := ResolveTable(key)
	if err != nil {
		return false, err
	}
	var pgStr string
	switch table {
	case Undefined:
		return false, unsupportedTableTypeErr
	case KVStore:
		pgStr = hasKVPgStr
	case Headers:
		pgStr = hasHeaderPgStr
	case Hashes:
		pgStr = hasHashPgStr
	case Bodies:
		pgStr = hasBodyPgStr
	case Receipts:
		pgStr = hasReceiptPgStr
	case TDs:
		pgStr = hasTDPgStr
	case BloomBits:
		pgStr = hasBloomBitsPgStr
	case TxLookUps:
		pgStr = hasTxLookupPgStr
	case Preimages:
		pgStr = hasPreimagePgStr
	case Numbers:
		pgStr = hasNumberPgStr
	case Configs:
		pgStr = hasConfigPgStr
	case BloomIndexes:
		pgStr = hasBloomIndexPgStr
	case TxMeta:
		pgStr = hasTxMetaPgStr
	}
	var exists bool
	return exists, d.db.Get(&exists, pgStr, key)
}

// Get satisfies the ethdb.KeyValueReader interface
// Get retrieves the given key if it's present in the key-value data store
func (d *Database) Get(key []byte) ([]byte, error) {
	table, err := ResolveTable(key)
	if err != nil {
		return nil, err
	}
	var pgStr string
	switch table {
	case Undefined:
		return nil, unsupportedTableTypeErr
	case KVStore:
		pgStr = getKVPgStr
	case Headers:
		pgStr = getHeaderPgStr
	case Hashes:
		pgStr = getHashPgStr
	case Bodies:
		pgStr = getBodyPgStr
	case Receipts:
		pgStr = getReceiptPgStr
	case TDs:
		pgStr = getTDPgStr
	case BloomBits:
		pgStr = getBloomBitsPgStr
	case TxLookUps:
		pgStr = getTxLookupPgStr
	case Preimages:
		pgStr = getPreimagePgStr
	case Numbers:
		pgStr = getNumberPgStr
	case Configs:
		pgStr = getConfigPgStr
	case BloomIndexes:
		pgStr = getBloomIndexPgStr
	case TxMeta:
		pgStr = getTxMetaPgStr
	}
	var data []byte
	return data, d.db.Get(&data, pgStr, key)
}

// Put satisfies the ethdb.KeyValueWriter interface
// Put inserts the given value into the key-value data store
// Key is expected to be the keccak256 hash of value
func (d *Database) Put(key []byte, value []byte) error {
	prefix, table, num, fk, hash, err := ResolvePutKey(key, value)
	if err != nil {
		return err
	}
	var pgStr string
	args := make([]interface{}, 0, 4)
	args = append(args, key, value)
	switch table {
	case Undefined:
		return unsupportedTableTypeErr
	case KVStore:
		pgStr = putKVPgStr
		args = append(args, prefix)
	case Headers:
		pgStr = putHeaderPgStr
		args = append(args, num, hash)
	case Hashes:
		pgStr = putHashPgStr
		args = append(args, fk)
	case Bodies:
		pgStr = putBodyPgStr
		args = append(args, fk)
	case Receipts:
		pgStr = putReceiptPgStr
		args = append(args, fk)
	case TDs:
		pgStr = putTDPgStr
		args = append(args, fk)
	case BloomBits:
		pgStr = putBloomBitsPgStr
	case TxLookUps:
		pgStr = putTxLookupPgStr
	case Preimages:
		pgStr = putPreimagePgStr
	case Numbers:
		pgStr = putNumberPgStr
		args = append(args, fk)
	case Configs:
		pgStr = putConfigPgStr
	case BloomIndexes:
		pgStr = putBloomIndexPgStr
	case TxMeta:
		pgStr = putTxMetaPgStr
	}
	if _, err = d.db.Exec(pgStr, args...); err != nil {
		return err
	}
	return nil
}

// Delete satisfies the ethdb.KeyValueWriter interface
// Delete removes the key from the key-value data store
func (d *Database) Delete(key []byte) error {
	table, err := ResolveTable(key)
	if err != nil {
		return err
	}
	var pgStr string
	switch table {
	case Undefined:
		return unsupportedTableTypeErr
	case KVStore:
		pgStr = deleteKVPgStr
	case Headers:
		pgStr = deleteHeaderPgStr
	case Hashes:
		pgStr = deleteHashPgStr
	case Bodies:
		pgStr = deleteBodyPgStr
	case Receipts:
		pgStr = deleteReceiptPgStr
	case TDs:
		pgStr = deleteTDPgStr
	case BloomBits:
		pgStr = deleteBloomBitsPgStr
	case TxLookUps:
		pgStr = deleteTxLookupPgStr
	case Preimages:
		pgStr = deletePreimagePgStr
	case Numbers:
		pgStr = deleteNumberPgStr
	case Configs:
		pgStr = deleteConfigPgStr
	case BloomIndexes:
		pgStr = deleteBloomIndexPgStr
	case TxMeta:
		pgStr = deleteTxMetaPgStr
	}
	_, err = d.db.Exec(pgStr, key)
	return err
}

// DatabaseProperty enum type
type DatabaseProperty int

const (
	Unknown DatabaseProperty = iota
	Size
	Idle
	InUse
	MaxIdleClosed
	MaxLifetimeClosed
	MaxOpenConnections
	OpenConnections
	WaitCount
	WaitDuration
)

// DatabasePropertyFromString helper function
func DatabasePropertyFromString(property string) (DatabaseProperty, error) {
	switch strings.ToLower(property) {
	case "size":
		return Size, nil
	case "idle":
		return Idle, nil
	case "inuse":
		return InUse, nil
	case "maxidleclosed":
		return MaxIdleClosed, nil
	case "maxlifetimeclosed":
		return MaxLifetimeClosed, nil
	case "maxopenconnections":
		return MaxOpenConnections, nil
	case "openconnections":
		return OpenConnections, nil
	case "waitcount":
		return WaitCount, nil
	case "waitduration":
		return WaitDuration, nil
	default:
		return Unknown, fmt.Errorf("unknown database property")
	}
}

// Stat satisfies the ethdb.Stater interface
// Stat returns a particular internal stat of the database
func (d *Database) Stat(property string) (string, error) {
	prop, err := DatabasePropertyFromString(property)
	if err != nil {
		return "", err
	}
	switch prop {
	case Size:
		var byteSize string
		return byteSize, d.db.Get(&byteSize, dbSizePgStr)
	case Idle:
		return string(d.db.Stats().Idle), nil
	case InUse:
		return string(d.db.Stats().InUse), nil
	case MaxIdleClosed:
		return string(d.db.Stats().MaxIdleClosed), nil
	case MaxLifetimeClosed:
		return string(d.db.Stats().MaxLifetimeClosed), nil
	case MaxOpenConnections:
		return string(d.db.Stats().MaxOpenConnections), nil
	case OpenConnections:
		return string(d.db.Stats().OpenConnections), nil
	case WaitCount:
		return string(d.db.Stats().WaitCount), nil
	case WaitDuration:
		return d.db.Stats().WaitDuration.String(), nil
	default:
		return "", fmt.Errorf("unhandled database property")
	}
}

// Compact satisfies the ethdb.Compacter interface
// Compact flattens the underlying data store for the given key range
func (d *Database) Compact(start []byte, limit []byte) error {
	return nil
}

// NewBatch satisfies the ethdb.Batcher interface
// NewBatch creates a write-only database that buffers changes to its host db
// until a final write is called
func (d *Database) NewBatch() ethdb.Batch {
	return NewBatch(d.db, nil)
}

// NewIterator creates a binary-alphabetical iterator over the entire keyspace
// contained within the key-value database.
func (d *Database) NewIterator() ethdb.Iterator {
	return NewIterator(nil, nil, d.db)
}

// NewIteratorWithStart creates a binary-alphabetical iterator over a subset of
// database content starting at a particular initial key (or after, if it does
// not exist).
func (d *Database) NewIteratorWithStart(start []byte) ethdb.Iterator {
	return NewIterator(start, nil, d.db)
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (d *Database) NewIteratorWithPrefix(prefix []byte) ethdb.Iterator {
	return NewIterator(nil, prefix, d.db)
}

// Close satisfies the io.Closer interface
// Close closes the db connection
func (d *Database) Close() error {
	return d.db.DB.Close()
}

// ExposeDB satisfies Exposer interface
func (d *Database) ExposeDB() interface{} {
	return d.db
}