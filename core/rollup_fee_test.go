package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"math/rand"
	"testing"
)

var feeTests = map[string]struct {
	dataLen        int
	gasUsed        uint64
	gasLimit       uint64
	dataPrice      int64
	executionPrice int64
}{
	"simple":               {100, 10, 20, 20, 30},
	"zero gas used":        {1000, 0, 10, 20, 30},
	"zero data price":      {100, 0, 10, 0, 30},
	"zero execution price": {10000, 0, 10, 0, 0},
}

func TestCalculateRollupFee(t *testing.T) {
	for name, tt := range feeTests {
		t.Run(name, func(t *testing.T) {
			tx := L1Tx{
				From:     randomAddress(),
				To:       randomAddress(),
				Gas:      tt.gasLimit,
				GasPrice: big.NewInt(tt.executionPrice),
				Data:     randomBytes(tt.dataLen),
			}

			fee, err := CalculateRollupFee(tx, tt.gasUsed, big.NewInt(tt.dataPrice), big.NewInt(tt.executionPrice))
			if err != nil {
				t.Fatal(err)
			}

			data, err := rlp.EncodeToBytes(tx)
			if err != nil {
				t.Fatal(err)
			}
			dataFee := uint64((SIGNATURE_SIZE + len(data)) * int(tt.dataPrice))
			executionFee := uint64(tt.executionPrice) * tt.gasUsed
			expectedFee := dataFee + executionFee
			if fee.Cmp(big.NewInt(int64(expectedFee))) != 0 {
				t.Errorf("rollup fee check failed: expected %d, got %s", expectedFee, fee.String())
			}
		})
	}
}

type L1Tx struct {
	From     *common.Address
	To       *common.Address
	Gas      uint64
	GasPrice *big.Int
	Value    *big.Int
	Data     []byte
}

func randomAddress() *common.Address {
	randAddr := make([]byte, 20)
	rand.Read(randAddr)
	addr := common.BytesToAddress(randAddr)
	return &addr
}

func randomBytes(length int) []byte {
	randBytes := make([]byte, length)
	rand.Read(randBytes)
	return randBytes
}
