// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package portdelegate

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ContextABI is the input ABI used to generate the binding from.
const ContextABI = "[]"

// Context is an auto generated Go binding around an Ethereum contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableFuncSigs = map[string]string{
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"f2fde38b": "transferOwnership(address)",
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ownable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ownable *OwnableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableOwnershipTransferred, error) {
	event := new(OwnableOwnershipTransferred)
	if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PortMemorizerABI is the input ABI used to generate the binding from.
const PortMemorizerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nebula\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"RequestCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"drop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initializer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nebula\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rqId\",\"type\":\"uint256\"}],\"name\":\"nextRq\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"persist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rqId\",\"type\":\"uint256\"}],\"name\":\"prevRq\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestsQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"swapStatus\",\"outputs\":[{\"internalType\":\"enumPortMemorizer.RequestStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"unwrapRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"foreignAddress\",\"type\":\"bytes32\"},{\"internalType\":\"enumPortMemorizer.RequestStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PortMemorizerFuncSigs maps the 4-byte function signature to its string representation.
var PortMemorizerFuncSigs = map[string]string{
	"00b52395": "drop(bytes)",
	"9ce110d7": "initializer()",
	"4ecde849": "nebula()",
	"29db9a47": "nextRq(uint256)",
	"8da5cb5b": "owner()",
	"d3f5d56e": "persist(bytes)",
	"a6f3ade9": "prevRq(uint256)",
	"715018a6": "renounceOwnership()",
	"56dcda94": "requestsQueue()",
	"0872512b": "swapStatus(uint256)",
	"9d76ea58": "tokenAddress()",
	"f2fde38b": "transferOwnership(address)",
	"d99c2a72": "unwrapRequests(uint256)",
}

// PortMemorizerBin is the compiled bytecode used for deploying new contracts.
var PortMemorizerBin = "0x608060405234801561001057600080fd5b50604051610b87380380610b878339818101604052604081101561003357600080fd5b50805160209091015160006100466100cc565b600080546001600160a01b0319166001600160a01b0383169081178255604051929350917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350600980546001600160a01b03199081163317909155600180546001600160a01b03948516908316179055600280549290931691161790556100d0565b3390565b610aa8806100df6000396000f3fe608060405234801561001057600080fd5b50600436106100ce5760003560e01c80638da5cb5b1161008c578063a6f3ade911610066578063a6f3ade914610217578063d3f5d56e14610234578063d99c2a72146102a4578063f2fde38b146102f0576100ce565b80638da5cb5b146101ff5780639ce110d7146102075780639d76ea581461020f576100ce565b8062b52395146100d35780630872512b1461014557806329db9a47146101835780634ecde849146101b257806356dcda94146101d6578063715018a6146101f7575b600080fd5b610143600480360360208110156100e957600080fd5b81019060208101813564010000000081111561010457600080fd5b82018360208201111561011657600080fd5b8035906020019184600183028401116401000000008311171561013857600080fd5b509092509050610316565b005b6101626004803603602081101561015b57600080fd5b503561045d565b6040518082600281111561017257fe5b815260200191505060405180910390f35b6101a06004803603602081101561019957600080fd5b5035610472565b60408051918252519081900360200190f35b6101ba610484565b604080516001600160a01b039092168252519081900360200190f35b6101de610493565b6040805192835260208301919091528051918290030190f35b61014361049c565b6101ba61055a565b6101ba610569565b6101ba610578565b6101a06004803603602081101561022d57600080fd5b5035610587565b6101436004803603602081101561024a57600080fd5b81019060208101813564010000000081111561026557600080fd5b82018360208201111561027757600080fd5b8035906020019184600183028401116401000000008311171561029957600080fd5b509092509050610599565b6102c1600480360360208110156102ba57600080fd5b503561083d565b604051808481526020018381526020018260028111156102dd57fe5b8152602001935050505060405180910390f35b6101436004803603602081101561030657600080fd5b50356001600160a01b0316610861565b6009546001600160a01b03163314610365576040805162461bcd60e51b815260206004820152600d60248201526c1858d8d95cdcc819195b9a5959609a1b604482015290519081900360640190fd5b6000806103ac84848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250602091506109759050565b600081815260046020818152604080842084815560018101859055600201805460ff1990811690915560038352818520805490911690558051639d6ad84b60e01b8152600593810193909352602483018590525195019492935073__$8eef590145d27c17e68c918edbd96fcba5$__92639d6ad84b92604480840193919291829003018186803b15801561043f57600080fd5b505af4158015610453573d6000803e3d6000fd5b5050505050505050565b60036020526000908152604090205460ff1681565b60009081526007602052604090205490565b6001546001600160a01b031681565b60055460065482565b6104a46109b6565b6001600160a01b03166104b561055a565b6001600160a01b031614610510576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031690565b6009546001600160a01b031681565b6002546001600160a01b031681565b60009081526008602052604090205490565b6009546001600160a01b031633146105e8576040805162461bcd60e51b815260206004820152600d60248201526c1858d8d95cdcc819195b9a5959609a1b604482015290519081900360640190fd5b60008061062f84848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250602091506109759050565b9050602082019150600061067d85858080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250879250602091506109759050565b905060208301925060006106cb86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250889250602091506109759050565b60001b9050602084019350600061071987878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508992506109ba915050565b9050600185019450604051806060016040528084815260200183815260200182600281111561074457fe5b81525060046000868152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff0219169083600281111561079257fe5b021790555050506000848152600360205260409020805482919060ff191660018360028111156107be57fe5b021790555060408051632941b65560e21b81526005600482015260248101869052905173__$8eef590145d27c17e68c918edbd96fcba5$__9163a506d954916044808301926000929190829003018186803b15801561081c57600080fd5b505af4158015610830573d6000803e3d6000fd5b5050505050505050505050565b60046020526000908152604090208054600182015460029092015490919060ff1683565b6108696109b6565b6001600160a01b031661087a61055a565b6001600160a01b0316146108d5576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6001600160a01b03811661091a5760405162461bcd60e51b8152600401808060200182810382526026815260200180610a4d6026913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b600080835b8385018110156109ad5785818151811061099057fe5b60209101015160f81c61010092909202919091019060010161097a565b50949350505050565b3390565b6000808383815181106109c957fe5b016020015160f81c9050806109e2576000915050610a46565b80600114156109f5576001915050610a46565b8060021415610a08576002915050610a46565b6040805162461bcd60e51b815260206004820152600e60248201526d696e76616c69642073746174757360901b604482015290519081900360640190fd5b9291505056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373a2646970667358221220ab8ef646ad780a3d84b414897bf58ae70e6cd9899b37269dcc077470a9924ec564736f6c63430007050033"

// DeployPortMemorizer deploys a new Ethereum contract, binding an instance of PortMemorizer to it.
func DeployPortMemorizer(auth *bind.TransactOpts, backend bind.ContractBackend, _nebula common.Address, _tokenAddress common.Address) (common.Address, *types.Transaction, *PortMemorizer, error) {
	parsed, err := abi.JSON(strings.NewReader(PortMemorizerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	queueLibAddr, _, _, _ := DeployQueueLib(auth, backend)
	PortMemorizerBin = strings.Replace(PortMemorizerBin, "__$8eef590145d27c17e68c918edbd96fcba5$__", queueLibAddr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PortMemorizerBin), backend, _nebula, _tokenAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PortMemorizer{PortMemorizerCaller: PortMemorizerCaller{contract: contract}, PortMemorizerTransactor: PortMemorizerTransactor{contract: contract}, PortMemorizerFilterer: PortMemorizerFilterer{contract: contract}}, nil
}

// PortMemorizer is an auto generated Go binding around an Ethereum contract.
type PortMemorizer struct {
	PortMemorizerCaller     // Read-only binding to the contract
	PortMemorizerTransactor // Write-only binding to the contract
	PortMemorizerFilterer   // Log filterer for contract events
}

// PortMemorizerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PortMemorizerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PortMemorizerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PortMemorizerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PortMemorizerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PortMemorizerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PortMemorizerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PortMemorizerSession struct {
	Contract     *PortMemorizer    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PortMemorizerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PortMemorizerCallerSession struct {
	Contract *PortMemorizerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// PortMemorizerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PortMemorizerTransactorSession struct {
	Contract     *PortMemorizerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// PortMemorizerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PortMemorizerRaw struct {
	Contract *PortMemorizer // Generic contract binding to access the raw methods on
}

// PortMemorizerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PortMemorizerCallerRaw struct {
	Contract *PortMemorizerCaller // Generic read-only contract binding to access the raw methods on
}

// PortMemorizerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PortMemorizerTransactorRaw struct {
	Contract *PortMemorizerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPortMemorizer creates a new instance of PortMemorizer, bound to a specific deployed contract.
func NewPortMemorizer(address common.Address, backend bind.ContractBackend) (*PortMemorizer, error) {
	contract, err := bindPortMemorizer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PortMemorizer{PortMemorizerCaller: PortMemorizerCaller{contract: contract}, PortMemorizerTransactor: PortMemorizerTransactor{contract: contract}, PortMemorizerFilterer: PortMemorizerFilterer{contract: contract}}, nil
}

// NewPortMemorizerCaller creates a new read-only instance of PortMemorizer, bound to a specific deployed contract.
func NewPortMemorizerCaller(address common.Address, caller bind.ContractCaller) (*PortMemorizerCaller, error) {
	contract, err := bindPortMemorizer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PortMemorizerCaller{contract: contract}, nil
}

// NewPortMemorizerTransactor creates a new write-only instance of PortMemorizer, bound to a specific deployed contract.
func NewPortMemorizerTransactor(address common.Address, transactor bind.ContractTransactor) (*PortMemorizerTransactor, error) {
	contract, err := bindPortMemorizer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PortMemorizerTransactor{contract: contract}, nil
}

// NewPortMemorizerFilterer creates a new log filterer instance of PortMemorizer, bound to a specific deployed contract.
func NewPortMemorizerFilterer(address common.Address, filterer bind.ContractFilterer) (*PortMemorizerFilterer, error) {
	contract, err := bindPortMemorizer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PortMemorizerFilterer{contract: contract}, nil
}

// bindPortMemorizer binds a generic wrapper to an already deployed contract.
func bindPortMemorizer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PortMemorizerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PortMemorizer *PortMemorizerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PortMemorizer.Contract.PortMemorizerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PortMemorizer *PortMemorizerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PortMemorizer.Contract.PortMemorizerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PortMemorizer *PortMemorizerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PortMemorizer.Contract.PortMemorizerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PortMemorizer *PortMemorizerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PortMemorizer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PortMemorizer *PortMemorizerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PortMemorizer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PortMemorizer *PortMemorizerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PortMemorizer.Contract.contract.Transact(opts, method, params...)
}

// Initializer is a free data retrieval call binding the contract method 0x9ce110d7.
//
// Solidity: function initializer() view returns(address)
func (_PortMemorizer *PortMemorizerCaller) Initializer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "initializer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Initializer is a free data retrieval call binding the contract method 0x9ce110d7.
//
// Solidity: function initializer() view returns(address)
func (_PortMemorizer *PortMemorizerSession) Initializer() (common.Address, error) {
	return _PortMemorizer.Contract.Initializer(&_PortMemorizer.CallOpts)
}

// Initializer is a free data retrieval call binding the contract method 0x9ce110d7.
//
// Solidity: function initializer() view returns(address)
func (_PortMemorizer *PortMemorizerCallerSession) Initializer() (common.Address, error) {
	return _PortMemorizer.Contract.Initializer(&_PortMemorizer.CallOpts)
}

// Nebula is a free data retrieval call binding the contract method 0x4ecde849.
//
// Solidity: function nebula() view returns(address)
func (_PortMemorizer *PortMemorizerCaller) Nebula(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "nebula")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Nebula is a free data retrieval call binding the contract method 0x4ecde849.
//
// Solidity: function nebula() view returns(address)
func (_PortMemorizer *PortMemorizerSession) Nebula() (common.Address, error) {
	return _PortMemorizer.Contract.Nebula(&_PortMemorizer.CallOpts)
}

// Nebula is a free data retrieval call binding the contract method 0x4ecde849.
//
// Solidity: function nebula() view returns(address)
func (_PortMemorizer *PortMemorizerCallerSession) Nebula() (common.Address, error) {
	return _PortMemorizer.Contract.Nebula(&_PortMemorizer.CallOpts)
}

// NextRq is a free data retrieval call binding the contract method 0x29db9a47.
//
// Solidity: function nextRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerCaller) NextRq(opts *bind.CallOpts, rqId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "nextRq", rqId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRq is a free data retrieval call binding the contract method 0x29db9a47.
//
// Solidity: function nextRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerSession) NextRq(rqId *big.Int) (*big.Int, error) {
	return _PortMemorizer.Contract.NextRq(&_PortMemorizer.CallOpts, rqId)
}

