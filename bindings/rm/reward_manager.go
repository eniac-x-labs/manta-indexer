// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rm

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
	_ = abi.ConvertType
)

// RewardManagerMetaData contains all meta data concerning the RewardManager contract.
var RewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_delegationManager\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"},{\"name\":\"_stragegyManager\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"_rewardTokenAddress\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"_stakePercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegationManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rewardManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payFeeManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorClaimReward\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorRewards\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payFee\",\"inputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"payFeeManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rewardTokenAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakeHolderClaimReward\",\"inputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakePercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyStakeRewards\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateStakePercent\",\"inputs\":[{\"name\":\"_stakePercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorAndStakeReward\",\"inputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"stakerFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorClaimReward\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeHolderClaimReward\",\"inputs\":[{\"name\":\"stakeHolder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressInsufficientBalance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// RewardManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RewardManagerMetaData.ABI instead.
var RewardManagerABI = RewardManagerMetaData.ABI

// RewardManager is an auto generated Go binding around an Ethereum contract.
type RewardManager struct {
	RewardManagerCaller     // Read-only binding to the contract
	RewardManagerTransactor // Write-only binding to the contract
	RewardManagerFilterer   // Log filterer for contract events
}

// RewardManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardManagerSession struct {
	Contract     *RewardManager    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardManagerCallerSession struct {
	Contract *RewardManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RewardManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardManagerTransactorSession struct {
	Contract     *RewardManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RewardManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardManagerRaw struct {
	Contract *RewardManager // Generic contract binding to access the raw methods on
}

// RewardManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardManagerCallerRaw struct {
	Contract *RewardManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RewardManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardManagerTransactorRaw struct {
	Contract *RewardManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardManager creates a new instance of RewardManager, bound to a specific deployed contract.
func NewRewardManager(address common.Address, backend bind.ContractBackend) (*RewardManager, error) {
	contract, err := bindRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardManager{RewardManagerCaller: RewardManagerCaller{contract: contract}, RewardManagerTransactor: RewardManagerTransactor{contract: contract}, RewardManagerFilterer: RewardManagerFilterer{contract: contract}}, nil
}

// NewRewardManagerCaller creates a new read-only instance of RewardManager, bound to a specific deployed contract.
func NewRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*RewardManagerCaller, error) {
	contract, err := bindRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardManagerCaller{contract: contract}, nil
}

// NewRewardManagerTransactor creates a new write-only instance of RewardManager, bound to a specific deployed contract.
func NewRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardManagerTransactor, error) {
	contract, err := bindRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardManagerTransactor{contract: contract}, nil
}

// NewRewardManagerFilterer creates a new log filterer instance of RewardManager, bound to a specific deployed contract.
func NewRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardManagerFilterer, error) {
	contract, err := bindRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardManagerFilterer{contract: contract}, nil
}

// bindRewardManager binds a generic wrapper to an already deployed contract.
func bindRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardManager *RewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardManager.Contract.RewardManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardManager *RewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.Contract.RewardManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardManager *RewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardManager.Contract.RewardManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardManager *RewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardManager *RewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardManager *RewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardManager.Contract.contract.Transact(opts, method, params...)
}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_RewardManager *RewardManagerCaller) DelegationManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "delegationManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_RewardManager *RewardManagerSession) DelegationManager() (common.Address, error) {
	return _RewardManager.Contract.DelegationManager(&_RewardManager.CallOpts)
}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_RewardManager *RewardManagerCallerSession) DelegationManager() (common.Address, error) {
	return _RewardManager.Contract.DelegationManager(&_RewardManager.CallOpts)
}

// OperatorRewards is a free data retrieval call binding the contract method 0x41a2b8d6.
//
// Solidity: function operatorRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerCaller) OperatorRewards(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "operatorRewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperatorRewards is a free data retrieval call binding the contract method 0x41a2b8d6.
//
// Solidity: function operatorRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerSession) OperatorRewards(arg0 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.OperatorRewards(&_RewardManager.CallOpts, arg0)
}

