package rawdb

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestReadWriteEth1HeaderHash(t *testing.T) {
	eth1HeaderHashes := []common.Hash{
		common.Hash{},
		common.HexToHash("0x00000000000000000004aba4215a0a494d25bc31f01014fa41bb74f082824146"),
		common.HexToHash("0xad8db6fc74e6f7e80af434588a2d98b03db03aeca0810635c0d3c6173cc027dd"),
	}

	db := NewMemoryDatabase()
	for _, hash := range eth1HeaderHashes {
		WriteHeadEth1HeaderHash(db, hash)
		got := ReadHeadEth1HeaderHash(db)
		if !bytes.Equal(hash.Bytes(), got.Bytes()) {
			t.Fatal("Header hash mismatch")
		}
	}
}

func TestReadWriteEth1HeaderHeight(t *testing.T) {
	eth1HeaderHeights := []uint64{
		1,
		1 << 2,
		1 << 8,
		1 << 16,
	}

	db := NewMemoryDatabase()
	for _, height := range eth1HeaderHeights {
		WriteHeadEth1HeaderHeight(db, height)
		got := ReadHeadEth1HeaderHeight(db)
		if height != got {
			t.Fatal("Header height mismatch")
		}
	}
}
