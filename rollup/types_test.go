package rollup

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestCanonicalChainBatchContext(t *testing.T) {
	tests := []struct {
		input  ctcBatchContext
		expect []byte
	}{
		{
			input: ctcBatchContext{
				NumSequencedTransactions:       big.NewInt(0),
				NumSubsequentQueueTransactions: big.NewInt(0),
				Timestamp:                      big.NewInt(0),
				BlockNumber:                    big.NewInt(0),
			},
			expect: hexutil.MustDecode("0x00"),
		},
		{
			input: ctcBatchContext{
				NumSequencedTransactions:       big.NewInt(10),
				NumSubsequentQueueTransactions: big.NewInt(11),
				Timestamp:                      big.NewInt(12),
				BlockNumber:                    big.NewInt(13),
			},
			expect: hexutil.MustDecode("0x00"),
		},
	}

	for _, test := range tests {
		buf := new(bytes.Buffer)
		err := test.input.Encode(buf)
		if err != nil {
			t.Fatal(err)
		}
		reader := bytes.NewReader(buf.Bytes())
		ctcCtx := ctcBatchContext{}
		err = ctcCtx.Decode(reader)
		if err != nil {
			t.Fatal(err)
		}

		if !isCtcBatchContextEqual(&test.input, &ctcCtx) {
			t.Fatal(err)
		}
	}
}

func isCtcBatchContextEqual(one *ctcBatchContext, two *ctcBatchContext) bool {
	if one == nil && two == nil {
		return true
	}
	if !bytes.Equal(one.NumSequencedTransactions.Bytes(), two.NumSequencedTransactions.Bytes()) {
		return false
	}
	if !bytes.Equal(one.NumSubsequentQueueTransactions.Bytes(), two.NumSubsequentQueueTransactions.Bytes()) {
		return false
	}
	if !bytes.Equal(one.Timestamp.Bytes(), two.Timestamp.Bytes()) {
		return false
	}
	if !bytes.Equal(one.BlockNumber.Bytes(), two.BlockNumber.Bytes()) {
		return false
	}
	return true
}

func TestSequencerBatchCalldata(t *testing.T) {
	tests := []struct {
		input  appendSequencerBatchCallData
		expect []byte
	}{
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(0),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x1234"),
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(0),
						Timestamp:                      big.NewInt(0),
						BlockNumber:                    big.NewInt(0),
					},
				},
				ShouldStartAtBatch:    big.NewInt(1234),
				TotalElementsToAppend: big.NewInt(1),
			},
			expect: hexutil.MustDecode("0x00000004d2000001000001000001000000000000000000000000000000021234"),
		},
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602820447),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x12"),
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602820447),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x1234"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(2),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602820447),
						BlockNumber:                    big.NewInt(0),
					},
				},
				ShouldStartAtBatch:    big.NewInt(0),
				TotalElementsToAppend: big.NewInt(3),
			},
			expect: hexutil.MustDecode("0x0000000000000003000001000002000001005f89195f0000000000000001120000021234"),
		},
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x12"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x1234"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x123434"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x12343434"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{NumSequencedTransactions: big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
				},
				ShouldStartAtBatch:    big.NewInt(0),
				TotalElementsToAppend: big.NewInt(8),
			},
			expect: hexutil.MustDecode("0x0000000000000008000004000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c00000112000002123400000312343400000412343434"),
		},
	}

	for _, test := range tests {
		buf := new(bytes.Buffer)
		err := test.input.Encode(buf)
		if err != nil {
			t.Fatalf("Cannot encode appendSequencerBatchCallData: %s\n", err.Error())
		}

		if !bytes.Equal(test.expect, buf.Bytes()) {
			t.Fatalf("Serialization mismatch: expect\nexpect:\n%x\ngot:\n%x", test.expect, buf.Bytes())
		}

		reader := bytes.NewReader(buf.Bytes())
		cd := appendSequencerBatchCallData{}
		err = cd.Decode(reader)
		if err != nil {
			t.Fatalf("Error decoding: %s", err.Error())
		}

		if !isSequencerBatchCalldataEqual(&test.input, &cd) {
			t.Fatalf("Deserialization result mismatch:\nexpect:\n%#v\ngot:\n%#v", test.input, cd)
		}
	}
}

