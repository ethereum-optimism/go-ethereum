package diffdb

import (
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
	db.SetDiffKey(big.NewInt(1), common.Address{0x1, 0x2}, common.Hash{0x12, 0x13}, false)
	db.SetDiffKey(big.NewInt(1), addr, hashes[0], false)
	db.SetDiffKey(big.NewInt(1), addr, hashes[1], false)
	db.SetDiffKey(big.NewInt(1), addr, hashes[2], false)
	db.SetDiffKey(big.NewInt(1), common.Address{0x2}, common.Hash{0x99}, false)
	db.SetDiffKey(big.NewInt(2), common.Address{0x2}, common.Hash{0x98}, true)

	diff, _ := db.GetDiff(big.NewInt(1))
	for i := range hashes {
		if hashes[i] != diff[addr][i].Key {
			t.Fatalf("Did not match")
		}
	}

	diff, _ = db.GetDiff(big.NewInt(2))
	if diff[common.Address{0x2}][0].Mutated != true {
		t.Fatalf("Did not match mutated")
	}
}
