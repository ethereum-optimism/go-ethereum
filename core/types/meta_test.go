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
			sighashType: GetSighashEthSign(),
		},
		{
			txid:        nil,
			msgSender:   &addr,
			sighashType: GetSighashEthSign(),
		},
		{
			txid:        &txid,
			msgSender:   nil,
			sighashType: GetSighashEthSign(),
		},
	}

	txMetaSighashEncodeTests = []struct {
		input  *SignatureHashType
		output *SignatureHashType
	}{
		{
			input:  nil,
			output: GetSighashEIP155(),
		},
		{
			input:  GetSighashEIP155(),
			output: GetSighashEIP155(),
		},
		{
			input:  GetSighashEthSign(),
			output: GetSighashEthSign(),
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