// NextRq is a free data retrieval call binding the contract method 0x29db9a47.
//
// Solidity: function nextRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerCallerSession) NextRq(rqId *big.Int) (*big.Int, error) {
	return _PortMemorizer.Contract.NextRq(&_PortMemorizer.CallOpts, rqId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PortMemorizer *PortMemorizerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PortMemorizer *PortMemorizerSession) Owner() (common.Address, error) {
	return _PortMemorizer.Contract.Owner(&_PortMemorizer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PortMemorizer *PortMemorizerCallerSession) Owner() (common.Address, error) {
	return _PortMemorizer.Contract.Owner(&_PortMemorizer.CallOpts)
}

// PrevRq is a free data retrieval call binding the contract method 0xa6f3ade9.
//
// Solidity: function prevRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerCaller) PrevRq(opts *bind.CallOpts, rqId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "prevRq", rqId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PrevRq is a free data retrieval call binding the contract method 0xa6f3ade9.
//
// Solidity: function prevRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerSession) PrevRq(rqId *big.Int) (*big.Int, error) {
	return _PortMemorizer.Contract.PrevRq(&_PortMemorizer.CallOpts, rqId)
}

// PrevRq is a free data retrieval call binding the contract method 0xa6f3ade9.
//
// Solidity: function prevRq(uint256 rqId) view returns(uint256)
func (_PortMemorizer *PortMemorizerCallerSession) PrevRq(rqId *big.Int) (*big.Int, error) {
	return _PortMemorizer.Contract.PrevRq(&_PortMemorizer.CallOpts, rqId)
}

