package tests

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/crypto"
)

var KEY = common.FromHex("0102030000000000000000000000000000000000000000000000000000000000")
var VALUE1 = common.FromHex("0405060000000000000000000000000000000000000000000000000000000000")
var VALUE2 = common.FromHex("0708090000000000000000000000000000000000000000000000000000000000")
var INIT_CODE = common.FromHex("608060405234801561001057600080fd5b5060405161026b38038061026b8339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506101d7806100946000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633408f73a1461003b578063d3404b6d14610045575b600080fd5b61004361004f565b005b61004d6100fa565b005b600060e060405180807f6f766d534c4f4144282900000000000000000000000000000000000000000000815250600a0190506040518091039020901c905060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060405136600082378260181c81538260101c60018201538260081c60028201538260038201536040516207a1208136846000875af160008114156100f657600080fd5b3d82f35b600060e060405180807f6f766d5353544f52452829000000000000000000000000000000000000000000815250600b0190506040518091039020901c905060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060405136600082378260181c81538260101c60018201538260081c600282015382600382015360008036836000865af1600081141561019c57600080fd5b5050505056fea265627a7a7231582047df4ba501514f65ab1e6f8215402e9240cb0cf954d608cdc4158258f468b12364736f6c634300050c0032000000000000000000000000fdfef9d10d929cb3905c71400ce6be1990ea0f34")

func mstoreBytes(bytes []byte, offset int) []byte {
	output := []byte{}
	for i := 0; i < len(bytes); i += vm.WORD_SIZE {
		end := i + vm.WORD_SIZE
		if end > len(bytes) {
			end = len(bytes)
		}
		output = append(output, byte(vm.PUSH32))
		output = append(output, common.RightPadBytes(bytes[i:end], vm.WORD_SIZE)...)
		output = append(output, pushN(int64(offset+i)))
		output = append(output, int64ToBytes(int64(offset+i))...)
		output = append(output, byte(vm.MSTORE))
	}
	return output
}

func call(addr common.Address, value int64, inOffset int64, inSize int64, retOffset int64, retSize int64) []byte {
	output := []byte{}
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(retSize))
	output = append(output, int64ToBytes(retSize)...)
	output = append(output, pushN(retOffset))
	output = append(output, int64ToBytes(retOffset)...)
	output = append(output, pushN(inSize))
	output = append(output, int64ToBytes(inSize)...)
	output = append(output, pushN(inOffset))
	output = append(output, int64ToBytes(inOffset)...)
	output = append(output, pushN(value))
	output = append(output, int64ToBytes(value)...)
	output = append(output, []byte{
		byte(vm.PUSH20)}...)
	output = append(output, addr.Bytes()...)
	output = append(output, []byte{
		byte(vm.GAS),
		byte(vm.CALL),
	}...)
	return output
}

func int64ToBytes(n int64) []byte {
	if bytes.Equal(big.NewInt(n).Bytes(), []byte{}) {
		return []byte{0}
	} else {
		return big.NewInt(n).Bytes()
	}
}
func pushN(n int64) byte {
	return byte(int(vm.PUSH1) + byteLength(n) - 1)
}
func byteLength(n int64) int {
	if bytes.Equal(big.NewInt(n).Bytes(), []byte{}) {
    return 1
  } else {
		return len(big.NewInt(n).Bytes())
 }
}

func mockPurityChecker(pure bool) []byte {
  var pureByte byte

  if pure {
    pureByte = 1
  } else {
    pureByte = 0
  }

	return []byte{
		byte(vm.PUSH1),
    pureByte,
		byte(vm.PUSH1),
    0,
		byte(vm.MSTORE8),
		byte(vm.PUSH1),
    1,
		byte(vm.PUSH1),
    0,
		byte(vm.RETURN),
  }
}

func TestCreateImpure(t *testing.T) {
  vm.PurityCheckerAddress = common.HexToAddress("0x0A")
	aliceAddr := common.HexToAddress("0x00")
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	codeAddr := common.HexToAddress("0xC0")
	initCode := INIT_CODE
	code := mstoreBytes(vm.OvmCREATEMethodId, 0)
	code = append(code, mstoreBytes(initCode, 4)...)
	code = append(code,
		call(
			vm.OvmContractAddress,
			0,
			0,
			int64(len(initCode))+4,
			0,
			32)...)
	code = append(code, []byte{
		byte(vm.PUSH1), 0,
		byte(vm.MSTORE8),
		byte(vm.PUSH1), 1,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}...)

	state.SetCode(codeAddr, code)
	state.SetCode(vm.PurityCheckerAddress, mockPurityChecker(false))

	returnVal, _, err := runtime.Call(codeAddr, nil, &runtime.Config{State: state, Debug: true, Origin: aliceAddr})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	expectedVal := common.LeftPadBytes([]byte{0}, 1)
	if !bytes.Equal(expectedVal, returnVal) {
		t.Errorf("Expected %020x; got %020x", expectedVal, returnVal)
	}
}

func TestCreate(t *testing.T) {
  vm.PurityCheckerAddress = common.HexToAddress("0x0A")
	aliceAddr := common.HexToAddress("0x00")
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	codeAddr := common.HexToAddress("0xC0")
	initCode := INIT_CODE
	code := mstoreBytes(vm.OvmCREATEMethodId, 0)
	code = append(code, mstoreBytes(initCode, 4)...)
	code = append(code,
		call(
			vm.OvmContractAddress,
			0,
			0,
			int64(len(initCode))+4,
			0,
			32)...)
	code = append(code, []byte{
		byte(vm.POP),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}...)

	state.SetCode(codeAddr, code)
	state.SetCode(vm.PurityCheckerAddress, mockPurityChecker(true))

	returnVal, _, err := runtime.Call(codeAddr, nil, &runtime.Config{State: state, Debug: true, Origin: aliceAddr})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	expectedVal := common.LeftPadBytes(crypto.CreateAddress(vm.ContractCreatorAddress, 0).Bytes(), vm.WORD_SIZE)
	if !bytes.Equal(expectedVal, returnVal) {
		t.Errorf("Expected %020x; got %020x", expectedVal, returnVal)
	}
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
