// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package addressmanager

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

// LibAddressManagerABI is the input ABI used to generate the binding from.
const LibAddressManagerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"getAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// LibAddressManagerBin is the compiled bytecode used for deploying new contracts.
var LibAddressManagerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a361056a806100696000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063715018a61461005c5780638da5cb5b146100665780639b2ea4bd1461008a578063bf40fac11461013b578063f2fde38b146101e1575b600080fd5b610064610207565b005b61006e6102b0565b604080516001600160a01b039092168252519081900360200190f35b610064600480360360408110156100a057600080fd5b8101906020810181356401000000008111156100bb57600080fd5b8201836020820111156100cd57600080fd5b803590602001918460018302840111640100000000831117156100ef57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550505090356001600160a01b031691506102bf9050565b61006e6004803603602081101561015157600080fd5b81019060208101813564010000000081111561016c57600080fd5b82018360208201111561017e57600080fd5b803590602001918460018302840111640100000000831117156101a057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610362945050505050565b610064600480360360208110156101f757600080fd5b50356001600160a01b0316610391565b6000546001600160a01b03163314610266576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031681565b6000546001600160a01b0316331461031e576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b806001600061032c85610490565b815260200190815260200160002060006101000a8154816001600160a01b0302191690836001600160a01b031602179055505050565b60006001600061037184610490565b81526020810191909152604001600020546001600160a01b031692915050565b6000546001600160a01b031633146103f0576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6001600160a01b0381166104355760405162461bcd60e51b815260040180806020018281038252602d815260200180610508602d913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6000816040516020018082805190602001908083835b602083106104c55780518252601f1990920191602091820191016104a6565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405160208183030381529060405280519060200120905091905056fe4f776e61626c653a206e6577206f776e65722063616e6e6f7420626520746865207a65726f2061646472657373a26469706673582212204367ffc2e6671623708150e2d0cff4c12cf566722a26b4748555d789953e2d2264736f6c63430007000033"

// DeployLibAddressManager deploys a new Ethereum contract, binding an instance of LibAddressManager to it.
func DeployLibAddressManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LibAddressManager, error) {
	parsed, err := abi.JSON(strings.NewReader(LibAddressManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LibAddressManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LibAddressManager{LibAddressManagerCaller: LibAddressManagerCaller{contract: contract}, LibAddressManagerTransactor: LibAddressManagerTransactor{contract: contract}, LibAddressManagerFilterer: LibAddressManagerFilterer{contract: contract}}, nil
}

// LibAddressManager is an auto generated Go binding around an Ethereum contract.
type LibAddressManager struct {
	LibAddressManagerCaller     // Read-only binding to the contract
	LibAddressManagerTransactor // Write-only binding to the contract
	LibAddressManagerFilterer   // Log filterer for contract events
}

// LibAddressManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type LibAddressManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LibAddressManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LibAddressManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LibAddressManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LibAddressManagerSession struct {
	Contract     *LibAddressManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LibAddressManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LibAddressManagerCallerSession struct {
	Contract *LibAddressManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// LibAddressManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LibAddressManagerTransactorSession struct {
	Contract     *LibAddressManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// LibAddressManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type LibAddressManagerRaw struct {
	Contract *LibAddressManager // Generic contract binding to access the raw methods on
}

// LibAddressManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LibAddressManagerCallerRaw struct {
	Contract *LibAddressManagerCaller // Generic read-only contract binding to access the raw methods on
}

// LibAddressManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LibAddressManagerTransactorRaw struct {
	Contract *LibAddressManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLibAddressManager creates a new instance of LibAddressManager, bound to a specific deployed contract.
func NewLibAddressManager(address common.Address, backend bind.ContractBackend) (*LibAddressManager, error) {
	contract, err := bindLibAddressManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LibAddressManager{LibAddressManagerCaller: LibAddressManagerCaller{contract: contract}, LibAddressManagerTransactor: LibAddressManagerTransactor{contract: contract}, LibAddressManagerFilterer: LibAddressManagerFilterer{contract: contract}}, nil
}

// NewLibAddressManagerCaller creates a new read-only instance of LibAddressManager, bound to a specific deployed contract.
func NewLibAddressManagerCaller(address common.Address, caller bind.ContractCaller) (*LibAddressManagerCaller, error) {
	contract, err := bindLibAddressManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LibAddressManagerCaller{contract: contract}, nil
}

// NewLibAddressManagerTransactor creates a new write-only instance of LibAddressManager, bound to a specific deployed contract.
func NewLibAddressManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*LibAddressManagerTransactor, error) {
	contract, err := bindLibAddressManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LibAddressManagerTransactor{contract: contract}, nil
}

