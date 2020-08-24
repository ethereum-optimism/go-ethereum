package types

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	addr = common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	txid = hexutil.Uint64(0)

	txMetaSerializationTests = []struct {
		txid        *hexutil.Uint64
		msgSender   *common.Address
		sighashType *SignatureHashType
	}{
		{
			txid:        &txid,
			msgSender:   &addr,
			sighashType: &SighashEthSign,
		},
		{
			txid:        nil,
			msgSender:   &addr,
			sighashType: &SighashEthSign,
		},
		{
			txid:        &txid,
			msgSender:   nil,
			sighashType: &SighashEthSign,
		},
	}

	txMetaSighashEncodeTests = []struct {
		input  *SignatureHashType
		output *SignatureHashType
	}{
		{
			input:  nil,
			output: &SighashEIP155,
		},
		{
			input:  &SighashEIP155,
			output: &SighashEIP155,
		},
		{
			input:  &SighashEthSign,
			output: &SighashEthSign,
		},
	}

	addr1 = common.HexToAddress("abee26c3644d908dc9f3dff1e203da36b50573a3")
	addr2 = common.HexToAddress("e5c99b740572c2dabf7e3a418c4e9df2f793a599")
	addr3 = common.HexToAddress("cc4d9fbd42d6523cf804f2b9f26da98a9d9fa205")
	txid1 = hexutil.Uint64(33)
	txid2 = hexutil.Uint64(923)
	txid3 = hexutil.Uint64(5190)

	blockMetaEncodeTests = []struct {
		txs []*TransactionMeta
	}{
		{
			txs: []*TransactionMeta{
				NewTransactionMeta(&txid1, &addr1, &SighashEIP155),
			},
		},
		{
			txs: []*TransactionMeta{
				NewTransactionMeta(&txid1, &addr1, &SighashEIP155),
				NewTransactionMeta(&txid2, &addr2, &SighashEthSign),
				NewTransactionMeta(&txid3, &addr3, &SighashEthSign),
			},
		},
		{
			txs: []*TransactionMeta{
				NewTransactionMeta(nil, nil, &SighashEIP155),
				NewTransactionMeta(nil, &addr2, &SighashEthSign),
				NewTransactionMeta(&txid3, nil, &SighashEthSign),
			},
		},
	}
)

func TestTransactionMetaEncode(t *testing.T) {
	for _, test := range txMetaSerializationTests {
		txmeta := NewTransactionMeta(test.txid, test.msgSender, test.sighashType)
		encoded := TxMetaEncode(txmeta)
		decoded, err := TxMetaDecode(encoded)

		if err != nil {
			t.Fatal(err)
		}

		if !isTxMetaEqual(txmeta, decoded) {
			t.Fatal("Encoding/decoding mismatch")
		}
	}
}

func TestTransactionSighashEncode(t *testing.T) {
	for _, test := range txMetaSighashEncodeTests {
		txmeta := NewTransactionMeta(&txid, &addr, test.input)
		encoded := TxMetaEncode(txmeta)
		decoded, err := TxMetaDecode(encoded)

		if err != nil {
			t.Fatal(err)
		}

		if *decoded.SignatureHashType != *test.output {
			t.Fatal("SighashTypes do not match")
		}
	}
}

func isTxMetaEqual(meta1 *TransactionMeta, meta2 *TransactionMeta) bool {
	if meta1.L1MessageSender == nil || meta2.L1MessageSender == nil {
		if meta1.L1MessageSender != meta2.L1MessageSender {
			return false
		}
	} else {
		if !bytes.Equal(meta1.L1MessageSender.Bytes(), meta2.L1MessageSender.Bytes()) {
			return false
		}
	}

	if meta1.L1RollupTxId == nil || meta2.L1RollupTxId == nil {
		if meta1.L1RollupTxId != meta2.L1RollupTxId {
			return false
		}
	} else {
		if *meta1.L1RollupTxId != *meta2.L1RollupTxId {
			return false
		}
	}

	if meta1.SignatureHashType == nil || meta2.SignatureHashType == nil {
		if meta1.SignatureHashType != meta2.SignatureHashType {
			return false
		}
	} else {
		if *meta1.SignatureHashType != *meta2.SignatureHashType {
			return false
		}
	}

	return true
}
