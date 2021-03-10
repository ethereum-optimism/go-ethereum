package gasprice

import (
	"context"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type L1Oracle struct {
}

func NewL1Oracle() *L1Oracle {
	return &L1Oracle{}
}

/// SuggestDataPrice returns the gas price which should be charged per byte of published
/// data by the sequencer.
func (gpo *L1Oracle) SuggestDataPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(100 * params.GWei), nil
}
