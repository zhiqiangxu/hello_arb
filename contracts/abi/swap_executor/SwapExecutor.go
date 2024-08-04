// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swap_executor

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

// SwapExecutorMetaData contains all meta data concerning the SwapExecutor contract.
var SwapExecutorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"BiswapCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPERATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"nomiswapCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"pancakeCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"enumExchange\",\"name\":\"exchange\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structPath[]\",\"name\":\"paths\",\"type\":\"tuple[]\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"enumExchange\",\"name\":\"exchange\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structPath[]\",\"name\":\"paths\",\"type\":\"tuple[]\"}],\"name\":\"swapFlashloan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swapV2Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608034620001ce57601f62001fbd38819003918201601f1916830192916001600160401b03841183851017620001d35780839260409586528339602092839181010312620001ce5751906001600160a01b03808316808403620001ce576200006733620001e9565b60016002551590811562000085575b8451611d879081620002368239f35b7f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b9299060009180835260018552868320338452855260ff87842054161562000181575b5033915416036200013f57620000ed5750620000e390620001e9565b3880808062000076565b60849083519062461bcd60e51b82526004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152fd5b60648285519062461bcd60e51b825280600483015260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152fd5b808352600185528683203384528552868320600160ff19825416179055339033907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d858a51a438620000c7565b600080fd5b634e487b7160e01b600052604160045260246000fd5b600080546001600160a01b039283166001600160a01b03198216811783556040519093909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a356fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a7146111565750806322109682146108a8578063248a9ca3146111275780632f2ff15d1461107257806336568abe14610fe05780635b3bc4fe146108a8578063715018a614610f8157806384800812146108a85780638da5cb5b14610f5857806391d1485414610f0b578063a217fddf14610eef578063a26c58c0146109d6578063adb574d4146108b1578063b2ff9f26146108a8578063beabacc8146107db578063bf21a2f114610235578063d547741f146101f4578063f2fde38b1461012c5763f5b541a6146100f1575b600080fd5b346100ec5760003660031901126100ec5760206040517f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b9298152f35b346100ec5760203660031901126100ec576101456111bf565b61014d611af6565b6001600160a01b039081169081156101a057600080546001600160a01b031981168417825560405191939192167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08484a3f35b60405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608490fd5b346100ec5760403660031901126100ec576102336004356102136111a9565b9080600052600160205261022e60016040600020015461172a565b611a7e565b005b346100ec57610243366116df565b9161024c611c9c565b6000546001600160a01b0316331480156107a3575b61026a90611b99565b8251610275816115f2565b9061028360405192836115d0565b808252610292601f19916115f2565b01366020830137819085936000945b86519363ffffffff8716948510156103ef57600491906001600160a01b0360206102cb888c611bc9565b51015116906060806102dd898d611bc9565b510151936001600160a01b0360406102f58b8f611bc9565b510151169060018060a01b0316109260405195868092630240bc6b60e21b82525afa9182156103e357610348946000906000946103ae575b506001600160701b03908116931690156103a8579190611cf0565b6103528484611bc9565b5261035d8383611bc9565b51926001600160a01b03906040906103759089611bc9565b510151169463ffffffff80911690811461039257600101946102a1565b634e487b7160e01b600052601160045260246000fd5b90611cf0565b90506103d391935060603d6060116103dc575b6103cb81836115d0565b810190611bfe565b5092908d61032d565b503d6103c1565b6040513d6000823e3d90fd5b8289878a87805194600019958681019081116103925761040f9083611bc9565b511061077a5781519182868101116103925760249460206001600160a01b03604061043c878b0186611bc9565b51015116604051978880926370a0823160e01b82523060048301525afa9586156103e357600096610746575b5081516001198101908111610392576001600160a01b039060409061048d9085611bc9565b510151166001600160a01b0382161115610730576104ad87850184611bc9565b5192600093915b6104c96104c38a880186611bc9565b51611c34565b60028110156106f15761070757909560405196879160a083019360018060a01b03166020840152604083015260806060830152845180935260c0820192602086019060005b81811061069957505050601f198284030160808301526020808251948581520191019260005b818110610680575050610550925003601f1981018752866115d0565b6001600160a01b036020610566868a0185611bc9565b51015116946bffffffffffffffffffffffff60a01b9586600354161760035560018060a01b03602061059a8a880186611bc9565b51015116803b156100ec576105cc94600080946040519788958694859363022c0d9f60e01b8552309160048601611c41565b03925af19081156103e35760249660209560409461060094610671575b506003541660035560018060a01b03940190611bc9565b51015116604051938480926370a0823160e01b82523060048301525afa80156103e35760009061063d575b610636925011611c6c565b6001600255005b506020823d602011610669575b81610657602093836115d0565b810103126100ec57610636915161062b565b3d915061064a565b61067a9061158a565b896105e9565b845183526020948501948a945090920191600101610534565b91809450949092945180519060028210156106f15782606060809260209460019652858060a01b03858201511685840152858060a01b036040820151166040840152015160608201520194019101918993949261050e565b634e487b7160e01b600052602160045260246000fd5b60405162461bcd60e51b81526020600482015260016024820152604560f81b6044820152606490fd5b61073c87850184611bc9565b51926000916104b4565b9095506020813d602011610772575b81610762602093836115d0565b810103126100ec57519487610468565b3d9150610755565b60405162461bcd60e51b81526020600482015260016024820152604160f81b6044820152606490fd5b503360009081527f31c1e66639f421f1853aeefe8ad6b62a3b96f3287efe23106923cd924aa025c2602052604090205460ff16610261565b346100ec5760603660031901126100ec576107f46111bf565b6107fc6111a9565b60443591610808611af6565b6001600160a01b039081168061083e5750600080938193829383918315610834575b1690f1156103e357005b6108fc925061082a565b60405163a9059cbb60e01b81526001600160a01b0393909316600484015260248301939093525090602090829060449082906000905af180156103e35761088157005b6102339060203d81116108a1575b61089981836115d0565b810190611b81565b503d61088f565b506100ec6111e9565b346100ec576020806003193601126100ec5767ffffffffffffffff6004358181116100ec57366023820112156100ec5780600401359182116100ec576024906005368385831b840101116100ec57610907611af6565b60005b84811061091357005b80821b83018401356001600160a01b03811691908290036100ec5761097f917f97667070c54ef182b0f5858b034beac1b6f3089aa2d3188bb1e8929f4fa9b92990816000526001808a526040600020826000528a5260ff6040600020541615610984575b505050611b72565b61090a565b82600052808a526040600020826000528a5260406000209060ff1982541617905533917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d6000604051a4878080610977565b346100ec576109e4366116df565b90916109ee611c9c565b8360018060a01b036000541633148015610eb7575b610a0c90611b99565b825115610ea15760208381015181015160405163a9059cbb60e01b81526001600160a01b0391821660048201526024810194909452909183916044918391600091165af180156103e357610e82575b5080516000198101908111610392576024906020906001600160a01b0390604090610a869086611bc9565b51015116604051928380926370a0823160e01b82523060048301525afa9081156103e357600091610e50575b506000935b8251851015610daf57600460606001600160a01b036020610ad88988611bc9565b5101511660405192838092630240bc6b60e21b82525afa9081156103e357600090600092610d8c575b506001600160701b0391821691166001600160a01b036040610b238988611bc9565b510151166001600160a01b0384161015610d86575b6001600160a01b036020610b4c8988611bc9565b5101516040516370a0823160e01b815291166004820152916020836024816001600160a01b0388165afa9283156103e357600093610d50575b50610b9382610ba994611bdd565b916060610ba08a89611bc9565b51015192611cf0565b906001600160a01b036040610bbe8887611bc9565b510151166001600160a01b039091161015610d485760005b8351600019810190811161039257861015610d425760018601808711610392576001600160a01b0390602090610c0c9087611bc9565b510151165b610c1e6104c38887611bc9565b60028110156106f157600103610cca576001600160a01b036020610c428988611bc9565b51015116803b156100ec576040516336cd320560e11b8152600481019390935260248301939093526001600160a01b03166044820152906000908290606490829084905af180156103e357610cbb575b505b610cb56001600160a01b036040610cab8786611bc9565b5101511694611b72565b93610ab7565b610cc49061158a565b84610c92565b6001600160a01b036020610cde8988611bc9565b5101511660405193610cef856115b4565b60008552813b156100ec5760008094610d1e6040519788968795869463022c0d9f60e01b865260048601611c41565b03925af180156103e357610d33575b50610c94565b610d3c9061158a565b84610d2d565b30610c11565b600090610bd6565b92506020833d602011610d7e575b81610d6b602093836115d0565b810103126100ec57915191610b93610b85565b3d9150610d5e565b90610b38565b9050610da7915060603d6060116103dc576103cb81836115d0565b509087610b01565b50925080516000198101908111610392576024916020916001600160a01b0391604091610ddb91611bc9565b51015116604051928380926370a0823160e01b82523060048301525afa9081156103e357600091610e1c575b5061063692610e1591611bdd565b1015611c6c565b90506020813d602011610e48575b81610e37602093836115d0565b810103126100ec5751610636610e07565b3d9150610e2a565b90506020813d602011610e7a575b81610e6b602093836115d0565b810103126100ec575184610ab2565b3d9150610e5e565b610e9a9060203d6020116108a15761089981836115d0565b5083610a5b565b634e487b7160e01b600052603260045260246000fd5b503360009081527f31c1e66639f421f1853aeefe8ad6b62a3b96f3287efe23106923cd924aa025c2602052604090205460ff16610a03565b346100ec5760003660031901126100ec57602060405160008152f35b346100ec5760403660031901126100ec57610f246111a9565b600435600052600160205260406000209060018060a01b0316600052602052602060ff604060002054166040519015158152f35b346100ec5760003660031901126100ec576000546040516001600160a01b039091168152602090f35b346100ec5760003660031901126100ec57610f9a611af6565b600080546001600160a01b0319811682556040519082906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08284a3f35b346100ec5760403660031901126100ec57610ff96111a9565b336001600160a01b038216036110155761023390600435611a7e565b60405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b6064820152608490fd5b346100ec5760403660031901126100ec5760043561108e6111a9565b8160005260016020526110a860016040600020015461172a565b81600052600160205260406000209060018060a01b0316908160005260205260ff60406000205416156110d757005b8160005260016020526040600020816000526020526040600020600160ff1982541617905533917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d6000604051a4005b346100ec5760203660031901126100ec5760043560005260016020526020600160406000200154604051908152f35b346100ec5760203660031901126100ec576004359063ffffffff60e01b82168092036100ec57602091637965db0b60e01b8114908115611198575b5015158152f35b6301ffc9a760e01b14905083611191565b602435906001600160a01b03821682036100ec57565b600435906001600160a01b03821682036100ec57565b35906001600160a01b03821682036100ec57565b50346100ec5760803660031901126100ec5760046001600160a01b038135818116036100ec5760643567ffffffffffffffff918282116100ec57366023830112156100ec57818401358381116100ec578201602494858201913683116100ec5783600354163303611563578460809103126100ec5785840135908382168092036100ec5760648501358681116100ec5783886112879288010161160a565b9560848601359081116100ec57850190836043830112156100ec5787820135936112b0856115f2565b906040956112c0875193846115d0565b80835260209460448685019260051b8201019283116100ec576044869101915b8383106115535750505050839688511561153f5783890151840151865163a9059cbb60e01b81529088166001600160a01b0316848201908152604492909201356020830152908490829081906040010381600080995af1801561145457611522575b5083965b885160001981019081116115105788101561150b578680876113688b8d611bc9565b51015116911610156114fa5761137e8782611bc9565b5184905b6001808a01808b116114e857906113b06104c38c8c8f9796956113a68c918a611bc9565b5101511696611bc9565b60028110156114d6578b8d8c938a931460001461146257906113d191611bc9565b5101511691823b1561145e5788516336cd320560e11b815286810191825260208201929092526001600160a01b039093166040840152918691839182908490829060600103925af1801561145457611445575b505b61143f86866114358a8c611bc9565b5101511697611b72565b96611346565b61144e9061158a565b38611424565b86513d87823e3d90fd5b8780fd5b9061146f91969596611bc9565b5101511690885193611480856115b4565b888552823b156114d2579088809493926114ae8c519788968795869463022c0d9f60e01b86528d8601611c41565b03925af18015611454576114c3575b50611426565b6114cc9061158a565b386114bd565b8880fd5b634e487b7160e01b8952602187528d89fd5b634e487b7160e01b8852601186528c88fd5b6115048782611bc9565b5184611382565b848651f35b634e487b7160e01b8652601184528a86fd5b61153890843d86116108a15761089981836115d0565b5038611342565b89603284634e487b7160e01b600052526000fd5b82358152918101918691016112e0565b60405162461bcd60e51b8152602081840152600181890152602360f91b6044820152606490fd5b67ffffffffffffffff811161159e57604052565b634e487b7160e01b600052604160045260246000fd5b6020810190811067ffffffffffffffff82111761159e57604052565b90601f8019910116810190811067ffffffffffffffff82111761159e57604052565b67ffffffffffffffff811161159e5760051b60200190565b81601f820112156100ec57803590611621826115f2565b92604092611631845195866115d0565b808552602091828087019260071b850101938185116100ec578301915b84831061165e5750505050505090565b60809081848403126100ec57865191820182811067ffffffffffffffff8211176116ca57875283359060028210156100ec57828692608094526116a28387016111d5565b838201526116b18987016111d5565b898201526060808701359082015281520192019161164e565b60246000634e487b7160e01b81526041600452fd5b9060806003198301126100ec576004356001600160a01b03811681036100ec579160243591604435916064359067ffffffffffffffff82116100ec576117279160040161160a565b90565b60008181526001602091818352604093848220338352845260ff858320541615611755575050505050565b33855193606085019267ffffffffffffffff9386811085821117611962578852602a86528686019288368537865115611a2257603084538651831015611a22576078602188015360295b8381116119b8575061197657908751936080850190858210908211176119625788526042845286840194606036873784511561194e5760308653845182101561194e5790607860218601536041915b8183116118e05750505061189e5761189a93869361187e9361186f6048946118469a519a8b957f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008c8801525180926037880190611a36565b8401917001034b99036b4b9b9b4b733903937b6329607d1b603784015251809386840190611a36565b010360288101875201856115d0565b5192839262461bcd60e51b845260048401526024830190611a59565b0390fd5b60648587519062461bcd60e51b825280600483015260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152fd5b909192600f8116601081101561193a576f181899199a1a9b1b9c1cb0b131b232b360811b901a6119108588611b61565b5360041c928015611926576000190191906117ee565b634e487b7160e01b82526011600452602482fd5b634e487b7160e01b83526032600452602483fd5b634e487b7160e01b81526032600452602490fd5b634e487b7160e01b86526041600452602486fd5b60648789519062461bcd60e51b825280600483015260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152fd5b90600f81166010811015611a0e576f181899199a1a9b1b9c1cb0b131b232b360811b901a6119e6838a611b61565b5360041c9080156119fa576000190161179f565b634e487b7160e01b87526011600452602487fd5b634e487b7160e01b88526032600452602488fd5b634e487b7160e01b86526032600452602486fd5b60005b838110611a495750506000910152565b8181015183820152602001611a39565b90602091611a7281518092818552858086019101611a36565b601f01601f1916010190565b906000918083526001602052604083209160018060a01b03169182845260205260ff604084205416611aaf57505050565b8083526001602052604083208284526020526040832060ff1981541690557ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b3393604051a4565b6000546001600160a01b03163303611b0a57565b606460405162461bcd60e51b815260206004820152602060248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152fd5b8060001904821181151516610392570290565b908151811015610ea1570160200190565b60001981146103925760010190565b908160209103126100ec575180151581036100ec5790565b15611ba057565b60405162461bcd60e51b81526020600482015260016024820152604f60f81b6044820152606490fd5b8051821015610ea15760209160051b010190565b9190820391821161039257565b51906001600160701b03821682036100ec57565b908160609103126100ec57611c1281611bea565b916040611c2160208401611bea565b92015163ffffffff811681036100ec5790565b5160028110156106f15790565b909260809261172795948352602083015260018060a01b031660408201528160608201520190611a59565b15611c7357565b60405162461bcd60e51b81526020600482015260016024820152603360f81b6044820152606490fd5b6002805414611cab5760028055565b60405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606490fd5b929091926127109182039082821161039257611d1691611d0f91611b4e565b9384611b4e565b9180600019048211811515166103925702918201809211610392578115611d3b570490565b634e487b7160e01b600052601260045260246000fdfea26469706673582212208bc1d584af407d2cd2b341db2b0566e78951fb9e029d93384eabec619693e5ed64736f6c63430008100033",
}

