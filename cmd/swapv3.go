package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
)

var SwapV3Cmd = cli.Command{
	Name:  "swapv3",
	Usage: "swapv3 actions",
	Subcommands: []cli.Command{
		swapV3CountCmd,
		swapV3Price2TickCmd,
		swapV3Tick2PriceCmd,
		swapV3Float2QCmd,
		swapV3Q2FloatCmd,
		swapV3Tick2BitmapCmd,
	},
}

var swapV3CountCmd = cli.Command{
	Name:   "count",
	Usage:  "count pairs",
	Action: swapV3Count,
	Flags: []cli.Flag{
		flag.NodeRPCFlag,
		flag.FactoryFlag,
	},
}

var swapV3Price2TickCmd = cli.Command{
	Name:   "price2tick",
	Usage:  "calculate tick by price",
	Action: swapV3Price2Tick,
	Flags: []cli.Flag{
		flag.PriceFlag,
	},
}

var swapV3Tick2PriceCmd = cli.Command{
	Name:   "tick2price",
	Usage:  "calculate price by tick",
	Action: swapV3Tick2Price,
	Flags: []cli.Flag{
		flag.TickFlag,
	},
}

var swapV3Tick2BitmapCmd = cli.Command{
	Name:   "tick2bm",
	Usage:  "calculate bitmap by tick",
	Action: swapV3Tick2Bitmap,
	Flags: []cli.Flag{
		flag.TickFlag,
	},
}

var swapV3Float2QCmd = cli.Command{
	Name:   "f2q",
	Usage:  "convert float to Q64.96",
	Action: swapV3Float2Q,
	Flags: []cli.Flag{
		flag.NumbersFlag,
	},
}

var swapV3Q2FloatCmd = cli.Command{
	Name:   "q2f",
	Usage:  "convert Q64.96 to float",
	Action: swapV3Q2Float,
	Flags: []cli.Flag{
		flag.NumbersFlag,
	},
}

func swapV3Count(ctx *cli.Context) (err error) {

	client := graphql.NewClient(ctx.String(flag.NodeRPCFlag.Name), nil)

	factory := ctx.String(flag.FactoryFlag.Name)

	var query struct {
		Factory struct {
			PoolCount graphql.String `graphql:"poolCount"`
		} `graphql:"factory(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(factory),
	}

	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		return
	}

	result, err := json.MarshalIndent(query, "", "  ")
	fmt.Println(string(result))
	return
}

func swapV3Price2Tick(ctx *cli.Context) (err error) {

	f, err := strconv.ParseFloat(ctx.String(flag.PriceFlag.Name), 10)
	if err != nil {
		return
	}

	fmt.Println(int(math.Log(f) / math.Log(1.0001)))
	return
}

func swapV3Tick2Price(ctx *cli.Context) (err error) {
	t, err := strconv.ParseInt(ctx.String(flag.TickFlag.Name), 10, 64)
	if err != nil {
		return
	}

	fmt.Println(math.Pow(1.0001, float64(t)))
	return
}

func swapV3Float2Q(ctx *cli.Context) (err error) {
	ns := strings.Split(ctx.String(flag.NumbersFlag.Name), ",")
	for _, n := range ns {
		var f float64
		f, err = strconv.ParseFloat(n, 10)
		if err != nil {
			return
		}
		fmt.Println(n, big.NewFloat(f*float64(1<<96)))
	}

	return
}

func swapV3Q2Float(ctx *cli.Context) (err error) {
	ns := strings.Split(ctx.String(flag.NumbersFlag.Name), ",")
	for _, n := range ns {
		var f float64
		f, err = strconv.ParseFloat(n, 10)
		if err != nil {
			return
		}
		fmt.Println(n, f/(1<<96))
	}
	return
}

func swapV3Tick2Bitmap(ctx *cli.Context) (err error) {
	t, err := strconv.ParseInt(ctx.String(flag.TickFlag.Name), 10, 64)
	if err != nil {
		return
	}

	wordPos := t >> 8
	bitPos := t % 256

	fmt.Println("wordPos", wordPos, "bitPos", bitPos)
	return
}
