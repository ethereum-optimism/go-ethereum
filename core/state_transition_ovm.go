package core

import (
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
		evm.Context.OvmStateManager.Address,
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

	from := msg.From()
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
	)

	return outmsg, nil
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
