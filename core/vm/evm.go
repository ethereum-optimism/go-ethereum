// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

// run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func run(evm *EVM, contract *Contract, input []byte, readOnly bool) ([]byte, error) {
	if UsingOVM {
		// OVM_ENABLED

		// Some simple logging here. First, check to see if we know about the address we're
		// interacting with and try to log the input data.
		var isUnknown = true
		for name, account := range OvmStateDump.Accounts {
			if contract.Address() == account.Address {
				isUnknown = false
				abi := &(account.ABI)
				method, err := abi.MethodById(input)
				if err != nil {
					log.Debug("Calling Known Contract", "Name", name, "ERROR", err)
				} else {
					log.Debug("Calling Known Contract", "Name", name, "Method", method.RawName)
				}
			}
		}

		// We don't know the contract, so print some generic information.
		if isUnknown {
			log.Debug("Calling Unknown Contract", "Address", contract.Address().Hex())
		}

		// Temporary: Safety checker always returns true.
		// ovmTODO: Remove this.
		if contract.Address() == OvmStateDump.Accounts["OVM_SafetyChecker"].Address {
			return AbiBytesTrue, nil
		}

		// If we're calling the state manager, we want to use our native implementation instead.
		if contract.Address() == OvmStateManager.Address {
			return callStateManager(input, evm, contract)
		}
	}

	if contract.CodeAddr != nil {
		precompiles := PrecompiledContractsHomestead
		if evm.chainRules.IsByzantium {
			precompiles = PrecompiledContractsByzantium
		}
		if evm.chainRules.IsIstanbul {
			precompiles = PrecompiledContractsIstanbul
		}
		if p := precompiles[*contract.CodeAddr]; p != nil {
			return RunPrecompiledContract(p, input, contract)
		}
	}

	for _, interpreter := range evm.interpreters {
		if interpreter.CanRun(contract.Code) {
			if evm.interpreter != interpreter {
				// Ensure that the interpreter pointer is set back
				// to its current value upon return.
				defer func(i Interpreter) {
					evm.interpreter = i
				}(evm.interpreter)
				evm.interpreter = interpreter
			}
			return interpreter.Run(contract, input, readOnly)
		}
	}

	return nil, ErrNoCompatibleInterpreter
}

