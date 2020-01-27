package tests

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

var KEY = common.FromHex("0102030000000000000000000000000000000000000000000000000000000000")
var VALUE1 = common.FromHex("0405060000000000000000000000000000000000000000000000000000000000")
var VALUE2 = common.FromHex("0708090000000000000000000000000000000000000000000000000000000000")

func mstoreBytes(bytes []byte, offset int) []byte {
	output := make([]byte, len(bytes)*5)
	for i, b := range bytes {
		output[i*5] = byte(vm.PUSH1)
		output[i*5+1] = b
		output[i*5+2] = byte(vm.PUSH1)
		output[i*5+3] = byte(offset + i)
		output[i*5+4] = byte(vm.MSTORE8)
	}
	return output
}

func call(addr common.Address, value uint, inOffset uint, inSize uint, retOffset uint, retSize uint) []byte {
	output := []byte{
		byte(vm.PUSH1), 0,
		byte(vm.PUSH1), 0,
		byte(vm.PUSH1), byte(retSize),
		byte(vm.PUSH1), byte(retOffset),
		byte(vm.PUSH1), byte(inSize),
		byte(vm.PUSH1), byte(inOffset),
		byte(vm.PUSH1), byte(value),
	}
	output = append(output, []byte{
		byte(vm.PUSH20)}...)
	output = append(output, addr.Bytes()...)
	output = append(output, []byte{
		byte(vm.GAS),
		byte(vm.CALL),
	}...)
	return output
}

func TestSloadAndStore(t *testing.T) {
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	codeAddr := common.HexToAddress("0xC0")
	code := mstoreBytes(vm.OvmSSTOREMethodId, 0)
	code = append(code, mstoreBytes(KEY, 4)...)
	code = append(code, mstoreBytes(VALUE1, 36)...)
	code = append(code,
		call(
			vm.OvmContractAddress,
			0,
			0,
			68,
			0,
			0)...)
	code = append(code, mstoreBytes(vm.OvmSLOADMethodId, 0)...)
	code = append(code, mstoreBytes(KEY, 4)...)
	code = append(code,
		call(
			vm.OvmContractAddress,
			0,
			0,
			36,
			0,
			32)...)
	code = append(code, []byte{
		byte(vm.POP),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}...)

	state.SetCode(codeAddr, code)

	returnValue, _, err := runtime.Call(codeAddr, nil, &runtime.Config{State: state, Debug: true})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	if !bytes.Equal(VALUE1, returnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE1, returnValue)
	}
}

func TestSstoreDoesntOverwrite(t *testing.T) {
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	aliceAddr := common.HexToAddress("0x0a")
	bobAddr := common.HexToAddress("0x0b")
	store1CodeAddr := common.HexToAddress("0xC1")
	store2CodeAddr := common.HexToAddress("0xC2")
	loadCodeAddr := common.HexToAddress("0xC3")
	store1Code := storeCode(KEY, VALUE1)
	store2Code := storeCode(KEY, VALUE2)
	loadCode := mstoreBytes(vm.OvmSLOADMethodId, 0)
	loadCode = append(loadCode, mstoreBytes(KEY, 4)...)
	loadCode = append(loadCode,
		call(
			vm.OvmContractAddress,
			0,
			0,
			36,
			0,
			32)...)
	loadCode = append(loadCode, []byte{
		byte(vm.POP),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}...)
	state.SetCode(store1CodeAddr, store1Code)
	state.SetCode(store2CodeAddr, store2Code)
	state.SetCode(loadCodeAddr, loadCode)
	_, _, err := runtime.Call(store1CodeAddr, nil, &runtime.Config{State: state, Origin: aliceAddr, Debug: true})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	_, _, err = runtime.Call(store2CodeAddr, nil, &runtime.Config{State: state, Origin: bobAddr, Debug: true})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	aliceReturnValue, _, err := runtime.Call(loadCodeAddr, nil, &runtime.Config{State: state, Origin: aliceAddr, Debug: true})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}
	bobReturnValue, _, err := runtime.Call(loadCodeAddr, nil, &runtime.Config{State: state, Origin: bobAddr, Debug: true})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	if !bytes.Equal(VALUE1, aliceReturnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE1, aliceReturnValue)
	}
	if !bytes.Equal(VALUE2, bobReturnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE2, bobReturnValue)
	}
}

func storeCode(key []byte, value []byte) []byte {
	storeCode := mstoreBytes(vm.OvmSSTOREMethodId, 0)
	storeCode = append(storeCode, mstoreBytes(key, 4)...)
	storeCode = append(storeCode, mstoreBytes(value, 36)...)
	storeCode = append(storeCode,
		call(
			vm.OvmContractAddress,
			0,
			0,
			68,
			0,
			0)...)
	storeCode = append(storeCode, []byte{
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}...)
	return storeCode
}