func TestCTCTransactionDeserialization(t *testing.T) {
	// Use a test vector generated by the javascript
	raw := hexutil.MustDecode("0x0011111111111111111111111111111111111111111111111111111111111111112222222222222222222222222222222222222222222222222222222222222222010001f4000064000064121212121212121212121212121212121212121299999999999999999999")
	// Expect these values
	sig := hexutil.MustDecode("0x1111111111111111111111111111111111111111111111111111111111111111222222222222222222222222222222222222222222222222222222222222222201")
	addr := common.HexToAddress("0x1212121212121212121212121212121212121212")
	data := hexutil.MustDecode("0x99999999999999999999")

	tx := CTCTransaction{}
	err := tx.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if tx.typ != CTCTransactionTypeEIP155 {
		t.Fatal("Wrong type decoded")
	}

	eip155, ok := tx.tx.(*CTCTxEIP155)
	if !ok {
		t.Fatal("Wrong type decoded")
	}
	if !bytes.Equal(eip155.Signature[:], sig) {
		t.Fatal("Wrong Signature decoded")
	}
	if eip155.gasLimit != 500 {
		t.Fatal("Wrong gas limit decoded")
	}
	if eip155.gasPrice != 100 {
		t.Fatal("Wrong gas price decoded")
	}
	if !bytes.Equal(eip155.target.Bytes(), addr.Bytes()) {
		t.Fatal("Wrong target decoded")
	}
	if !bytes.Equal(eip155.data, data) {
		t.Fatal("Wrong data decoded")
	}
	length, err := tx.Len()
	if err != nil {
		t.Fatal(err)
	}

	// Reserialize the struct and make sure that it matches the
	// original raw data
	encoded := make([]byte, length)
	err = tx.Encode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(raw, encoded) {
		t.Fatal("Serialization/Deserialization mismatch")
	}
}

func TestCTCTransactionSerialization(t *testing.T) {
	raw := hexutil.MustDecode("0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301")
	var sig [65]byte
	copy(sig[:], raw)

	txs := []CTCTransaction{
		{
			typ: CTCTransactionTypeEIP155,
			tx: &CTCTxEIP155{
				Signature: sig,
				gasLimit:  620,
				gasPrice:  (1 << 16) + 8,
				nonce:     (1 << 18) + 32,
				target:    common.Address{},
				data:      []byte("abcdef"),
			},
		},
		{
			typ: CTCTransactionTypeEIP155,
			tx: &CTCTxEIP155{
				Signature: sig,
				gasLimit:  89,
				gasPrice:  45,
				nonce:     (1 << 21) + 14,
				target:    common.HexToAddress("0x5769785087b1b64e4cbd9a38d48a1ca35a2fd75cf5cd941d75b2e2fbc6018e8a"),
				data:      raw,
			},
		},
		{
			typ: CTCTransactionTypeEIP155,
			tx: &CTCTxEIP155{
				Signature: sig,
				gasLimit:  (1 << 20) + 45,
				gasPrice:  20,
				nonce:     (1 << 12) + 99,
				target:    common.HexToAddress("0x5769785087b1b64e4cbd9a38d48a1ca35a2fd75cf5cd941d75b2e2fbc6018e8a"),
				data:      []byte("foobarbazlolololololololol"),
			},
		},
		{
			typ: CTCTransactionTypeEOA,
			tx: &CTCTxCreateEOA{
				Signature: sig,
				Hash:      common.HexToHash("0xffcfa4cf82b5326f382ece74ba547c368677f922b8f652f5b370a38eccf5f8e1"),
			},
		},
	}

	for _, tx := range txs {
		length, err := tx.Len()
		if err != nil {
			t.Fatal("Cannot read legnth")
		}
		slice := make([]byte, length)
		err = tx.Encode(slice)
		if err != nil {
			t.Fatalf("Cannot encode slice: %s", err)
		}

		decoded := CTCTransaction{}
		err = decoded.Decode(slice)
		if err != nil {
			t.Fatalf("Cannot decode ctc tx: %s", err)
		}
		if tx.typ != decoded.typ {
			t.Fatal("Invalid types")
		}
		rt := reflect.TypeOf(tx.tx)
		if tx.typ == CTCTransactionTypeEIP155 {
			if rt != reflect.TypeOf(&CTCTxEIP155{}) {
				t.Fatal("Invalid type")
			}
			typ, ok := tx.tx.(*CTCTxEIP155)
			if !ok {
				t.Fatal("Cannot type cast")
			}
			got, ok := decoded.tx.(*CTCTxEIP155)
			if !ok {
				t.Fatal("Cannot type cast")
			}

			if !bytes.Equal(typ.Signature[:], got.Signature[:]) {
				t.Fatalf("Signature Serialization mismatch\ngot:\n%x\nexpected:\n%x\n", got.Signature, typ.Signature)
			}
			if typ.gasLimit != got.gasLimit {
				t.Fatalf("Gas limit mismatch\ngot:\n%d\nexpected:\n%d\n", got.gasLimit, typ.gasLimit)
			}
			if typ.gasPrice != got.gasPrice {
				t.Fatalf("Gas price mismatch\ngot:\n%d\nexpected:\n%d\n", got.gasPrice, typ.gasPrice)
			}
			if typ.nonce != got.nonce {
				t.Fatalf("Nonce mismatch\ngot:\n%d\nexpected:\n%d\n", got.nonce, typ.nonce)
			}
			if !bytes.Equal(typ.target.Bytes(), got.target.Bytes()) {
				t.Fatalf("Target mismatch\ngot:\n%x\nexpected:\n%x\n", got.target, typ.target)
			}
			if !bytes.Equal(typ.data, got.data) {
				t.Fatalf("Data mismatch\ngot:\n%x\nexpected:\n%x\n", got.data, typ.data)
			}
		} else if tx.typ == CTCTransactionTypeEOA {
			if rt != reflect.TypeOf(&CTCTxCreateEOA{}) {
				t.Fatal("Invalid type")
			}
			typ, ok := tx.tx.(*CTCTxCreateEOA)
			if !ok {
				t.Fatal("Cannot type cast")
			}
			got, ok := decoded.tx.(*CTCTxCreateEOA)
			if !ok {
				t.Fatal("Cannot type cast")
			}

			if !bytes.Equal(typ.Signature[:], got.Signature[:]) {
				t.Fatalf("Signature Serialization mismatch\ngot:\n%x\nexpected:\n%x\n", got.Signature, typ.Signature)
			}
			if !bytes.Equal(typ.Hash[:], got.Hash[:]) {
				t.Fatalf("Hash mismatch\ngot:\n%x\nexpected:\n%x\n", got.Hash, typ.Hash)
			}
		}
	}
}