// SwapExecutorABI is the input ABI used to generate the binding from.
// Deprecated: Use SwapExecutorMetaData.ABI instead.
var SwapExecutorABI = SwapExecutorMetaData.ABI

// SwapExecutorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SwapExecutorMetaData.Bin instead.
var SwapExecutorBin = SwapExecutorMetaData.Bin

// DeploySwapExecutor deploys a new Ethereum contract, binding an instance of SwapExecutor to it.
func DeploySwapExecutor(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address) (common.Address, *types.Transaction, *SwapExecutor, error) {
	parsed, err := SwapExecutorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SwapExecutorBin), backend, newOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapExecutor{SwapExecutorCaller: SwapExecutorCaller{contract: contract}, SwapExecutorTransactor: SwapExecutorTransactor{contract: contract}, SwapExecutorFilterer: SwapExecutorFilterer{contract: contract}}, nil
}

// SwapExecutor is an auto generated Go binding around an Ethereum contract.
type SwapExecutor struct {
	SwapExecutorCaller     // Read-only binding to the contract
	SwapExecutorTransactor // Write-only binding to the contract
	SwapExecutorFilterer   // Log filterer for contract events
}

// SwapExecutorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwapExecutorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapExecutorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapExecutorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapExecutorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapExecutorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapExecutorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapExecutorSession struct {
	Contract     *SwapExecutor     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapExecutorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapExecutorCallerSession struct {
	Contract *SwapExecutorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SwapExecutorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapExecutorTransactorSession struct {
	Contract     *SwapExecutorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SwapExecutorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwapExecutorRaw struct {
	Contract *SwapExecutor // Generic contract binding to access the raw methods on
}

// SwapExecutorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapExecutorCallerRaw struct {
	Contract *SwapExecutorCaller // Generic read-only contract binding to access the raw methods on
}

// SwapExecutorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapExecutorTransactorRaw struct {
	Contract *SwapExecutorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapExecutor creates a new instance of SwapExecutor, bound to a specific deployed contract.
func NewSwapExecutor(address common.Address, backend bind.ContractBackend) (*SwapExecutor, error) {
	contract, err := bindSwapExecutor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapExecutor{SwapExecutorCaller: SwapExecutorCaller{contract: contract}, SwapExecutorTransactor: SwapExecutorTransactor{contract: contract}, SwapExecutorFilterer: SwapExecutorFilterer{contract: contract}}, nil
}

// NewSwapExecutorCaller creates a new read-only instance of SwapExecutor, bound to a specific deployed contract.
func NewSwapExecutorCaller(address common.Address, caller bind.ContractCaller) (*SwapExecutorCaller, error) {
	contract, err := bindSwapExecutor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorCaller{contract: contract}, nil
}

// NewSwapExecutorTransactor creates a new write-only instance of SwapExecutor, bound to a specific deployed contract.
func NewSwapExecutorTransactor(address common.Address, transactor bind.ContractTransactor) (*SwapExecutorTransactor, error) {
	contract, err := bindSwapExecutor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorTransactor{contract: contract}, nil
}

// NewSwapExecutorFilterer creates a new log filterer instance of SwapExecutor, bound to a specific deployed contract.
func NewSwapExecutorFilterer(address common.Address, filterer bind.ContractFilterer) (*SwapExecutorFilterer, error) {
	contract, err := bindSwapExecutor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorFilterer{contract: contract}, nil
}

// bindSwapExecutor binds a generic wrapper to an already deployed contract.
func bindSwapExecutor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapExecutorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapExecutor *SwapExecutorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapExecutor.Contract.SwapExecutorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapExecutor *SwapExecutorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapExecutorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapExecutor *SwapExecutorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapExecutorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapExecutor *SwapExecutorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapExecutor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapExecutor *SwapExecutorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapExecutor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapExecutor *SwapExecutorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapExecutor.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _SwapExecutor.Contract.DEFAULTADMINROLE(&_SwapExecutor.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _SwapExecutor.Contract.DEFAULTADMINROLE(&_SwapExecutor.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorCaller) OPERATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "OPERATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorSession) OPERATORROLE() ([32]byte, error) {
	return _SwapExecutor.Contract.OPERATORROLE(&_SwapExecutor.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_SwapExecutor *SwapExecutorCallerSession) OPERATORROLE() ([32]byte, error) {
	return _SwapExecutor.Contract.OPERATORROLE(&_SwapExecutor.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SwapExecutor *SwapExecutorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SwapExecutor *SwapExecutorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _SwapExecutor.Contract.GetRoleAdmin(&_SwapExecutor.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_SwapExecutor *SwapExecutorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _SwapExecutor.Contract.GetRoleAdmin(&_SwapExecutor.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SwapExecutor *SwapExecutorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SwapExecutor *SwapExecutorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _SwapExecutor.Contract.HasRole(&_SwapExecutor.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_SwapExecutor *SwapExecutorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _SwapExecutor.Contract.HasRole(&_SwapExecutor.CallOpts, role, account)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SwapExecutor *SwapExecutorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SwapExecutor *SwapExecutorSession) Owner() (common.Address, error) {
	return _SwapExecutor.Contract.Owner(&_SwapExecutor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SwapExecutor *SwapExecutorCallerSession) Owner() (common.Address, error) {
	return _SwapExecutor.Contract.Owner(&_SwapExecutor.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SwapExecutor *SwapExecutorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SwapExecutor.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SwapExecutor *SwapExecutorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SwapExecutor.Contract.SupportsInterface(&_SwapExecutor.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SwapExecutor *SwapExecutorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SwapExecutor.Contract.SupportsInterface(&_SwapExecutor.CallOpts, interfaceId)
}

// BiswapCall is a paid mutator transaction binding the contract method 0x5b3bc4fe.
//
// Solidity: function BiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactor) BiswapCall(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "BiswapCall", sender, amount0, amount1, data)
}

// BiswapCall is a paid mutator transaction binding the contract method 0x5b3bc4fe.
//
// Solidity: function BiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorSession) BiswapCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.BiswapCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// BiswapCall is a paid mutator transaction binding the contract method 0x5b3bc4fe.
//
// Solidity: function BiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) BiswapCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.BiswapCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// AddOperator is a paid mutator transaction binding the contract method 0xadb574d4.
//
// Solidity: function addOperator(address[] operators) returns()
func (_SwapExecutor *SwapExecutorTransactor) AddOperator(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "addOperator", operators)
}

// AddOperator is a paid mutator transaction binding the contract method 0xadb574d4.
//
// Solidity: function addOperator(address[] operators) returns()
func (_SwapExecutor *SwapExecutorSession) AddOperator(operators []common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.AddOperator(&_SwapExecutor.TransactOpts, operators)
}

// AddOperator is a paid mutator transaction binding the contract method 0xadb574d4.
//
// Solidity: function addOperator(address[] operators) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) AddOperator(operators []common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.AddOperator(&_SwapExecutor.TransactOpts, operators)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.GrantRole(&_SwapExecutor.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.GrantRole(&_SwapExecutor.TransactOpts, role, account)
}

// NomiswapCall is a paid mutator transaction binding the contract method 0x22109682.
//
// Solidity: function nomiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactor) NomiswapCall(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "nomiswapCall", sender, amount0, amount1, data)
}

// NomiswapCall is a paid mutator transaction binding the contract method 0x22109682.
//
// Solidity: function nomiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorSession) NomiswapCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.NomiswapCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// NomiswapCall is a paid mutator transaction binding the contract method 0x22109682.
//
// Solidity: function nomiswapCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) NomiswapCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.NomiswapCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactor) PancakeCall(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "pancakeCall", sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorSession) PancakeCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.PancakeCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// PancakeCall is a paid mutator transaction binding the contract method 0x84800812.
//
// Solidity: function pancakeCall(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) PancakeCall(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.PancakeCall(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SwapExecutor *SwapExecutorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SwapExecutor *SwapExecutorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SwapExecutor.Contract.RenounceOwnership(&_SwapExecutor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SwapExecutor *SwapExecutorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SwapExecutor.Contract.RenounceOwnership(&_SwapExecutor.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.RenounceRole(&_SwapExecutor.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.RenounceRole(&_SwapExecutor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.RevokeRole(&_SwapExecutor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.RevokeRole(&_SwapExecutor.TransactOpts, role, account)
}

// Swap is a paid mutator transaction binding the contract method 0xa26c58c0.
//
// Solidity: function swap(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorTransactor) Swap(opts *bind.TransactOpts, tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "swap", tokenIn, amountIn, amountOutMin, paths)
}

// Swap is a paid mutator transaction binding the contract method 0xa26c58c0.
//
// Solidity: function swap(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorSession) Swap(tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.Contract.Swap(&_SwapExecutor.TransactOpts, tokenIn, amountIn, amountOutMin, paths)
}

// Swap is a paid mutator transaction binding the contract method 0xa26c58c0.
//
// Solidity: function swap(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) Swap(tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.Contract.Swap(&_SwapExecutor.TransactOpts, tokenIn, amountIn, amountOutMin, paths)
}

// SwapFlashloan is a paid mutator transaction binding the contract method 0xbf21a2f1.
//
// Solidity: function swapFlashloan(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorTransactor) SwapFlashloan(opts *bind.TransactOpts, tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "swapFlashloan", tokenIn, amountIn, amountOutMin, paths)
}