// Context provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY

	// OVM_ADDITION
	EthCallSender         *common.Address
	OriginalTargetAddress *common.Address
	OriginalTargetResult  []byte
	OriginalTargetReached bool
}

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	vmConfig Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreters []Interpreter
	interpreter  Interpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(ctx Context, statedb StateDB, chainConfig *params.ChainConfig, vmConfig Config) *EVM {
	evm := &EVM{
		Context:      ctx,
		StateDB:      statedb,
		vmConfig:     vmConfig,
		chainConfig:  chainConfig,
		chainRules:   chainConfig.Rules(ctx.BlockNumber),
		interpreters: make([]Interpreter, 0, 1),
	}

	if chainConfig.IsEWASM(ctx.BlockNumber) {
		// to be implemented by EVM-C and Wagon PRs.
		// if vmConfig.EWASMInterpreter != "" {
		//  extIntOpts := strings.Split(vmConfig.EWASMInterpreter, ":")
		//  path := extIntOpts[0]
		//  options := []string{}
		//  if len(extIntOpts) > 1 {
		//    options = extIntOpts[1..]
		//  }
		//  evm.interpreters = append(evm.interpreters, NewEVMVCInterpreter(evm, vmConfig, options))
		// } else {
		// 	evm.interpreters = append(evm.interpreters, NewEWASMInterpreter(evm, vmConfig))
		// }
		panic("No supported ewasm interpreter yet.")
	}

	// vmConfig.EVMInterpreter will be used by EVM-C, it won't be checked here
	// as we always want to have the built-in EVM as the failover option.
	evm.interpreters = append(evm.interpreters, NewEVMInterpreter(evm, vmConfig))
	evm.interpreter = evm.interpreters[0]

	return evm
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// Cancelled returns true if Cancel has been called
func (evm *EVM) Cancelled() bool {
	return atomic.LoadInt32(&evm.abort) == 1
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() Interpreter {
	return evm.interpreter
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	var isTarget = false
	if UsingOVM {
		// OVM_ENABLED
		if evm.depth == 0 {
			// We're inside a new transaction, so make sure to wipe these variables beforehand.
			evm.Context.OriginalTargetAddress = nil
			evm.Context.OriginalTargetResult = []byte("00")
			evm.Context.OriginalTargetReached = false
		}

		if caller.Address() == OvmExecutionManager.Address &&
			!strings.HasPrefix(strings.ToLower(addr.Hex()), "0xdeaddeaddeaddeaddeaddeaddeaddeaddead") &&
			!strings.HasPrefix(strings.ToLower(addr.Hex()), "0x000000000000000000000000000000000000") &&
			!strings.HasPrefix(strings.ToLower(addr.Hex()), "0x420000000000000000000000000000000000") &&
			evm.Context.OriginalTargetAddress == nil {
			// Whew. Okay, so: we consider ourselves to be at a "target" as long as we were called
			// by the execution manager, and we're not a precompile or "dead" address.
			evm.Context.OriginalTargetAddress = &addr
			evm.Context.OriginalTargetReached = true
			isTarget = true
		}
	}

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	if !UsingOVM {
		// OVM_DISABLED
		// Fail if we're trying to transfer more than the available balance
		if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			return nil, gas, ErrInsufficientBalance
		}
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	if !evm.StateDB.Exist(addr) {
		precompiles := PrecompiledContractsHomestead
		if evm.chainRules.IsByzantium {
			precompiles = PrecompiledContractsByzantium
		}
		if evm.chainRules.IsIstanbul {
			precompiles = PrecompiledContractsIstanbul
		}
		if precompiles[addr] == nil && evm.chainRules.IsEIP158 && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}
		evm.StateDB.CreateAccount(addr)
	}

	if !UsingOVM {
		// OVM_DISABLED
		evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)
	}

	var prevCode []byte
	if UsingOVM {
		// OVM_ENABLED
		if evm.Context.EthCallSender != nil && *evm.Context.EthCallSender == addr {
			// We have to handle eth_call in a special manner as it doesn't stem from a signed
			// transaction and therefore can't go through the standard EOA contract. When we detect
			// this case, we temporarily insert some mock code that allows the user to pass through
			// the EOA without a signature. EthCallSender should *never* be made non-nil outside
			// of an eth_call. We store the old code here so that we're able to reset the code to
			// the original code for the remainder of the transaction.
			prevCode = evm.StateDB.GetCode(addr)
			evm.StateDB.SetCode(addr, common.FromHex(OvmStateDump.Accounts["mockOVM_ECDSAContractAccount"].Code))
		}
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	if UsingOVM {
		// OVM_ENABLED
		if evm.Context.EthCallSender != nil && *evm.Context.EthCallSender == addr {
			// Reset the code once it's been loaded into the contract object.
			evm.StateDB.SetCode(addr, prevCode)
		}
	}

	// Even if the account has no code, we need to continue because it might be a precompile
	start := time.Now()

	// Capture the tracer start/end events in debug mode
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

		defer func() { // Lazy evaluation of the parameters
			evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
		}()
	}

	ret, err = run(evm, contract, input, false)

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	if UsingOVM {
		// OVM_ENABLED

		if isTarget {
			// If this was our target contract, store the result so that it can be later re-inserted
			// into the user-facing return data (as seen below).
			evm.Context.OriginalTargetResult = ret
		}

		if evm.depth == 0 {
			// We're back at the root-level message call, so we'll need to modify the return data
			// sent to us by the OVM_ExecutionManager to instead be the intended return data.

			if !evm.Context.OriginalTargetReached {
				// If we didn't get to the target contract, then our execution somehow failed
				// (perhaps due to insufficient gas). Just return an error that represents this.
				ret = common.FromHex("0x")
				err = ErrOvmExecutionFailed
			} else if len(evm.Context.OriginalTargetResult) >= 96 {
				// We expect that EOA contracts return at least 96 bytes of data, where the first
				// 32 bytes are the boolean success value and the next 64 bytes are unnecessary
				// ABI encoding data. The actual return data starts at the 96th byte and can be
				// empty.
				success := evm.Context.OriginalTargetResult[:32]
				ret = evm.Context.OriginalTargetResult[96:]

				if !bytes.Equal(success, AbiBytesTrue) && !bytes.Equal(success, AbiBytesFalse) {
					// If the first 32 bytes not either are the ABI encoding of "true" or "false",
					// then the user hasn't correctly ABI encoded the result. We return the null
					// hex string as a default here (an annoying default that would convince most
					// people to just use the standard form).
					ret = common.FromHex("0x")
				} else if bytes.Equal(success, AbiBytesFalse) {
					// If the first 32 bytes are the ABI encoding of "false", then we need to add an
					// artificial error that represents the revert.
					err = errExecutionReverted

					// We also currently need to add an extra four empty bytes to the return data
					// to appease ethers.js. Our return correctly inserts the four specific bytes
					// that represent a "string error" to clients, but somehow the returndata size
					// is a multiple of 32 (when we expect size % 32 == 4). ethers.js checks that
					// [size % 32 == 4] before trying to decode a string error result. Adding these
					// four empty bytes tricks ethers into correctly decoding the error string.
					// ovmTODO: Figure out how to actually deal with this.
					// ovmTODO: This may actually be completely broken if the first four bytes of
					// the return data are **not** the specific "string error" bytes.
					ret = append(ret, make([]byte, 4)...)
				}
			} else {
				// User hasn't conformed the standard format, just return "null" for the success
				// (with no return data) to convince them to use the standard.
				ret = common.FromHex("0x")
			}

			log.Debug("Reached the end of an OVM execution", "Return Data", hexutil.Encode(ret), "Error", err)
		}
	}

	return ret, contract.Gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	if !UsingOVM {
		// OVM_DISABLED
		// Fail if we're trying to transfer more than the available balance
		if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
			return nil, gas, ErrInsufficientBalance
		}
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, contract, input, false)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, contract, input, false)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	// We do an AddBalance of zero here, just in order to trigger a touch.
	// This doesn't matter on Mainnet, where all empties are gone at the time of Byzantium,
	// but is the correct thing to do and matters on other networks, in tests, and potential
	// future scenarios
	evm.StateDB.AddBalance(addr, bigZero)

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in Homestead this also counts for code storage gas errors.
	ret, err = run(evm, contract, input, true)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

