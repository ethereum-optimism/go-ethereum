package diffdb

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func TestInMemoryDb(t *testing.T) {
	db, err := NewDiffDb("whatever")
	if err != nil {
		t.Fatal(err)
	}

	hashes := []common.Hash{
		common.Hash{0x1},
		common.Hash{0x2},
		common.Hash{0x0},
	}
	addr := common.Address{0x1}
	db.SetDiffKey(big.NewInt(1), common.Address{0x1, 0x2}, common.Hash{0x12, 0x13})
	db.SetDiffKey(big.NewInt(1), addr, hashes[0])
	db.SetDiffKey(big.NewInt(1), addr, hashes[1])
	db.SetDiffKey(big.NewInt(1), addr, hashes[2])
	db.SetDiffKey(big.NewInt(1), common.Address{0x2}, common.Hash{0x99})
	db.SetDiffKey(big.NewInt(2), common.Address{0x2}, common.Hash{0x98})

	diff, _ := db.GetDiff(big.NewInt(1))
	for i := range hashes {
		if hashes[i] != diff[addr][i] {
			t.Fatalf("Did not match")
		}
	}

	diff, _ = db.GetDiff(big.NewInt(5))
	fmt.Println(diff)
}
