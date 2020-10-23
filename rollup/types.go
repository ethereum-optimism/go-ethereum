package rollup

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

const (
	MinTxBytes               = uint64(100)
	MinTxGas                 = MinTxBytes*params.TxDataNonZeroGasEIP2028 + params.SstoreSetGas
	TransitionBatchGasBuffer = uint64(1_000_000)
)

type BlockStore interface {
	GetBlockByNumber(number uint64) *types.Block
}

type Transition struct {
	transaction *types.Transaction
	postState   common.Hash
}

func newTransition(tx *types.Transaction, postState common.Hash) *Transition {
	return &Transition{
		transaction: tx,
		postState:   postState,
	}
}

type TransitionBatch struct {
	transitions []*Transition
}

func NewTransitionBatch(defaultSize int) *TransitionBatch {
	return &TransitionBatch{transitions: make([]*Transition, 0, defaultSize)}
}

// addBlock adds a Geth Block to the TransitionBatch. This is just its transaction and state root.
func (r *TransitionBatch) addBlock(block *types.Block) {
	r.transitions = append(r.transitions, newTransition(block.Transactions()[0], block.Root()))
}

// Implement the Sort interface for []types.Log by Index
type LogsByIndex []types.Log

func (l LogsByIndex) Len() int           { return len(l) }
func (l LogsByIndex) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l LogsByIndex) Less(i, j int) bool { return l[i].Index < l[i].Index }

type ctcBatchContext struct {
	NumSequencedTransactions       *big.Int
	NumSubsequentQueueTransactions *big.Int
	Timestamp                      *big.Int
	BlockNumber                    *big.Int
}

type chainElement struct {
	IsSequenced bool
	Timestamp   *big.Int // Origin Sequencer
	BlockNumber *big.Int // Origin Sequencer
	TxData      []byte   // Origin Sequencer
}

type appendSequencerBatchCallData struct {
	ShouldStartAtBatch    *big.Int
	TotalElementsToAppend *big.Int
	Contexts              []ctcBatchContext
	ChainElements         []chainElement
}

