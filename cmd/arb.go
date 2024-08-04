package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/ethclient"

	"os/signal"

	"net/http"

	_ "github.com/zhiqiangxu/util/monitor"

	zlog "github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	sclient "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/clients"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/litenode"
)

var ArbCmd = cli.Command{
	Name:  "arb",
	Usage: "arb actions",
	Subcommands: []cli.Command{
		arbFullCmd,
		arbP2PCmd,
		arbPrepareCmd,
		arbGenCPKFileCmd,
		arbStatCmd,
	},
}

var arbPrepareCmd = cli.Command{
	Name:   "prepare",
	Usage:  "arb prepare",
	Action: arbPrepare,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ReuseFlag,
		flag.DebugPortFlag,
	},
}

var arbStatCmd = cli.Command{
	Name:   "stat",
	Usage:  "query arb stat info",
	Action: arbStat,
	Flags: []cli.Flag{
		flag.NodeRPCsFlag,
		flag.OptionalHeightFlag,
	},
}

var arbGenCPKFileCmd = cli.Command{
	Name:   "gencpk",
	Usage:  "generate cpk file",
	Action: arbGenCPKFile,
	Flags: []cli.Flag{
		flag.PKFlag,
		flag.OutFlag,
	},
}

var arbCPKFileCmd = cli.Command{
	Name:   "cpk",
	Usage:  "dump cpk file",
	Action: arbCPKFile,
	Flags: []cli.Flag{
		flag.InFlag,
	},
}

var arbFullCmd = cli.Command{
	Name:   "full",
	Usage:  "arb full stack",
	Action: arbFull,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.DebugPortFlag,
	},
}

var arbP2PCmd = cli.Command{
	Name:   "p2p",
	Usage:  "arb p2p only",
	Action: arbP2P,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.DebugPortFlag,
	},
}

var cipherPK = "448e1b2ceb7d7359f7b1b99439b2c6efc795f34f93e5f764b88d15ace6d40ee4"

func init() {

	// a hack to make test find artifact
	if strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], "__debug_bin") {
		_, b, _, _ := runtime.Caller(0)
		os.Chdir(filepath.Dir(filepath.Dir(b)))
	}

	cipherKey, err := crypto.HexToECDSA(cipherPK)
	if err != nil {
		panic(err)
	}
	eciespk := ecies.ImportECDSA(cipherKey)

	for _, config := range arbConfig {
		if config.Bot.Eth != nil && config.Bot.Eth.CPKFile != "" {
			{
				var (
					cpks []string
				)
				cpkBytes, err := os.ReadFile(config.Bot.Eth.CPKFile)
				if err == nil {
					err = json.Unmarshal(cpkBytes, &cpks)
				}
				if err != nil {
					panic(fmt.Sprintf("read cpkFile failed:%v", err))
				}

				for _, cpk := range cpks {
					var plain []byte
					plain, err = ethEciesDecryptImpl(eciespk, cpk)
					if err != nil {
						panic(err)
					}

					var realPK *ecdsa.PrivateKey
					realPK, err = crypto.HexToECDSA(string(plain))
					if err != nil {
						panic(err)
					}

					config.Bot.Eth.PKs = append(config.Bot.Eth.PKs, realPK)
				}
			}
		}
	}

}