func isSequencerBatchCalldataEqual(one, two *appendSequencerBatchCallData) bool {
	if one == nil && two == nil {
		return true
	}
	if xornil(one.ShouldStartAtBatch, two.ShouldStartAtBatch) {
		return false
	}
	if one.ShouldStartAtBatch != nil && two.ShouldStartAtBatch != nil {
		if !bytes.Equal(one.ShouldStartAtBatch.Bytes(), two.ShouldStartAtBatch.Bytes()) {
			return false
		}
	}
	if xornil(one.TotalElementsToAppend, two.TotalElementsToAppend) {
		return false
	}
	if one.TotalElementsToAppend != nil && two.TotalElementsToAppend != nil {
		if !bytes.Equal(one.TotalElementsToAppend.Bytes(), two.TotalElementsToAppend.Bytes()) {
			return false
		}
	}
	if len(one.Contexts) != len(two.Contexts) {
		return false
	}
	for i, oneCtx := range one.Contexts {
		twoCtx := two.Contexts[i]
		if !isCtcBatchContextEqual(&oneCtx, &twoCtx) {
			return false
		}
	}
	if len(one.ChainElements) != len(two.ChainElements) {
		return false
	}
	for i, oneEl := range one.ChainElements {
		twoEl := two.ChainElements[i]
		if oneEl.IsSequenced != twoEl.IsSequenced {
			return false
		}
		if xornil(oneEl.Timestamp, twoEl.Timestamp) {
			return false
		}
		if oneEl.Timestamp != nil && twoEl.Timestamp != nil {
			if !bytes.Equal(oneEl.Timestamp.Bytes(), twoEl.Timestamp.Bytes()) {
				return false
			}
		}
		if xornil(oneEl.BlockNumber, twoEl.BlockNumber) {
			return false
		}
		if oneEl.BlockNumber != nil && twoEl.BlockNumber != nil {
			if !bytes.Equal(oneEl.BlockNumber.Bytes(), twoEl.BlockNumber.Bytes()) {
				return false
			}
		}
		if !bytes.Equal(oneEl.TxData, twoEl.TxData) {
			return false
		}
	}
	return true
}

func xornil(one, two interface{}) bool {
	if one == nil && two != nil {
		return true
	}
	if two == nil && one != nil {
		return true
	}

	return false
}
