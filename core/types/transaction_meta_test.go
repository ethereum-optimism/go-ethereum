package types

import (
	"bytes"
	"math/big"
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
		sighashType SignatureHashType
		queueOrigin *big.Int
	}{
		{
			txid:        &txid,
			msgSender:   &addr,
			sighashType: SighashEthSign,
			queueOrigin: big.NewInt(2),
		},
		{
			txid:        nil,
			msgSender:   &addr,
			sighashType: SighashEthSign,
			queueOrigin: big.NewInt(2),
		},
		{
			txid:        &txid,
			msgSender:   nil,
			sighashType: SighashEthSign,
			queueOrigin: big.NewInt(2),
		},
		{
			txid:        &txid,
			msgSender:   &addr,
			sighashType: SighashEthSign,
			queueOrigin: nil,
		},
		{
			txid:        nil,
			msgSender:   nil,
			sighashType: SighashEthSign,
			queueOrigin: nil,
		},
		{
			txid:        &txid,
			msgSender:   &addr,
			sighashType: SighashEthSign,
			queueOrigin: big.NewInt(0),
			queueOrigin: QueueOriginSequencer,
		},
	}

	txMetaSighashEncodeTests = []struct {
		input  SignatureHashType
		output SignatureHashType
	}{
		{
			input:  SighashEIP155,
			output: SighashEIP155,
		},
		{
			input:  SighashEthSign,
			output: SighashEthSign,
		},
	}
)

func TestTransactionMetaEncode(t *testing.T) {
	for _, test := range txMetaSerializationTests {
		txmeta := NewTransactionMeta(test.txid, test.msgSender, test.sighashType)
		txmeta.QueueOrigin = test.queueOrigin

		txmeta := NewTransactionMeta(test.txid, test.msgSender, test.queueOrigin, test.sighashType)
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
	queueOrigin := QueueOriginSequencer
	for _, test := range txMetaSighashEncodeTests {
		txmeta := NewTransactionMeta(&txid, &addr, queueOrigin, test.input)
		encoded := TxMetaEncode(txmeta)
		decoded, err := TxMetaDecode(encoded)

		if err != nil {
			t.Fatal(err)
		}

		if decoded.SignatureHashType != test.output {
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

	if meta1.SignatureHashType != meta2.SignatureHashType {
		return false
	}

	if meta1.QueueOrigin == nil || meta2.QueueOrigin == nil {
		// Note: this only works because it is the final comparison
		if meta1.QueueOrigin == nil && meta2.QueueOrigin == nil {
			return true
		}
	}

	if !bytes.Equal(meta1.QueueOrigin.Bytes(), meta2.QueueOrigin.Bytes()) {
		return false
	}

	return true
}
