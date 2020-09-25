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
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/jmoiron/sqlx"
)

const (
	hasPgStr    = "SELECT exists(SELECT 1 FROM eth.kvstore WHERE eth_key = $1)"
	getPgStr    = "SELECT eth_data FROM eth.kvstore WHERE eth_key = $1"
	putPgStr    = "INSERT INTO eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3) ON CONFLICT (eth_key) DO NOTHING"
	deletePgStr = "DELETE FROM eth.kvstore WHERE eth_key = $1"
	dbSizePgStr = "SELECT pg_total_relation_size('eth.kvstore')"
)

// Database is the type that satisfies the ethdb.Database and ethdb.KeyValueStore interfaces for PG-IPFS Ethereum data using a direct Postgres connection
type Database struct {
	db        *DB
	ancientTx *sqlx.Tx
}

// NewKeyValueStore returns a ethdb.KeyValueStore interface for PG-IPFS
func NewKeyValueStore(db *DB) ethdb.KeyValueStore {
	return &Database{
		db: db,
	}
}

// NewDatabase returns a ethdb.Database interface for PG-IPFS
func NewDatabase(db *DB) ethdb.Database {
	return &Database{
		db: db,
	}
}

// Has satisfies the ethdb.KeyValueReader interface
// Has retrieves if a key is present in the key-value data store
// Has uses the eth.key_preimages table
func (d *Database) Has(key []byte) (bool, error) {
	var exists bool
	return exists, d.db.Get(&exists, hasPgStr, key)
}

// Get satisfies the ethdb.KeyValueReader interface
// Get retrieves the given key if it's present in the key-value data store
// Get uses the eth.key_preimages table
func (d *Database) Get(key []byte) ([]byte, error) {
	var data []byte
	return data, d.db.Get(&data, getPgStr, key)
}

// Put satisfies the ethdb.KeyValueWriter interface
// Put inserts the given value into the key-value data store
// Key is expected to be the keccak256 hash of value
// Put inserts the keccak256 key into the eth.key_preimages table
func (d *Database) Put(key []byte, value []byte) error {
	dsKey, prefix, err := ResolveKeyPrefix(key)
	if err != nil {
		return err
	}
	if _, err = d.db.Exec(putPgStr, dsKey, value, prefix); err != nil {
		return err
	}
	return nil
}

// Delete satisfies the ethdb.KeyValueWriter interface
// Delete removes the key from the key-value data store
// Delete uses the eth.key_preimages table
func (d *Database) Delete(key []byte) error {
	_, err := d.db.Exec(deletePgStr, key)
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