type codeAndHash struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash) Hash() common.Hash {
	if c.hash == (common.Hash{}) {
		c.hash = crypto.Keccak256Hash(c.code)
	}
	return c.hash
}

// create creates a new contract using code as deployment code.
func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth
	}
	if !UsingOVM {
		// OVM_DISABLED
		if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
			return nil, common.Address{}, gas, ErrInsufficientBalance
		}
	}

	// Ensure there's no existing contract already at the designated address
	contractHash := evm.StateDB.GetCodeHash(address)
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(address)
	if !UsingOVM {
		// OVM_DISABLED
		evm.Transfer(evm.StateDB, caller.Address(), address, value)
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, codeAndHash.code, gas, value)
	}
	start := time.Now()

	ret, err := run(evm, contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := evm.chainRules.IsEIP158 && len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas)) {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err

}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	if !UsingOVM {
		// OVM_DISABLED
		contractAddr = crypto.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	} else {
		// OVM_ENABLED
		if caller.Address() != OvmExecutionManager.Address {
			log.Error("Creation called by non-Execution Manager contract! This should never happen.", "Offending address", caller.Address().Hex())
			return nil, caller.Address(), 0, ErrOvmCreationFailed
		}

		// ovmADDRESS will be set by the execution manager to the target address whenever it's
		// about to create a new contract. This value is currently stored at the [15] storage slot.
		// Can pull this specific storage slot to get the address that the execution manager is
		// trying to create to, and create to it.
		slot := common.HexToHash(strconv.FormatInt(15, 16))
		contractAddr = common.BytesToAddress(evm.StateDB.GetState(OvmExecutionManager.Address, slot).Bytes())

		log.Debug("[EM] Creating contract.", "New contract address", contractAddr.Hex(), "Caller Addr", caller.Address().Hex(), "Caller nonce", evm.StateDB.GetNonce(caller.Address()))
	}

	return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr)
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	codeAndHash := &codeAndHash{code: code}
	if !UsingOVM {
		// OVM_DISABLED
		contractAddr = crypto.CreateAddress2(caller.Address(), common.BigToHash(salt), codeAndHash.Hash().Bytes())
	} else {
		// OVM_ENABLED
		if caller.Address() != OvmExecutionManager.Address {
			log.Error("Creation called by non-Execution Manager contract! This should never happen.", "Offending address", caller.Address().Hex())
			return nil, caller.Address(), 0, ErrOvmCreationFailed
		}

		// Same logic here as in Create, as seen above.
		slot := common.HexToHash(strconv.FormatInt(15, 16))
		contractAddr = common.BytesToAddress(evm.StateDB.GetState(OvmExecutionManager.Address, slot).Bytes())

		log.Debug("[EM] Creating contract [create2].", "New contract address", contractAddr.Hex(), "Caller Addr", caller.Address().Hex(), "Caller nonce", evm.StateDB.GetNonce(caller.Address()))
	}

	return evm.create(caller, codeAndHash, gas, endowment, contractAddr)
}

// ChainConfig returns the environment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig { return evm.chainConfig }
