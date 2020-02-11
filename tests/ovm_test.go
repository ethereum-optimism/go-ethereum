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
	"github.com/ethereum/go-ethereum/params"
)

var KEY = common.FromHex("0102030000000000000000000000000000000000000000000000000000000000")
var VALUE1 = common.FromHex("0405060000000000000000000000000000000000000000000000000000000000")
var VALUE2 = common.FromHex("0708090000000000000000000000000000000000000000000000000000000000")
var INIT_CODE = common.FromHex("608060405234801561001057600080fd5b5060405161026b38038061026b8339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506101d7806100946000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633408f73a1461003b578063d3404b6d14610045575b600080fd5b61004361004f565b005b61004d6100fa565b005b600060e060405180807f6f766d534c4f4144282900000000000000000000000000000000000000000000815250600a0190506040518091039020901c905060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060405136600082378260181c81538260101c60018201538260081c60028201538260038201536040516207a1208136846000875af160008114156100f657600080fd5b3d82f35b600060e060405180807f6f766d5353544f52452829000000000000000000000000000000000000000000815250600b0190506040518091039020901c905060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060405136600082378260181c81538260101c60018201538260081c600282015382600382015360008036836000865af1600081141561019c57600080fd5b5050505056fea265627a7a7231582047df4ba501514f65ab1e6f8215402e9240cb0cf954d608cdc4158258f468b12364736f6c634300050c0032000000000000000000000000fdfef9d10d929cb3905c71400ce6be1990ea0f34")
var chainConfig params.ChainConfig

func init() {
	chainConfig = params.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      new(big.Int),
		ByzantiumBlock:      new(big.Int),
		ConstantinopleBlock: new(big.Int),
		DAOForkBlock:        new(big.Int),
		DAOForkSupport:      false,
		EIP150Block:         new(big.Int),
		EIP155Block:         new(big.Int),
		EIP158Block:         new(big.Int),
	}
}

func TestCreate(t *testing.T) {
	vm.PurityCheckerAddress = common.HexToAddress("0x0A")
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	codeAddr := common.HexToAddress("0xC0")
	state, _ := state.New(common.Hash{}, db)
	createCode := ovmMethodId("CREATE")
	createCode = append(createCode, INIT_CODE...)
	state.SetCode(vm.PurityCheckerAddress, mockPurityChecker(true))
	state.SetCode(codeAddr, callCode(vm.ExecutionManagerAddress))
	returnVal, _ := call(t, state, codeAddr, createCode)
	expectedVal := crypto.CreateAddress(codeAddr, 0).Bytes()
	if !bytes.Equal(expectedVal, returnVal) {
		t.Errorf("Expected %020x; got %020x", expectedVal, returnVal)
	}
}

func TestCreateImpure(t *testing.T) {
	vm.PurityCheckerAddress = common.HexToAddress("0x0A")
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	codeAddr := common.HexToAddress("0xC0")
	state, _ := state.New(common.Hash{}, db)
	createCode := ovmMethodId("CREATE")
	createCode = append(createCode, INIT_CODE...)
	state.SetCode(vm.PurityCheckerAddress, mockPurityChecker(false))
	state.SetCode(codeAddr, callCode(vm.ExecutionManagerAddress))
	returnVal, _ := call(t, state, codeAddr, createCode)
	expectedVal := []byte{}
	if !bytes.Equal(expectedVal, returnVal) {
		t.Errorf("Expected %x; got %x", expectedVal, returnVal)
	}
}

func TestCreate2(t *testing.T) {
	vm.PurityCheckerAddress = common.HexToAddress("0x0A")
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	codeAddr := common.HexToAddress("0xC0")
	state, _ := state.New(common.Hash{}, db)
	createCode := ovmMethodId("CREATE2")
	createCode = append(createCode, INIT_CODE...)
	state.SetCode(vm.PurityCheckerAddress, mockPurityChecker(true))
	state.SetCode(codeAddr, callCode(vm.ExecutionManagerAddress))
	returnVal, _ := call(t, state, codeAddr, createCode)
	salt := [32]byte{}
	expectedVal := crypto.CreateAddress2(codeAddr, salt, crypto.Keccak256(INIT_CODE)).Bytes()
	if !bytes.Equal(expectedVal, returnVal) {
		t.Errorf("Expected %020x; got %020x", expectedVal, returnVal)
	}
}