// OperatorRewards is a free data retrieval call binding the contract method 0x41a2b8d6.
//
// Solidity: function operatorRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerCallerSession) OperatorRewards(arg0 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.OperatorRewards(&_RewardManager.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardManager *RewardManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardManager *RewardManagerSession) Owner() (common.Address, error) {
	return _RewardManager.Contract.Owner(&_RewardManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardManager *RewardManagerCallerSession) Owner() (common.Address, error) {
	return _RewardManager.Contract.Owner(&_RewardManager.CallOpts)
}

// PayFeeManager is a free data retrieval call binding the contract method 0x057e3f24.
//
// Solidity: function payFeeManager() view returns(address)
func (_RewardManager *RewardManagerCaller) PayFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "payFeeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayFeeManager is a free data retrieval call binding the contract method 0x057e3f24.
//
// Solidity: function payFeeManager() view returns(address)
func (_RewardManager *RewardManagerSession) PayFeeManager() (common.Address, error) {
	return _RewardManager.Contract.PayFeeManager(&_RewardManager.CallOpts)
}

// PayFeeManager is a free data retrieval call binding the contract method 0x057e3f24.
//
// Solidity: function payFeeManager() view returns(address)
func (_RewardManager *RewardManagerCallerSession) PayFeeManager() (common.Address, error) {
	return _RewardManager.Contract.PayFeeManager(&_RewardManager.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_RewardManager *RewardManagerCaller) RewardManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "rewardManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_RewardManager *RewardManagerSession) RewardManager() (common.Address, error) {
	return _RewardManager.Contract.RewardManager(&_RewardManager.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_RewardManager *RewardManagerCallerSession) RewardManager() (common.Address, error) {
	return _RewardManager.Contract.RewardManager(&_RewardManager.CallOpts)
}

// RewardTokenAddress is a free data retrieval call binding the contract method 0x125f9e33.
//
// Solidity: function rewardTokenAddress() view returns(address)
func (_RewardManager *RewardManagerCaller) RewardTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "rewardTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardTokenAddress is a free data retrieval call binding the contract method 0x125f9e33.
//
// Solidity: function rewardTokenAddress() view returns(address)
func (_RewardManager *RewardManagerSession) RewardTokenAddress() (common.Address, error) {
	return _RewardManager.Contract.RewardTokenAddress(&_RewardManager.CallOpts)
}

// RewardTokenAddress is a free data retrieval call binding the contract method 0x125f9e33.
//
// Solidity: function rewardTokenAddress() view returns(address)
func (_RewardManager *RewardManagerCallerSession) RewardTokenAddress() (common.Address, error) {
	return _RewardManager.Contract.RewardTokenAddress(&_RewardManager.CallOpts)
}

// StakePercent is a free data retrieval call binding the contract method 0x34a6cdc5.
//
// Solidity: function stakePercent() view returns(uint256)
func (_RewardManager *RewardManagerCaller) StakePercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "stakePercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakePercent is a free data retrieval call binding the contract method 0x34a6cdc5.
//
// Solidity: function stakePercent() view returns(uint256)
func (_RewardManager *RewardManagerSession) StakePercent() (*big.Int, error) {
	return _RewardManager.Contract.StakePercent(&_RewardManager.CallOpts)
}

// StakePercent is a free data retrieval call binding the contract method 0x34a6cdc5.
//
// Solidity: function stakePercent() view returns(uint256)
func (_RewardManager *RewardManagerCallerSession) StakePercent() (*big.Int, error) {
	return _RewardManager.Contract.StakePercent(&_RewardManager.CallOpts)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_RewardManager *RewardManagerCaller) StrategyManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "strategyManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_RewardManager *RewardManagerSession) StrategyManager() (common.Address, error) {
	return _RewardManager.Contract.StrategyManager(&_RewardManager.CallOpts)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_RewardManager *RewardManagerCallerSession) StrategyManager() (common.Address, error) {
	return _RewardManager.Contract.StrategyManager(&_RewardManager.CallOpts)
}

