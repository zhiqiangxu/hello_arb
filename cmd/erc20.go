package cmd

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
	wrapper_abi "github.com/zhiqiangxu/arbbot/contracts/abi/wrapper"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/multicall"
)

var Erc20Cmd = cli.Command{
	Name:  "erc20",
	Usage: "erc20 actions",
	Subcommands: []cli.Command{
		erc20InfoCmd,
		erc20BalanceCmd,
		erc20AllowanceCmd,
		erc20TransferFromCmd,
		erc20WrapperCmd,
	},
}

var erc20InfoCmd = cli.Command{
	Name:   "info",
	Usage:  "dump erc20 info",
	Action: erc20Info,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
	},
}

var erc20BalanceCmd = cli.Command{
	Name:   "balance",
	Usage:  "show balance info",
	Action: erc20Balance,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.AccountFlag,
	},
}

var erc20AllowanceCmd = cli.Command{
	Name:   "allowance",
	Usage:  "show allowance info",
	Action: erc20Allowance,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.OwnerFlag,
		flag.SpenderFlag,
	},
}

var erc20TransferFromCmd = cli.Command{
	Name:   "transferFrom",
	Usage:  "call transferFrom",
	Action: erc20TransferFrom,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.OwnerFlag,
		flag.SpenderFlag,
		flag.AmountFlag,
	},
}

var erc20WrapperCmd = cli.Command{
	Name:   "wrapper",
	Usage:  "call wrapper deposit and show ratio",
	Action: erc20Wrapper,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.AccountFlag,
	},
}

func erc20Info(ctx *cli.Context) (err error) {
	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	invokes := []multicall.Invoke{
		{
			Contract: contract,
			Name:     "symbol",
		},
		{
			Contract: contract,
			Name:     "name",
		},
		{
			Contract: contract,
			Name:     "decimals",
		},
	}

	result := make([]interface{}, len(invokes))
	height, err := multicall.Do(context.Background(), client, &tokenAB, invokes, result)
	if err != nil {
		return
	}

	fmt.Println("height", height, "result", result)
	return
}

func erc20Balance(ctx *cli.Context) (err error) {
	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	account := common.HexToAddress(ctx.String(flag.AccountFlag.Name))

	tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	invokes := []multicall.Invoke{
		{
			Contract: contract,
			Name:     "balanceOf",
			Args:     []interface{}{account},
		},
	}

	result := make([]interface{}, len(invokes))
	height, err := multicall.Do(context.Background(), client, &tokenAB, invokes, result)
	if err != nil {
		return
	}

	fmt.Println("height", height, "result", result)
	return
}

func erc20Allowance(ctx *cli.Context) (err error) {
	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	invokes := []multicall.Invoke{
		{
			Contract: contract,
			Name:     "allowance",
			Args: []interface{}{
				common.HexToAddress(ctx.String(flag.OwnerFlag.Name)),
				common.HexToAddress(ctx.String(flag.SpenderFlag.Name)),
			},
		},
	}

	result := make([]interface{}, len(invokes))
	height, err := multicall.Do(context.Background(), client, &tokenAB, invokes, result)
	if err != nil {
		return
	}

	fmt.Println("height", height, "result", result)
	return
}

func erc20TransferFrom(ctx *cli.Context) (err error) {
	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	amount, ok := new(big.Int).SetString(ctx.String(flag.AmountFlag.Name), 10)
	if !ok {
		err = fmt.Errorf("invalid amount:%s", ctx.String(flag.AmountFlag.Name))
		return
	}
	spender := common.HexToAddress(ctx.String(flag.SpenderFlag.Name))

	args := []interface{}{
		common.HexToAddress(ctx.String(flag.OwnerFlag.Name)),
		spender,
		amount,
	}

	method := tokenAB.Methods["transferFrom"]

	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return
	}
	resultBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{From: spender, To: &contract, Data: append(method.ID, packed...)}, nil)
	if err != nil {
		return
	}

	result, err := method.Outputs.Unpack(resultBytes)
	if err != nil {
		return
	}

	fmt.Println("result", result[0])

	return
}

func erc20Wrapper(ctx *cli.Context) (err error) {

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	account := common.HexToAddress(ctx.String(flag.AccountsFlag.Name))

	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return
	}

	if balance.Uint64() == 0 {
		err = fmt.Errorf("account has no balance")
		return
	}

	ab, _ := abi.JSON(strings.NewReader(wrapper_abi.WrapperMetaData.ABI))

	method := ab.Constructor

	wrapper := common.HexToAddress(ctx.String(flag.ContractFlag.Name))
	packed, err := method.Inputs.Pack(wrapper)
	if err != nil {
		return
	}

	resultBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{From: account, Value: big.NewInt(1), Data: append(common.FromHex(wrapper_abi.WrapperMetaData.Bin), packed...)}, nil)
	if err != nil {
		return
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return
	}
	arguments := abi.Arguments{
		{Type: uint256Type, Name: "Height"},
	}

	result, err := arguments.Unpack(resultBytes)
	if err != nil {
		return
	}

	fmt.Println("1 vs", result[0].(*big.Int))

	return
}

func fetchTokens(clients []*ethclient.Client, tokens []common.Address, tokenAB *abi.ABI) (tokenInfos []*defi.Token, err error) {

	tokenInfos = make([]*defi.Token, len(tokens))
	batch := 20
	_, err = multicall.DoSliceCvtConcurrent(context.Background(), clients, tokenAB, len(tokens), batch, func(i int) []multicall.Invoke {
		return []multicall.Invoke{
			{
				Contract: tokens[i],
				Name:     "symbol",
				Args:     []interface{}{},
			},
			{
				Contract: tokens[i],
				Name:     "name",
				Args:     []interface{}{},
			},
			{
				Contract: tokens[i],
				Name:     "decimals",
				Args:     []interface{}{},
			},
		}
	}, func(from, to int, result []interface{}) (err error) {
		for i := from; i < to; i++ {
			symbol := result[3*(i-from)]
			name := result[3*(i-from)+1]
			decimals := result[3*(i-from)+2]
			if symbol == nil || name == nil || decimals == nil {
				// auto blacklist
				panic(fmt.Sprintf("bogus token:%v symbol:%v name:%v decimals:%d\n", tokens[i], symbol, name, decimals))
			}

			tokenInfos[i] = &defi.Token{
				Address:  tokens[i],
				Symbol:   result[3*(i-from)].(string),
				Name:     result[3*(i-from)+1].(string),
				Decimals: result[3*(i-from)+2].(uint8),
			}
		}
		return
	}, nil, nil)

	return
}