func arbP2P(ctx *cli.Context) (err error) {
	networks := strings.Split(ctx.String(flag.NetworkFlag.Name), ",")

	nodes, err := startLiteNodes(networks)
	if err != nil {
		return
	}
	defer func() {
		for _, node := range nodes {
			node.Stop()
		}
	}()

	zlog.Info().Msg("p2p started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	node := nodes[0]
	for {
		time.Sleep(time.Second * 10)
		zlog.Info().Int("peers", node.Eth.PeerCount()).Msg("stat")
		for i, peer := range node.Eth.AllPeers() {
			zlog.Info().Int("idx", i).Str("enode", peer.Node().URLv4()).Str("id", peer.ID()).Msg("peer info")
		}

		select {
		case <-signalChan:
			return
		default:
		}

	}

}

func startLiteNodes(networks []string) (nodes []*litenode.Node, err error) {
	zlog.Info().Strs("networks", networks).Msg("startLiteNodes")

	for _, network := range networks {
		if config, ok := arbConfig[network]; !ok {
			err = fmt.Errorf("network %s not supported yet", network)
			return
		} else {

			node := litenode.New(config.Lite)
			err = node.Start()
			if err != nil {
				return
			}
			nodes = append(nodes, node)

		}
	}
	return
}

func startArbbots(networks []string, lites []*litenode.Node) (robots []*bot.Bot, err error) {
	zlog.Info().Strs("networks", networks).Msg("startArbbots")

	if len(networks) != len(lites) {
		panic("#networks != #nodes")
	}

	for i, network := range networks {
		config, ok := arbConfig[network]
		if !ok {
			err = fmt.Errorf("network %s not supported yet", network)
			return
		}

		robot := bot.New(config.Bot, lites[i])
		err = robot.Start()
		if err != nil {
			return
		}
		robots = append(robots, robot)

	}
	return
}

func arbGenCPKFile(ctx *cli.Context) (err error) {

	cipherKey, err := crypto.HexToECDSA(cipherPK)
	if err != nil {
		return
	}
	eciespk := ecies.ImportECDSA(cipherKey)

	if ctx.String(flag.PKFlag.Name) == "" {
		err = fmt.Errorf("empty pk")
		return
	}
	var cpks []string
	pks := strings.Split(ctx.String(flag.PKFlag.Name), ",")
	for _, pk := range pks {
		var cpk string
		cpk, err = ethEciesEncryptImpl(eciespk, []byte(pk))
		if err != nil {
			return
		}
		cpks = append(cpks, cpk)
	}

	cpksBytes, err := json.Marshal(cpks)
	if err != nil {
		return
	}

	err = os.WriteFile(ctx.String(flag.OutFlag.Name), cpksBytes, defaultPermission)
	return
}

func arbCPKFile(ctx *cli.Context) (err error) {
	var cipherKey *ecdsa.PrivateKey
	cipherKey, err = crypto.HexToECDSA(cipherPK)
	if err != nil {
		return
	}
	eciespk := ecies.ImportECDSA(cipherKey)

	cpksBytes, err := os.ReadFile(ctx.String(flag.InFlag.Name))
	if err != nil {
		return
	}

	var cpks []string
	err = json.Unmarshal(cpksBytes, &cpks)
	if err != nil {
		return
	}

	var pks []string
	for _, cpk := range cpks {
		var plain []byte
		plain, err = ethEciesDecryptImpl(eciespk, cpk)
		if err != nil {
			err = fmt.Errorf("ethEciesDecryptImpl failed:%v", err)
			return
		}
		pks = append(pks, string(plain))
	}

	fmt.Println(pks)
	return
}

type RouterInfo struct {
	ExchangeName string
	Exchange     uint8
	Contract     common.Address
	Fee          uint64
}

func arbPrepare(ctx *cli.Context) (err error) {
	dp := ctx.Int(flag.DebugPortFlag.Name)
	if dp > 0 {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf("localhost:%d", dp), nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	networks := strings.Split(ctx.String(flag.NetworkFlag.Name), ",")

	type routerInfo struct {
		exchangeName string
		exchange     uint8
		contract     string
		fee          uint64
	}

	routerInfo2Canonical := func(routers []routerInfo) (result []RouterInfo) {
		for _, router := range routers {
			result = append(result, RouterInfo{ExchangeName: router.exchangeName, Exchange: router.exchange, Fee: router.fee, Contract: common.HexToAddress(router.contract)})
		}
		return
	}

	type networkPrepareConfig struct {
		routers    []routerInfo
		stables    []string
		threshold  uint
		tokenBatch int
	}
	config := map[string]networkPrepareConfig{
		"ok": {
			routers: []routerInfo{
				{contract: "0x865bfde337C8aFBffF144Ff4C29f9404EBb22b15", fee: 30, exchangeName: "cherryswap"}, // cherryswap
				{contract: "0xc3364A27f56b95f4bEB0742a7325D67a04D80942", fee: 30, exchangeName: "kswap"},      // kswap
				{contract: "0x069A306A638ac9d3a68a6BD8BE898774C073DCb3", fee: 30, exchangeName: "jswap"},      // jswap
			},
			stables: []string{
				"0x382bb369d343125bfb2117af9c149795c6c65c50", // usdt
				"0xc946daf81b08146b1c7a8da2a851ddf2b3eaaf85", // usdc
			},
			threshold: 10000, // pool tvl
		},
		"bsc": {
			routers: []routerInfo{
				{contract: "0x10ed43c718714eb63d5aa57b78b54704e256024e", fee: 25, exchangeName: "pancakeswap"},  // pancakeswap
				{contract: "0x3a6d8ca21d1cf76f653a67577fa0d27453350dd8" /*fee: 10,*/, exchangeName: "biswap"},   // biswap
				{contract: "0xd654953d746f0b114d1f85332dc43446ac79413d" /*fee: 10,*/, exchangeName: "nomiswap"}, // nomiswap
				{contract: "0x7dae51bd3e3376b8c7c4900e9107f12be3af1ba8" /*fee: 30,*/, exchangeName: "mdex"},     // mdex
			},
			stables: []string{
				"0x55d398326f99059fF775485246999027B3197955", // usdt
				"0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3", // dai
				"0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d", // usdc
				"0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56", // busd
			},
			threshold: 10000, // pool tvl
		},
		"eth": {
			routers: []routerInfo{
				{contract: "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", fee: 20},
			},
			stables: []string{
				"0xdac17f958d2ee523a2206206994597c13d831ec7", // usdt
			},
			threshold: 10000, // pool tvl
		},
		"fsn": {
			routers: []routerInfo{
				{contract: "0x447cb6966470c25447501e855fc8c0712349ad00", fee: 30},
			},
			stables: []string{
				"0x9636d3294e45823ec924c8d89dd1f1dffcf044e6", // usdt
				"0x6b52048a01c41d1625a6893c80fbe4aa2c22bb54", // usdc
				"0x947250c8664600b7cb18b0de73e592ed78598b8f", // dai
			},
			threshold:  1000, // pool tvl
			tokenBatch: 50,
		},
	}

	root, err := os.Getwd()
	if err != nil {
		return
	}

	for _, network := range networks {
		prepareConfig, ok := config[network]
		if !ok {
			panic(fmt.Sprintf("network %s not exists", network))
		}
		arbNetworkConfig, ok := arbConfig[network]
		if !ok {
			panic(fmt.Sprintf("network %s not exists", network))
		}

		var stableList []common.Address
		for _, contract := range prepareConfig.stables {
			stableList = append(stableList, common.HexToAddress(contract))
		}
		var arbTokens []common.Address
		for _, tokenInfo := range arbNetworkConfig.Bot.Eth.Arb.Tokens {
			arbTokens = append(arbTokens, tokenInfo.Token)
		}

		all := root + "/cmd/data/" + network + ".json"
		if !ctx.Bool(flag.ReuseFlag.Name) {
			fmt.Println("dump all pairs start")
			start := time.Now()
			err = swapPairsImpl(arbNetworkConfig.Bot.Eth.RPCs, routerInfo2Canonical(prepareConfig.routers), all, prepareConfig.tokenBatch)
			fmt.Println("dump all pairs end, took", time.Since(start), "err", err)
			if err != nil {
				return
			}
		}

		arbitraged := root + "/cmd/data/" + network + ".arb.json"
		err = swapFilterArbitragedImpl(network, all, arbitraged)
		if err != nil {
			return
		}

		verified := root + "/cmd/artifact/" + network + ".verified.json"
		fmt.Println("verify pool start")
		start := time.Now()
		err = swapVerifyPoolImpl(network, arbitraged, verified)
		fmt.Println("verify pool end, took", time.Since(start), "err", err)
		if err != nil {
			return
		}

		priced := root + "/cmd/data/" + network + ".priced.json"
		fmt.Println("calc tvl start")
		start = time.Now()
		err = swapTVLImpl(verified, stableList, prepareConfig.threshold, priced)
		fmt.Println("calc tvl end, took", time.Since(start), "err", err)
		if err != nil {
			return
		}

		fmt.Println("calc profit start")
		start = time.Now()
		err = swapProfitImpl(priced, arbTokens)
		fmt.Println("calc profit end, took", time.Since(start), "err", err)
		if err != nil {
			return
		}

	}

	return
}

func arbFull(ctx *cli.Context) (err error) {
	dp := ctx.Int(flag.DebugPortFlag.Name)
	if dp > 0 {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf("localhost:%d", dp), nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	networks := strings.Split(ctx.String(flag.NetworkFlag.Name), ",")

	lites, err := startLiteNodes(networks)
	if err != nil {
		return
	}

	robots, err := startArbbots(networks, lites)
	if err != nil {
		return
	}

	go func() {
		for {
			for i, node := range lites {
				zlog.Info().Int("node", i).Int("#peers", node.Eth.PeerCount()).Msg("p2p")
			}
			time.Sleep(time.Second)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	for _, node := range lites {
		node.Stop()
	}

	for _, bot := range robots {
		bot.Stop()
	}
	return
}

var (
	swapV2Topic common.Hash
	swapV2Event abi.Arguments

	swapV3Topic common.Hash
	swapV3Event abi.Arguments
)

func init() {
	var (
		err error
		sig string
	)
	sig, swapV2Event, err = homeabi.ParseEvent("event Swap(address indexed sender,uint amount0In,uint amount1In,uint amount0Out,uint amount1Out,address indexed to)")
	if err != nil {
		panic(err)
	}
	swapV2Topic = homeabi.SigToTopic(sig)

	sig, swapV3Event, err = homeabi.ParseEvent("event Swap(address indexed sender,address indexed recipient,int256 amount0,int256 amount1,uint160 sqrtPriceX96,uint128 liquidity,int24 tick)")
	if err != nil {
		panic(err)
	}
	swapV3Topic = homeabi.SigToTopic(sig)
}

func arbStat(ctx *cli.Context) (err error) {
	rpcs := strings.Split(ctx.String(flag.NodeRPCsFlag.Name), ",")
	var clients []*ethclient.Client
	for _, rpc := range rpcs {
		var client *ethclient.Client
		client, err = ethclient.Dial(rpc)
		if err != nil {
			return
		}
		clients = append(clients, client)
	}

	height := ctx.Uint64(flag.OptionalHeightFlag.Name)
	if height == 0 {
		height, err = clients[0].BlockNumber(context.Background())
		if err != nil {
			return
		}
	}

	sclients := sclient.Parallel(clients)
	findTarget := func(height *big.Int, pairs []common.Address, index uint) (found int, targetHash common.Hash, targetIndex uint, err error) {
		filter := ethereum.FilterQuery{
			FromBlock: height,
			ToBlock:   height,
			Addresses: pairs,
			Topics:    [][]common.Hash{{swapV2Topic, swapV3Topic}},
		}

		logs, err := sclients.FilterLogs(context.Background(), filter)
		if err != nil {
			fmt.Println("FilterLogs", err)
		}

		for _, l := range logs {
			if l.TxIndex >= index {
				break
			}
			found++
			targetIndex = l.TxIndex
			targetHash = l.TxHash
		}
		return
	}
	handleTx := func(receipt *types.Receipt) {
		if receipt.Status == types.ReceiptStatusFailed {
			return
		}
		var swaps []*SwapLog
		for _, l := range receipt.Logs {
			if l.Removed {
				return
			}

			if len(l.Topics) == 0 {
				return
			}

			if l.Topics[0] == swapV2Topic || l.Topics[0] == swapV3Topic {
				tokens, err := sclients.PoolTokens(context.Background(), l.Address)
				if err != nil {
					fmt.Println("PoolTokens", err)
					return
				}

				swap, err := parseSwapFromLog(l, tokens)
				if err != nil {
					fmt.Println("parseSwapFromLog", err)
					return
				}

				swaps = append(swaps, swap)
			}
		}

		if len(swaps) < 2 {
			return
		}

		if swaps[0].From == swaps[len(swaps)-1].To && swaps[len(swaps)-1].AmountOut.Cmp(swaps[0].AmountIn) > 0 {
			fmt.Printf("arb tx:%s, token:%s, delta:%d, idx:%d\n", receipt.TxHash, defi.FastAddrHex(swaps[0].From), new(big.Int).Sub(swaps[len(swaps)-1].AmountOut, swaps[0].AmountIn), receipt.TransactionIndex)
		} else {
			return
		}
		if receipt.TransactionIndex == 0 {
			return
		}

		pairs := make([]common.Address, 0, len(swaps))
		for _, swap := range swaps {
			pairs = append(pairs, swap.Pool)
		}
		found, targetHash, targetIndex, err := findTarget(receipt.BlockNumber, pairs, receipt.TransactionIndex)
		if err != nil {
			fmt.Println("findTarget", err)
			return
		}
		if found == 0 {
			fmt.Println("no target found")
			return
		}

		fmt.Printf("target tx:%s, idx:%d found:%d\n", targetHash, targetIndex, found)

	}
	handleHeight := func(height uint64) (err error) {
		start := time.Now()
		fmt.Println("height", height)
		block, _ := sclients.BlockByNumber(context.Background(), big.NewInt(int64(height)))
		fmt.Println("#tx", len(block.Transactions()))
		for _, tx := range block.Transactions() {
			receipt, _ := sclients.TransactionReceipt(context.Background(), tx.Hash())
			handleTx(receipt)
		}

		fmt.Println("took", time.Since(start))
		return
	}
	for {
		err = handleHeight(height)
		if err != nil {
			fmt.Println("handleHeight", err)
			time.Sleep(time.Second)
			continue
		}
		height++
	}
}

type SwapLog struct {
	Pool      common.Address
	From      common.Address
	To        common.Address
	AmountIn  *big.Int
	AmountOut *big.Int
	Dir       bool
}

func parseSwapFromLog(l *types.Log, tokens sclient.PoolTokens) (swap *SwapLog, err error) {
	switch l.Topics[0] {
	case swapV2Topic:
		v := make(map[string]interface{})
		err = swapV2Event.UnpackIntoMap(v, l.Data)
		if err != nil {
			return
		}
		if v["amount0Out"].(*big.Int).Cmp(v["amount0In"].(*big.Int)) > 0 {
			amountIn := new(big.Int).Sub(v["amount1In"].(*big.Int), v["amount1Out"].(*big.Int))
			amountOut := new(big.Int).Sub(v["amount0Out"].(*big.Int), v["amount0In"].(*big.Int))
			swap = &SwapLog{Pool: l.Address, From: tokens.Token1, To: tokens.Token0, AmountIn: amountIn, AmountOut: amountOut}
		} else {
			amountIn := new(big.Int).Sub(v["amount0In"].(*big.Int), v["amount0Out"].(*big.Int))
			amountOut := new(big.Int).Sub(v["amount1Out"].(*big.Int), v["amount1In"].(*big.Int))
			swap = &SwapLog{Pool: l.Address, From: tokens.Token0, To: tokens.Token1, Dir: true, AmountIn: amountIn, AmountOut: amountOut}
		}

	case swapV3Topic:
		v := make(map[string]interface{})
		err = swapV3Event.UnpackIntoMap(v, l.Data)
		if err != nil {
			return
		}
		if v["amount0"].(*big.Int).Sign() > 0 {
			amountIn := v["amount0"].(*big.Int)
			amountOut := new(big.Int).Abs(v["amount1"].(*big.Int))
			swap = &SwapLog{Pool: l.Address, From: tokens.Token0, To: tokens.Token1, Dir: true, AmountIn: amountIn, AmountOut: amountOut}
		} else {
			amountIn := v["amount1"].(*big.Int)
			amountOut := new(big.Int).Abs(v["amount0"].(*big.Int))
			swap = &SwapLog{Pool: l.Address, From: tokens.Token0, To: tokens.Token1, AmountIn: amountIn, AmountOut: amountOut}
		}
	default:
		err = fmt.Errorf("not swap log")
		return
	}
	return
}