// StrategyStakeRewards is a free data retrieval call binding the contract method 0x65d4dab3.
//
// Solidity: function strategyStakeRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerCaller) StrategyStakeRewards(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "strategyStakeRewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StrategyStakeRewards is a free data retrieval call binding the contract method 0x65d4dab3.
//
// Solidity: function strategyStakeRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerSession) StrategyStakeRewards(arg0 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.StrategyStakeRewards(&_RewardManager.CallOpts, arg0)
}

// StrategyStakeRewards is a free data retrieval call binding the contract method 0x65d4dab3.
//
// Solidity: function strategyStakeRewards(address ) view returns(uint256)
func (_RewardManager *RewardManagerCallerSession) StrategyStakeRewards(arg0 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.StrategyStakeRewards(&_RewardManager.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _rewardManager, address _payFeeManager) returns()
func (_RewardManager *RewardManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _rewardManager common.Address, _payFeeManager common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "initialize", initialOwner, _rewardManager, _payFeeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _rewardManager, address _payFeeManager) returns()
func (_RewardManager *RewardManagerSession) Initialize(initialOwner common.Address, _rewardManager common.Address, _payFeeManager common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.Initialize(&_RewardManager.TransactOpts, initialOwner, _rewardManager, _payFeeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _rewardManager, address _payFeeManager) returns()
func (_RewardManager *RewardManagerTransactorSession) Initialize(initialOwner common.Address, _rewardManager common.Address, _payFeeManager common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.Initialize(&_RewardManager.TransactOpts, initialOwner, _rewardManager, _payFeeManager)
}

// OperatorClaimReward is a paid mutator transaction binding the contract method 0x2cb08272.
//
// Solidity: function operatorClaimReward() returns(bool)
func (_RewardManager *RewardManagerTransactor) OperatorClaimReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "operatorClaimReward")
}

// OperatorClaimReward is a paid mutator transaction binding the contract method 0x2cb08272.
//
// Solidity: function operatorClaimReward() returns(bool)
func (_RewardManager *RewardManagerSession) OperatorClaimReward() (*types.Transaction, error) {
	return _RewardManager.Contract.OperatorClaimReward(&_RewardManager.TransactOpts)
}

// OperatorClaimReward is a paid mutator transaction binding the contract method 0x2cb08272.
//
// Solidity: function operatorClaimReward() returns(bool)
func (_RewardManager *RewardManagerTransactorSession) OperatorClaimReward() (*types.Transaction, error) {
	return _RewardManager.Contract.OperatorClaimReward(&_RewardManager.TransactOpts)
}

// PayFee is a paid mutator transaction binding the contract method 0x0adfcd81.
//
// Solidity: function payFee(address strategy, address operator, uint256 baseFee) returns()
func (_RewardManager *RewardManagerTransactor) PayFee(opts *bind.TransactOpts, strategy common.Address, operator common.Address, baseFee *big.Int) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "payFee", strategy, operator, baseFee)
}

// PayFee is a paid mutator transaction binding the contract method 0x0adfcd81.
//
// Solidity: function payFee(address strategy, address operator, uint256 baseFee) returns()
func (_RewardManager *RewardManagerSession) PayFee(strategy common.Address, operator common.Address, baseFee *big.Int) (*types.Transaction, error) {
	return _RewardManager.Contract.PayFee(&_RewardManager.TransactOpts, strategy, operator, baseFee)
}

// PayFee is a paid mutator transaction binding the contract method 0x0adfcd81.
//
// Solidity: function payFee(address strategy, address operator, uint256 baseFee) returns()
func (_RewardManager *RewardManagerTransactorSession) PayFee(strategy common.Address, operator common.Address, baseFee *big.Int) (*types.Transaction, error) {
	return _RewardManager.Contract.PayFee(&_RewardManager.TransactOpts, strategy, operator, baseFee)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardManager *RewardManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardManager *RewardManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _RewardManager.Contract.RenounceOwnership(&_RewardManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardManager *RewardManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RewardManager.Contract.RenounceOwnership(&_RewardManager.TransactOpts)
}

// StakeHolderClaimReward is a paid mutator transaction binding the contract method 0x0058091d.
//
// Solidity: function stakeHolderClaimReward(address strategy) returns(bool)
func (_RewardManager *RewardManagerTransactor) StakeHolderClaimReward(opts *bind.TransactOpts, strategy common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "stakeHolderClaimReward", strategy)
}

