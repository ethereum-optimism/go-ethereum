package core

import (
	"math/big"
)

var ROLLUP_BASE_TX_SIZE int = 96

/// CalculateFee calculates the fee that must be paid to the Rollup sequencer, taking into
/// account the cost of publishing data to L1.
/// Returns: (ROLLUP_BASE_TX_SIZE + len(data)) * dataPrice + executionPrice * gasUsed
func CalculateRollupFee(data []byte, gasUsed uint64, dataPrice, executionPrice *big.Int) *big.Int {
	dataLen := int64(ROLLUP_BASE_TX_SIZE + len(data))
	// get the data fee
	dataFee := new(big.Int).Mul(dataPrice, big.NewInt(dataLen))
	executionFee := new(big.Int).Mul(executionPrice, big.NewInt(int64(gasUsed)))
	fee := new(big.Int).Add(dataFee, executionFee)
	return fee
}
