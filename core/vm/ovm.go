package vm

import (
	"bytes"
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrImpureInitcode = errors.New("initCode is impure")
	OvmCREATEMethodId        = crypto.Keccak256([]byte("ovmCREATE()"))[0:4]
	OvmCREATE2MethodId        = crypto.Keccak256([]byte("ovmCREATE2()"))[0:4]
	OvmSLOADMethodId         = crypto.Keccak256([]byte("ovmSLOAD()"))[0:4]
	OvmSSTOREMethodId        = crypto.Keccak256([]byte("ovmSSTORE()"))[0:4]
	OvmContractAddress       = common.HexToAddress(os.Getenv("EXECUTION_MANAGER_ADDRESS"))
	ContractAddress          = common.HexToAddress(os.Getenv("EXECUTION_MANAGER_ADDRESS"))
	PurityCheckerAddress     = common.HexToAddress(os.Getenv("PURITY_CHECKER_ADDRESS"))
	ContractCreatorAddress   = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

var ExecutionManagerAddress = common.HexToAddress(os.Getenv("EXECUTION_MANAGER_ADDRESS"))

type ovmOperation func(*EVM, ContractRef, *Contract, []byte) ([]byte, error)

var funcs = map[string]ovmOperation{
	"ovmSSTORE()": sStore,
	"ovmSLOAD()":  sLoad,
	"ovmCREATE()": create,
	"ovmCREATE2()": create2,
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

func runOvmOperation(input []byte, evm *EVM, caller ContractRef, contract *Contract) (ret []byte, err error) {
	var methodId [4]byte
	copy(methodId[:], input[:4])
	return methodIds[methodId](evm, caller, contract, input)
}

func create(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	initCode := input[4:]
	gas := contract.Gas
	if evm.chainRules.IsEIP150 {
		gas -= gas / 64
	}
	contract.UseGas(gas)

	if isPure(evm, caller, gas, initCode) {
		_, address, _, _ := evm.Create(caller, initCode, contract.Gas, big.NewInt(0))
		return address.Bytes(), nil
	} else {
		return nil, ErrImpureInitcode
	}
}

func create2(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	initCode := input[4:]
	gas := contract.Gas
	if evm.chainRules.IsEIP150 {
		gas -= gas / 64
	}
	contract.UseGas(gas)

	if isPure(evm, caller, gas, initCode) {
		_, address, _, _ := evm.Create2(caller, initCode, contract.Gas, big.NewInt(0), big.NewInt(0))
		return address.Bytes(), nil
	} else {
		return nil, ErrImpureInitcode
	}
}

func isPure(evm *EVM, caller ContractRef, gas uint64, code []byte) bool {
	returnValue, _, _ := evm.Call(caller, PurityCheckerAddress, code, gas, big.NewInt(0))
	return bytes.Equal(returnValue, []byte{1})
}

func sLoad(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := evm.StateDB.GetState(caller.Address(), key)
	return val.Bytes(), nil
}
func sStore(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := common.BytesToHash(input[36:68])
	evm.StateDB.SetState(caller.Address(), key, val)
	return []byte{}, nil
}