// RequestsQueue is a free data retrieval call binding the contract method 0x56dcda94.
//
// Solidity: function requestsQueue() view returns(bytes32 first, bytes32 last)
func (_PortMemorizer *PortMemorizerCaller) RequestsQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "requestsQueue")

	outstruct := new(struct {
		First [32]byte
		Last  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.First = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Last = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// RequestsQueue is a free data retrieval call binding the contract method 0x56dcda94.
//
// Solidity: function requestsQueue() view returns(bytes32 first, bytes32 last)
func (_PortMemorizer *PortMemorizerSession) RequestsQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _PortMemorizer.Contract.RequestsQueue(&_PortMemorizer.CallOpts)
}

// RequestsQueue is a free data retrieval call binding the contract method 0x56dcda94.
//
// Solidity: function requestsQueue() view returns(bytes32 first, bytes32 last)
func (_PortMemorizer *PortMemorizerCallerSession) RequestsQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _PortMemorizer.Contract.RequestsQueue(&_PortMemorizer.CallOpts)
}

// SwapStatus is a free data retrieval call binding the contract method 0x0872512b.
//
// Solidity: function swapStatus(uint256 ) view returns(uint8)
func (_PortMemorizer *PortMemorizerCaller) SwapStatus(opts *bind.CallOpts, arg0 *big.Int) (uint8, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "swapStatus", arg0)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SwapStatus is a free data retrieval call binding the contract method 0x0872512b.
//
// Solidity: function swapStatus(uint256 ) view returns(uint8)
func (_PortMemorizer *PortMemorizerSession) SwapStatus(arg0 *big.Int) (uint8, error) {
	return _PortMemorizer.Contract.SwapStatus(&_PortMemorizer.CallOpts, arg0)
}

