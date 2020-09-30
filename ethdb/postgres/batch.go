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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/jmoiron/sqlx"
)

// Batch is the type that satisfies the ethdb.Batch interface for PG-IPFS Ethereum data using a direct Postgres connection
type Batch struct {
	db          *sqlx.DB
	tx          *sqlx.Tx
	valueSize   int
	replayCache map[string][]byte
}

// NewBatch returns a ethdb.Batch interface for PG-IPFS
func NewBatch(db *sqlx.DB, tx *sqlx.Tx) ethdb.Batch {
	b := &Batch{
		db:          db,
		tx:          tx,
		replayCache: make(map[string][]byte),
	}
	if tx == nil {
		b.Reset()
	}
	return b
}

// Put satisfies the ethdb.Batch interface
// Put inserts the given value into the key-value data store
// Key is expected to be the keccak256 hash of value
func (b *Batch) Put(key []byte, value []byte) (err error) {
	prefix, table, num, fk, err := ResolvePutKey(key, value)
	if err != nil {
		return err
	}
	var pgStr string
	args := make([]interface{}, 0, 3)
	args = append(args, key, value)
	switch table {
	case Undefined:
		return unsupportedTableTypeErr
	case KVStore:
		pgStr = putKVPgStr
		args = append(args, prefix)
	case Headers:
		pgStr = putHeaderPgStr
		args = append(args, num)
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
	if _, err = b.tx.Exec(pgStr, args...); err != nil {
		return err
	}
	b.valueSize += len(value)
	b.replayCache[common.Bytes2Hex(key)] = value
	return nil
}

// Delete satisfies the ethdb.Batch interface
// Delete removes the key from the key-value data store
func (b *Batch) Delete(key []byte) (err error) {
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
	_, err = b.tx.Exec(pgStr, key)
	if err != nil {
		return err
	}
	delete(b.replayCache, common.Bytes2Hex(key))
	return nil
}

// ValueSize satisfies the ethdb.Batch interface
// ValueSize retrieves the amount of data queued up for writing
// The returned value is the total byte length of all data queued to write
func (b *Batch) ValueSize() int {
	return b.valueSize
}

// Write satisfies the ethdb.Batch interface
// Write flushes any accumulated data to disk
// Reset should be called after every write
func (b *Batch) Write() error {
	if b.tx == nil {
		return nil
	}
	if err := b.tx.Commit(); err != nil {
		return err
	}
	b.replayCache = nil
	return nil
}

// Replay satisfies the ethdb.Batch interface
// Replay replays the batch contents
func (b *Batch) Replay(w ethdb.KeyValueWriter) error {
	if b.tx != nil {
		b.tx.Rollback()
		b.tx = nil
	}
	for key, value := range b.replayCache {
		if err := w.Put(common.Hex2Bytes(key), value); err != nil {
			return err
		}
	}
	b.replayCache = nil
	return nil
}

// Reset satisfies the ethdb.Batch interface
// Reset resets the batch for reuse
// This should be called after every write
func (b *Batch) Reset() {
	var err error
	b.tx, err = b.db.Beginx()
	if err != nil {
		panic(err)
	}
	b.replayCache = make(map[string][]byte)
	b.valueSize = 0
}
