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
	"github.com/ethereum/go-ethereum/ethdb"
)

const (
	initPgStr = `SELECT eth_key, eth_data
				FROM eth.kvstore 
				WHERE eth_key = $1`
	nextPgStr = `SELECT eth_key, eth_data
				FROM eth.kvstore
				WHERE eth_key > $1
				ORDER BY eth_key LIMIT 1`
	nextPgStrWithPrefix = `SELECT eth_key, eth_data
				FROM eth.kvstore
				WHERE eth_key > $1
				AND prefix = $2 
				ORDER BY eth_key LIMIT 1`
)

type nextModel struct {
	Key   []byte `db:"eth_key"`
	Value []byte `db:"eth_data"`
}

// Iterator is the type that satisfies the ethdb.Iterator interface for PG-IPFS Ethereum data using a direct Postgres connection
type Iterator struct {
	db                               *DB
	currentKey, prefix, currentValue []byte
	err                              error
	init                             bool
}

// NewIterator returns an ethdb.Iterator interface for PG-IPFS
func NewIterator(start, prefix []byte, db *DB) ethdb.Iterator {
	return &Iterator{
		db:         db,
		prefix:     prefix,
		currentKey: start,
		init:       start != nil,
	}
}

// Next satisfies the ethdb.Iterator interface
// Next moves the iterator to the next key/value pair
// It returns whether the iterator is exhausted
func (i *Iterator) Next() bool {
	next := new(nextModel)
	if i.init {
		i.init = false
		if err := i.db.Get(next, initPgStr, i.currentKey); err != nil {
			i.currentKey, i.currentValue, i.err = nil, nil, err
			return false
		}
	} else if i.prefix != nil {
		if err := i.db.Get(next, nextPgStrWithPrefix, i.currentKey, i.prefix); err != nil {
			i.currentKey, i.currentValue, i.err = nil, nil, err
			return false
		}
	} else {
		if err := i.db.Get(next, nextPgStr, i.currentKey); err != nil {
			i.currentKey, i.currentValue, i.err = nil, nil, err
			return false
		}
	}
	i.currentKey, i.currentValue, i.err = next.Key, next.Value, nil
	return true
}

// Error satisfies the ethdb.Iterator interface
// Error returns any accumulated error
// Exhausting all the key/value pairs is not considered to be an error
func (i *Iterator) Error() error {
	return i.err
}

// Key satisfies the ethdb.Iterator interface
// Key returns the key of the current key/value pair, or nil if done
// The caller should not modify the contents of the returned slice
// and its contents may change on the next call to Next
func (i *Iterator) Key() []byte {
	return i.currentKey
}

// Value satisfies the ethdb.Iterator interface
// Value returns the value of the current key/value pair, or nil if done
// The caller should not modify the contents of the returned slice
// and its contents may change on the next call to Next
func (i *Iterator) Value() []byte {
	return i.currentValue
}

// Release satisfies the ethdb.Iterator interface
// Release releases associated resources
// Release should always succeed and can be called multiple times without causing error
func (i *Iterator) Release() {
	i.db, i.currentKey, i.currentValue, i.err, i.prefix = nil, nil, nil, nil, nil
}
