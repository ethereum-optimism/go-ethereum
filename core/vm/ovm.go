package vm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrImpureInitcode       = errors.New("initCode is impure")
	ExecutionManagerAddress = common.HexToAddress(os.Getenv("EXECUTION_MANAGER_ADDRESS"))
	PurityCheckerAddress    = common.HexToAddress(os.Getenv("PURITY_CHECKER_ADDRESS"))
	WORD_SIZE               = 32
)

type ovmOperation func(*EVM, ContractRef, *Contract, []byte) ([]byte, error)
type methodId [4]byte

var funcs = map[string]ovmOperation{
	"SSTORE":  sStore,
	"SLOAD":   sLoad,
	"CREATE":  create,
	"CREATE2": create2,
}
var methodIds map[[4]byte]ovmOperation

func init() {
	methodIds = make(map[[4]byte]ovmOperation, len(funcs))
	for methodName, f := range funcs {
		methodIds[OvmMethodId(methodName)] = f
	}
}

func OvmMethodId(methodName string) [4]byte {
	var methodId [4]byte
	var fullMethodName = "ovm" + methodName + "()"
	copy(methodId[:], crypto.Keccak256([]byte(fullMethodName)))
	return methodId
}

func isOvmOperation(contract *Contract, input []byte) bool {
	if contract.Address() != ExecutionManagerAddress {
		// fmt.Printf("%020x == %020x\n", contract.Address(), ExecutionManagerAddress)
		return false
	} else {
		// fmt.Printf("%020x == %020x\n", contract.Address(), ExecutionManagerAddress)
	}
	if len(input) < 4 {
		return false
	}
	for methodId := range methodIds {
		// fmt.Printf("MethodId: %x\n", input[:4])
		if bytes.Equal(input[0:4], methodId[:]) {
			// fmt.Println("isOvmOperation!")
			return true
		} else {
			// fmt.Printf("not MethodId: %x\n", methodId)
		}
	}
	// fmt.Println("not methodId")
	return false
}

func runOvmOperation(input []byte, evm *EVM, caller ContractRef, contract *Contract) (ret []byte, err error) {
	fmt.Println("runOvmOperation")
	var methodId [4]byte
	copy(methodId[:], input[:4])
	return methodIds[methodId](evm, caller, contract, input)
}

func create(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	fmt.Println("ovmCREATE")
	initCode := input[4:]
	gas := contract.Gas
	if evm.chainRules.IsEIP150 {
		gas -= gas / 64
	}
	contract.UseGas(gas)

	// contractRef := &Contract{self: AccountRef(contract.Address())}
	// if isPure(evm, contractRef, gas, initCode) {
	_, address, _, _ := evm.Create(caller, initCode, contract.Gas, big.NewInt(0))
	// fmt.Printf("Pure %x\n", address.Bytes())
	// emitEvent(evm, contract, "test", []byte{})
	// emitActiveContract(address)
	// emitCreatedContract(address,a
	emitActiveContract(evm, contract, address)
	emitCreatedContract(evm, contract, address, address, [32]byte{1})
	return address.Bytes(), nil
	// return nil, nil
	// } else {
	//   return nil, ErrImpureInitcode
	// }
}

func create2(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	initCode := input[4:]
	gas := contract.Gas
	if evm.chainRules.IsEIP150 {
		gas -= gas / 64
	}
	contract.UseGas(gas)

	// if isPure(evm, caller, gas, initCode) {
	_, address, _, _ := evm.Create2(caller, initCode, contract.Gas, big.NewInt(0), big.NewInt(0))
	return address.Bytes(), nil
	// } else {
	// 	return nil, ErrImpureInitcode
	// }
}

func sLoad(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := evm.StateDB.GetState(caller.Address(), key)
	fmt.Printf("%x > %x\n", key, val)
	return val.Bytes(), nil
}
func sStore(evm *EVM, caller ContractRef, contract *Contract, input []byte) (ret []byte, err error) {
	key := common.BytesToHash(input[4:36])
	val := common.BytesToHash(input[36:68])
	evm.StateDB.SetState(caller.Address(), key, val)
	fmt.Printf("%x = %x\n", key, val)
	return []byte{}, nil
}

func isPure(evm *EVM, caller ContractRef, gas uint64, code []byte) bool {
	return true
}

func emitActiveContract(
	evm *EVM,
	contract *Contract,
	contractAddress common.Address,
) {
	typ, _ := abi.NewType("(address)", "", []abi.ArgumentMarshaling{})
	data, _ := typ.Pack(reflect.ValueOf(contractAddress))
	emitEvent(evm, contract, "ActiveContract(address)", data)
}

func emitCreatedContract(
	evm *EVM,
	contract *Contract,
	ovmContractAddress common.Address,
	codeContractAddress common.Address,
	codeContractHash [32]byte,
) {
	typ, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "a", Type: "address"},
		{Name: "b", Type: "address"},
		{Name: "c", Type: "bytes32"},
	})
	data, _ := typ.Pack(reflect.ValueOf(struct {
		A common.Address
		B common.Address
		C [32]byte
	}{
		ovmContractAddress,
		codeContractAddress,
		codeContractHash,
	}))

	emitEvent(evm, contract, "CreatedContract(address,address,bytes32)", data)
}

func emitEvent(evm *EVM, contract *Contract, topic string, data []byte) {
	fmt.Printf("topic %s\n", topic)
	fmt.Printf("hash: %x\n", crypto.Keccak256([]byte(topic)))
	// fmt.Printf("topic2 %s\n", "CreatedContract(address,address,bytes32)")
	topics := []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(topic)))}
	evm.StateDB.AddLog(&types.Log{
		Address:     contract.Address(),
		Topics:      topics,
		Data:        data,
		BlockNumber: evm.BlockNumber.Uint64(),
	})
}
