// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package postgres

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
)

var (
	testTx = types.Transaction{}
	testTxHash = testTx.Hash()
	testTxRLP, _ = rlp.EncodeToBytes(testTx)
	testNumber = big.NewInt(1337)
	testHeader = types.Header{Number: testNumber}
	testHeaderHash = testHeader.Hash()
	testHeaderRLP, _ = rlp.EncodeToBytes(testHeader)
	testHeaderKey = headerKey(testNumber.Uint64(), testHeaderHash)
	testTDKey = headerTDKey(testNumber.Uint64(), testHeaderHash)
	testBodyKey = blockBodyKey(testNumber.Uint64(), testHeaderHash)
	testReceiptsKey = blockReceiptsKey(testNumber.Uint64(), testHeaderHash)
	testHeaderHashkey = headerHashKey(testNumber.Uint64())
	testNumberKey = headerNumberKey(testHeaderHash)
	testTxLookupKey = txLookupKey(testTxHash)
	testTxMetaKey = txMetaKey(testTxHash)
	testBBKey = bloomBitsKey(1, 10, testTxHash)
	testConfigKey = configKey(testHeaderHash)
	testBBIKeyComponent = []byte{1,2,3,4,5}
	testBBIKey = append(append(bloomBitsIndexPrefix, prefixDelineation...), testBBIKeyComponent...)
	testStateLeafNodeRLP, _ = rlp.EncodeToBytes([]interface{}{
		[]byte{1,2,3,4,5},
		[]{1,2,3,4,5},
	})
	testStateLeafNodeHash = crypto.Keccak256Hash(testStateLeafNodeRLP)
	testPreimageKey = preimageKey(testStateLeafNodeHash)
)

func TestDatabaseGet(t *testing.T) {

}

func TestDatabasePut(t *testing.T) {

}

func TestDatabaseHas(t *testing.T) {

}

func TestDatabaseDelete(t *testing.T) {

}
