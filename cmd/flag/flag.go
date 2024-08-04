package flag

import "github.com/urfave/cli"

var NetworkFlag = cli.StringFlag{
	Name:     "network",
	Usage:    "specify network",
	Required: true,
}

var ForkFlag = cli.StringFlag{
	Name:     "fork",
	Usage:    "specify fork",
	Required: true,
}

var TxFlag = cli.StringFlag{
	Name:     "tx",
	Usage:    "specify tx hash",
	Required: true,
}

var MethodFlag = cli.StringFlag{
	Name:     "method",
	Usage:    "specify method",
	Required: true,
}

var EventFlag = cli.StringFlag{
	Name:     "event",
	Usage:    "specify event",
	Required: true,
}

var NumbersFlag = cli.StringFlag{
	Name:     "numbers",
	Usage:    "specify numbers",
	Required: true,
}

var HeightFlag = cli.Uint64Flag{
	Name:     "height",
	Usage:    "specify height",
	Required: true,
}

var DebugPortFlag = cli.IntFlag{
	Name:  "dp",
	Usage: "specify debug port",
}

var PKFlag = cli.StringFlag{
	Name:  "pk",
	Usage: "specify pk",
}

var NodeRPCFlag = cli.StringFlag{
	Name:     "node_rpc",
	Usage:    "specify one node rpc addr",
	Required: true,
}

var ListenFlag = cli.StringFlag{
	Name:     "listen",
	Usage:    "specify listen addr",
	Required: true,
}

var NodeRPCsFlag = cli.StringFlag{
	Name:     "node_rpcs",
	Usage:    "specify one node rpc addr",
	Required: true,
}

var ContractFlag = cli.StringFlag{
	Name:     "contract",
	Usage:    "specify contract addr",
	Required: true,
}

var AccountFlag = cli.StringFlag{
	Name:     "account",
	Usage:    "specify account addr",
	Required: true,
}

var AccountsFlag = cli.StringFlag{
	Name:     "accounts",
	Usage:    "specify account addr list seperated by comma",
	Required: true,
}

var WtokenFlag = cli.StringFlag{
	Name:     "wtoken",
	Usage:    "specify wrap token addr",
	Required: true,
}

var OwnerFlag = cli.StringFlag{
	Name:     "owner",
	Usage:    "specify owner addr",
	Required: true,
}

var SpenderFlag = cli.StringFlag{
	Name:     "spender",
	Usage:    "specify spender addr",
	Required: true,
}

var AmountFlag = cli.StringFlag{
	Name:     "amount",
	Usage:    "specify amount",
	Required: true,
}

var AmountInFlag = cli.StringFlag{
	Name:     "amount_in",
	Usage:    "specify amount in",
	Required: true,
}

var AmountOutFlag = cli.StringFlag{
	Name:     "amount_in",
	Usage:    "specify amount out",
	Required: true,
}

var SlotFlag = cli.StringFlag{
	Name:     "slot",
	Usage:    "specify slot",
	Required: true,
}

var FactoryFlag = cli.StringFlag{
	Name:     "factory",
	Usage:    "specify factory addr",
	Required: true,
}

var FeesFlag = cli.StringFlag{
	Name:     "fees",
	Usage:    "specify fees",
	Required: true,
}

var FeeFlag = cli.IntFlag{
	Name:     "fee",
	Usage:    "specify fee",
	Required: true,
}

var ExsFlag = cli.StringFlag{
	Name:     "exs",
	Usage:    "specify exs",
	Required: true,
}

var RouterFlag = cli.StringFlag{
	Name:     "router",
	Usage:    "specify router addr",
	Required: true,
}

var VerboseFlag = cli.BoolFlag{
	Name:  "verbose",
	Usage: "verbose or not",
}

var TokensFlag = cli.StringFlag{
	Name:     "tokens",
	Usage:    "specify tokens",
	Required: true,
}

var IndexFlag = cli.IntFlag{
	Name:     "index",
	Usage:    "specify index",
	Required: true,
}

var HexFlag = cli.BoolFlag{
	Name:  "hex",
	Usage: "in hex",
}

var PairFlag = cli.StringFlag{
	Name:  "pair",
	Usage: "specify pair",
}

var InFlag = cli.StringFlag{
	Name:     "in",
	Usage:    "in file",
	Required: true,
}

var PKFileFlag = cli.StringFlag{
	Name:     "pk_file",
	Usage:    "pk file",
	Required: true,
}

var DryRunFlag = cli.BoolFlag{
	Name:  "dry",
	Usage: "dry run",
}

var PoolsFlag = cli.StringFlag{
	Name:     "pools",
	Usage:    "pools file",
	Required: true,
}

var PathFlag = cli.StringFlag{
	Name:     "path",
	Usage:    "path file",
	Required: true,
}

var WrapPathFlag = cli.StringFlag{
	Name:     "wrap_path",
	Usage:    "wrap path",
	Required: true,
}

var TxDataFlag = cli.StringFlag{
	Name:     "tx_data",
	Usage:    "tx data",
	Required: true,
}

var GasPriceFlag = cli.StringFlag{
	Name:     "gas_price",
	Usage:    "specify gas price",
	Required: true,
}

var PriceFlag = cli.StringFlag{
	Name:     "price",
	Usage:    "specify price",
	Required: true,
}

var TickFlag = cli.StringFlag{
	Name:     "tick",
	Usage:    "specify tick",
	Required: true,
}

var GasLimitFlag = cli.Uint64Flag{
	Name:     "gas_limit",
	Usage:    "specify gas limit",
	Required: true,
}

var FromFlag = cli.StringFlag{
	Name:     "from",
	Usage:    "from",
	Required: true,
}

var ToFlag = cli.StringFlag{
	Name:     "to",
	Usage:    "to",
	Required: true,
}

var ThresholdFlag = cli.UintFlag{
	Name:     "shreshold",
	Usage:    "specify shreshold",
	Required: true,
}

var OutFlag = cli.StringFlag{
	Name:     "out",
	Usage:    "out file",
	Required: true,
}

var StableFlag = cli.StringFlag{
	Name:     "stable",
	Usage:    "stable coins",
	Required: true,
}

var ArbFlag = cli.StringFlag{
	Name:     "arb",
	Usage:    "arb coins",
	Required: true,
}

var ReuseFlag = cli.BoolFlag{
	Name:  "reuse",
	Usage: "reuse flag",
}

var InvokesFlag = cli.StringFlag{
	Name:     "invokes",
	Usage:    "invokes flag",
	Required: true,
}

var AuthFlag = cli.StringFlag{
	Name:     "auth",
	Usage:    "specify auth text",
	Required: true,
}

var WalletFlag = cli.StringFlag{
	Name:     "wallet",
	Usage:    "specify wallet path",
	Required: true,
}

var OptionalOutFlag = OutFlag
var OptionalInFlag = InFlag
var OptionalThresholdFlag = ThresholdFlag
var OptionalMethodFlag = MethodFlag
var OptionalHeightFlag = HeightFlag
var OptionalPKFlag = PKFlag
var OptionalContractFlag = ContractFlag
var OptionalListenFlag = ListenFlag

func init() {
	OptionalOutFlag.Required = false
	OptionalInFlag.Required = false
	OptionalThresholdFlag.Required = false
	OptionalMethodFlag.Required = false
	OptionalHeightFlag.Required = false
	OptionalPKFlag.Required = false
	OptionalContractFlag.Required = false
	OptionalListenFlag.Required = false
}