// SwapStatus is a free data retrieval call binding the contract method 0x0872512b.
//
// Solidity: function swapStatus(uint256 ) view returns(uint8)
func (_PortMemorizer *PortMemorizerCallerSession) SwapStatus(arg0 *big.Int) (uint8, error) {
	return _PortMemorizer.Contract.SwapStatus(&_PortMemorizer.CallOpts, arg0)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_PortMemorizer *PortMemorizerCaller) TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "tokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_PortMemorizer *PortMemorizerSession) TokenAddress() (common.Address, error) {
	return _PortMemorizer.Contract.TokenAddress(&_PortMemorizer.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_PortMemorizer *PortMemorizerCallerSession) TokenAddress() (common.Address, error) {
	return _PortMemorizer.Contract.TokenAddress(&_PortMemorizer.CallOpts)
}

// UnwrapRequests is a free data retrieval call binding the contract method 0xd99c2a72.
//
// Solidity: function unwrapRequests(uint256 ) view returns(uint256 amount, bytes32 foreignAddress, uint8 status)
func (_PortMemorizer *PortMemorizerCaller) UnwrapRequests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount         *big.Int
	ForeignAddress [32]byte
	Status         uint8
}, error) {
	var out []interface{}
	err := _PortMemorizer.contract.Call(opts, &out, "unwrapRequests", arg0)

	outstruct := new(struct {
		Amount         *big.Int
		ForeignAddress [32]byte
		Status         uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ForeignAddress = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Status = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// UnwrapRequests is a free data retrieval call binding the contract method 0xd99c2a72.
//
// Solidity: function unwrapRequests(uint256 ) view returns(uint256 amount, bytes32 foreignAddress, uint8 status)
func (_PortMemorizer *PortMemorizerSession) UnwrapRequests(arg0 *big.Int) (struct {
	Amount         *big.Int
	ForeignAddress [32]byte
	Status         uint8
}, error) {
	return _PortMemorizer.Contract.UnwrapRequests(&_PortMemorizer.CallOpts, arg0)
}

// UnwrapRequests is a free data retrieval call binding the contract method 0xd99c2a72.
//
// Solidity: function unwrapRequests(uint256 ) view returns(uint256 amount, bytes32 foreignAddress, uint8 status)
func (_PortMemorizer *PortMemorizerCallerSession) UnwrapRequests(arg0 *big.Int) (struct {
	Amount         *big.Int
	ForeignAddress [32]byte
	Status         uint8
}, error) {
	return _PortMemorizer.Contract.UnwrapRequests(&_PortMemorizer.CallOpts, arg0)
}

// Drop is a paid mutator transaction binding the contract method 0x00b52395.
//
// Solidity: function drop(bytes value) returns()
func (_PortMemorizer *PortMemorizerTransactor) Drop(opts *bind.TransactOpts, value []byte) (*types.Transaction, error) {
	return _PortMemorizer.contract.Transact(opts, "drop", value)
}

// Drop is a paid mutator transaction binding the contract method 0x00b52395.
//
// Solidity: function drop(bytes value) returns()
func (_PortMemorizer *PortMemorizerSession) Drop(value []byte) (*types.Transaction, error) {
	return _PortMemorizer.Contract.Drop(&_PortMemorizer.TransactOpts, value)
}

// Drop is a paid mutator transaction binding the contract method 0x00b52395.
//
// Solidity: function drop(bytes value) returns()
func (_PortMemorizer *PortMemorizerTransactorSession) Drop(value []byte) (*types.Transaction, error) {
	return _PortMemorizer.Contract.Drop(&_PortMemorizer.TransactOpts, value)
}

// Persist is a paid mutator transaction binding the contract method 0xd3f5d56e.
//
// Solidity: function persist(bytes value) returns()
func (_PortMemorizer *PortMemorizerTransactor) Persist(opts *bind.TransactOpts, value []byte) (*types.Transaction, error) {
	return _PortMemorizer.contract.Transact(opts, "persist", value)
}

// Persist is a paid mutator transaction binding the contract method 0xd3f5d56e.
//
// Solidity: function persist(bytes value) returns()
func (_PortMemorizer *PortMemorizerSession) Persist(value []byte) (*types.Transaction, error) {
	return _PortMemorizer.Contract.Persist(&_PortMemorizer.TransactOpts, value)
}

// Persist is a paid mutator transaction binding the contract method 0xd3f5d56e.
//
// Solidity: function persist(bytes value) returns()
func (_PortMemorizer *PortMemorizerTransactorSession) Persist(value []byte) (*types.Transaction, error) {
	return _PortMemorizer.Contract.Persist(&_PortMemorizer.TransactOpts, value)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PortMemorizer *PortMemorizerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PortMemorizer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PortMemorizer *PortMemorizerSession) RenounceOwnership() (*types.Transaction, error) {
	return _PortMemorizer.Contract.RenounceOwnership(&_PortMemorizer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PortMemorizer *PortMemorizerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PortMemorizer.Contract.RenounceOwnership(&_PortMemorizer.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PortMemorizer *PortMemorizerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PortMemorizer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PortMemorizer *PortMemorizerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PortMemorizer.Contract.TransferOwnership(&_PortMemorizer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PortMemorizer *PortMemorizerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PortMemorizer.Contract.TransferOwnership(&_PortMemorizer.TransactOpts, newOwner)
}

// PortMemorizerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PortMemorizer contract.
type PortMemorizerOwnershipTransferredIterator struct {
	Event *PortMemorizerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PortMemorizerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PortMemorizerOwnershipTransferred)
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
		it.Event = new(PortMemorizerOwnershipTransferred)
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
func (it *PortMemorizerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PortMemorizerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PortMemorizerOwnershipTransferred represents a OwnershipTransferred event raised by the PortMemorizer contract.
type PortMemorizerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PortMemorizer *PortMemorizerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PortMemorizerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PortMemorizer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PortMemorizerOwnershipTransferredIterator{contract: _PortMemorizer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PortMemorizer *PortMemorizerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PortMemorizerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PortMemorizer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PortMemorizerOwnershipTransferred)
				if err := _PortMemorizer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PortMemorizer *PortMemorizerFilterer) ParseOwnershipTransferred(log types.Log) (*PortMemorizerOwnershipTransferred, error) {
	event := new(PortMemorizerOwnershipTransferred)
	if err := _PortMemorizer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PortMemorizerRequestCreatedIterator is returned from FilterRequestCreated and is used to iterate over the raw logs and unpacked data for RequestCreated events raised by the PortMemorizer contract.
type PortMemorizerRequestCreatedIterator struct {
	Event *PortMemorizerRequestCreated // Event containing the contract specifics and raw log

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
func (it *PortMemorizerRequestCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PortMemorizerRequestCreated)
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
		it.Event = new(PortMemorizerRequestCreated)
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
func (it *PortMemorizerRequestCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PortMemorizerRequestCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PortMemorizerRequestCreated represents a RequestCreated event raised by the PortMemorizer contract.
type PortMemorizerRequestCreated struct {
	Arg0 *big.Int
	Arg1 common.Address
	Arg2 [32]byte
	Arg3 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRequestCreated is a free log retrieval operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_PortMemorizer *PortMemorizerFilterer) FilterRequestCreated(opts *bind.FilterOpts) (*PortMemorizerRequestCreatedIterator, error) {

	logs, sub, err := _PortMemorizer.contract.FilterLogs(opts, "RequestCreated")
	if err != nil {
		return nil, err
	}
	return &PortMemorizerRequestCreatedIterator{contract: _PortMemorizer.contract, event: "RequestCreated", logs: logs, sub: sub}, nil
}

// WatchRequestCreated is a free log subscription operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_PortMemorizer *PortMemorizerFilterer) WatchRequestCreated(opts *bind.WatchOpts, sink chan<- *PortMemorizerRequestCreated) (event.Subscription, error) {

	logs, sub, err := _PortMemorizer.contract.WatchLogs(opts, "RequestCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PortMemorizerRequestCreated)
				if err := _PortMemorizer.contract.UnpackLog(event, "RequestCreated", log); err != nil {
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

// ParseRequestCreated is a log parse operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_PortMemorizer *PortMemorizerFilterer) ParseRequestCreated(log types.Log) (*PortMemorizerRequestCreated, error) {
	event := new(PortMemorizerRequestCreated)
	if err := _PortMemorizer.contract.UnpackLog(event, "RequestCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QueueLibABI is the input ABI used to generate the binding from.
const QueueLibABI = "[]"

// QueueLibFuncSigs maps the 4-byte function signature to its string representation.
var QueueLibFuncSigs = map[string]string{
	"9d6ad84b": "drop(QueueLib.Queue storage,bytes32)",
	"a506d954": "push(QueueLib.Queue storage,bytes32)",
}

// QueueLibBin is the compiled bytecode used for deploying new contracts.
var QueueLibBin = "0x610198610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100405760003560e01c80639d6ad84b14610045578063a506d95414610077575b600080fd5b81801561005157600080fd5b506100756004803603604081101561006857600080fd5b50803590602001356100a7565b005b81801561008357600080fd5b506100756004803603604081101561009a57600080fd5b5080359060200135610114565b6000818152600383016020908152604080832054600286019092529091205481156100e457600082815260028501602052604090208190556100e8565b8084555b8015610106576000818152600385016020526040902082905561010e565b600184018290555b50505050565b8154610129578082556001820181905561015e565b600182018054600090815260028401602081815260408084208690558454868552600388018352818520559190528120558190555b505056fea2646970667358221220ae08655007e0e16a4521804785e2701ac72d60f9722060485bec517ea1de18cc64736f6c63430007050033"

// DeployQueueLib deploys a new Ethereum contract, binding an instance of QueueLib to it.
func DeployQueueLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *QueueLib, error) {
	parsed, err := abi.JSON(strings.NewReader(QueueLibABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(QueueLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &QueueLib{QueueLibCaller: QueueLibCaller{contract: contract}, QueueLibTransactor: QueueLibTransactor{contract: contract}, QueueLibFilterer: QueueLibFilterer{contract: contract}}, nil
}

// QueueLib is an auto generated Go binding around an Ethereum contract.
type QueueLib struct {
	QueueLibCaller     // Read-only binding to the contract
	QueueLibTransactor // Write-only binding to the contract
	QueueLibFilterer   // Log filterer for contract events
}

// QueueLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type QueueLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QueueLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QueueLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QueueLibSession struct {
	Contract     *QueueLib         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QueueLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QueueLibCallerSession struct {
	Contract *QueueLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// QueueLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QueueLibTransactorSession struct {
	Contract     *QueueLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// QueueLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type QueueLibRaw struct {
	Contract *QueueLib // Generic contract binding to access the raw methods on
}

// QueueLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QueueLibCallerRaw struct {
	Contract *QueueLibCaller // Generic read-only contract binding to access the raw methods on
}

// QueueLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QueueLibTransactorRaw struct {
	Contract *QueueLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQueueLib creates a new instance of QueueLib, bound to a specific deployed contract.
func NewQueueLib(address common.Address, backend bind.ContractBackend) (*QueueLib, error) {
	contract, err := bindQueueLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QueueLib{QueueLibCaller: QueueLibCaller{contract: contract}, QueueLibTransactor: QueueLibTransactor{contract: contract}, QueueLibFilterer: QueueLibFilterer{contract: contract}}, nil
}

// NewQueueLibCaller creates a new read-only instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibCaller(address common.Address, caller bind.ContractCaller) (*QueueLibCaller, error) {
	contract, err := bindQueueLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QueueLibCaller{contract: contract}, nil
}

// NewQueueLibTransactor creates a new write-only instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibTransactor(address common.Address, transactor bind.ContractTransactor) (*QueueLibTransactor, error) {
	contract, err := bindQueueLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QueueLibTransactor{contract: contract}, nil
}

// NewQueueLibFilterer creates a new log filterer instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibFilterer(address common.Address, filterer bind.ContractFilterer) (*QueueLibFilterer, error) {
	contract, err := bindQueueLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QueueLibFilterer{contract: contract}, nil
}

// bindQueueLib binds a generic wrapper to an already deployed contract.
func bindQueueLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(QueueLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueueLib *QueueLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueueLib.Contract.QueueLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueueLib *QueueLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueueLib.Contract.QueueLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueueLib *QueueLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueueLib.Contract.QueueLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueueLib *QueueLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueueLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueueLib *QueueLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueueLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueueLib *QueueLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueueLib.Contract.contract.Transact(opts, method, params...)
}