func (c *ctcBatchContext) Encode(w io.Writer) error {
	elements := [][]byte{
		common.LeftPadBytes(c.NumSequencedTransactions.Bytes(), 3),
		common.LeftPadBytes(c.NumSubsequentQueueTransactions.Bytes(), 3),
		common.LeftPadBytes(c.Timestamp.Bytes(), 5),
		common.LeftPadBytes(c.BlockNumber.Bytes(), 5),
	}
	for _, element := range elements {
		_, err := w.Write(element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ctcBatchContext) Len() int {
	return 3 + 3 + 5 + 5
}

func (c *ctcBatchContext) Decode(r io.ReaderAt) error {
	offset := int64(0)
	elements := [][]byte{
		make([]byte, 3),
		make([]byte, 3),
		make([]byte, 5),
		make([]byte, 5),
	}

	for i, element := range elements {
		sr := io.NewSectionReader(r, offset, int64(len(element)))
		off, err := sr.Read(element)
		if err != nil {
			return err
		}

		switch i {
		case 0:
			c.NumSequencedTransactions = new(big.Int).SetBytes(element)
		case 1:
			c.NumSubsequentQueueTransactions = new(big.Int).SetBytes(element)
		case 2:
			c.Timestamp = new(big.Int).SetBytes(element)
		case 3:
			c.BlockNumber = new(big.Int).SetBytes(element)
		}
		offset += int64(off)
	}

	return nil
}

func (c *appendSequencerBatchCallData) Encode(w io.Writer) error {
	contexts := new(bytes.Buffer)
	for _, context := range c.Contexts {
		buf := new(bytes.Buffer)
		err := context.Encode(buf)
		if err != nil {
			return err
		}
		contexts.Write(buf.Bytes())
	}
	transactions := new(bytes.Buffer)
	for _, el := range c.ChainElements {
		if !el.IsSequenced {
			continue
		}
		header := make([]byte, 0, 4)
		buf := bytes.NewBuffer(header)
		err := binary.Write(buf, binary.BigEndian, uint32(len(el.TxData)))
		if err != nil {
			return err
		}
		_ = buf.Next(1) // Move forward a byte
		transactions.Write(buf.Bytes())
		transactions.Write(el.TxData)
	}
	elements := [][]byte{
		common.LeftPadBytes(c.ShouldStartAtBatch.Bytes(), 5),
		common.LeftPadBytes(c.TotalElementsToAppend.Bytes(), 3),
		common.LeftPadBytes(new(big.Int).SetUint64(uint64(len(c.Contexts))).Bytes(), 3),
		contexts.Bytes(),
		transactions.Bytes(),
	}
	for _, element := range elements {
		_, err := w.Write(element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *appendSequencerBatchCallData) Decode(r io.ReaderAt) error {
	offset := int64(0)
	elements := [][]byte{
		make([]byte, 5),
		make([]byte, 3),
		make([]byte, 3), // context count
	}
	ctxCount := new(big.Int)
	for i, element := range elements {
		sr := io.NewSectionReader(r, offset, int64(len(element)))
		off, err := sr.Read(element)
		if err != nil {
			return err
		}

		switch i {
		case 0:
			a.ShouldStartAtBatch = new(big.Int).SetBytes(element)
		case 1:
			a.TotalElementsToAppend = new(big.Int).SetBytes(element)
		case 2:
			ctxCount.SetBytes(element)
		}
		offset += int64(off)
	}

	a.Contexts = make([]ctcBatchContext, ctxCount.Uint64())
	for i := uint64(0); i < ctxCount.Uint64(); i++ {
		batchCtx := ctcBatchContext{}
		sr := io.NewSectionReader(r, offset, offset+int64(batchCtx.Len()))
		err := batchCtx.Decode(sr)
		if err != nil {
			return fmt.Errorf("Cannot decode batch context: %w", err)
		}
		a.Contexts[i] = batchCtx
		offset += int64(batchCtx.Len())
	}

	txCount := uint64(0)
	for _, ctx := range a.Contexts {
		txCount += ctx.NumSequencedTransactions.Uint64()
		txCount += ctx.NumSubsequentQueueTransactions.Uint64()
	}

	if txCount != a.TotalElementsToAppend.Uint64() {
		return errors.New("Incorrect number of elements")
	}

	a.ChainElements = []chainElement{}
	for _, ctx := range a.Contexts {
		timestamp := ctx.Timestamp
		blockNumber := ctx.BlockNumber
		for i := uint64(0); i < ctx.NumSequencedTransactions.Uint64(); i++ {
			header := make([]byte, 3)
			sr := io.NewSectionReader(r, offset, offset+3)
			off, err := sr.Read(header)

			if err != nil {
				return fmt.Errorf("Cannot read tx header: %w", err)
			}
			offset += int64(off)

			var sizeHi uint16
			var sizeLo uint8
			hr := bytes.NewReader(header)
			err = binary.Read(hr, binary.BigEndian, &sizeHi)
			if err != nil {
				return fmt.Errorf("Cannot read tx header hi bits: %w", err)
			}
			err = binary.Read(hr, binary.BigEndian, &sizeLo)
			if err != nil {
				return fmt.Errorf("Cannot read tx header lo bits: %w", err)
			}

			size := (sizeHi << 8) | uint16(sizeLo)
			tx := make([]byte, size)
			tsr := io.NewSectionReader(r, offset, offset+int64(size))
			off, err = tsr.Read(tx)
			if err != nil {
				return fmt.Errorf("Cannot read tx: %w", err)
			}
			offset += int64(off)

			element := chainElement{
				IsSequenced: true,
				Timestamp:   timestamp,
				BlockNumber: blockNumber,
				TxData:      tx,
			}
			a.ChainElements = append(a.ChainElements, element)
		}

		for i := uint64(0); i < ctx.NumSubsequentQueueTransactions.Uint64(); i++ {
			element := chainElement{
				IsSequenced: false,
				Timestamp:   nil,
				BlockNumber: nil,
				TxData:      []byte{},
			}
			a.ChainElements = append(a.ChainElements, element)
		}
	}
	return nil
}

// Canonical Chain Transaction Serialization
type CTCTransactionType uint8

const (
	CTCTransactionTypeEOA    CTCTransactionType = 0
	CTCTransactionTypeEIP155 CTCTransactionType = 1
)

type CTCTransaction struct {
	typ CTCTransactionType
	tx  Ser
}

func (c *CTCTransaction) Len() (int, error) {
	if c.tx == nil {
		return int(^uint(0) >> 1), errors.New("Cannot compute length")
	}
	length, err := c.tx.Len()
	return 1 + length, err
}

func (c *CTCTransaction) Encode(b []byte) error {
	length, _ := c.Len()
	if len(b) < length {
		return errors.New("Encoding overflow")
	}
	b[0] = uint8(c.typ)
	err := c.tx.Encode(b[1:])
	if err != nil {
		return fmt.Errorf("Cannot encode ctc tx: %w", err)
	}

	return nil
}

func (c *CTCTransaction) Decode(b []byte) error {
	// only care about the first byte, the other decode methods
	// will handle length checks
	if len(b) < 1 {
		return errors.New("CTCTransaction Decoding overflow")
	}
	c.typ = CTCTransactionType(b[0])
	switch c.typ {
	case CTCTransactionTypeEOA:
		tx := CTCTxCreateEOA{}
		err := tx.Decode(b[1:])
		if err != nil {
			return fmt.Errorf("Cannot decode EOA ctc tx %x: %w", b, err)
		}
		c.tx = &tx
	case CTCTransactionTypeEIP155:
		tx := CTCTxEIP155{}
		err := tx.Decode(b[1:])
		if err != nil {
			return fmt.Errorf("Cannot decode EIP155 ctc tx %x: %w", b, err)
		}
		c.tx = &tx
	}
	return nil
}

type Ser interface {
	Encode([]byte) error
	Decode([]byte) error
	Len() (int, error)
}

type CTCTxCreateEOA struct {
	Signature [65]byte
	Hash      [32]byte
}

func (c *CTCTxCreateEOA) Len() (int, error) { return 65 + 32, nil }
func (c *CTCTxCreateEOA) Encode(b []byte) error {
	length, _ := c.Len()
	if len(b) < length {
		return errors.New("CTCTxCreateEOA encoding overflow")
	}
	copy(b[:65], c.Signature[:])
	copy(b[65:65+32], c.Hash[:])
	return nil
}
func (c *CTCTxCreateEOA) Decode(b []byte) error {
	length, _ := c.Len()
	if len(b) < length {
		return errors.New("CTCTxCreateEOA decoding overflow")
	}
	copy(c.Signature[:], b[:65])
	copy(c.Hash[:], b[65:])
	return nil
}

type CTCTxEIP155 struct {
	Signature [65]byte
	gasLimit  uint16 // uint16
	gasPrice  uint8  // uint8
	nonce     uint32 // uint24
	target    common.Address
	data      []byte
}

func (c *CTCTxEIP155) Len() (int, error) {
	return 65 + 2 + 1 + 3 + 20 + len(c.data), nil
}

func (c *CTCTxEIP155) Encode(b []byte) error {
	length, _ := c.Len()
	if len(b) < length {
		return errors.New("CTCTxEIP155 encoding overflow")
	}
	copy(b[:65], c.Signature[:])
	binary.BigEndian.PutUint16(b[65:65+2], c.gasLimit)
	b[67] = c.gasPrice
	hi := c.nonce & 0x00FF0000 >> 16
	lo := c.nonce & 0x0000FFFF
	b[68] = uint8(hi)
	binary.BigEndian.PutUint16(b[69:69+2], uint16(lo))
	copy(b[71:71+20], c.target.Bytes())
	copy(b[91:], c.data)
	return nil
}

func (c *CTCTxEIP155) Decode(b []byte) error {
	length, _ := c.Len()
	if len(b) < length {
		return errors.New("CTCTxEIP155 decoding overflow")
	}
	copy(c.Signature[:], b[:65])
	c.gasLimit = binary.BigEndian.Uint16(b[65 : 65+2])
	c.gasPrice = b[67]
	hi := uint32(b[68])
	lo := uint32(binary.BigEndian.Uint16(b[69 : 69+2]))
	c.nonce = (hi << 16) | lo
	c.target = common.BytesToAddress(b[71 : 71+20])
	c.data = b[91:]
	return nil
}

func isCtcTxEqual(a, b *types.Transaction) bool {
	if a.To() == nil && b.To() != nil {
		return false
	}
	if b.To() != nil && a.To() == nil {
		return false
	}
	if !bytes.Equal(a.To().Bytes(), b.To().Bytes()) {
		return false
	}
	if !bytes.Equal(a.Data(), b.Data()) {
		return false
	}
	if a.L1MessageSender() == nil && b.L1MessageSender() != nil {
		return false
	}
	if a.L1MessageSender() != nil && b.L1MessageSender() == nil {
		return false
	}
	if !bytes.Equal(a.L1MessageSender().Bytes(), b.L1MessageSender().Bytes()) {
		return false
	}
	if a.Gas() != b.Gas() {
		return false
	}
	return true
}
