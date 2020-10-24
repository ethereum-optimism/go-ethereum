package core

import (
	"fmt"
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)


var GodAddress = common.HexToAddress("0x444400000000000000000000000000000000000")
var NullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var NullHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")

type OVMTransaction struct {
	Timestamp *big.Int 			"json:\"timestamp\""
	BlockNumber *big.Int 		"json:\"blockNumber\""
	L1QueueOrigin uint8 		"json:\"l1QueueOrigin\""
	L1TxOrigin common.Address 	"json:\"l1TxOrigin\""
	Entrypoint common.Address 	"json:\"entrypoint\""
	GasLimit *big.Int 			"json:\"gasLimit\""
	Data []uint8 				"json:\"data\""
}

func NeedsEOACreate(tx *types.Transaction, signer types.Signer, statedb *state.StateDB) (bool, error) {
	msg, err := tx.AsMessage(signer)
	if err != nil {
		return false, err
	}

	if (msg.From() == GodAddress) {
		return false, nil
	} else {
		return statedb.GetCodeSize(msg.From()) == 0, nil
	}
}

func ToOvmMessage(tx *types.Transaction, signer types.Signer, eoa bool) (Message, error) {
	msg, err := tx.AsMessage(signer)
	if err != nil {
		return nil, err
	}

	var inputmsg Message
	if (msg.From() == GodAddress) {
		inputmsg = msg
	} else {
		inputmsg, err = encodeCompressedMessage(msg, tx, signer, eoa)
		if err != nil {
			return nil, err
		}
	}

	// Take the compressed message and encode it into something the execution
	// manager can understand (i.e., the input to "run").
	data, err := encodeExecutionManagerRun(inputmsg)
	if err != nil {
		return nil, err
	}

	outputmsg, err := modMessage(
		inputmsg,
		inputmsg.From(),
		&vm.OVMExecutionManager.Address,
		data,
	)
	if err != nil {
		return nil, err
	}

	return outputmsg, nil
}

func encodeCompressedMessage(
	msg Message,
	tx *types.Transaction,
	signer types.Signer,
	eoa bool,
) (Message, error) {
	v, r, s := tx.RawSignatureValues()
	v = new(big.Int).Mod(v, big.NewInt(256))
	var data = new(bytes.Buffer)

	var sigtype uint8
	if eoa {
		sigtype = 0
	} else {
		sigtype = getSignatureType(msg)
	}

	var target common.Address
	if tx.To() == nil {
		target = NullAddress
	} else {
		target = *tx.To()
	}

	// Signature type
	data.WriteByte(byte(sigtype)) 									// 1 byte: 00 == EOACreate, 01 == EIP 155, 02 == ETH Sign Message

	// Signature data
	data.Write(v.FillBytes(make([]byte, 1, 1)))										// 1 byte: Signature `v` parameter
	data.Write(r.FillBytes(make([]byte, 32, 32)))									// 32 bytes: Signature `r` parameter
	data.Write(s.FillBytes(make([]byte, 32, 32)))									// 32 bytes: Signature `s` parameter

	if (sigtype == 0) {
		// EOACreate: Encode the transaction hash.
		data.Write(tx.Hash().Bytes())												// 32 bytes: Transaction hash
	} else {
		// EIP 155 or ETH Sign Message: Encode the full transaction data.
		data.Write(big.NewInt(int64(msg.Nonce())).FillBytes(make([]byte, 2, 2))) 	// 2 bytes: Nonce
		data.Write(big.NewInt(int64(msg.Gas())).FillBytes(make([]byte, 3, 3)))	 	// 3 bytes: Gas limit
		data.Write(msg.GasPrice().FillBytes(make([]byte, 1, 1)))				 	// 1 byte: Gas price
		data.Write(tx.ChainId().FillBytes(make([]byte, 4, 4)))		 				// 4 bytes: Chain ID
		data.Write(target.Bytes())											 		// 20 bytes: Target address
		data.Write(msg.Data())													 	// ?? bytes: Transaction data
	}

	fmt.Printf("ORIGINAL DATA: %x\n", data.Bytes())

	outmsg, err := modMessage(
		msg,
		GodAddress,
		&vm.OVMSequencerMessageDecompressor.Address,
		data.Bytes(),
	)
	if err != nil {
		return nil, err
	}

	return outmsg, nil
}

func encodeExecutionManagerRun(
	msg Message,
) ([]byte, error) {
	ovmtx := OVMTransaction{
		big.NewInt(0), // TODO
		big.NewInt(0), // TODO
		uint8(msg.QueueOrigin().Uint64()),
		*msg.L1MessageSender(),
		*msg.To(),
		big.NewInt(int64(msg.Gas())),
		msg.Data(),
	}

	var abi = vm.OVMExecutionManager.ABI
	var args = []interface{}{
		ovmtx,
		vm.OVMStateManager.Address,
	}

	ret, err := abi.Pack("run", args...)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func modMessage(
	msg Message,
	from common.Address,
	to *common.Address,
	data []byte,
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
		msg.Gas(),
		msg.GasPrice(),
		data,
		false,
		msg.L1MessageSender(),
		msg.L1RollupTxId(),
		queueOrigin,
		msg.SignatureHashType(),
	)

	return outmsg, nil
}

func getQueueOrigin(
	queueOrigin *big.Int,
) (types.QueueOrigin, error) {
	if (queueOrigin.Cmp(big.NewInt(0)) == 0) {
		return types.QueueOriginSequencer, nil
	} else if (queueOrigin.Cmp(big.NewInt(1)) == 0) {
		return types.QueueOriginL1ToL2, nil
	} else {
		return types.QueueOriginSequencer, fmt.Errorf("Invalid queue origin: %d\n", queueOrigin)
	}
}

func getSignatureType(
	msg Message,
) uint8 {
	if (msg.SignatureHashType() == 0) {
		return 1
	} else if (msg.SignatureHashType() == 1) {
		return 2
	} else {
		return 0
	}
}
