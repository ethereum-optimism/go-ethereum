package vm

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

type stateManagerFunction func(*EVM, *Contract, map[string]interface{}) ([]interface{}, error)

var funcs = map[string]stateManagerFunction{
	"owner":                                    owner,
	"setAccountNonce":                          setAccountNonce,
	"getAccountNonce":                          getAccountNonce,
	"getAccountEthAddress":                     getAccountEthAddress,
	"getContractStorage":                       getContractStorage,
	"putContractStorage":                       putContractStorage,
	"isAuthenticated":                          nativeFunctionTrue,
	"hasAccount":                               nativeFunctionTrue,
	"hasEmptyAccount":                          nativeFunctionTrue,
	"hasContractStorage":                       nativeFunctionTrue,
	"testAndSetAccountLoaded":                  nativeFunctionTrue,
	"testAndSetAccountChanged":                 nativeFunctionTrue,
	"testAndSetContractStorageLoaded":          nativeFunctionTrue,
	"testAndSetContractStorageChanged":         nativeFunctionTrue,
	"incrementTotalUncommittedAccounts":        nativeFunctionVoid,
	"incrementTotalUncommittedContractStorage": nativeFunctionVoid,
	"initPendingAccount":                       nativeFunctionVoid,
	"commitPendingAccount":                     nativeFunctionVoid,
}

func callStateManager(input []byte, evm *EVM, contract *Contract) (ret []byte, err error) {
	rawabi := OvmStateManager.ABI
	abi := &rawabi

	method, err := abi.MethodById(input)
	if err != nil {
		return nil, err
	}

	var inputArgs = make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(inputArgs, input[4:])
	if err != nil {
		return nil, err
	}

	fn, exist := funcs[method.RawName]
	if !exist {
		return nil, fmt.Errorf("Native OVM_StateManager function not found for method '%s'", method.RawName)
	}

	outputArgs, err := fn(evm, contract, inputArgs)
	if err != nil {
		return nil, err
	}

	returndata, err := method.Outputs.PackValues(outputArgs)
	if err != nil {
		return nil, err
	}

	return returndata, nil
}

func owner(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	origin := evm.Context.Origin

	return []interface{}{origin}, nil
}

func setAccountNonce(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address := args["_address"].(common.Address)
	nonce := args["_nonce"].(*big.Int)

	evm.StateDB.SetNonce(address, nonce.Uint64())

	return []interface{}{}, nil
}

func getAccountNonce(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address := args["_address"].(common.Address)

	nonce := evm.StateDB.GetNonce(address)
	return []interface{}{new(big.Int).SetUint64(reflect.ValueOf(nonce).Uint())}, nil
}

func getAccountEthAddress(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address := args["_address"].(common.Address)

	return []interface{}{address}, nil
}

func getContractStorage(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address := args["_contract"].(common.Address)
	key := toHash(args["_key"])

	val := evm.StateDB.GetState(address, key)

	return []interface{}{val}, nil
}

func putContractStorage(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address := args["_contract"].(common.Address)
	key := toHash(args["_key"])
	val := toHash(args["_value"])

	evm.StateDB.SetState(address, key, val)

	return []interface{}{}, nil
}

func nativeFunctionTrue(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	return []interface{}{true}, nil
}

func nativeFunctionVoid(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	return []interface{}{}, nil
}

func toHash(arg interface{}) common.Hash {
	b := [32]byte(arg.([32]uint8))
	return common.BytesToHash(b[:])
}
