package vm

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
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
	address, ok := args["_address"].(common.Address)
	if !ok {
		return nil, errors.New("Could not parse address arg in setAccountNonce")
	}
	nonce, ok := args["_nonce"].(*big.Int)
	if !ok {
		return nil, errors.New("Could not parse nonce arg in setAccountNonce")
	}
	evm.StateDB.SetNonce(address, nonce.Uint64())
	return []interface{}{}, nil
}

func getAccountNonce(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address, ok := args["_address"].(common.Address)
	if !ok {
		return nil, errors.New("Could not parse address arg in getAccountNonce")
	}
	nonce := evm.StateDB.GetNonce(address)
	return []interface{}{new(big.Int).SetUint64(reflect.ValueOf(nonce).Uint())}, nil
}

func getAccountEthAddress(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address, ok := args["_address"].(common.Address)
	if !ok {
		return nil, errors.New("Could not parse address arg in getAccountEthAddress")
	}
	return []interface{}{address}, nil
}

func getContractStorage(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address, ok := args["_contract"].(common.Address)
	if !ok {
		return nil, errors.New("Could not parse contract arg in getContractStorage")
	}
	_key, ok := args["_key"]
	if !ok {
		return nil, errors.New("Could not parse key arg in getContractStorage")
	}
	key := toHash(_key)
	val := evm.StateDB.GetState(address, key)
	log.Debug("Got contract storage", "address", address.Hex(), "key", key.Hex(), "val", val.Hex())
	return []interface{}{val}, nil
}

func putContractStorage(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	address, ok := args["_contract"].(common.Address)
	if !ok {
		return nil, errors.New("Could not parse address arg in putContractStorage")
	}
	_key, ok := args["_key"]
	if !ok {
		return nil, errors.New("Could not parse key arg in putContractStorage")
	}
	key := toHash(_key)
	_value, ok := args["_value"]
	if !ok {
		return nil, errors.New("Could not parse value arg in putContractStorage")
	}
	val := toHash(_value)

	// save the block number and address with modified key if it's not an eth_call
	if evm.Context.EthCallSender == nil {
		// save the value before
		before := evm.StateDB.GetState(address, key)
		evm.StateDB.SetState(address, key, val)
		err := evm.StateDB.SetDiffKey(
			evm.Context.BlockNumber,
			address,
			key,
			before != val,
		)
		if err != nil {
			log.Error("error", err)
		}
	} else {
		// otherwise just do the db update
		evm.StateDB.SetState(address, key, val)
	}

	log.Debug("Put contract storage", "address", address.Hex(), "key", key.Hex(), "val", val.Hex())
	return []interface{}{}, nil
}

func nativeFunctionTrue(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	return []interface{}{true}, nil
}

func nativeFunctionVoid(evm *EVM, contract *Contract, args map[string]interface{}) ([]interface{}, error) {
	return []interface{}{}, nil
}

func toHash(arg interface{}) common.Hash {
	b := arg.([32]uint8)
	return common.BytesToHash(b[:])
}
