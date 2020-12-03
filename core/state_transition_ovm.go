package core

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

type ovmTransaction struct {
	Timestamp     *big.Int       "json:\"timestamp\""
	BlockNumber   *big.Int       "json:\"blockNumber\""
	L1QueueOrigin uint8          "json:\"l1QueueOrigin\""
	L1TxOrigin    common.Address "json:\"l1TxOrigin\""
	Entrypoint    common.Address "json:\"entrypoint\""
	GasLimit      *big.Int       "json:\"gasLimit\""
	Data          []uint8        "json:\"data\""
}

func toExecutionManagerRun(evm *vm.EVM, msg Message) (Message, error) {
	tx := ovmTransaction{
		evm.Context.Time,
		evm.Context.BlockNumber, // TODO (what's the correct block number?)
		uint8(msg.QueueOrigin().Uint64()),
		*msg.L1MessageSender(),
		*msg.To(),
		big.NewInt(int64(msg.Gas())),
		msg.Data(),
	}

	var abi = evm.Context.OvmExecutionManager.ABI
	var args = []interface{}{
		tx,
		evm.OvmStateManager.Address,
	}

	ret, err := abi.Pack("run", args...)
	if err != nil {
		return nil, err
	}

	outputmsg, err := modMessage(
		msg,
		msg.From(),
		&evm.Context.OvmExecutionManager.Address,
		ret,
		evm.Context.GasLimit,
	)
	if err != nil {
		return nil, err
	}

	return outputmsg, nil
}

func asOvmMessage(tx *types.Transaction, signer types.Signer, decompressor common.Address) (Message, error) {
	msg, err := tx.AsMessage(signer)
	if err != nil {
		return msg, err
	}

	// Queue origin L1ToL2 transactions do not go through the
	// sequencer entrypoint. The calldata is expected to be in the
	// correct format when deserialized from the EVM events, see
	// rollup/sync_service.go.
	qo := msg.QueueOrigin()
	if qo != nil && qo.Uint64() == uint64(types.QueueOriginL1ToL2) {
		return msg, nil
	}

	v, r, s := tx.RawSignatureValues()

	// V parameter here will include the chain ID, so we need to recover the original V. If the V
	// does not equal zero or one, we have an invalid parameter and need to throw an error.
	// TODO: the chainid needs to be pulled in from config
	v = big.NewInt(int64(v.Uint64() - 35 - 2*420))
	if v.Uint64() != 0 && v.Uint64() != 1 {
		return msg, fmt.Errorf("invalid signature v parameter")
	}

	// Since we use a fixed encoding, we need to insert some placeholder address to represent that
	// the user wants to create a contract (in this case, the zero address).
	var target common.Address
	if tx.To() == nil {
		target = ZeroAddress
	} else {
		target = *tx.To()
	}

	// Sequencer uses a custom encoding structure --
	// We originally receive sequencer transactions encoded in this way, but we decode them before
	// inserting into Geth so we can make transactions easily parseable. However, this means that
	// we need to re-encode the transactions before executing them.
	var data = new(bytes.Buffer)
	data.WriteByte(getSignatureType(msg))                    // 1 byte: 00 == EIP 155, 02 == ETH Sign Message
	data.Write(fillBytes(r, 32))                             // 32 bytes: Signature `r` parameter
	data.Write(fillBytes(s, 32))                             // 32 bytes: Signature `s` parameter
	data.Write(fillBytes(v, 1))                              // 1 byte: Signature `v` parameter
	data.Write(fillBytes(big.NewInt(int64(msg.Gas())), 3))   // 3 bytes: Gas limit
	data.Write(fillBytes(msg.GasPrice(), 3))                 // 3 bytes: Gas price
	data.Write(fillBytes(big.NewInt(int64(msg.Nonce())), 3)) // 3 bytes: Nonce
	data.Write(target.Bytes())                               // 20 bytes: Target address
	data.Write(msg.Data())                                   // ?? bytes: Transaction data

	// Sequencer transactions get sent to the "sequencer entrypoint," a contract that decompresses
	// the incoming transaction data.
	outmsg, err := modMessage(
		msg,
		msg.From(),
		&decompressor,
		data.Bytes(),
		msg.Gas(),
	)

	if err != nil {
		return msg, err
	}

	return outmsg, nil
}

func EncodeFakeMessage(
	msg Message,
	account abi.ABI,
) (Message, error) {
	var input = []interface{}{
		big.NewInt(int64(msg.Gas())),
		msg.To(),
		msg.Data(),
	}

	output, err := account.Pack("qall", input...)
	if err != nil {
		return nil, err
	}

	var from = msg.From()
	return modMessage(
		msg,
		from,
		&from,
		output,
		msg.Gas(),
	)
}

func modMessage(
	msg Message,
	from common.Address,
	to *common.Address,
	data []byte,
	gasLimit uint64,
) (Message, error) {
	queueOrigin, err := getQueueOrigin(msg.QueueOrigin())
	if err != nil {
		return nil, err
	}

	outmsg := types.NewMessage(
		from,
		to,
		msg.Nonce(),
		msg.Value(),
		gasLimit,
		msg.GasPrice(),
		data,
		false,
		msg.L1MessageSender(),
		msg.L1BlockNumber(),
		queueOrigin,
		msg.SignatureHashType(),
	)

	return outmsg, nil
}

func getSignatureType(
	msg Message,
) uint8 {
	if msg.SignatureHashType() == 0 {
		return 0
	} else if msg.SignatureHashType() == 1 {
		return 2
	} else {
		return 1
	}
}

func getQueueOrigin(
	queueOrigin *big.Int,
) (types.QueueOrigin, error) {
	if queueOrigin.Cmp(big.NewInt(0)) == 0 {
		return types.QueueOriginSequencer, nil
	} else if queueOrigin.Cmp(big.NewInt(1)) == 0 {
		return types.QueueOriginL1ToL2, nil
	} else if queueOrigin.Cmp(big.NewInt(2)) == 0 {
		return types.QueueOriginL1ToL2, nil
	} else {
		return types.QueueOriginSequencer, fmt.Errorf("invalid queue origin: %d", queueOrigin)
	}
}

func fillBytes(x *big.Int, size int) []byte {
	b := x.Bytes()
	switch {
	case len(b) > size:
		panic("math/big: value won't fit requested size")
	case len(b) == size:
		return b
	default:
		buf := make([]byte, size)
		copy(buf[size-len(b):], b)
		return buf
	}
}
