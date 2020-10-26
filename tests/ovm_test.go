package tests

import (
	"fmt"
	"testing"
	"math/big"
	"bytes"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

var nulladdress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var nullhash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
var emaddress = common.HexToAddress("0xdeaddeaddeaddeaddeaddeaddeaddeaddead000b")

var addresses []common.Address
var keys []*ecdsa.PrivateKey

var hashes = []common.Hash{
	common.HexToHash("0x1020304050600000000000000000000000000000000000000000000000000000"),
	common.HexToHash("0x6050403020100000000000000000000000000000000000000000000000000000"),
	common.HexToHash("0x6050403020101020304050600000000000000000000000000000000000000000"),
}

type ovmResult struct {
	Success bool		`abi:"_success"`
	Data	[]byte		`abi:"_returndata"`
}

type ovmSimpleContractTest struct {
	fn		string
	in		[]interface{}
	out		[]interface{}
	success bool
}

func init() {
	key, _ := crypto.GenerateKey()
	addresses = append(addresses, crypto.PubkeyToAddress(key.PublicKey))
	keys = append(keys, key)
}

func TestContextFunctions(t *testing.T) {
	var ovmContextFunctionTests = []ovmSimpleContractTest{
		{
			fn: "ovmCALLER",
			in: []interface{}{},
			out: []interface{}{addresses[0]},
			success: true,
		},
		{
			fn: "ovmADDRESS",
			in: []interface{}{},
			out: []interface{}{crypto.CreateAddress(addresses[0], 0)},
			success: true,
		},
		{
			fn: "ovmTIMESTAMP",
			in: []interface{}{},
			out: []interface{}{big.NewInt(0)},
			success: true,
		},
		{
			fn: "ovmNUMBER",
			in: []interface{}{},
			out: []interface{}{big.NewInt(0)},
			success: true,
		},
		{
			fn: "ovmSLOAD",
			in: []interface{}{hashes[0]},
			out: []interface{}{nullhash},
			success: true,
		},
		{
			fn: "ovmSSTORE",
			in: []interface{}{hashes[0], hashes[1]},
			out: []interface{}{},
			success: true,
		},
		{
			fn: "ovmSLOAD",
			in: []interface{}{hashes[0]},
			out: []interface{}{hashes[1]},
			success: true,
		},
		{
			fn: "doReturn",
			in: []interface{}{},
			out: []interface{}{common.FromHex("0x420adfadf1234789098484848069")},
			success: true,
		},
		{
			fn: "doCallToReturn",
			in: []interface{}{},
			out: []interface{}{common.FromHex("0x420adfadf1234789098484848069")},
			success: true,
		},
		{
			fn: "doCREATE",
			in: []interface{}{big.NewInt(1234)},
			out: []interface{}{},
			success: true,
		},
		{
			fn: "doRevert",
			in: []interface{}{},
			out: []interface{}{[]byte("this is a revert message")},
			success: false,
		},
	}

	runSimpleContractTest(t, ovmContextFunctionTests)
}

func runSimpleContractTest(t *testing.T, tests []ovmSimpleContractTest) {
	ovm, statedb := initOVM()
	
	signer := types.NewEIP155Signer(big.NewInt(420))

	// Deploy the contract first.
	runOVMTransaction(ovm, addresses[0], nulladdress, common.FromHex(OVMSimpleContract.Bytecode), signer, keys[0], statedb)
	addr := crypto.CreateAddress(addresses[0], 0)

	fmt.Printf("original: %s\n", addresses[0].Hex())
	fmt.Printf("created: %s\n", addr.Hex())

	// Now run each of the test parameters.
	for _, tst := range tests {
		tx, err := OVMSimpleContract.ABI.Pack(tst.fn, tst.in...)
		if err != nil {
			panic(fmt.Errorf("Could not encode test input data: %v\n", err))
		}

		out := runOVMTransaction(ovm, addresses[0], addr, tx, signer, keys[0], statedb)
		var args ovmResult
		err = vm.OVMStateDump.Accounts["mockOVM_ECDSAContractAccount"].ABI.Unpack(&args, "execute", out)
		if err != nil {
			panic(fmt.Errorf("Could not unpack result: %v\n", err))
		}

		if args.Success != tst.success {
			t.Errorf("Expected success status does not match actual return status.")
		}

		ret := args.Data
		if !tst.success {
			ret = ret[4:]
		}

		expected, err := OVMSimpleContract.ABI.Methods[tst.fn].Outputs.Pack(tst.out...)
		if err != nil {
			t.Errorf("Could not pack output values for '%s': %v\n", tst.fn, err)
		}

		if !bytes.Equal(ret, expected) {
			t.Errorf("Expected %x; got %x\n", expected, ret)
		}
	}
}

func runOVMTransaction(evm *vm.EVM, from common.Address, to common.Address, tx []byte, signer types.Signer, key *ecdsa.PrivateKey, statedb *state.StateDB) []byte {
	st := makeStateTransition(evm, from, to, tx, signer, key, statedb)
	ret, _, failed, err := st.TransitionDb()

	if err != nil {
		panic(fmt.Errorf("Running state transition failed: %v\n", err))
	}

	if failed {
		panic(fmt.Errorf("Running state transition failed (revert error): %x\n", ret))
	}

	return ret
}

func initOVM() (*vm.EVM, *state.StateDB) {
	db := rawdb.NewMemoryDatabase()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(db))
	if err != nil {
		panic(fmt.Errorf("Unable to initialize OVM StateDB: %v", err))
	}

	core.ApplyOvmStateToState(statedb)
	return vm.NewEVM(vm.Context{
		BlockNumber: big.NewInt(1),
	}, statedb, params.TestChainConfig, vm.Config{}), statedb
}

func makeStateTransition(
	evm *vm.EVM,
	from common.Address,
	to common.Address,
	data []byte,
	signer types.Signer,
	key *ecdsa.PrivateKey,
	statedb *state.StateDB,
) *core.StateTransition {
	tx, _ := types.SignTx(types.NewTransaction(
		0,
		to,
		big.NewInt(0),
		10000000,
		big.NewInt(0),
		data,
		&from,
		nil,
		types.QueueOriginSequencer,
		0,
	), signer, key)

	needseoa, err := core.NeedsEOACreate(tx, signer, statedb)
	if err != nil {
		panic(fmt.Errorf("Unable to check if needs EOA: %v", err))
	}

	gaspool := new(core.GasPool).AddGas(50000000)
	evm.Context.Origin = core.GodAddress
	evm.Context.L1MessageSender = &from

	if needseoa {
		fmt.Printf("\n-------- GOING INTO EOA CREATE --------\n")
		eoamsg, err := core.ToOvmMessage(tx, signer, true)
		if err != nil {
			panic(fmt.Errorf("Unable to create EOACreate message: %v", err))
		}

		st := core.NewStateTransition(
			evm,
			eoamsg,
			gaspool,
		)

		st.TransitionDb()
	}

	if to == nulladdress {
		fmt.Printf("\n-------- GOING INTO CONTRACT CREATE --------\n")
	} else {
		fmt.Printf("\n-------- GOING INTO CONTRACT CALL --------\n")
	}

	msg, err := core.ToOvmMessage(tx, signer, false)
	if err != nil {
		panic(fmt.Errorf("Unable to create OVM message: %v", err))
	}

	return core.NewStateTransition(
		evm,
		msg,
		gaspool,
	)
}
