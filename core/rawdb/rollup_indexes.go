package rawdb

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

// ReadHeadEth1HeaderHash reads the last processed Eth1 header hash
func ReadHeadEth1HeaderHash(db ethdb.KeyValueReader) common.Hash {
	data, _ := db.Get(headEth1HeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadEth1HeaderHash writes the last processed Eth1 header hash
func WriteHeadEth1HeaderHash(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Put(headEth1HeaderKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last eth1 header hash", "err", err)
	}
}

// ReadHeadEth1HeightKey reads the last processed Eth1 header height
func ReadHeadEth1HeaderHeight(db ethdb.KeyValueReader) uint64 {
	data, _ := db.Get(headEth1HeightKey)
	if len(data) == 0 {
		return 0
	}
	return new(big.Int).SetBytes(data).Uint64()
}

// WriteHeadEth1HeightKey writes the last processed Eth1 header height
func WriteHeadEth1HeaderHeight(db ethdb.KeyValueWriter, height uint64) {
	if err := db.Put(headEth1HeightKey, new(big.Int).SetUint64(height).Bytes()); err != nil {
		log.Crit("Failed to store eth1 header height", "err", err)
	}
}