// NewLibAddressManagerFilterer creates a new log filterer instance of LibAddressManager, bound to a specific deployed contract.
func NewLibAddressManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*LibAddressManagerFilterer, error) {
	contract, err := bindLibAddressManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LibAddressManagerFilterer{contract: contract}, nil
}

// bindLibAddressManager binds a generic wrapper to an already deployed contract.
func bindLibAddressManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LibAddressManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibAddressManager *LibAddressManagerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibAddressManager.Contract.LibAddressManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibAddressManager *LibAddressManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibAddressManager.Contract.LibAddressManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibAddressManager *LibAddressManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibAddressManager.Contract.LibAddressManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LibAddressManager *LibAddressManagerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LibAddressManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LibAddressManager *LibAddressManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibAddressManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LibAddressManager *LibAddressManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LibAddressManager.Contract.contract.Transact(opts, method, params...)
}

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string _name) constant returns(address)
func (_LibAddressManager *LibAddressManagerCaller) GetAddress(opts *bind.CallOpts, _name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LibAddressManager.contract.Call(opts, out, "getAddress", _name)
	return *ret0, err
}

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string _name) constant returns(address)
func (_LibAddressManager *LibAddressManagerSession) GetAddress(_name string) (common.Address, error) {
	return _LibAddressManager.Contract.GetAddress(&_LibAddressManager.CallOpts, _name)
}

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string _name) constant returns(address)
func (_LibAddressManager *LibAddressManagerCallerSession) GetAddress(_name string) (common.Address, error) {
	return _LibAddressManager.Contract.GetAddress(&_LibAddressManager.CallOpts, _name)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LibAddressManager *LibAddressManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LibAddressManager.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LibAddressManager *LibAddressManagerSession) Owner() (common.Address, error) {
	return _LibAddressManager.Contract.Owner(&_LibAddressManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_LibAddressManager *LibAddressManagerCallerSession) Owner() (common.Address, error) {
	return _LibAddressManager.Contract.Owner(&_LibAddressManager.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LibAddressManager *LibAddressManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LibAddressManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LibAddressManager *LibAddressManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _LibAddressManager.Contract.RenounceOwnership(&_LibAddressManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LibAddressManager *LibAddressManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LibAddressManager.Contract.RenounceOwnership(&_LibAddressManager.TransactOpts)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string _name, address _address) returns()
func (_LibAddressManager *LibAddressManagerTransactor) SetAddress(opts *bind.TransactOpts, _name string, _address common.Address) (*types.Transaction, error) {
	return _LibAddressManager.contract.Transact(opts, "setAddress", _name, _address)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string _name, address _address) returns()
func (_LibAddressManager *LibAddressManagerSession) SetAddress(_name string, _address common.Address) (*types.Transaction, error) {
	return _LibAddressManager.Contract.SetAddress(&_LibAddressManager.TransactOpts, _name, _address)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string _name, address _address) returns()
func (_LibAddressManager *LibAddressManagerTransactorSession) SetAddress(_name string, _address common.Address) (*types.Transaction, error) {
	return _LibAddressManager.Contract.SetAddress(&_LibAddressManager.TransactOpts, _name, _address)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_LibAddressManager *LibAddressManagerTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _LibAddressManager.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_LibAddressManager *LibAddressManagerSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _LibAddressManager.Contract.TransferOwnership(&_LibAddressManager.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_LibAddressManager *LibAddressManagerTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _LibAddressManager.Contract.TransferOwnership(&_LibAddressManager.TransactOpts, _newOwner)
}

// LibAddressManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LibAddressManager contract.
type LibAddressManagerOwnershipTransferredIterator struct {
	Event *LibAddressManagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LibAddressManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LibAddressManagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LibAddressManagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LibAddressManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LibAddressManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LibAddressManagerOwnershipTransferred represents a OwnershipTransferred event raised by the LibAddressManager contract.
type LibAddressManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LibAddressManager *LibAddressManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LibAddressManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LibAddressManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LibAddressManagerOwnershipTransferredIterator{contract: _LibAddressManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LibAddressManager *LibAddressManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LibAddressManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LibAddressManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LibAddressManagerOwnershipTransferred)
				if err := _LibAddressManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LibAddressManager *LibAddressManagerFilterer) ParseOwnershipTransferred(log types.Log) (*LibAddressManagerOwnershipTransferred, error) {
	event := new(LibAddressManagerOwnershipTransferred)
	if err := _LibAddressManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