func TestSloadAndStore(t *testing.T) {
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	codeAddr := common.HexToAddress("0xC0")
	storeCode := ovmMethodId("SSTORE")
	storeCode = append(storeCode, KEY...)
	storeCode = append(storeCode, VALUE1...)
	loadCode := ovmMethodId("SLOAD")
	loadCode = append(loadCode, KEY...)
	state.SetCode(codeAddr, callCode(vm.ExecutionManagerAddress))
	call(t, state, codeAddr, storeCode)
	returnValue, _ := call(t, state, codeAddr, loadCode)
	if !bytes.Equal(VALUE1, returnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE1, returnValue)
	}
}

func call(t *testing.T, state *state.StateDB, address common.Address, callData []byte) ([]byte, error) {
	returnValue, _, err := runtime.Call(address, callData, &runtime.Config{
		State:       state,
		ChainConfig: &chainConfig,
	})

	return returnValue, err
}

func TestSstoreDoesntOverwrite(t *testing.T) {
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, _ := state.New(common.Hash{}, db)
	codeAddr1 := common.HexToAddress("0xC1")
	codeAddr2 := common.HexToAddress("0xC2")
	storeCode1 := ovmMethodId("SSTORE")
	storeCode1 = append(storeCode1, KEY...)
	storeCode1 = append(storeCode1, VALUE1...)
	storeCode2 := ovmMethodId("SSTORE")
	storeCode2 = append(storeCode2, KEY...)
	storeCode2 = append(storeCode2, VALUE2...)
	loadCode := ovmMethodId("SLOAD")
	loadCode = append(loadCode, KEY...)
	state.SetCode(codeAddr1, callCode(vm.ExecutionManagerAddress))
	state.SetCode(codeAddr2, callCode(vm.ExecutionManagerAddress))

	state.SetCode(vm.ExecutionManagerAddress, []byte{})
	call(t, state, codeAddr1, storeCode1)
	call(t, state, codeAddr2, storeCode2)
	code1ReturnValue, _ := call(t, state, codeAddr1, loadCode)
	code2ReturnValue, _ := call(t, state, codeAddr2, loadCode)

	if !bytes.Equal(VALUE1, code1ReturnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE1, code1ReturnValue)
	}
	if !bytes.Equal(VALUE2, code2ReturnValue) {
		t.Errorf("Expected %020x; got %020x", VALUE2, code2ReturnValue)
	}
}

/*
  callCode generates EVM bytecode which makes a single CALL with call data as
  it's input.
*/
func callCode(addr common.Address) []byte {
	output := []byte{}
	output = append(output, []byte{
		byte(vm.CALLDATASIZE),
		byte(vm.PUSH1), 0,
		byte(vm.PUSH1), 0,
		byte(vm.CALLDATACOPY),
	}...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, byte(vm.CALLDATASIZE))
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, pushN(0))
	output = append(output, int64ToBytes(0)...)
	output = append(output, []byte{
		byte(vm.PUSH20)}...)
	output = append(output, addr.Bytes()...)
	output = append(output, []byte{
		byte(vm.GAS),
		byte(vm.CALL),
		byte(vm.POP),
		byte(vm.RETURNDATASIZE),
		byte(vm.PUSH1), 0,
		byte(vm.PUSH1), 0,
		byte(vm.RETURNDATACOPY),
		byte(vm.RETURNDATASIZE),
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
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

func ovmMethodId(methodName string) []byte {
	fixedBytes := vm.OvmMethodId(methodName)
	return fixedBytes[:]
}
