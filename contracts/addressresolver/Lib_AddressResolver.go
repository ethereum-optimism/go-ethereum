// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package addressresolver

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LibAddressResolverABI is the input ABI used to generate the binding from.
const LibAddressResolverABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_libAddressManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LibAddressResolverBin is the compiled bytecode used for deploying new contracts.
var LibAddressResolverBin = "0x608060405234801561001057600080fd5b506040516102693803806102698339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b0319909216919091179055610204806100656000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063461a447814610030575b600080fd5b6100d66004803603602081101561004657600080fd5b81019060208101813564010000000081111561006157600080fd5b82018360208201111561007357600080fd5b8035906020019184600183028401116401000000008311171561009557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506100f2945050505050565b604080516001600160a01b039092168252519081900360200190f35b6000805460405163bf40fac160e01b81526020600482018181528551602484015285516001600160a01b039094169363bf40fac19387938392604490920191908501908083838b5b8381101561015257818101518382015260200161013a565b50505050905090810190601f16801561017f5780820380516001836020036101000a031916815260200191505b509250505060206040518083038186803b15801561019c57600080fd5b505afa1580156101b0573d6000803e3d6000fd5b505050506040513d60208110156101c657600080fd5b50519291505056fea26469706673582212202d65863aab3b960819b4723b506367e35ce09bc37fd925a075175883011ce4d864736f6c63430007000033"

// DeployLibAddressResolver deploys a new Ethereum contract, binding an instance of LibAddressResolver to it.
func DeployLibAddressResolver(auth *bind.TransactOpts, backend bind.ContractBackend, _libAddressManager common.Address) (common.Address, *types.Transaction, *LibAddressResolver, error) {
	parsed, err := abi.JSON(strings.NewReader(LibAddressResolverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LibAddressResolverBin), backend, _libAddressManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LibAddressResolver{LibAddressResolverCaller: LibAddressResolverCaller{contract: contract}, LibAddressResolverTransactor: LibAddressResolverTransactor{contract: contract}, LibAddressResolverFilterer: LibAddressResolverFilterer{contract: contract}}, nil
}

// LibAddressResolver is an auto generated Go binding around an Ethereum contract.
type LibAddressResolver struct {
	LibAddressResolverCaller     // Read-only binding to the contract
	LibAddressResolverTransactor // Write-only binding to the contract
	LibAddressResolverFilterer   // Log filterer for contract events
}

// LibAddressResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type LibAddressResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LibAddressResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LibAddressResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LibAddressResolverSession struct {
	Contract     *LibAddressResolver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LibAddressResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LibAddressResolverCallerSession struct {
	Contract *LibAddressResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LibAddressResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LibAddressResolverTransactorSession struct {
	Contract     *LibAddressResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LibAddressResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type LibAddressResolverRaw struct {
	Contract *LibAddressResolver // Generic contract binding to access the raw methods on
}

// LibAddressResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LibAddressResolverCallerRaw struct {
	Contract *LibAddressResolverCaller // Generic read-only contract binding to access the raw methods on
}

// LibAddressResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LibAddressResolverTransactorRaw struct {
	Contract *LibAddressResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLibAddressResolver creates a new instance of LibAddressResolver, bound to a specific deployed contract.
func NewLibAddressResolver(address common.Address, backend bind.ContractBackend) (*LibAddressResolver, error) {
	contract, err := bindLibAddressResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LibAddressResolver{LibAddressResolverCaller: LibAddressResolverCaller{contract: contract}, LibAddressResolverTransactor: LibAddressResolverTransactor{contract: contract}, LibAddressResolverFilterer: LibAddressResolverFilterer{contract: contract}}, nil
}

// NewLibAddressResolverCaller creates a new read-only instance of LibAddressResolver, bound to a specific deployed contract.
func NewLibAddressResolverCaller(address common.Address, caller bind.ContractCaller) (*LibAddressResolverCaller, error) {
	contract, err := bindLibAddressResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LibAddressResolverCaller{contract: contract}, nil
}

// NewLibAddressResolverTransactor creates a new write-only instance of LibAddressResolver, bound to a specific deployed contract.
func NewLibAddressResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*LibAddressResolverTransactor, error) {
	contract, err := bindLibAddressResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LibAddressResolverTransactor{contract: contract}, nil
}

// NewLibAddressResolverFilterer creates a new log filterer instance of LibAddressResolver, bound to a specific deployed contract.
func NewLibAddressResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*LibAddressResolverFilterer, error) {
	contract, err := bindLibAddressResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LibAddressResolverFilterer{contract: contract}, nil
}

// bindLibAddressResolver binds a generic wrapper to an already deployed contract.
func bindLibAddressResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LibAddressResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibAddressResolver *LibAddressResolverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibAddressResolver.Contract.LibAddressResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibAddressResolver *LibAddressResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibAddressResolver.Contract.LibAddressResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibAddressResolver *LibAddressResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibAddressResolver.Contract.LibAddressResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibAddressResolver *LibAddressResolverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibAddressResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibAddressResolver *LibAddressResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibAddressResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibAddressResolver *LibAddressResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibAddressResolver.Contract.contract.Transact(opts, method, params...)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_LibAddressResolver *LibAddressResolverCaller) Resolve(opts *bind.CallOpts, _name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LibAddressResolver.contract.Call(opts, out, "resolve", _name)
	return *ret0, err
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_LibAddressResolver *LibAddressResolverSession) Resolve(_name string) (common.Address, error) {
	return _LibAddressResolver.Contract.Resolve(&_LibAddressResolver.CallOpts, _name)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_LibAddressResolver *LibAddressResolverCallerSession) Resolve(_name string) (common.Address, error) {
	return _LibAddressResolver.Contract.Resolve(&_LibAddressResolver.CallOpts, _name)
}
