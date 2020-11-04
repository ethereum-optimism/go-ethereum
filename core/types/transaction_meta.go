/**
 * Optimism 2020 Copyright
 */

package types

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type QueueOrigin int64

const (
	// Possible `queue_origin` values
	QueueOriginSequencer QueueOrigin = 0
	QueueOriginL1ToL2    QueueOrigin = 1
)

//go:generate gencodec -type TransactionMeta -out gen_tx_meta_json.go

type TransactionMeta struct {
	L1BlockNumber     *big.Int          `json:"l1BlockNumber"`
	L1MessageSender   *common.Address   `json:"l1MessageSender" gencodec:"required"`
	SignatureHashType SignatureHashType `json:"signatureHashType" gencodec:"required"`
	QueueOrigin       *big.Int          `json:"queueOrigin" gencodec:"required"`
	Index             *uint64           `json:"index" gencodec:"required"`
}

// NewTransactionMeta creates a TransactionMeta
func NewTransactionMeta(l1BlockNumber *big.Int, l1MessageSender *common.Address, sighashType SignatureHashType, queueOrigin QueueOrigin) *TransactionMeta {
	return &TransactionMeta{
		L1BlockNumber:     l1BlockNumber,
		L1MessageSender:   l1MessageSender,
		SignatureHashType: sighashType,
		QueueOrigin:       big.NewInt(int64(queueOrigin)),
	}
}

// TxMetaDecode deserializes bytes as a TransactionMeta struct.
// The schema is:
//   varbytes(SignatureHashType) ||
//   varbytes(L1BlockNumber) ||
//   varbytes(L1MessageSender) ||
//   varbytes(QueueOrigin)
func TxMetaDecode(input []byte) (*TransactionMeta, error) {
	var err error
	meta := TransactionMeta{}
	b := bytes.NewReader(input)

	sb, err := common.ReadVarBytes(b, 0, 1024, "SignatureHashType")
	if err != nil {
		return nil, err
	}

	var sighashType SignatureHashType
	binary.Read(bytes.NewReader(sb), binary.LittleEndian, &sighashType)
	meta.SignatureHashType = sighashType

	lb, err := common.ReadVarBytes(b, 0, 1024, "l1BlockNumber")
	if err != nil {
		return nil, err
	}
	if !isNullValue(lb) {
		l1BlockNumber := new(big.Int).SetBytes(lb)
		meta.L1BlockNumber = l1BlockNumber
	}

	mb, err := common.ReadVarBytes(b, 0, 1024, "L1MessageSender")
	if err != nil {
		return nil, err
	}
	if !isNullValue(mb) {
		var l1MessageSender common.Address
		binary.Read(bytes.NewReader(mb), binary.LittleEndian, &l1MessageSender)
		meta.L1MessageSender = &l1MessageSender
	}

	qo, err := common.ReadVarBytes(b, 0, 1024, "QueueOrigin")
	if err != nil {
		return nil, err
	}
	if !isNullValue(qo) {
		queueOrigin := new(big.Int).SetBytes(qo)
		meta.QueueOrigin = queueOrigin
	}

	return &meta, nil
}

// TxMetaEncode serializes the TransactionMeta as bytes.
func TxMetaEncode(meta *TransactionMeta) []byte {
	b := new(bytes.Buffer)

	s := new(bytes.Buffer)
	binary.Write(s, binary.LittleEndian, &meta.SignatureHashType)
	common.WriteVarBytes(b, 0, s.Bytes())

	L1BlockNumber := meta.L1BlockNumber
	if L1BlockNumber == nil {
		common.WriteVarBytes(b, 0, getNullValue())
	} else {
		l := new(bytes.Buffer)
		binary.Write(l, binary.LittleEndian, L1BlockNumber.Bytes())
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

	queueOrigin := meta.QueueOrigin
	if queueOrigin == nil {
		common.WriteVarBytes(b, 0, getNullValue())
	} else {
		q := new(bytes.Buffer)
		binary.Write(q, binary.LittleEndian, queueOrigin.Bytes())
		common.WriteVarBytes(b, 0, q.Bytes())
	}

	return b.Bytes()
}

// This may collide with a uint8
func isNullValue(b []byte) bool {
	nullValue := []byte{0x00}
	return bytes.Equal(b, nullValue)
}

func getNullValue() []byte {
	return []byte{0x00}
}