// StakeHolderClaimReward is a paid mutator transaction binding the contract method 0x0058091d.
//
// Solidity: function stakeHolderClaimReward(address strategy) returns(bool)
func (_RewardManager *RewardManagerSession) StakeHolderClaimReward(strategy common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.StakeHolderClaimReward(&_RewardManager.TransactOpts, strategy)
}

// StakeHolderClaimReward is a paid mutator transaction binding the contract method 0x0058091d.
//
// Solidity: function stakeHolderClaimReward(address strategy) returns(bool)
func (_RewardManager *RewardManagerTransactorSession) StakeHolderClaimReward(strategy common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.StakeHolderClaimReward(&_RewardManager.TransactOpts, strategy)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardManager *RewardManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardManager *RewardManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.TransferOwnership(&_RewardManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardManager *RewardManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.TransferOwnership(&_RewardManager.TransactOpts, newOwner)
}

// UpdateStakePercent is a paid mutator transaction binding the contract method 0x86d48fb5.
//
// Solidity: function updateStakePercent(uint256 _stakePercent) returns()
func (_RewardManager *RewardManagerTransactor) UpdateStakePercent(opts *bind.TransactOpts, _stakePercent *big.Int) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "updateStakePercent", _stakePercent)
}

// UpdateStakePercent is a paid mutator transaction binding the contract method 0x86d48fb5.
//
// Solidity: function updateStakePercent(uint256 _stakePercent) returns()
func (_RewardManager *RewardManagerSession) UpdateStakePercent(_stakePercent *big.Int) (*types.Transaction, error) {
	return _RewardManager.Contract.UpdateStakePercent(&_RewardManager.TransactOpts, _stakePercent)
}

// UpdateStakePercent is a paid mutator transaction binding the contract method 0x86d48fb5.
//
// Solidity: function updateStakePercent(uint256 _stakePercent) returns()
func (_RewardManager *RewardManagerTransactorSession) UpdateStakePercent(_stakePercent *big.Int) (*types.Transaction, error) {
	return _RewardManager.Contract.UpdateStakePercent(&_RewardManager.TransactOpts, _stakePercent)
}

// RewardManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RewardManager contract.
type RewardManagerInitializedIterator struct {
	Event *RewardManagerInitialized // Event containing the contract specifics and raw log

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
func (it *RewardManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerInitialized)
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
		it.Event = new(RewardManagerInitialized)
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
func (it *RewardManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardManagerInitialized represents a Initialized event raised by the RewardManager contract.
type RewardManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RewardManager *RewardManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*RewardManagerInitializedIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RewardManagerInitializedIterator{contract: _RewardManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RewardManager *RewardManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RewardManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardManagerInitialized)
				if err := _RewardManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RewardManager *RewardManagerFilterer) ParseInitialized(log types.Log) (*RewardManagerInitialized, error) {
	event := new(RewardManagerInitialized)
	if err := _RewardManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardManagerOperatorAndStakeRewardIterator is returned from FilterOperatorAndStakeReward and is used to iterate over the raw logs and unpacked data for OperatorAndStakeReward events raised by the RewardManager contract.
type RewardManagerOperatorAndStakeRewardIterator struct {
	Event *RewardManagerOperatorAndStakeReward // Event containing the contract specifics and raw log

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
func (it *RewardManagerOperatorAndStakeRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerOperatorAndStakeReward)
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
		it.Event = new(RewardManagerOperatorAndStakeReward)
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
func (it *RewardManagerOperatorAndStakeRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardManagerOperatorAndStakeRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardManagerOperatorAndStakeReward represents a OperatorAndStakeReward event raised by the RewardManager contract.
type RewardManagerOperatorAndStakeReward struct {
	Strategy    common.Address
	Operator    common.Address
	StakerFee   *big.Int
	OperatorFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAndStakeReward is a free log retrieval operation binding the contract event 0x80531474d393ce4cf1f6eff52666de82a1d793bd3644571afcd02e839460cdc8.
//
// Solidity: event OperatorAndStakeReward(address strategy, address operator, uint256 stakerFee, uint256 operatorFee)
func (_RewardManager *RewardManagerFilterer) FilterOperatorAndStakeReward(opts *bind.FilterOpts) (*RewardManagerOperatorAndStakeRewardIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "OperatorAndStakeReward")
	if err != nil {
		return nil, err
	}
	return &RewardManagerOperatorAndStakeRewardIterator{contract: _RewardManager.contract, event: "OperatorAndStakeReward", logs: logs, sub: sub}, nil
}

// WatchOperatorAndStakeReward is a free log subscription operation binding the contract event 0x80531474d393ce4cf1f6eff52666de82a1d793bd3644571afcd02e839460cdc8.
//
// Solidity: event OperatorAndStakeReward(address strategy, address operator, uint256 stakerFee, uint256 operatorFee)
func (_RewardManager *RewardManagerFilterer) WatchOperatorAndStakeReward(opts *bind.WatchOpts, sink chan<- *RewardManagerOperatorAndStakeReward) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "OperatorAndStakeReward")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardManagerOperatorAndStakeReward)
				if err := _RewardManager.contract.UnpackLog(event, "OperatorAndStakeReward", log); err != nil {
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

// ParseOperatorAndStakeReward is a log parse operation binding the contract event 0x80531474d393ce4cf1f6eff52666de82a1d793bd3644571afcd02e839460cdc8.
//
// Solidity: event OperatorAndStakeReward(address strategy, address operator, uint256 stakerFee, uint256 operatorFee)
func (_RewardManager *RewardManagerFilterer) ParseOperatorAndStakeReward(log types.Log) (*RewardManagerOperatorAndStakeReward, error) {
	event := new(RewardManagerOperatorAndStakeReward)
	if err := _RewardManager.contract.UnpackLog(event, "OperatorAndStakeReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardManagerOperatorClaimRewardIterator is returned from FilterOperatorClaimReward and is used to iterate over the raw logs and unpacked data for OperatorClaimReward events raised by the RewardManager contract.
type RewardManagerOperatorClaimRewardIterator struct {
	Event *RewardManagerOperatorClaimReward // Event containing the contract specifics and raw log

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
func (it *RewardManagerOperatorClaimRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerOperatorClaimReward)
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
		it.Event = new(RewardManagerOperatorClaimReward)
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
func (it *RewardManagerOperatorClaimRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardManagerOperatorClaimRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardManagerOperatorClaimReward represents a OperatorClaimReward event raised by the RewardManager contract.
type RewardManagerOperatorClaimReward struct {
	Operator common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorClaimReward is a free log retrieval operation binding the contract event 0xf358f13bbeaf315167e16e69fa1b06ca23f46b7fcc119b65b735a427539b332b.
//
// Solidity: event OperatorClaimReward(address operator, uint256 amount)
func (_RewardManager *RewardManagerFilterer) FilterOperatorClaimReward(opts *bind.FilterOpts) (*RewardManagerOperatorClaimRewardIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "OperatorClaimReward")
	if err != nil {
		return nil, err
	}
	return &RewardManagerOperatorClaimRewardIterator{contract: _RewardManager.contract, event: "OperatorClaimReward", logs: logs, sub: sub}, nil
}

// WatchOperatorClaimReward is a free log subscription operation binding the contract event 0xf358f13bbeaf315167e16e69fa1b06ca23f46b7fcc119b65b735a427539b332b.
//
// Solidity: event OperatorClaimReward(address operator, uint256 amount)
func (_RewardManager *RewardManagerFilterer) WatchOperatorClaimReward(opts *bind.WatchOpts, sink chan<- *RewardManagerOperatorClaimReward) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "OperatorClaimReward")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardManagerOperatorClaimReward)
				if err := _RewardManager.contract.UnpackLog(event, "OperatorClaimReward", log); err != nil {
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

// ParseOperatorClaimReward is a log parse operation binding the contract event 0xf358f13bbeaf315167e16e69fa1b06ca23f46b7fcc119b65b735a427539b332b.
//
// Solidity: event OperatorClaimReward(address operator, uint256 amount)
func (_RewardManager *RewardManagerFilterer) ParseOperatorClaimReward(log types.Log) (*RewardManagerOperatorClaimReward, error) {
	event := new(RewardManagerOperatorClaimReward)
	if err := _RewardManager.contract.UnpackLog(event, "OperatorClaimReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RewardManager contract.
type RewardManagerOwnershipTransferredIterator struct {
	Event *RewardManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewardManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerOwnershipTransferred)
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
		it.Event = new(RewardManagerOwnershipTransferred)
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
func (it *RewardManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardManagerOwnershipTransferred represents a OwnershipTransferred event raised by the RewardManager contract.
type RewardManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RewardManager *RewardManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardManagerOwnershipTransferredIterator{contract: _RewardManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RewardManager *RewardManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardManagerOwnershipTransferred)
				if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_RewardManager *RewardManagerFilterer) ParseOwnershipTransferred(log types.Log) (*RewardManagerOwnershipTransferred, error) {
	event := new(RewardManagerOwnershipTransferred)
	if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardManagerStakeHolderClaimRewardIterator is returned from FilterStakeHolderClaimReward and is used to iterate over the raw logs and unpacked data for StakeHolderClaimReward events raised by the RewardManager contract.
type RewardManagerStakeHolderClaimRewardIterator struct {
	Event *RewardManagerStakeHolderClaimReward // Event containing the contract specifics and raw log

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
func (it *RewardManagerStakeHolderClaimRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerStakeHolderClaimReward)
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
		it.Event = new(RewardManagerStakeHolderClaimReward)
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
func (it *RewardManagerStakeHolderClaimRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardManagerStakeHolderClaimRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardManagerStakeHolderClaimReward represents a StakeHolderClaimReward event raised by the RewardManager contract.
type RewardManagerStakeHolderClaimReward struct {
	StakeHolder common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakeHolderClaimReward is a free log retrieval operation binding the contract event 0x05c713f9714a7f88571f4dd4709e245825e594762ffafd8b649886a952836185.
//
// Solidity: event StakeHolderClaimReward(address stakeHolder, uint256 amount)
func (_RewardManager *RewardManagerFilterer) FilterStakeHolderClaimReward(opts *bind.FilterOpts) (*RewardManagerStakeHolderClaimRewardIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "StakeHolderClaimReward")
	if err != nil {
		return nil, err
	}
	return &RewardManagerStakeHolderClaimRewardIterator{contract: _RewardManager.contract, event: "StakeHolderClaimReward", logs: logs, sub: sub}, nil
}

// WatchStakeHolderClaimReward is a free log subscription operation binding the contract event 0x05c713f9714a7f88571f4dd4709e245825e594762ffafd8b649886a952836185.
//
// Solidity: event StakeHolderClaimReward(address stakeHolder, uint256 amount)
func (_RewardManager *RewardManagerFilterer) WatchStakeHolderClaimReward(opts *bind.WatchOpts, sink chan<- *RewardManagerStakeHolderClaimReward) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "StakeHolderClaimReward")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardManagerStakeHolderClaimReward)
				if err := _RewardManager.contract.UnpackLog(event, "StakeHolderClaimReward", log); err != nil {
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

// ParseStakeHolderClaimReward is a log parse operation binding the contract event 0x05c713f9714a7f88571f4dd4709e245825e594762ffafd8b649886a952836185.
//
// Solidity: event StakeHolderClaimReward(address stakeHolder, uint256 amount)
func (_RewardManager *RewardManagerFilterer) ParseStakeHolderClaimReward(log types.Log) (*RewardManagerStakeHolderClaimReward, error) {
	event := new(RewardManagerStakeHolderClaimReward)
	if err := _RewardManager.contract.UnpackLog(event, "StakeHolderClaimReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
