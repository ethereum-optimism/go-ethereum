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

// probably over leveldb?
type DiffDb struct {
}

func (diff *DiffDb) GetDiff(*big.Int) ([]Diff, error) {
	return nil, errors.New("lol'")
}

func NewDiffDb(path string) *DiffDb {
	return &DiffDb{}
}
