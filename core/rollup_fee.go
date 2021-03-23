package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

/// Standard ECDSA signature length
var SIGNATURE_SIZE int = 65

/// CalculateRollupFee calculates the fee that must be paid to the Rollup sequencer, taking into
/// account the cost of publishing data to L1.
///
/// Returns: (SIGNATURE_SIZE + len(data)) * dataPrice + executionPrice * gasUsed
func CalculateRollupFee(tx interface{}, gasUsed uint64, dataPrice, executionPrice *big.Int) (*big.Int, error) {
	// RLP encode the transaction to get its serialized representation
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	// Get the length and add the signature length, since it is not included when
	// submitting a transaction for gas estimation
	dataLen := int64(SIGNATURE_SIZE + len(data))

	// dataFee = dataPrice * dataLen
	dataFee := new(big.Int).Mul(dataPrice, big.NewInt(dataLen))
	// executionFee = executionPrice * gasUsed
	executionFee := new(big.Int).Mul(executionPrice, new(big.Int).SetUint64(gasUsed))
	// total fee = dataFee + executionFee
	fee := new(big.Int).Add(dataFee, executionFee)

	return fee, nil
}