// SwapFlashloan is a paid mutator transaction binding the contract method 0xbf21a2f1.
//
// Solidity: function swapFlashloan(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorSession) SwapFlashloan(tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapFlashloan(&_SwapExecutor.TransactOpts, tokenIn, amountIn, amountOutMin, paths)
}

// SwapFlashloan is a paid mutator transaction binding the contract method 0xbf21a2f1.
//
// Solidity: function swapFlashloan(address tokenIn, uint256 amountIn, uint256 amountOutMin, (uint8,address,address,uint256)[] paths) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) SwapFlashloan(tokenIn common.Address, amountIn *big.Int, amountOutMin *big.Int, paths []Path) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapFlashloan(&_SwapExecutor.TransactOpts, tokenIn, amountIn, amountOutMin, paths)
}

// SwapV2Call is a paid mutator transaction binding the contract method 0xb2ff9f26.
//
// Solidity: function swapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactor) SwapV2Call(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "swapV2Call", sender, amount0, amount1, data)
}

// SwapV2Call is a paid mutator transaction binding the contract method 0xb2ff9f26.
//
// Solidity: function swapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorSession) SwapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapV2Call(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// SwapV2Call is a paid mutator transaction binding the contract method 0xb2ff9f26.
//
// Solidity: function swapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) SwapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _SwapExecutor.Contract.SwapV2Call(&_SwapExecutor.TransactOpts, sender, amount0, amount1, data)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address token, address to, uint256 value) returns()
func (_SwapExecutor *SwapExecutorTransactor) Transfer(opts *bind.TransactOpts, token common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "transfer", token, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address token, address to, uint256 value) returns()
func (_SwapExecutor *SwapExecutorSession) Transfer(token common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapExecutor.Contract.Transfer(&_SwapExecutor.TransactOpts, token, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address token, address to, uint256 value) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) Transfer(token common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapExecutor.Contract.Transfer(&_SwapExecutor.TransactOpts, token, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SwapExecutor *SwapExecutorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SwapExecutor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SwapExecutor *SwapExecutorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.TransferOwnership(&_SwapExecutor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SwapExecutor *SwapExecutorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SwapExecutor.Contract.TransferOwnership(&_SwapExecutor.TransactOpts, newOwner)
}

// SwapExecutorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SwapExecutor contract.
type SwapExecutorOwnershipTransferredIterator struct {
	Event *SwapExecutorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SwapExecutorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapExecutorOwnershipTransferred)
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
		it.Event = new(SwapExecutorOwnershipTransferred)
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
func (it *SwapExecutorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapExecutorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapExecutorOwnershipTransferred represents a OwnershipTransferred event raised by the SwapExecutor contract.
type SwapExecutorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SwapExecutor *SwapExecutorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SwapExecutorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SwapExecutor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorOwnershipTransferredIterator{contract: _SwapExecutor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SwapExecutor *SwapExecutorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SwapExecutorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SwapExecutor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapExecutorOwnershipTransferred)
				if err := _SwapExecutor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SwapExecutor *SwapExecutorFilterer) ParseOwnershipTransferred(log types.Log) (*SwapExecutorOwnershipTransferred, error) {
	event := new(SwapExecutorOwnershipTransferred)
	if err := _SwapExecutor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapExecutorRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the SwapExecutor contract.
type SwapExecutorRoleAdminChangedIterator struct {
	Event *SwapExecutorRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *SwapExecutorRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapExecutorRoleAdminChanged)
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
		it.Event = new(SwapExecutorRoleAdminChanged)
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
func (it *SwapExecutorRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapExecutorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapExecutorRoleAdminChanged represents a RoleAdminChanged event raised by the SwapExecutor contract.
type SwapExecutorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SwapExecutor *SwapExecutorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*SwapExecutorRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _SwapExecutor.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorRoleAdminChangedIterator{contract: _SwapExecutor.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SwapExecutor *SwapExecutorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *SwapExecutorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _SwapExecutor.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapExecutorRoleAdminChanged)
				if err := _SwapExecutor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_SwapExecutor *SwapExecutorFilterer) ParseRoleAdminChanged(log types.Log) (*SwapExecutorRoleAdminChanged, error) {
	event := new(SwapExecutorRoleAdminChanged)
	if err := _SwapExecutor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapExecutorRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the SwapExecutor contract.
type SwapExecutorRoleGrantedIterator struct {
	Event *SwapExecutorRoleGranted // Event containing the contract specifics and raw log

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
func (it *SwapExecutorRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapExecutorRoleGranted)
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
		it.Event = new(SwapExecutorRoleGranted)
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
func (it *SwapExecutorRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapExecutorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapExecutorRoleGranted represents a RoleGranted event raised by the SwapExecutor contract.
type SwapExecutorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*SwapExecutorRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SwapExecutor.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorRoleGrantedIterator{contract: _SwapExecutor.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *SwapExecutorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SwapExecutor.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapExecutorRoleGranted)
				if err := _SwapExecutor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) ParseRoleGranted(log types.Log) (*SwapExecutorRoleGranted, error) {
	event := new(SwapExecutorRoleGranted)
	if err := _SwapExecutor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapExecutorRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the SwapExecutor contract.
type SwapExecutorRoleRevokedIterator struct {
	Event *SwapExecutorRoleRevoked // Event containing the contract specifics and raw log

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
func (it *SwapExecutorRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapExecutorRoleRevoked)
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
		it.Event = new(SwapExecutorRoleRevoked)
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
func (it *SwapExecutorRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapExecutorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapExecutorRoleRevoked represents a RoleRevoked event raised by the SwapExecutor contract.
type SwapExecutorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*SwapExecutorRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SwapExecutor.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &SwapExecutorRoleRevokedIterator{contract: _SwapExecutor.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *SwapExecutorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SwapExecutor.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapExecutorRoleRevoked)
				if err := _SwapExecutor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_SwapExecutor *SwapExecutorFilterer) ParseRoleRevoked(log types.Log) (*SwapExecutorRoleRevoked, error) {
	event := new(SwapExecutorRoleRevoked)
	if err := _SwapExecutor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

