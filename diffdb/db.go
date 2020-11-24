package diffdb

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// maybe map by address?
type Diff struct {
	Address common.Address
	Keys    []string
}

type DiffDb struct {
	// Todo: should this be go-ethereum's leveldb maybe?
	// db *leveldb.DB
	inner map[uint64]map[common.Address][]common.Hash
}

// Called by the OVM StateManager
func (diff DiffDb) SetDiffKey(block *big.Int, address common.Address, key common.Hash) {
	// instantiate the diff
	if diff.inner[block.Uint64()] == nil {
		diff.inner[block.Uint64()] = make(map[common.Address][]common.Hash)
	}

	// set the value
	diff.inner[block.Uint64()][address] = append(diff.inner[block.Uint64()][address], key)
}

/// Gets a list of diffs from the databse for the corresponding
func (diff *DiffDb) GetDiff(block *big.Int) (map[common.Address][]common.Hash, error) {
	res, ok := diff.inner[block.Uint64()]
	if !ok {
		return nil, errors.New("No diff was found for the provided block")
	}
	return res, nil
}

func NewDiffDb(path string) (*DiffDb, error) {
	// db, err := leveldb.OpenFile(path, nil)
	// if err != nil {
	//     return nil, err
	// }
	// return &DiffDb{ db: db }, nil
	diffdb := make(map[uint64]map[common.Address][]common.Hash)
	return &DiffDb{inner: diffdb}, nil
}