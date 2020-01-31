package vm

import (
	"bytes"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var ExecutionManagerAddress = common.HexToAddress(os.Getenv("EXECUTION_MANAGER_ADDRESS"))

type ovmOperation func(*EVM, ContractRef, []byte) ([]byte, error)

var funcs = map[string]ovmOperation{
	"ovmSSTORE()": sStore,
	"ovmSLOAD()":  sLoad,
	"ovmCREATE()": create,
}
var methodIds map[[4]byte]ovmOperation

func init() {
	methodIds = make(map[[4]byte]ovmOperation, 4)
	var methodId [4]byte
	for methodName, f := range funcs {
		copy(methodId[:], crypto.Keccak256([]byte(methodName)))
		methodIds[methodId] = f
	}
}

func isOvmOperation(contract *Contract, input []byte) bool {
	if contract.Address() != ExecutionManagerAddress {
		return false
	}
	if len(input) < 4 {
		return false
	}
	for methodId := range methodIds {
		if bytes.Equal(input[0:4], methodId[:]) {
			return true
		}
	}
	return false
}

func runOvmOperation(input []byte, evm *EVM, caller ContractRef) (ret []byte, err error) {
	var methodId [4]byte
	copy(methodId[:], input[:4])
	return methodIds[methodId](evm, caller, input)
}

func create(evm *EVM, caller ContractRef, input []byte) (ret []byte, err error) {
	return []byte{}, nil
}
func sLoad(evm *EVM, caller ContractRef, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := evm.StateDB.GetState(caller.Address(), key)
	return val.Bytes(), nil
}
func sStore(evm *EVM, caller ContractRef, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := common.BytesToHash(input[36:68])
	evm.StateDB.SetState(caller.Address(), key, val)
	return []byte{}, nil
}
