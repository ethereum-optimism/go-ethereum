/**
 * Optimism 2020 Copyright
 */

package types

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TransactionMeta struct {
	L1RollupTxId      *hexutil.Uint64    `json:"l1RollupTxId"`
	L1MessageSender   *common.Address    `json:"l1MessageSender"`
	SignatureHashType *SignatureHashType `json:"signatureHashType"`
}

func NewTransactionMeta(L1RollupTxId *hexutil.Uint64, L1MessageSender *common.Address, sighashType *SignatureHashType) *TransactionMeta {
	return &TransactionMeta{L1RollupTxId: L1RollupTxId, L1MessageSender: L1MessageSender, SignatureHashType: sighashType}
}

// TxMetaDecode deserializes bytes as a TransactionMeta struct.
// The schema is:
// varbytes(SignatureHashType) || varbytes(L1RollupTxId) || varbytes(L1MessageSender)
func TxMetaDecode(input []byte) (*TransactionMeta, error) {
	var err error
	meta := TransactionMeta{}

	b := bytes.NewReader(input)

	sb, err := common.ReadVarBytes(b, 0, 1024, "SignatureHashType")
	if err != nil {
		return &TransactionMeta{}, err
	}

	var sighashType SignatureHashType
	binary.Read(bytes.NewReader(sb), binary.LittleEndian, &sighashType)
	meta.SignatureHashType = &sighashType

	lb, err := common.ReadVarBytes(b, 0, 1024, "L1RollupTxId")
	if err != nil {
		return &TransactionMeta{}, err
	}

	if !isNullValue(lb) {
		var l1RollupTxId hexutil.Uint64
		binary.Read(bytes.NewReader(lb), binary.LittleEndian, l1RollupTxId)
		meta.L1RollupTxId = &l1RollupTxId
	}

	mb, err := common.ReadVarBytes(b, 0, 1024, "L1MessageSender")
	if err != nil {
		return &TransactionMeta{}, err
	}

	if !isNullValue(mb) {
		var l1MessageSender common.Address
		binary.Read(bytes.NewReader(mb), binary.LittleEndian, &l1MessageSender)
		meta.L1MessageSender = &l1MessageSender
	}

	return &meta, nil
}

// TxMetaEncode serializes the TransactionMeta as bytes.
func TxMetaEncode(meta *TransactionMeta) []byte {
	b := new(bytes.Buffer)

	// If the SignatureHashType is not explicitly defined, then it uses EIP155.
	sighashType := meta.SignatureHashType
	if sighashType == nil {
		sighashType = &SighashEIP155
	}

	s := new(bytes.Buffer)
	binary.Write(s, binary.LittleEndian, *sighashType)
	common.WriteVarBytes(b, 0, s.Bytes())

	L1RollupTxId := meta.L1RollupTxId
	if L1RollupTxId == nil {
		common.WriteVarBytes(b, 0, getNullValue())
	} else {
		l := new(bytes.Buffer)
		binary.Write(l, binary.LittleEndian, *L1RollupTxId)
		common.WriteVarBytes(b, 0, l.Bytes())
	}

	L1MessageSender := meta.L1MessageSender
	if L1MessageSender == nil {
		common.WriteVarBytes(b, 0, getNullValue())
	} else {
		l := new(bytes.Buffer)
		binary.Write(l, binary.LittleEndian, *L1MessageSender)
		common.WriteVarBytes(b, 0, l.Bytes())
	}

	return b.Bytes()
}

func isNullValue(b []byte) bool {
	nullValue := []byte{0x00}
	return bytes.Equal(b, nullValue)
}

func getNullValue() []byte {
	return []byte{0x00}
}
