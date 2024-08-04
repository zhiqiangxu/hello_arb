// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrapper

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// WrapperMetaData contains all meta data concerning the Wrapper contract.
var WrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wtoken\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"}]",
	Bin: "0x6080604052610250803803809161001582610193565b6080396020809112610162576080516001600160a01b03811690819003610162576040516370a0823160e01b8082523060048301529091908383602481855afa928315610186575b600093610167575b50813b15610162576100f0926100fe92856100dc9360405163b6b55f2560e01b81526000818061009d34600483019190602083019252565b038134885af18015610155575b61013c575b5060405190815230600482015291829060249082905afa90811561012f575b600091610102575b5061022c565b604051938401908152929182906020850190565b03601f1981018352826101e2565b5190f35b6101229150863d8811610128575b61011a81836101e2565b810190610205565b386100d6565b503d610110565b610137610214565b6100ce565b8061014961014f926101cf565b80610221565b386100af565b61015d610214565b6100aa565b600080fd5b61017f919350843d86116101285761011a81836101e2565b9138610065565b61018e610214565b61005d565b6080601f91909101601f19168101906001600160401b038211908210176101b957604052565b634e487b7160e01b600052604160045260246000fd5b6001600160401b0381116101b957604052565b601f909101601f19168101906001600160401b038211908210176101b957604052565b90816020910312610162575190565b506040513d6000823e3d90fd5b600091031261016257565b9190820391821161023957565b634e487b7160e01b600052601160045260246000fdfe",
}

// WrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use WrapperMetaData.ABI instead.
var WrapperABI = WrapperMetaData.ABI

// WrapperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WrapperMetaData.Bin instead.
var WrapperBin = WrapperMetaData.Bin

// DeployWrapper deploys a new Ethereum contract, binding an instance of Wrapper to it.
func DeployWrapper(auth *bind.TransactOpts, backend bind.ContractBackend, wtoken common.Address) (common.Address, *types.Transaction, *Wrapper, error) {
	parsed, err := WrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WrapperBin), backend, wtoken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Wrapper{WrapperCaller: WrapperCaller{contract: contract}, WrapperTransactor: WrapperTransactor{contract: contract}, WrapperFilterer: WrapperFilterer{contract: contract}}, nil
}

// Wrapper is an auto generated Go binding around an Ethereum contract.
type Wrapper struct {
	WrapperCaller     // Read-only binding to the contract
	WrapperTransactor // Write-only binding to the contract
	WrapperFilterer   // Log filterer for contract events
}

// WrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type WrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrapperSession struct {
	Contract     *Wrapper          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrapperCallerSession struct {
	Contract *WrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// WrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrapperTransactorSession struct {
	Contract     *WrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// WrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type WrapperRaw struct {
	Contract *Wrapper // Generic contract binding to access the raw methods on
}

// WrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrapperCallerRaw struct {
	Contract *WrapperCaller // Generic read-only contract binding to access the raw methods on
}

// WrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrapperTransactorRaw struct {
	Contract *WrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWrapper creates a new instance of Wrapper, bound to a specific deployed contract.
func NewWrapper(address common.Address, backend bind.ContractBackend) (*Wrapper, error) {
	contract, err := bindWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Wrapper{WrapperCaller: WrapperCaller{contract: contract}, WrapperTransactor: WrapperTransactor{contract: contract}, WrapperFilterer: WrapperFilterer{contract: contract}}, nil
}

// NewWrapperCaller creates a new read-only instance of Wrapper, bound to a specific deployed contract.
func NewWrapperCaller(address common.Address, caller bind.ContractCaller) (*WrapperCaller, error) {
	contract, err := bindWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrapperCaller{contract: contract}, nil
}

// NewWrapperTransactor creates a new write-only instance of Wrapper, bound to a specific deployed contract.
func NewWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*WrapperTransactor, error) {
	contract, err := bindWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrapperTransactor{contract: contract}, nil
}

// NewWrapperFilterer creates a new log filterer instance of Wrapper, bound to a specific deployed contract.
func NewWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*WrapperFilterer, error) {
	contract, err := bindWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrapperFilterer{contract: contract}, nil
}

// bindWrapper binds a generic wrapper to an already deployed contract.
func bindWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wrapper *WrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Wrapper.Contract.WrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wrapper *WrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wrapper.Contract.WrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wrapper *WrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wrapper.Contract.WrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wrapper *WrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Wrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wrapper *WrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wrapper *WrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wrapper.Contract.contract.Transact(opts, method, params...)
}

