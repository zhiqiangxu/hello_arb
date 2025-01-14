// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swap_verifier

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

// Path is an auto generated low-level Go binding around an user-defined struct.
type Path struct {
	Exchange uint8
	Pool     common.Address
	To       common.Address
	Fee      *big.Int
}

// SwapVerifierMetaData contains all meta data concerning the SwapVerifier contract.
var SwapVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"arbToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wrapToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"wrapPath\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"arbAmountIn\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amountOutMin\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"enumExchange\",\"name\":\"exchange\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structPath[][]\",\"name\":\"paths\",\"type\":\"tuple[][]\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"}]",
	Bin: "0x60808060405261127b8038038091610016826100cc565b833960e081830191126100b05761002b610186565b9061003461019c565b9261003d6101b2565b60e0516001600160401b039591908681116100b0578461005e918501610202565b91610100518781116100b05785610076918601610268565b93610120518881116100b0578661008e918301610268565b95610140519889116100b0576100ae986100a892016102c6565b95610743565b005b600080fd5b50634e487b7160e01b600052604160045260246000fd5b6080601f91909101601f19168101906001600160401b038211908210176100f257604052565b6100fa6100b5565b604052565b608081019081106001600160401b038211176100f257604052565b6001600160401b0381116100f257604052565b604081019081106001600160401b038211176100f257604052565b606081019081106001600160401b038211176100f257604052565b601f909101601f19168101906001600160401b038211908210176100f257604052565b608051906001600160a01b03821682036100b057565b60a051906001600160a01b03821682036100b057565b60c051906001600160a01b03821682036100b057565b51906001600160a01b03821682036100b057565b6020906001600160401b0381116101f5575b60051b0190565b6101fd6100b5565b6101ee565b81601f820112156100b057805191610219836101dc565b926102276040519485610163565b808452602092838086019260051b8201019283116100b0578301905b828210610251575050505090565b83809161025d846101c8565b815201910190610243565b81601f820112156100b05780519161027f836101dc565b9261028d6040519485610163565b808452602092838086019260051b8201019283116100b0578301905b8282106102b7575050505090565b815181529083019083016102a9565b9080601f830112156100b0578151916102de836101dc565b926040926102ee84519586610163565b818552602093848087019360051b850101938285116100b057858101935b85851061031d575050505050505090565b84516001600160401b0381116100b057820184603f820112156100b0578781015190610348826101dc565b9161035586519384610163565b808352858a84019160071b830101918783116100b05791868b94929593015b81811061038b57505082935081520194019361030c565b91935091936080828903126100b0578651906103a6826100ff565b82519060028210156100b057828d92608094526103c48386016101c8565b838201526103d38a86016101c8565b8a820152606080860151908201528152019101918a93919492610374565b604051906103fe8261012d565b6005825264746573743160d81b6020830152565b50634e487b7160e01b600052603260045260246000fd5b602090805115610437570190565b61043f610412565b0190565b604090805160011015610437570190565b6020918151811015610469575b60051b010190565b610471610412565b610461565b1561047d57565b60405162461bcd60e51b815260206004820152600d60248201526c14915457d25395905492505395609a1b6044820152606490fd5b50634e487b7160e01b600052601160045260246000fd5b60019060001981146104d9570190565b61043f6104b2565b90606482018092116104ef57565b6104f76104b2565b565b90600182018092116104ef57565b919082018092116104ef57565b6040519061052182610148565b600282526040366020840137565b906020828203126100b05781516001600160401b0381116100b0576105549201610268565b90565b9190949392946080830190835260209060808285015282518091528160a0850193019160005b82811061059e5750505050906060919460018060a01b031660408201520152565b83516001600160a01b03168552938101939281019260010161057d565b506040513d6000823e3d90fd5b60009103126100b057565b906105dd826101dc565b6105ea6040519182610163565b82815280926105fb601f19916101dc565b0190602036910137565b9061060f826101dc565b61061c6040519182610163565b828152809261062d601f19916101dc565b019060005b82811061063e57505050565b806060602080938501015201610632565b908160209103126100b0575180151581036100b05790565b919082519283825260005b848110610693575050826000602080949584010152601f8019910116010190565b602081830181015184830182015201610672565b906040820191604081528151809352606081019260208093019060005b81811061072d57505050818184039101528251908183528083019281808460051b8301019501936000915b8483106106ff5750505050505090565b909192939495848061071d600193601f198682030187528a51610667565b98019301930191949392906106ef565b82511515865294840194918401916001016106c4565b919590965093929361075433610b1b565b61075d34610a82565b60018060a01b039687811688841694818614996107798b610af2565b6107896107846103f1565b610ac4565b6107938951610a82565b61079d8851610a82565b6107a78551610a82565b6107ba6107b386610429565b5151610a82565b6107d1895189518091149081610a77575b50610476565b6000998a5b8a518c1015610803576107f76107fd916107f08e8e610454565b5190610507565b9b6104c9565b9a6107d6565b98999798908c156109bb57505050508091503b156100b057600060049160405192838092630d0e30db60e41b825234905af180156109ae575b610995575b505b61084d83516105d3565b946108588451610605565b9360005b815181101561096c5761092c9060206108d68161089581610886610880878c610454565b51610429565b5101516001600160a01b031690565b8b6108a08689610454565b5160405163a9059cbb60e01b81526001600160a01b03909316600484015260248301529092839190829060009082906044820190565b03925af1801561095f575b610931575b505061092761090b6108f88387610454565b516109038489610454565b519089610d06565b610915848b610454565b52610920838c610454565b9015159052565b6104c9565b61085c565b8161095092903d10610958575b6109488183610163565b81019061064f565b5089806108e6565b503d61093e565b6109676105bb565b6108e1565b6040516020810190610991816109838a8d866106a7565b03601f198101835282610163565b5190f35b806109a26109a89261011a565b806105c8565b85610841565b6109b66105bb565b61083c565b600094506109f595845115610a3e575b506109d5426104e1565b60405163fb3bdb4160e01b81529687958694859391309160048601610557565b03923491165af18015610a31575b610a0e575b50610843565b610a2a903d806000833e610a228183610163565b81019061052f565b5085610a08565b610a396105bb565b610a03565b9350610a64610a4b610514565b94610a5586610429565b6001600160a01b039091169052565b610a7188610a5586610443565b8b6109cb565b9050865114386107cb565b6104f7906040519063f82c50f160e01b6020830152602482015260248152610aa981610148565b600080916020815191016a636f6e736f6c652e6c6f675afa50565b610983610aa96104f79260405192839163104c13eb60e21b6020840152602060248401526044830190610667565b6104f790604051906332458eed60e01b60208301521515602482015260248152610aa981610148565b60405163161765e160e11b60208201526001600160a01b0390911660248083019190915281526104f790610aa981610148565b6000198101919082116104ef57565b919082039182116104ef57565b908160209103126100b0575190565b51906001600160701b03821682036100b057565b908160609103126100b057610ba181610b79565b916040610bb060208401610b79565b92015163ffffffff811681036100b05790565b60405190610bd08261012d565b600d82526c17d9d95d105b5bdd5b9d13dd5d609a1b6020830152565b60021115610bf657565b634e487b7160e01b600052602160045260246000fd5b516002811015610bf65790565b6040516020810191906000906001600160401b03841181851017610c45575b8360405281815292369037565b610c4d6100b5565b610c38565b909260809261055495948352602083015260018060a01b031660408201528160608201520190610667565b3d15610cc5573d906001600160401b038211610cb8575b60405191610cac601f8201601f191660200184610163565b82523d6000602084013e565b610cc06100b5565b610c94565b606090565b60405190610cd78261012d565b600382526266656560e81b6020830152565b60405190610cf68261012d565b60018252603360f81b6020830152565b92909192610d34610d28610d286040610886610d228951610b4e565b89610454565b6001600160a01b031690565b6040516370a0823160e01b81523060048201529390602090859060249082905afa938415611177575b600094611156575b506000915b85518310156110ca576004906060610d8d610d28610d286020610886898d610454565b604051630240bc6b60e21b815293849182905afa9182156110bd575b6000908193611088575b506001600160701b03908116921690610de6610dd46040610886888c610454565b6001600160a01b039081169083161090565b1561107c57610e8c610ec392935b89610e87610e796060610e6e8b610e6887610e456020610e1881610886878d610454565b6040516370a0823160e01b81526001600160a01b0390911660048201529182908e90829081906024820190565b03916001600160a01b03165afa90811561106f575b600091611040575b50610b5d565b95610454565b5101518885856111f1565b968793610e87610784610bc3565b611184565b610ea3610e9e6020610886888c610454565b610b1b565b610eb26040610886878b610454565b6001600160a01b0390811691161090565b156110385760005b610ed58751610b4e565b84101561103257610ef36020610886610eed876104f9565b8a610454565b848860006001610f0c610f068585610454565b51610c0c565b610f1581610bec565b03610fd9575060009492610f79610f3a610d28610d2860206108868b9a988b98610454565b6040516336cd320560e11b60208201908152602482019690965260448101969096526001600160a01b0390921660648601529093908160848101610983565b51925af1610f85610c7d565b905b15610fab5750610fa5610f9f60406108868589610454565b926104c9565b91610d6a565b94610fd3939450610fcb915091606092610fc6610784610cca565b610454565b510151610a82565b60009190565b9492611020610ff6610d28610d2860206108868b9a988b98610454565b94610983611002610c19565b604051948593602085019863022c0d9f60e01b8a5260248601610c52565b51925af161102c610c7d565b90610f87565b30610ef3565b600090610ecb565b611062915060203d602011611068575b61105a8183610163565b810190610b6a565b38610e62565b503d611050565b6110776105bb565b610e5a565b610e8c610ec392610df4565b90506110ad91925060603d6060116110b6575b6110a58183610163565b810190610b8d565b50919038610db3565b503d61109b565b6110c56105bb565b610da9565b509290506110f1610d28610d28604061088661112596986110eb8151610b4e565b90610454565b6040516370a0823160e01b815230600482015290602090829060249082905afa90811561106f576000916110405750610b5d565b90811061114a576111456040519160208352602083015260408201604052565b600191565b50600090610554610ce9565b61117091945060203d6020116110685761105a8183610163565b9238610d65565b61117f6105bb565b610d5d565b6104f79160405191637b3338ad60e11b602084015260248301526044820152604481526080810181811060018060401b038211176111c5575b604052610aa9565b6111cd6100b5565b6111bd565b80600019048211811515166111e5570290565b6111ed6104b2565b0290565b61120f61121691949293946127109384039084821161126d576111d2565b93846111d2565b918060001904821181151516611260575b02918201809211611253575b811561123d570490565b634e487b7160e01b600052601260045260246000fd5b61125b6104b2565b611233565b6112686104b2565b611227565b6112756104b2565b6111d256fe",
}

// SwapVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use SwapVerifierMetaData.ABI instead.
var SwapVerifierABI = SwapVerifierMetaData.ABI

// SwapVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SwapVerifierMetaData.Bin instead.
var SwapVerifierBin = SwapVerifierMetaData.Bin

// DeploySwapVerifier deploys a new Ethereum contract, binding an instance of SwapVerifier to it.
func DeploySwapVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, arbToken common.Address, wrapToken common.Address, router common.Address, wrapPath []common.Address, arbAmountIn []*big.Int, amountOutMin []*big.Int, paths [][]Path) (common.Address, *types.Transaction, *SwapVerifier, error) {
	parsed, err := SwapVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SwapVerifierBin), backend, arbToken, wrapToken, router, wrapPath, arbAmountIn, amountOutMin, paths)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapVerifier{SwapVerifierCaller: SwapVerifierCaller{contract: contract}, SwapVerifierTransactor: SwapVerifierTransactor{contract: contract}, SwapVerifierFilterer: SwapVerifierFilterer{contract: contract}}, nil
}

// SwapVerifier is an auto generated Go binding around an Ethereum contract.
type SwapVerifier struct {
	SwapVerifierCaller     // Read-only binding to the contract
	SwapVerifierTransactor // Write-only binding to the contract
	SwapVerifierFilterer   // Log filterer for contract events
}

// SwapVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwapVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapVerifierSession struct {
	Contract     *SwapVerifier     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapVerifierCallerSession struct {
	Contract *SwapVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SwapVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapVerifierTransactorSession struct {
	Contract     *SwapVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SwapVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwapVerifierRaw struct {
	Contract *SwapVerifier // Generic contract binding to access the raw methods on
}

// SwapVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapVerifierCallerRaw struct {
	Contract *SwapVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// SwapVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapVerifierTransactorRaw struct {
	Contract *SwapVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapVerifier creates a new instance of SwapVerifier, bound to a specific deployed contract.
func NewSwapVerifier(address common.Address, backend bind.ContractBackend) (*SwapVerifier, error) {
	contract, err := bindSwapVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapVerifier{SwapVerifierCaller: SwapVerifierCaller{contract: contract}, SwapVerifierTransactor: SwapVerifierTransactor{contract: contract}, SwapVerifierFilterer: SwapVerifierFilterer{contract: contract}}, nil
}

// NewSwapVerifierCaller creates a new read-only instance of SwapVerifier, bound to a specific deployed contract.
func NewSwapVerifierCaller(address common.Address, caller bind.ContractCaller) (*SwapVerifierCaller, error) {
	contract, err := bindSwapVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapVerifierCaller{contract: contract}, nil
}

// NewSwapVerifierTransactor creates a new write-only instance of SwapVerifier, bound to a specific deployed contract.
func NewSwapVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*SwapVerifierTransactor, error) {
	contract, err := bindSwapVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapVerifierTransactor{contract: contract}, nil
}

// NewSwapVerifierFilterer creates a new log filterer instance of SwapVerifier, bound to a specific deployed contract.
func NewSwapVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*SwapVerifierFilterer, error) {
	contract, err := bindSwapVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapVerifierFilterer{contract: contract}, nil
}

// bindSwapVerifier binds a generic wrapper to an already deployed contract.
func bindSwapVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapVerifier *SwapVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapVerifier.Contract.SwapVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapVerifier *SwapVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapVerifier.Contract.SwapVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapVerifier *SwapVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapVerifier.Contract.SwapVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapVerifier *SwapVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapVerifier *SwapVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapVerifier *SwapVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapVerifier.Contract.contract.Transact(opts, method, params...)
}

