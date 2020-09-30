// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package rawdb contains a collection of low level database accessors.
package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/metrics"
)

// The fields below define the low level database schema prefixing.
var (
	// prefixDelineation is used to delineate the key prefixes and suffixesg
	prefixDelineation = []byte("-fix-")

	//  numberDelineation is used to delineate the block number encoded in a key
	numberDelineation = []byte("-nmb-")

	// databaseVerisionKey tracks the current database version.
	databaseVerisionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest known header's hash.
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest known full block's hash.
	headBlockKey = []byte("LastBlock")

	// headFastBlockKey tracks the latest known incomplete block's hash during fast sync.
	headFastBlockKey = []byte("LastFast")

	// fastTrieProgressKey tracks the number of trie entries imported during fast sync.
	fastTrieProgressKey = []byte("TrieSync")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
	headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix  = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits

	// Optimism specific
	txMetaPrefix = []byte("x") // txMetaPrefix + hash -> transaction metadata

	// headEth1HeaderKey tracks the latest processed Eth1 Block
	headEth1HeaderKey = []byte("LastEth1Header")
	// headEth1HeightKey tracks the latest processed Eth1 Height
	headEth1HeightKey = []byte("LastEth1Height")

	PreimagePrefix = []byte("secure-key-")      // preimagePrefix + hash -> preimage

	configPrefix   = []byte("ethereum-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	preimageCounter    = metrics.NewRegisteredCounter("db/preimage/total", nil)
	preimageHitCounter = metrics.NewRegisteredCounter("db/preimage/hits", nil)
)

const (
	// freezerHeaderTable indicates the name of the freezer header table.
	freezerHeaderTable = "headers"

	// freezerHashTable indicates the name of the freezer canonical hash table.
	freezerHashTable = "hashes"

	// freezerBodiesTable indicates the name of the freezer block body table.
	freezerBodiesTable = "bodies"

	// freezerReceiptTable indicates the name of the freezer receipts table.
	freezerReceiptTable = "receipts"

	// freezerDifficultyTable indicates the name of the freezer total difficulty table.
	freezerDifficultyTable = "diffs"
)

// freezerNoSnappy configures whether compression is disabled for the ancient-tables.
// Hashes and difficulties don't compress well.
var freezerNoSnappy = map[string]bool{
	freezerHeaderTable:     false,
	freezerHashTable:       true,
	freezerBodiesTable:     false,
	freezerReceiptTable:    false,
	freezerDifficultyTable: true,
}

// LegacyTxLookupEntry is the legacy TxLookupEntry definition with some unnecessary
// fields.
type LegacyTxLookupEntry struct {
	BlockHash  common.Hash
	BlockIndex uint64
	Index      uint64
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKeyPrefix = headerPrefix + prefixDelineation + num (uint64 big endian) +  numberDelineation
func headerKeyPrefix(number uint64) []byte {
	return append(append(append(headerPrefix, prefixDelineation...), encodeBlockNumber(number)...),  numberDelineation...)
}

// headerKey = headerPrefix + prefixDelineation + num (uint64 big endian) +  numberDelineation + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(append(append(headerPrefix, prefixDelineation...), encodeBlockNumber(number)...),  numberDelineation...), hash.Bytes()...)
}

// headerTDKey = headerPrefix + prefixDelineation + num (uint64 big endian) + numberDelineation + hash + prefixDelineation + headerTDSuffix
func headerTDKey(number uint64, hash common.Hash) []byte {
	return append(append(headerKey(number, hash), prefixDelineation...), headerTDSuffix...)
}

// headerHashKey = headerPrefix + prefixDelineation + num (uint64 big endian) + numberDelineation + prefixDelineation + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(append(append(append(headerPrefix, prefixDelineation...), encodeBlockNumber(number)...), numberDelineation...), prefixDelineation...), headerHashSuffix...)
}

// headerNumberKey = headerNumberPrefix + prefixDelineation + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(append(headerNumberPrefix, prefixDelineation...), hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + prefixDelineation + num (uint64 big endian) +  numberDelineation + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(append(append(blockBodyPrefix, prefixDelineation...), encodeBlockNumber(number)...),  numberDelineation...), hash.Bytes()...)
}

// blockReceiptsKey = blockReceiptsPrefix + prefixDelineation + num (uint64 big endian) +  numberDelineation + hash
func blockReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(append(append(blockReceiptsPrefix, prefixDelineation...), encodeBlockNumber(number)...),  numberDelineation...), hash.Bytes()...)
}

// txLookupKey = txLookupPrefix + prefixDelineation + hash
func txLookupKey(hash common.Hash) []byte {
	return append(append(txLookupPrefix, prefixDelineation...), hash.Bytes()...)
}

// txMetaKey = txMetaPrefix + hash
func txMetaKey(hash common.Hash) []byte {
	return append(append(txMetaPrefix, prefixDelineation...), hash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + prefixDelineation + bit (uint16 big endian) + section (uint64 big endian) + hash
func bloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(append(bloomBitsPrefix, prefixDelineation...), make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[2:], uint16(bit))
	binary.BigEndian.PutUint64(key[4:], section)

	return key
}

// preimageKey = preimagePrefix + prefixDelineation + hash
func preimageKey(hash common.Hash) []byte {
	return append(append(PreimagePrefix, prefixDelineation...), hash.Bytes()...)
}

// configKey = configPrefix + prefixDelineation + hash
func configKey(hash common.Hash) []byte {
	return append(append(configPrefix, prefixDelineation...), hash.Bytes()...)
}
