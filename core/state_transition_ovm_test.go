package core

import (
	"testing"
	"math/big"
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

var addresses = []common.Address{
	common.HexToAddress("0x5DF4785003ee8A7D2Ed8A0f4A5A49E5399c6d5c6"),
	common.HexToAddress("0xE23E5067F37d5c959F6cC712A7B24823C9a1a7c3"),
	common.HexToAddress("0x84341Ad8Eedb81FdC8EBffa935c3B93dD578Ea5E"),
}

func TestNeedsEOACreate(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := types.NewEIP155Signer(big.NewInt(18))
	tx, err := types.SignTx(types.NewTransaction(0, addr, new(big.Int), 0, new(big.Int), nil, &addresses[0], nil, types.QueueOriginSequencer, types.SighashEIP155), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	needseoa, err := NeedsEOACreate(tx, signer, statedb)
	if needseoa != true {
		t.Errorf("Expected transaction to need EOACreate, but did not.")
	}

	code := vm.OVMStateDump.Accounts["mockOVM_ECDSAContractAccount"].Code
	statedb.SetCode(addr, common.FromHex(code))
	
	needseoa, err = NeedsEOACreate(tx, signer, statedb)
	if needseoa != false {
		t.Errorf("Expected transaction to not need EOACreate, but did.")
	}
}

func TestToOvmMessage(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	var nonce = uint64(1234)
	var gasLimit = uint64(4321)
	var data = common.FromHex("0x56785678567856785678567856785678")

	signer := types.NewEIP155Signer(big.NewInt(420))
	tx, err := types.SignTx(types.NewTransaction(nonce, addresses[0], new(big.Int), gasLimit, new(big.Int), data, &addr, nil, types.QueueOriginSequencer, types.SighashEIP155), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := ToOvmMessage(tx, signer, false)
	if err != nil {
		t.Fatal(err)
	}

	var decoded = make(map[string]interface{})
	err = vm.OVMExecutionManager.ABI.Methods["run"].Inputs.UnpackIntoMap(decoded, msg.Data()[4:])
	if err != nil {
		t.Fatal(err)
	}

	var runtx = decoded["_transaction"].(struct { Timestamp *big.Int "json:\"timestamp\""; BlockNumber *big.Int "json:\"blockNumber\""; L1QueueOrigin uint8 "json:\"l1QueueOrigin\""; L1TxOrigin common.Address "json:\"l1TxOrigin\""; Entrypoint common.Address "json:\"entrypoint\""; GasLimit *big.Int "json:\"gasLimit\""; Data []uint8 "json:\"data\"" })
	
	if runtx.L1QueueOrigin != 2 {
		t.Errorf("Expected L1QueueOrigin to be 2, but got: %d\n", runtx.L1QueueOrigin)
	}

	if runtx.L1TxOrigin != addr {
		t.Errorf("Expected L1QueueOrigin to be %s, but got: %s\n", addr.Hex(), runtx.L1TxOrigin.Hex())
	}

	if runtx.Entrypoint != vm.OVMSequencerMessageDecompressor.Address {
		t.Errorf("Expected Entrypoint to be %s, but got: %s\n", vm.OVMSequencerMessageDecompressor.Address, runtx.Entrypoint.Hex())
	}

	if runtx.GasLimit.Uint64() != gasLimit {
		t.Errorf("Expected GasLimit to be %d, but got: %d\n", gasLimit, runtx.GasLimit.Uint64())
	}

	var buf = new(bytes.Buffer)
	buf.Write(runtx.Data)
	
	var actual = buf.Next(1)
	var expected = common.FromHex("0x01")
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected sig type to be %x, but got: %x\n", expected, actual)
	}
}
