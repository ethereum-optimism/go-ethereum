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
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	_ "github.com/lib/pq" //postgres driver
)

type Table string

const (
	Undefined Table = "undefined"
	KVStore Table = "kvstore"
	Headers Table = "headers"
	Hashes Table = "hashes"
	Bodies Table = "bodies"
	Receipts Table = "receipts"
	TDs Table = "tds"
	BloomBits Table = "bloom_bits"
	TxLookUps Table = "tx_lookups"
	Preimages Table = "preimages"
	Numbers Table = "numbers"
	Configs Table = "configs"
	BloomIndexes Table = "bloom_indexes"
)

var (
	// prefixDelineation is used to delineate the key prefixes
	PrefixDelineation = []byte("-fix-")

	// numberDelineation is used to delineate the block number encoded in a key
	NumberDelineation = []byte("-nmb-")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
	headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix  = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits

	preimagePrefix = []byte("secure-key-")      // preimagePrefix + hash -> preimage
	configPrefix   = []byte("ethereum-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	bloomBitsIndexPrefix = []byte("iB") // bloomBitsIndexPrefix is the data table of a chain indexer to track its progress
)

// headerKey = headerPrefix + prefixDelineation + num (uint64 big endian) + numberDelineation + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(append(append(headerPrefix, PrefixDelineation...), encodeBlockNumber(number)...), NumberDelineation...), hash.Bytes()...)
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// decodeBlockNumber decodes a block number as big endian uint64
func decodeBlockNumber(enc []byte) uint64 {
	return binary.BigEndian.Uint64(enc)
}


// ResolvePutKey takes a key-value pair and returns:
// key prefix, table id, block number, header fk, error
func ResolvePutKey(key, val []byte) ([]byte, Table, uint64, []byte, error) {
	psk := bytes.Split(key, PrefixDelineation)
	l := len(psk)
	switch l {
	case 1:
		return nil, KVStore, 0, nil, nil
	case 2:
		bsk := bytes.Split(psk[1], NumberDelineation)
		if len(bsk) > 1 {
			num := decodeBlockNumber(bsk[0])
			switch prefix := psk[0]; {
			case bytes.Equal(prefix, headerPrefix):
				return psk[0], Headers, num, nil, nil
			case bytes.Equal(prefix, blockBodyPrefix):
				return psk[0], Bodies, num, headerKey(num, common.BytesToHash(bsk[1])), nil
			case bytes.Equal(prefix, blockReceiptsPrefix):
				return psk[0], Receipts, num, headerKey(num, common.BytesToHash(bsk[1])), nil
			}
		} else {
			switch prefix := psk[0]; {
			case bytes.Equal(prefix, headerNumberPrefix):
				num := decodeBlockNumber(val)
				return psk[0], Numbers, num, headerKey(num, common.BytesToHash(psk[1])), nil
			case bytes.Equal(prefix, txLookupPrefix):
				return psk[0], TxLookUps, 0, nil, nil
			case bytes.Equal(prefix, bloomBitsPrefix):
				return psk[0], BloomBits, 0, nil, nil
			case bytes.Equal(prefix, preimagePrefix):
				return psk[0], Preimages, 0, nil, nil
			case bytes.Equal(prefix, configPrefix):
				return psk[0], Configs, 0, nil, nil
			case bytes.Equal(prefix, bloomBitsIndexPrefix):
				return psk[0], BloomIndexes, 0, nil, nil
			}
		}
	case 3:
		bsk := bytes.Split(psk[1], NumberDelineation)
		num := decodeBlockNumber(bsk[0])
		switch suffix := psk[2]; {
		case bytes.Equal(suffix, headerTDSuffix):
			return psk[0], TDs, num, headerKey(num, common.BytesToHash(bsk[1])), nil
		case bytes.Equal(suffix, headerHashSuffix):
			return psk[0], Hashes, num, headerKey(num, common.BytesToHash(val)), nil
		}
	}
	return nil, Undefined, 0, nil, fmt.Errorf("unexpected number of key components: %d", l)
}

// ResolveTable returns the Table id from a given key
func ResolveTable(key []byte) (Table, error) {
	psk := bytes.Split(key, PrefixDelineation)
	l := len(psk)
	switch l {
	case 1:
		return KVStore, nil
	case 2:
		switch prefix := psk[0]; {
		case bytes.Equal(prefix, headerPrefix):
			return Headers, nil
		case bytes.Equal(prefix, blockBodyPrefix):
			return Bodies, nil
		case bytes.Equal(prefix, blockReceiptsPrefix):
			return Receipts, nil
		case bytes.Equal(prefix, headerNumberPrefix):
			return Numbers, nil
		case bytes.Equal(prefix, txLookupPrefix):
			return TxLookUps, nil
		case bytes.Equal(prefix, bloomBitsPrefix):
			return BloomBits, nil
		case bytes.Equal(prefix, preimagePrefix):
			return Preimages, nil
		case bytes.Equal(prefix, configPrefix):
			return Configs, nil
		case bytes.Equal(prefix, bloomBitsIndexPrefix):
			return BloomIndexes, nil
		}
	case 3:
		switch suffix := psk[2]; {
		case bytes.Equal(suffix, headerTDSuffix):
			return TDs, nil
		case bytes.Equal(suffix, headerHashSuffix):
			return Hashes, nil
		}
	}
	return Undefined, fmt.Errorf("unexpected number of key components: %d", l)
}