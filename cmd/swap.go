package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
	"github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
	abi2 "github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
	abi3 "github.com/zhiqiangxu/arbbot/contracts/abi/swap_verifier"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	txPkg "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/tx"
	"github.com/zhiqiangxu/multicall"
	"github.com/zhiqiangxu/util/parallel"
)

var SwapCmd = cli.Command{
	Name:  "swap",
	Usage: "swap actions",
	Subcommands: []cli.Command{
		swapCountCmd,
		swapPairsCmd,
		swapFactoryCmd,
		swapQueryCmd,
		swapPoolStatCmd,
		swapPairCmd,
		swapPairAtCmd,
		swapTVLCmd,
		swapProfitCmd,
		swapIfProfitCmd,
		swapVerifyCmd,
		swapVerifyPathCmd,
		swapVerifyPoolCmd,
		swapFilterArbitragedCmd,
		swapTokensCmd,
		swapFeeCmd,
		swapBiSwapFeeCmd,
		swapMdexFeeCmd,
		swapOutCmd,
		swapBlocksCmd,
		swapEstimateDeploySwapExecutorCmd,
		swapDeploySwapExecutorCmd,
		swapPercentCmd,
		swapExactTokensForTokensTxDataCmd,
		swapExactTokensForTokensDryRunCmd,
		swapBnbPriceCmd,
	},
}

var swapCountCmd = cli.Command{
	Name:   "count",
	Usage:  "count pairs",
	Action: swapCount,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.FactoryFlag,
	},
}

var swapPairsCmd = cli.Command{
	Name:   "pairs",
	Usage:  "dump pairs",
	Action: swapPairs,
	Flags: []cli.Flag{
		flag.NodeRPCsFlag,
		flag.RouterFlag,
		flag.FeesFlag,
		flag.ExsFlag,
		flag.OptionalOutFlag,
	},
}

var swapTVLCmd = cli.Command{
	Name:   "tvl",
	Usage:  "calculate tvl of pools",
	Action: swapTVL,
	Flags: []cli.Flag{
		flag.StableFlag,
		flag.InFlag,
		flag.OptionalOutFlag,
		flag.OptionalThresholdFlag,
	},
}

var swapIfProfitCmd = cli.Command{
	Name:   "if_profit",
	Usage:  "do swap if profitable",
	Action: swapIfProfit,
	Flags: []cli.Flag{
		flag.OptionalContractFlag,
		flag.NetworkFlag,
		flag.InFlag,
		flag.DryRunFlag,
		flag.PoolsFlag,
		flag.OptionalPKFlag,
	},
}

var swapVerifyCmd = cli.Command{
	Name:   "verify",
	Usage:  "verify swap path",
	Action: swapVerify,
	Flags: []cli.Flag{
		flag.AccountFlag,
		flag.TokensFlag,
		flag.RouterFlag,
		flag.NetworkFlag,
	},
}

var swapVerifyPoolCmd = cli.Command{
	Name:   "verify_pool",
	Usage:  "verify swap pool",
	Action: swapVerifyPool,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.NetworkFlag,
		flag.OptionalOutFlag,
	},
}

var swapFilterArbitragedCmd = cli.Command{
	Name:   "filter_arbitraged",
	Usage:  "filter arbitraged pools",
	Action: swapFilterArbitraged,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.NetworkFlag,
		flag.OptionalOutFlag,
	},
}

var swapVerifyPathCmd = cli.Command{
	Name:   "verify_path",
	Usage:  "verify swap path",
	Action: swapVerifyPath,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.PathFlag,
	},
}

var swapProfitCmd = cli.Command{
	Name:   "profit",
	Usage:  "find profitable swap paths",
	Action: swapProfit,
	Flags: []cli.Flag{
		flag.ArbFlag,
		flag.InFlag,
	},
}

var swapTokensCmd = cli.Command{
	Name:   "tokens",
	Usage:  "fetch token info with multicall and invokes",
	Action: swapTokens,
	Flags: []cli.Flag{
		flag.InvokesFlag,
		flag.NetworkFlag,
	},
}

var swapFeeCmd = cli.Command{
	Name:   "fee",
	Usage:  "calculate fee based on sync and swap event",
	Action: swapFee,
	Flags: []cli.Flag{
		flag.NumbersFlag,
	},
}

var swapBiSwapFeeCmd = cli.Command{
	Name:   "biswap_fee",
	Usage:  "get swap fee from biswap pair",
	Action: swapBiSwapFee,
	Flags: []cli.Flag{
		flag.PairFlag,
	},
}

var swapMdexFeeCmd = cli.Command{
	Name:   "mdex_fee",
	Usage:  "get swap fee from mdex pair",
	Action: swapMdexFee,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.PairFlag,
	},
}

var swapOutCmd = cli.Command{
	Name:   "out",
	Usage:  "calculate out of a single pair based on amount in, reserves and fee",
	Action: swapOut,
	Flags: []cli.Flag{
		flag.NumbersFlag,
	},
}

var swapBlocksCmd = cli.Command{
	Name:   "blocks",
	Usage:  "show blocks with swaps start from specific height",
	Action: swapBlocks,
	Flags: []cli.Flag{
		flag.HeightFlag,
		flag.NetworkFlag,
	},
}

var swapEstimateDeploySwapExecutorCmd = cli.Command{
	Name:   "est_se",
	Usage:  "estimate gas cost of deploy swap executor contract",
	Action: swapEstimateDeploySwapExecutor,
	Flags:  []cli.Flag{},
}

var swapDeploySwapExecutorCmd = cli.Command{
	Name:   "deploy_se",
	Usage:  "deploy swap executor",
	Action: swapDeploySwapExecutor,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.AccountFlag,
		flag.PKFileFlag,
	},
}

var swapPercentCmd = cli.Command{
	Name:   "percent",
	Usage:  "query percent from stash",
	Action: swapPercent,
}

var swapQueryCmd = cli.Command{
	Name:   "query",
	Usage:  "query pool info",
	Action: swapQuery,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.FactoryFlag,
		flag.TokensFlag,
	},
}

var swapPairAtCmd = cli.Command{
	Name:   "pair_at",
	Usage:  "query pool info by index",
	Action: swapPairAt,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.FactoryFlag,
		flag.IndexFlag,
	},
}

var swapPairCmd = cli.Command{
	Name:   "pair",
	Usage:  "show pool info",
	Action: swapPair,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.PairFlag,
	},
}

var swapPoolStatCmd = cli.Command{
	Name:   "pool_stat",
	Usage:  "show pool stat",
	Action: swapPoolStat,
	Flags: []cli.Flag{
		flag.InFlag,
	},
}

var swapFactoryCmd = cli.Command{
	Name:   "factory",
	Usage:  "fetch factory address from router or pair",
	Action: swapFactory,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
	},
}

var swapExactTokensForTokensTxDataCmd = cli.Command{
	Name:   "swapExactTokensForTokens_data",
	Usage:  "generate txData for swapExactTokensForTokens",
	Action: swapExactTokensForTokensTxData,
	Flags: []cli.Flag{
		flag.AmountInFlag,
		flag.AmountOutFlag,
		flag.PathFlag,
		flag.ToFlag,
	},
}

var swapExactTokensForTokensDryRunCmd = cli.Command{
	Name:   "swapExactTokensForTokens_dryrun",
	Usage:  "calculate amountOutMin for swapExactTokensForTokens",
	Action: swapExactTokensForTokensDryRun,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ForkFlag,
		flag.AmountInFlag,
		flag.FeeFlag,
		flag.WrapPathFlag,
		flag.PathFlag,
	},
}

var swapBnbPriceCmd = cli.Command{
	Name:   "bnb_price",
	Usage:  "calculate bnb price via wbnb/busd of pancakeswap",
	Action: swapBnbPrice,
	Flags:  []cli.Flag{},
}

func swapCount(ctx *cli.Context) (err error) {

	factorys := toEthAddrList(strings.Split(ctx.String(flag.FactoryFlag.Name), ","))

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	height, output := fetchSwapCount(client, factorys)

	fmt.Println("height", height, "output", output)
	return
}

func fetchSwapCount(client *ethclient.Client, factorys []common.Address) (height uint64, counts []*big.Int) {
	var invokes []multicall.Invoke
	for _, factory := range factorys {
		invokes = append(invokes, multicall.Invoke{
			Contract: factory,
			Name:     "allPairsLength",
			Args:     []interface{}{},
		})
	}

	counts = make([]*big.Int, len(invokes))

	ab, _ := abi.JSON(strings.NewReader(homeabi.CherryFactoryABI))
	height, err := multicall.Do(context.Background(), client, &ab, invokes, counts)
	if err != nil {
		panic(fmt.Sprintf("fetchSwapCount:%v", err))
	}
	return
}

func toEthAddrList(l []string) []common.Address {
	r := make([]common.Address, 0, len(l))
	for _, e := range l {
		r = append(r, common.HexToAddress(e))
	}
	return r
}

func swapPairsImpl(rpcs []string, routers []RouterInfo, out string, tokenBatch int) (err error) {
	if _, err = os.Stat(out); err == nil {
		fmt.Println("out already exists")
		return
	}

	var clients []*ethclient.Client
	for _, rpc := range rpcs {
		var client *ethclient.Client
		client, err = ethclient.Dial(rpc)
		if err != nil {
			return
		}
		clients = append(clients, client)
	}

	root, err := os.Getwd()
	if err != nil {
		return
	}
	stashDir := root + "/cmd/data/stash"
	err = os.MkdirAll(stashDir, 0777)
	if err != nil {
		return
	}

	var routerAddrs []common.Address
	for _, router := range routers {
		routerAddrs = append(routerAddrs, router.Contract)
	}
	_, factorys, err := swapFactoryImpl(rpcs[0], routerAddrs)
	if err != nil {
		return
	}
	if len(factorys) != len(routers) {
		err = fmt.Errorf("#factoryes != #routers")
		return
	}

	height, counts := fetchSwapCount(clients[0], factorys)
	factoryAB, _ := abi.JSON(strings.NewReader(homeabi.CherryFactoryABI))

	batch := 200

	pairsStash := stashDir + "/pairs.json"
	pairTokensStash := stashDir + "/pairTokens.json"
	blacklistStash := stashDir + "/blacklist.json"
	var pairs [][]common.Address
	if _, err = os.Stat(pairsStash); err == nil {
		pairsBytes, err := os.ReadFile(pairsStash)
		if err == nil {
			err = json.Unmarshal(pairsBytes, &pairs)
		}
		if err != nil {
			fmt.Println("error when reading pairs stash, ignored. err:", err)
			pairs = nil
		}
		// os.Remove(pairsStash)
	}

	if len(pairs) > 0 {
		goto CheckPairs
	}
	for i, factory := range factorys {
		fmt.Println("factory idx", i)
		count := counts[i].Uint64()
		factoryPairs := make([]common.Address, count)
		start := time.Now()
		_, err = multicall.DoSliceConcurrent(context.Background(), clients, &factoryAB, int(count), batch, func(i int) []multicall.Invoke {
			return []multicall.Invoke{{
				Contract: factory,
				Name:     "allPairs",
				Args:     []interface{}{big.NewInt(int64(i))},
			}}
		}, func(from, to int) {
			if from == 0 {
				fmt.Println("[allPairs] from", from, "to", to)
				return
			}

			took := time.Since(start)
			fmt.Println("[allPairs] from", from, "to", to, "took", took, "eta", time.Duration(float64(took)*float64(count-uint64(from))/float64((from))))
		}, func(subInvokes []multicall.Invoke, err error, client *ethclient.Client) {
			subInvokesBytes, _ := json.Marshal(subInvokes)

			fmt.Println("invoke err", err)
			fmt.Println("subInvokes", string(subInvokesBytes))
		}, factoryPairs)
		if err != nil {
			return
		}

		pairs = append(pairs, factoryPairs)
	}

CheckPairs:
	{
		type Position struct {
			i int
			j int
		}
		dedup := make(map[common.Address]Position)
		for i, factoryPairs := range pairs {
			for j, pair := range factoryPairs {
				if dup, ok := dedup[pair]; ok {
					panic(fmt.Sprintf("dup pair exists (%d,%d) vs (%d,%d)", i, j, dup.i, dup.j))
				}
				dedup[pair] = Position{i: i, j: j}
			}
		}

		fmt.Println("height", height, "#pairs", len(dedup))
	}

	if out != "" {

		batch := tokenBatch
		if batch == 0 {
			batch = 100
		}

		// fetch pair tokens
		pairTokens := make([][]struct {
			Token0 common.Address
			Token1 common.Address
		}, 0, len(pairs))
		if _, err = os.Stat(pairTokensStash); err == nil {
			pairTokensBytes, err := os.ReadFile(pairTokensStash)
			if err == nil {
				err = json.Unmarshal(pairTokensBytes, &pairTokens)
			}
			if err != nil {
				fmt.Println("error when reading pairs stash, ignored. err:", err)
				pairTokens = make([][]struct {
					Token0 common.Address
					Token1 common.Address
				}, 0, len(pairs))
			}
			// os.Remove(pairTokensStash)
		}

		if len(pairTokens) > 0 {
			goto AfterPairTokens
		}

		for _, factoryPairs := range pairs {
			factoryTokens := make([]struct {
				Token0 common.Address
				Token1 common.Address
			}, len(factoryPairs))

			start := time.Now()
			_, err = multicall.DoSliceCvtConcurrent(context.Background(), clients, &homeabi.PairABI, len(factoryPairs), batch, func(i int) []multicall.Invoke {
				return []multicall.Invoke{
					{
						Contract: factoryPairs[i],
						Name:     "token0",
						Args:     []interface{}{},
					},
					{
						Contract: factoryPairs[i],
						Name:     "token1",
						Args:     []interface{}{},
					},
				}
			}, func(from, to int, result []common.Address) error {
				if 2*(to-from) != len(result) {
					panic("bug")
				}
				for k := 0; k < len(result)/2; k++ {
					factoryTokens[from+k] = struct {
						Token0 common.Address
						Token1 common.Address
					}{result[2*k], result[2*k+1]}
				}
				return nil
			}, func(from, to int) {
				if from == 0 {
					fmt.Println("[token0/token1] from", from, "to", to)
					return
				}

				took := time.Since(start)
				fmt.Println("[token0/token1] from", from, "to", to, "took", took, "eta", time.Duration(float64(took)*float64(len(factoryPairs)-from)/float64((from))))
			}, func(subInvokes []multicall.Invoke, err error, client *ethclient.Client) {
				subInvokesBytes, _ := json.Marshal(subInvokes)

				fmt.Println("invoke err", err)
				fmt.Println("subInvokes", string(subInvokesBytes))
			})
			if err != nil {
				return
			}
			pairTokens = append(pairTokens, factoryTokens)
		}

	AfterPairTokens:
		// blacklist
		blackList := make(map[common.Address]bool)
		for _, addr := range []string{
			"0x8B381Acd499D6B59E2C7f620035Bb58035dA09D2", //ok
			"0x57c0eCa48b1Fb5fB9Aeb2c0863Dd43eb72fbfeE8", //ok
			"0xFb2373aA9024Ec0d4aee5f6ab9dA1A778E277bB7", //ok
			"0x950231109A1DbD8f12f769D383aD1EF6d197E6Ac", //ok
			"0x54b69c06682FF146CD94dAD8605E15055091B198", //bsc
			"0x200b8E2895E14c6c7314c00731466C4eEE0F3e5c", //bsc
			"0x6636F7B89f64202208f608DEFFa71293EEF7b466", //bsc
			"0xAc9D0458567C85570409D5dBa53bA1a416C7f8D1", //bsc
			"0x944D9Add7D922b7bF6091341E82eB7C5d4017390", //bsc
			"0x335E73C38AE9ff15DDFfBA54D01C4Cf89ccb1D66", //bsc
			"0x9104A2fa4cB44226f07257698BB52e38eE8F51e4", //bsc
			"0x1767221a8FF6d13B04d1F4a66671d9b0FD481d6b", //bsc
			"0x4bB9156C36f766FDD6668845692C77127F77F094", //bsc
			"0x6Da8cCcde628eA2D3dA35c72F5E4394A07776a05", //bsc
			"0x4FC04c80d3B0Ad2039f3E3074198822de2ECD22B", //bsc
			"0x89fdA91d5889cE11c7Fc0a8F172e7893Ec2bB863", //bsc
			"0xC2418c9d64e4EcDC3503F42bDC14392cCBC50949", //bsc
			"0x4D27f877F4c4D6EEc5Eb0bcf53574cFED3d02868", //bsc
			"0x19CC8355426D36862C86AF927C31ff13B96d7A8F", //bsc
			"0x278e5BDaa90D392a2a0EC6DBCA2eb7F35c591C26", //bsc
			"0x57022C78058b0FdC242DA5213bf4BCb12a005A1a", //bsc
			"0x2C25894c6CB8B2bCbCfAfe3D8210c5dfa597737E", //bsc
			"0x687292d11e0b7771B7Ea8591C727940CbCE390D1", //bsc
			"0x20476bcf77FCD9c436CB16f11794f0e5f5471DeD", //bsc
			"0x68B47Ae6a929A073b9Ff26B7850e982239cDDA00", //bsc
			"0xDEDEa58385D8DFf930d11e233497BbAcb6afbac1", //bsc
			"0x85ab2F3aCAFA1BDf0DeeF5066769bc6f43105627", //bsc
			"0x105f177ed8Bb81C160e7072B6520b83AdF7e036d", //bsc
		} {
			blackList[common.HexToAddress(addr)] = true
		}
		if _, err = os.Stat(blacklistStash); err == nil {
			newblackList := make(map[common.Address]bool)
			blacklistBytes, err := os.ReadFile(blacklistStash)
			if err == nil {
				err = json.Unmarshal(blacklistBytes, &newblackList)
			}
			if err != nil {
				fmt.Println("error when reading pairs stash, ignored. err:", err)
			} else {
				blackList = newblackList
			}
			// os.Remove(blacklistStash)
		}
		fmt.Println("#blackList", len(blackList))

		onErrCalled := false
		var lock sync.Mutex

		// fetch erc20 info
		var tokens []common.Address
		{
			uniqueTokens := make(map[common.Address]bool)
			for _, factoryTokens := range pairTokens {
				for _, tokens := range factoryTokens {
					uniqueTokens[tokens.Token0] = true
					uniqueTokens[tokens.Token1] = true
				}
			}

			for blacked := range blackList {
				delete(uniqueTokens, blacked)
			}

			for token := range uniqueTokens {
				tokens = append(tokens, token)
			}
		}

		tokenInfo := make(map[common.Address]defi.Token)
		fmt.Println("#tokens", len(tokens))
		tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))
		{
			batch := 20
			start := time.Now()
			_, err = multicall.DoSliceCvtConcurrent(context.Background(), clients, &tokenAB, len(tokens), batch, func(i int) []multicall.Invoke {
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
						fmt.Printf("bogus token:%v symbol:%v name:%v decimals:%d\n", tokens[i], symbol, name, decimals)
						blackList[tokens[i]] = true
						return
					}

					tokenInfo[tokens[i]] = defi.Token{
						Address:  tokens[i],
						Symbol:   result[3*(i-from)].(string),
						Name:     result[3*(i-from)+1].(string),
						Decimals: result[3*(i-from)+2].(uint8),
					}
				}
				return
			}, func(from, to int) {
				if from == 0 {
					fmt.Println("[symbol/name/decimals] from", from, "to", to)
					return
				}

				took := time.Since(start)
				fmt.Println("[symbol/name/decimals] from", from, "to", to, "took", took, "eta", time.Duration(float64(took)*float64(len(tokens)-from)/float64((from))))
			}, func(subInvokes []multicall.Invoke, err error, client *ethclient.Client) {
				badContracts := findBadContracts(client, &tokenAB, subInvokes)

				lock.Lock()
				onErrCalled = true
				for contract := range badContracts {
					blackList[contract] = true
				}
				lock.Unlock()
				subInvokesBytes, _ := json.Marshal(subInvokes)

				fmt.Println("invoke err", err)
				fmt.Println("subInvokes", string(subInvokesBytes))
				fmt.Println("#badContracts", len(badContracts))
			})
			if err != nil {
				fmt.Printf("trying to stash for %v\n", err)
				if onErrCalled {
					blackListBytes, serr := json.Marshal(blackList)
					if serr != nil {
						fmt.Println("stash fail1", serr)
						return
					}
					serr = os.WriteFile(blacklistStash, blackListBytes, defaultPermission)
					if serr != nil {
						fmt.Println("stash fail2", serr)
						return
					}
					fmt.Println("blacklistStash stored, #blackList", len(blackList))
				}
				pairsBytes, serr := json.Marshal(pairs)
				if serr != nil {
					fmt.Println("stash fail3", serr)
					return
				}
				pairTokensBytes, serr := json.Marshal(pairTokens)
				if serr != nil {
					fmt.Println("stash fail4", serr)
					return
				}

				serr = os.WriteFile(pairsStash, pairsBytes, defaultPermission)
				if serr != nil {
					fmt.Println("stash fail5", serr)
					return
				}
				serr = os.WriteFile(pairTokensStash, pairTokensBytes, defaultPermission)
				if serr != nil {
					fmt.Println("stash fail6", serr)
					return
				}

				fmt.Println("pairsStash and pairTokensStash stored")
				return
			}
		}

		// fetch reserves
		reserves := make([][]struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}, 0, len(pairs))

		for _, factoryPairs := range pairs {
			factoryReserves := make([]struct {
				Reserve0           *big.Int
				Reserve1           *big.Int
				BlockTimestampLast uint32
			}, len(factoryPairs))

			start := time.Now()
			_, err = multicall.DoSliceConcurrent(context.Background(), clients, &homeabi.PairABI, len(factoryPairs), batch, func(i int) []multicall.Invoke {
				return []multicall.Invoke{{
					Contract: factoryPairs[i],
					Name:     "getReserves",
					Args:     []interface{}{},
				}}
			}, func(from, to int) {

				if from == 0 {
					fmt.Println("[getReserves] from", from, "to", to)
					return
				}

				took := time.Since(start)
				fmt.Println("[getReserves] from", from, "to", to, "took", took, "eta", time.Duration(float64(took)*float64(len(factoryPairs)-from)/float64((from))))
			}, func(subInvokes []multicall.Invoke, err error, client *ethclient.Client) {
				subInvokesBytes, _ := json.Marshal(subInvokes)

				fmt.Println("invoke err", err)
				fmt.Println("subInvokes", string(subInvokesBytes))
			}, factoryReserves)
			if err != nil {
				return
			}

			reserves = append(reserves, factoryReserves)
		}

		var poolList []*defi.Pool
		var tokenList []*defi.Token
		now := time.Now()
		for i, factoryPairs := range pairs {
			factoryReserves := reserves[i]
			factoryTokens := pairTokens[i]
			if len(factoryPairs) != len(factoryReserves) || len(factoryPairs) != len(factoryTokens) {
				panic("bug")
			}

			fmt.Printf("#pairs of %s:%d\n", routers[i].ExchangeName, len(factoryPairs))

			var realFee []uint32
			switch routers[i].ExchangeName {
			case "biswap", "nomiswap":

				realFee = make([]uint32, len(factoryPairs))
				// fetch fee for biswap
				ab, err := homeabi.ParseFunctionAsABI("function swapFee() returns (uint32 fee)")
				if err != nil {
					panic(err)
				}
				_, err = multicall.DoSliceConcurrent(context.Background(), clients, &ab, len(factoryPairs), 100, func(idx int) []multicall.Invoke {
					return []multicall.Invoke{{Contract: factoryPairs[idx], Name: "swapFee", Args: []interface{}{}}}
				}, nil, nil, realFee)
				if err != nil {
					panic(err)
				}

				// 千分比转万分比
				for j, fee := range realFee {
					realFee[j] = fee * 10
				}
			case "mdex":
				realFeeBig := make([]*big.Int, len(factoryPairs))
				// fetch fee for biswap
				ab, err := homeabi.ParseFunctionAsABI("function getPairFees(address pair) returns (uint256 fee)")
				if err != nil {
					panic(err)
				}
				_, err = multicall.DoSliceConcurrent(context.Background(), clients, &ab, len(factoryPairs), 100, func(idx int) []multicall.Invoke {
					return []multicall.Invoke{{Contract: factorys[i], Name: "getPairFees", Args: []interface{}{factoryPairs[idx]}}}
				}, nil, nil, realFeeBig)
				if err != nil {
					panic(err)
				}

				realFee = make([]uint32, len(factoryPairs))
				for j := range realFee {
					realFee[j] = uint32(realFeeBig[j].Uint64())
				}
			}

			for j, pair := range factoryPairs {
				pairTokens := factoryTokens[j]
				if blackList[pairTokens.Token0] || blackList[pairTokens.Token1] {
					continue
				}
				pairReserves := factoryReserves[j]
				poolTokens := []common.Address{
					pairTokens.Token0,
					pairTokens.Token1,
				}
				poolReserves := []*big.Int{
					pairReserves.Reserve0,
					pairReserves.Reserve1,
				}

				swapFee := routers[i].Fee
				if len(realFee) > 0 {
					swapFee = uint64(realFee[j])
				}
				pool := &defi.Pool{
					Exchange: routers[i].Exchange, Address: pair, Factory: factorys[i], Router: routerAddrs[i], SwapFee: swapFee, Tokens: poolTokens, Timestamp: now.Unix(),
				}
				pool.Reserves.Store(&poolReserves)
				poolList = append(poolList, pool)
			}
		}

		for addr := range tokenInfo {
			if blackList[addr] {
				continue
			}
			token := tokenInfo[addr]
			tokenList = append(tokenList, &token)
		}

		err = swapOutputData(poolList, tokenList, out)
		if err != nil {
			return
		}
	}

	return
}

func swapPairs(ctx *cli.Context) (err error) {
	routers := toEthAddrList(strings.Split(ctx.String(flag.RouterFlag.Name), ","))

	feesStrs := strings.Split(ctx.String(flag.FeesFlag.Name), ",")
	if len(routers) != len(feesStrs) {
		panic("#factorys != #feesStrs")
	}
	var (
		fees []uint64
		fee  int
	)
	for _, feeStr := range feesStrs {
		fee, err = strconv.Atoi(feeStr)
		if err != nil {
			return
		}
		fees = append(fees, uint64(fee))
	}

	exsStrs := strings.Split(ctx.String(flag.ExsFlag.Name), ",")
	if len(routers) != len(exsStrs) {
		panic("#factorys != #exsStrs")
	}
	var (
		exs []uint8
		ex  int
	)
	for _, exsStr := range exsStrs {
		ex, err = strconv.Atoi(exsStr)
		if err != nil {
			return
		}
		exs = append(exs, uint8(ex))
	}

	var routerInfos []RouterInfo
	for i := 0; i < len(routers); i++ {
		routerInfos = append(routerInfos, RouterInfo{Exchange: exs[i], Fee: fees[i], Contract: routers[i]})
	}

	err = swapPairsImpl(strings.Split(ctx.String(flag.NodeRPCsFlag.Name), ","), routerInfos, ctx.String(flag.OptionalOutFlag.Name), 0)

	return

}

func swapGetPair(client *ethclient.Client, factory, tokens0, token1 common.Address) (pair common.Address, err error) {
	invokes := []multicall.Invoke{
		{
			Contract: factory,
			Name:     "getPair",
			Args:     []interface{}{tokens0, token1},
		},
	}
	factoryAB, _ := abi.JSON(strings.NewReader(homeabi.CherryFactoryABI))
	result := make([]common.Address, 1)
	_, err = multicall.Do(context.Background(), client, &factoryAB, invokes, result)
	if err != nil {
		return
	}

	pair = result[0]
	return
}

func swapPairAt(ctx *cli.Context) (err error) {

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	factoryAB, _ := abi.JSON(strings.NewReader(homeabi.CherryFactoryABI))

	factory := common.HexToAddress(ctx.String(flag.FactoryFlag.Name))
	invokes := []multicall.Invoke{{
		Contract: factory,
		Name:     "allPairs",
		Args:     []interface{}{big.NewInt(int64(ctx.Int(flag.IndexFlag.Name)))},
	}}

	result := make([]common.Address, len(invokes))
	_, err = multicall.Do(context.Background(), client, &factoryAB, invokes, result)
	if err != nil {
		return
	}

	err = swapPairImpl(client, result[0], common.Address{}, common.Address{})
	return
}

func networkToClient(ctx *cli.Context) (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(ctx.String(flag.NetworkFlag.Name))
	if err != nil {
		config, ok := arbConfig[ctx.String(flag.NetworkFlag.Name)]
		if !ok {
			err = fmt.Errorf("invalid network:%s", ctx.String(flag.NetworkFlag.Name))
			return
		}
		client, err = ethclient.Dial(config.Bot.Eth.RPCs[0])
		if err != nil {
			return
		}
	}
	return
}

func swapQuery(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)

	tokens := toEthAddrList(strings.Split(ctx.String(flag.TokensFlag.Name), ","))
	if len(tokens) != 2 {
		panic("#tokens != 2")
	}
	if bytes.Compare(tokens[0][:], tokens[1][:]) > 0 {
		tokens[0], tokens[1] = tokens[1], tokens[0]
	}

	pair, err := swapGetPair(client, common.HexToAddress(ctx.String(flag.FactoryFlag.Name)), tokens[0], tokens[1])
	if err != nil {
		return
	}

	swapPairImpl(client, pair, tokens[0], tokens[1])
	return
}

func swapPair(ctx *cli.Context) (err error) {

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}
	pair := common.HexToAddress(ctx.String(flag.PairFlag.Name))

	err = swapPairImpl(client, pair, common.Address{}, common.Address{})
	return
}

func swapPoolStat(ctx *cli.Context) (err error) {
	pools, _, err := swapReadPoolsAndTokens(ctx.String(flag.InFlag.Name))
	if err != nil {
		return
	}

	byFee := make(map[uint64]int)
	byFactory := make(map[common.Address]int)
	for _, pool := range pools {
		byFee[pool.SwapFee] += 1
		byFactory[pool.Factory] += 1
	}

	fmt.Println(byFee)
	fmt.Println(byFactory)
	return
}

func swapGetReserves(client *ethclient.Client, pairAB *abi.ABI, pairs []common.Address) (reserves []struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, err error) {
	if pairAB == nil {
		var ab abi.ABI
		ab, err = abi.JSON(strings.NewReader(homeabi.CherryPairABI))
		if err != nil {
			return
		}
		pairAB = &ab
	}

	reserves = make([]struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}, len(pairs))

	invokeFunc := func(i int) []multicall.Invoke {
		return []multicall.Invoke{
			{
				Contract: pairs[i],
				Name:     "getReserves",
				Args:     []interface{}{},
			},
		}
	}
	_, err = multicall.DoSlice(context.Background(), client, pairAB, len(pairs), 20, invokeFunc, nil, reserves)
	if err != nil {
		return
	}

	return
}

func swapPairImpl(client *ethclient.Client, pair, token0, token1 common.Address) (err error) {
	fmt.Println("pair", pair)

	pairAB, _ := abi.JSON(strings.NewReader(homeabi.CherryPairABI))

	reserves, err := swapGetReserves(client, &pairAB, []common.Address{pair})
	if err != nil {
		return
	}

	reserve := reserves[0]

	fmt.Printf("-----\nReserve0:	%d (%x)\nReserve1:	%d (%x)\nBlockTimestampLast:	%d (%x)\n------\n", reserve.Reserve0, reserve.Reserve0, reserve.Reserve1, reserve.Reserve1, reserve.BlockTimestampLast, reserve.BlockTimestampLast)

	if token0 == (common.Address{}) || (token1 == common.Address{}) {
		invokes := []multicall.Invoke{
			{
				Contract: pair,
				Name:     "token0",
				Args:     []interface{}{},
			},
			{
				Contract: pair,
				Name:     "token1",
				Args:     []interface{}{},
			},
		}
		result := make([]common.Address, len(invokes))
		_, err = multicall.Do(context.Background(), client, &pairAB, invokes, result)
		if err != nil {
			return
		}
		token0 = result[0]
		token1 = result[1]
	}

	fmt.Printf("pair tokens:\n------------\ntoken0:	%v\ntoken1:	%v\n--------------\n", token0, token1)

	{

		tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))
		invokes := []multicall.Invoke{
			{
				Contract: token0,
				Name:     "symbol",
				Args:     []interface{}{},
			},
			{
				Contract: token0,
				Name:     "name",
				Args:     []interface{}{},
			},
			{
				Contract: token0,
				Name:     "decimals",
				Args:     []interface{}{},
			},
			{
				Contract: token1,
				Name:     "symbol",
				Args:     []interface{}{},
			},
			{
				Contract: token1,
				Name:     "name",
				Args:     []interface{}{},
			},
			{
				Contract: token1,
				Name:     "decimals",
				Args:     []interface{}{},
			},
		}
		result := make([]interface{}, len(invokes))
		_, err = multicall.Do(context.Background(), client, &tokenAB, invokes, result)
		if err != nil {
			return
		}

		symbol0 := result[0].(string)
		name0 := result[1].(string)
		decimals0 := result[2].(uint8)

		symbol1 := result[3].(string)
		name1 := result[4].(string)
		decimals1 := result[5].(uint8)

		fmt.Println(
			"Symbol0", symbol0,
			"Name0", name0,
			"Decimals0", decimals0,
			"Reserve0", big.NewInt(0).Div(reserve.Reserve0, unitReserve(decimals0)),
			"Reserve0 (raw)", reserve.Reserve0,
		)
		fmt.Println(
			"Symbol1", symbol1,
			"Name1", name1,
			"Decimals1", decimals1,
			"Reserve1", big.NewInt(0).Div(reserve.Reserve1, unitReserve(decimals1)),
			"Reserve1 (raw)", reserve.Reserve1,
		)
	}

	return
}

func swapTVL(ctx *cli.Context) (err error) {
	stables := toEthAddrList(strings.Split(ctx.String(flag.StableFlag.Name), ","))

	err = swapTVLImpl(ctx.String(flag.InFlag.Name), stables, ctx.Uint(flag.OptionalThresholdFlag.Name), ctx.String(flag.OptionalOutFlag.Name))
	return
}

func swapIfProfit(ctx *cli.Context) (err error) {
	// load pool map
	poolMap := make(map[common.Address]*defi.Pool)
	tokenMap := make(map[common.Address]*defi.Token)

	inBytes, err := os.ReadFile(ctx.String(flag.PoolsFlag.Name))
	if err != nil {
		return
	}

	var (
		pools  []*defi.Pool
		tokens []*defi.Token
	)
	data := []interface{}{&pools, &tokens}
	err = json.Unmarshal(inBytes, &data)
	if err != nil {
		return
	}

	for _, pool := range pools {
		poolMap[pool.Address] = pool
	}
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}

	// initialize others
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	var arbResult defi.ArbResult
	arbResultBytes, err := os.ReadFile(ctx.String(flag.InFlag.Name))
	if err != nil {
		return
	}
	err = json.Unmarshal(arbResultBytes, &arbResult)
	if err != nil {
		return
	}

	// fetch latest reserves
	{
		var (
			invokes []multicall.Invoke
		)
		for _, swapPath := range arbResult.SwapPaths.Swaps {
			invokes = append(invokes,
				multicall.Invoke{
					Contract: swapPath.Pool,
					Name:     "getReserves",
					Args:     []interface{}{},
				})
		}
		reserves := make([]struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}, len(invokes))

		pairAB, _ := abi.JSON(strings.NewReader(homeabi.CherryPairABI))
		_, err = multicall.Do(context.Background(), client, &pairAB, invokes, reserves)
		if err != nil {
			return
		}

		for i := 0; i < len(reserves); i++ {
			pool := poolMap[arbResult.SwapPaths.Swaps[i].Pool]
			reserve := []*big.Int{reserves[i].Reserve0, reserves[i].Reserve1}
			pool.Reserves.Store(&reserve)
		}

		amountIn := defi.SwapArbAmountIn(arbResult.SwapPaths, poolMap, nil)
		if amountIn == nil || amountIn.Sign() <= 0 {
			fmt.Println("negative amount", amountIn)
			return
		}
		profit := arbResult.SwapPaths.Profit(amountIn, poolMap, nil)
		if profit.Sign() <= 0 {
			fmt.Println("negative profit", profit)
			return
		}

		token := tokenMap[arbResult.SwapPaths.From]

		fmt.Printf("token:\t%s\nreal profit:\t%v\norigin profit:\t%v\namountIn:\t%v\n", token.Name, token.ToUnit(profit), token.ToUnit(arbResult.Profit), token.ToUnit(amountIn))

		fmt.Println("swap path:", arbResult.SwapPaths.String(tokenMap, poolMap))

		if ctx.Bool(flag.DryRunFlag.Name) {
			return
		}
	}

	if ctx.String(flag.OptionalPKFlag.Name) == "" {
		err = fmt.Errorf("pk required")
		return
	}
	pk, err := crypto.HexToECDSA(flag.OptionalPKFlag.Name)
	if err != nil {
		return
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return
	}

	opt, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return
	}
	dispatcher := txPkg.NewDispatcher([]*ethclient.Client{client}, nil, []*bind.TransactOpts{opt})
	dispatcher.Start()

	if ctx.String(flag.OptionalContractFlag.Name) == "" {
		err = fmt.Errorf("swapexecutor contract required")
		return
	}
	contract := common.HexToAddress(ctx.String(flag.OptionalContractFlag.Name))
	swapExecutorABI, err := abi.JSON(strings.NewReader(abi2.SwapExecutorMetaData.ABI))
	if err != nil {
		return
	}

	txData, err := defi.SwapArbConstructTxData(&swapExecutorABI, &arbResult, arbResult.Amount)
	if err != nil {
		return
	}
	dispatcher.Dispatch(&contract, 0, nil, txData, true)
	return
}

func swapVerify(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	router := common.HexToAddress(ctx.String(flag.RouterFlag.Name))
	tokens := toEthAddrList(strings.Split(ctx.String(flag.TokensFlag.Name), ","))

	account := common.HexToAddress(ctx.String(flag.AccountFlag.Name))

	deadline := time.Now().Unix() + 30

	args := []interface{}{big.NewInt(100000), big.NewInt(0), tokens, account, big.NewInt(deadline)}
	ab, _ := abi.JSON(strings.NewReader(homeabi.CherryRouterABI))
	method := ab.Methods["swapExactTokensForTokens"]
	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return
	}

	resultBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{From: account, To: &router, Data: append(method.ID, packed...)}, nil)
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

func swapFilterArbitraged(ctx *cli.Context) (err error) {
	err = swapFilterArbitragedImpl(ctx.String(flag.NetworkFlag.Name), ctx.String(flag.InFlag.Name), ctx.String(flag.OptionalOutFlag.Name))
	return
}

func swapRepairPoolsAndTokens(network string, pools []*defi.Pool, tokens []*defi.Token) (repairedPool []*defi.Pool, repairedTokens []*defi.Token, err error) {
	config, ok := arbConfig[network]
	if !ok {
		err = fmt.Errorf("invalid network:%s", network)
		return
	}
	var clients []*ethclient.Client
	for _, rpc := range config.Bot.Eth.RPCs {
		var client *ethclient.Client
		client, err = ethclient.Dial(rpc)
		if err != nil {
			return
		}
		clients = append(clients, client)
	}

	// handle nil reserve
	// TODO this is no longer needed if swapPairsImpl is perfected
	{
		nilPools := make(map[common.Address]int)

		for i, pool := range pools {
			poolReserves := *pool.Reserves.Load()
			if poolReserves == nil || poolReserves[0] == nil || poolReserves[1] == nil {
				nilPools[pool.Address] = i
				continue
			}
		}

		if len(nilPools) > 0 {
			fmt.Println("#nilPools", len(nilPools))

			nilPoolsList := make([]common.Address, 0, len(nilPools))
			for nilPool := range nilPools {
				nilPoolsList = append(nilPoolsList, nilPool)
			}
			pairAB, _ := abi.JSON(strings.NewReader(homeabi.CherryPairABI))
			var reserves []struct {
				Reserve0           *big.Int
				Reserve1           *big.Int
				BlockTimestampLast uint32
			}
			reserves, err = swapGetReserves(clients[0], &pairAB, nilPoolsList)
			if err != nil {
				return
			}

			for i, reserve := range reserves {
				poolAddr := nilPoolsList[i]
				if reserve.Reserve0 == nil || reserve.Reserve1 == nil {
					panic(fmt.Sprintf("reserve repair failed for pool:%v", poolAddr))
				}
				poolIdx, ok := nilPools[poolAddr]
				if !ok {
					panic("bug")
				}
				pool := pools[poolIdx]
				poolReserve := []*big.Int{reserve.Reserve0, reserve.Reserve1}
				pool.Reserves.Store(&poolReserve)
			}

			fmt.Println("nilPools repaired")
		}

	}

	tokenMap := make(map[common.Address]*defi.Token)
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}
	// handle missing token
	// TODO this is no longer needed if swapPairsImpl is perfected
	missingTokens := make(map[common.Address]bool)
	for _, pool := range pools {
		for _, token := range pool.Tokens {
			if tokenMap[token] == nil {
				missingTokens[token] = true
			}
		}
	}

	if len(missingTokens) > 0 {
		fmt.Println("#missing tokens", len(missingTokens))
		missingTokensList := make([]common.Address, 0, len(missingTokens))
		for token := range missingTokens {
			missingTokensList = append(missingTokensList, token)
		}
		tokenAB, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))
		var tokenInfos []*defi.Token
		tokenInfos, err = fetchTokens(clients, missingTokensList, &tokenAB)
		if err != nil {
			err = fmt.Errorf("failed to fetch missing tokens:%v", err)
			return
		}
		for _, token := range tokenInfos {
			if token == nil {
				panic("bug in fetchTokens")
			}

			tokens = append(tokens, token)
		}
		fmt.Println("missing tokens repaired")
	}

	repairedPool = pools
	repairedTokens = tokens
	return
}

func swapFilterArbitragedImpl(network, in, optionalOut string) (err error) {
	config, ok := arbConfig[network]
	if !ok {
		err = fmt.Errorf("invalid network:%s", network)
		return
	}

	arbMap := make(map[common.Address]bool)
	for _, arbToken := range config.Bot.Eth.Arb.Tokens {
		arbMap[arbToken.Token] = true
	}

	pools, tokens, err := swapReadPoolsAndTokens(in)
	if err != nil {
		return
	}

	pools, tokens, err = swapRepairPoolsAndTokens(network, pools, tokens)
	if err != nil {
		return
	}

	poolMap := make(map[common.Address]*defi.Pool)
	tokenMap := make(map[common.Address]*defi.Token)

	for _, pool := range pools {
		poolMap[pool.Address] = pool
	}
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}

	_, token2way, token3way := defi.SwapArbEssentials(poolMap, tokenMap, arbMap)

	usedPools := make(map[common.Address]bool)
	usedTokens := make(map[common.Address]bool)
	circles := 0
	handleFunc := func(tokenNway map[common.Address][]*defi.Swaps) {
		for _, swapPaths := range tokenNway {
			circles += len(swapPaths)
			for _, swapPath := range swapPaths {
				for _, swap := range swapPath.Swaps {
					usedPools[swap.Pool] = true
					usedTokens[swap.From] = true
				}
			}
		}
	}
	handleFunc(token2way)
	handleFunc(token3way)

	fmt.Println("#circles", circles)

	var arbPools []*defi.Pool
	var arbTokens []*defi.Token
	for _, pool := range pools {
		if !usedPools[pool.Address] {
			continue
		}
		arbPools = append(arbPools, pool)
	}
	for _, token := range tokens {
		if !usedTokens[token.Address] {
			continue
		}
		arbTokens = append(arbTokens, token)
	}

	if optionalOut == "" {
		optionalOut = in
	}
	err = swapOutputData(arbPools, arbTokens, optionalOut)
	return
}

func swapVerifyPoolImpl(network, in, optionalOut string) (err error) {
	config, ok := arbConfig[network]
	if !ok {
		err = fmt.Errorf("invalid network:%s", network)
		return
	}
	var clients []*ethclient.Client
	for _, rpc := range config.Bot.Eth.RPCs {
		var client *ethclient.Client
		client, err = ethclient.Dial(rpc)
		if err != nil {
			return
		}
		clients = append(clients, client)
	}

	pools, tokens, err := swapReadPoolsAndTokens(in)
	if err != nil {
		return
	}

	poolMap := make(map[common.Address]*defi.Pool)
	tokenMap := make(map[common.Address]*defi.Token)

	for _, pool := range pools {
		poolMap[pool.Address] = pool
	}
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}

	router := config.Bot.Eth.Verify.Router
	wtoken := config.Bot.Eth.Verify.Wtoken
	arbTokens := make([]common.Address, len(config.Bot.Eth.Arb.Tokens))
	for _, arbToken := range config.Bot.Eth.Arb.Tokens {
		arbTokens = append(arbTokens, arbToken.Token)
	}
	arbAccount := config.Bot.Eth.Verify.Account

	arbTokenMap := make(map[common.Address]bool)
	for _, token := range arbTokens {
		arbTokenMap[token] = true
	}
	balance, err := clients[0].BalanceAt(context.Background(), arbAccount, nil)
	if err != nil {
		return
	}
	arbValue := new(big.Int).Quo(balance, big.NewInt(2))

	_, token2way, token3way := defi.SwapArbEssentials(poolMap, tokenMap, arbTokenMap)

	ab, _ := abi.JSON(strings.NewReader(abi3.SwapVerifierMetaData.ABI))
	// pairAB, _ := abi.JSON(strings.NewReader(homeabi.CherryPairABI))

	var lock sync.Mutex
	validpairToken2Swaps := make(map[string][]*defi.Swaps)
	var (
		invalid, valid, reverted int64
	)

	var flatSwapPaths []*defi.Swaps
	for _, swapPaths := range token2way {
		flatSwapPaths = append(flatSwapPaths, swapPaths...)
	}
	for _, swapPaths := range token3way {
		flatSwapPaths = append(flatSwapPaths, swapPaths...)
	}
	start := time.Now()
	n := len(clients)
	unit := 5
	total := len(flatSwapPaths)

	parallel.All(context.Background(), total, unit, n, func(ctx context.Context, workerIdx, from, to int) (err error) {

		var (
			swapPathBatch             []*defi.Swaps
			arbAmountIn, amountOutMin []*big.Int
		)
		for i := from; i < to; i++ {
			swapPaths := flatSwapPaths[i]
			swapPathBatch = append(swapPathBatch, swapPaths)
			oneAmountIn := defi.SwapArbAmountIn(swapPaths, poolMap, nil)
			tokenInfo := tokenMap[swapPaths.From]
			oneToken := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tokenInfo.Decimals)), nil)
			if oneAmountIn == nil || oneAmountIn.Cmp(oneToken) < 0 {
				oneAmountIn = oneToken
			}
			arbAmountIn = append(arbAmountIn, oneAmountIn)
			amountOutMin = append(amountOutMin, big.NewInt(0))
		}

		result, err := swapVerifyPathImpl(clients[workerIdx%n], &ab, swapPathBatch, arbAmountIn, amountOutMin, arbAccount, router, wtoken, arbValue, nil)
		if err != nil {
			if strings.Contains(err.Error(), "CallContract:execution reverted") {
				err = nil
				atomic.AddInt64(&reverted, 1)
			} else if strings.Contains(err.Error(), "out of gas") {
				err = nil
				atomic.AddInt64(&reverted, 1)
			}
			return
		}

		for i, success := range result.Success {
			if success {
				atomic.AddInt64(&valid, 1)
				key := defi.PairTokenKey(swapPathBatch[i].Swaps[0].Pool, swapPathBatch[i].From)

				lock.Lock()
				validpairToken2Swaps[key] = append(validpairToken2Swaps[key], swapPathBatch[i])
				lock.Unlock()
			} else {
				atomic.AddInt64(&invalid, 1)

				errStr, err := defi.UnpackVerifyError(result.Reason[i])
				if err != nil {
					errStr = fmt.Sprintf("UnpackVerifyError failed:%v, reason:%s", err, string(result.Reason[i]))
				}

				fmt.Printf("swapVerifyPathImpl ng:[%s] paths:[%s]\n", errStr, swapPathBatch[i].String(tokenMap, poolMap))
			}
		}
		return
	}, func(from, to int) {
		if from > 0 && from%(unit*n) == 0 {
			fmt.Printf("dispatching %d - %d\n", from, to)
			took := time.Since(start)
			fmt.Println(
				"took", took, "eta", time.Duration(float64(took)*float64(total-int(from))/float64(from)),
				"#valid", atomic.LoadInt64(&valid), "#invalid", atomic.LoadInt64(&invalid), "#reverted", atomic.LoadInt64(&reverted),
			)
		}
	}, 3, time.Second*2)

	validPools := make(map[common.Address]*defi.Pool)
	validTokens := make(map[common.Address]*defi.Token)
	for _, swapPaths := range validpairToken2Swaps {
		for _, swapPath := range swapPaths {
			for _, swap := range swapPath.Swaps {
				validPools[swap.Pool] = poolMap[swap.Pool]
				validTokens[swap.From] = tokenMap[swap.From]
			}
		}
	}

	fmt.Println("#invalid", invalid, "#valid", valid, "#reverted", atomic.LoadInt64(&reverted), "#pool", len(validPools), "#token", len(validTokens))

	if optionalOut != "" {
		var outPools []*defi.Pool
		var outTokens []*defi.Token
		for _, pool := range validPools {
			outPools = append(outPools, pool)
		}
		for _, token := range validTokens {
			outTokens = append(outTokens, token)
		}

		err = swapOutputData(outPools, outTokens, optionalOut)
	}
	return
}
func swapVerifyPool(ctx *cli.Context) (err error) {
	err = swapVerifyPoolImpl(ctx.String(flag.NetworkFlag.Name), ctx.String(flag.InFlag.Name), ctx.String(flag.OptionalOutFlag.Name))
	return
}

func swapVerifyPath(ctx *cli.Context) (err error) {
	network := ctx.String(flag.NetworkFlag.Name)
	config, ok := arbConfig[network]
	if !ok {
		err = fmt.Errorf("invalid network:%s", network)
		return
	}

	var result defi.ArbResult
	resultBytes, err := os.ReadFile(ctx.String(flag.PathFlag.Name))
	if err != nil {
		return
	}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return
	}

	hops := len(result.SwapPaths.Swaps)
	if result.SwapPaths.From != result.SwapPaths.Swaps[hops-1].To {
		err = fmt.Errorf("not circular swaps")
		return
	}

	ab, _ := abi.JSON(strings.NewReader(abi3.SwapVerifierMetaData.ABI))

	client, err := ethclient.Dial(config.Bot.Eth.RPCs[0])
	if err != nil {
		return
	}
	wtoken := config.Bot.Eth.Verify.Wtoken
	router := config.Bot.Eth.Verify.Router
	account := config.Bot.Eth.Verify.Account
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return
	}
	value := new(big.Int).Quo(balance, big.NewInt(2))

	verifyResult, err := swapVerifyPathImpl(client, &ab, []*defi.Swaps{result.SwapPaths}, []*big.Int{result.Amount}, []*big.Int{result.Amount}, account, router, wtoken, value, nil)
	if err != nil {
		return
	}

	if verifyResult.Success[0] {
		fmt.Println("OK", "in", result.Amount, "out", new(big.Int).SetBytes(verifyResult.Reason[0]))
	} else {
		fmt.Println("NG", string(verifyResult.Reason[0]))

		pathsBytes, _ := json.Marshal(result.SwapPaths)
		fmt.Println(wtoken, router, string(pathsBytes))
	}

	return
}

func swapVerifyPathImpl(client *ethclient.Client, ab *abi.ABI, swapPathBatch []*defi.Swaps, arbAmountIn, amountOutMin []*big.Int, account, router, wtoken common.Address, value *big.Int, wrapPath []common.Address) (result defi.SwapVerifierResult, err error) {
	if ab == nil {
		var verifyABI abi.ABI
		verifyABI, err = abi.JSON(strings.NewReader(abi3.SwapVerifierMetaData.ABI))
		if err != nil {
			return
		}
		ab = &verifyABI
	}

	var (
		paths [][]abi2.Path
	)
	for _, swapPath := range swapPathBatch {
		paths = append(paths, swapPath.AbiPath())
	}

	method := ab.Constructor
	packed, err := method.Inputs.Pack(swapPathBatch[0].From, wtoken, router, wrapPath, arbAmountIn, amountOutMin, paths)
	if err != nil {
		panic(err)
	}
	if value == nil {
		var balance *big.Int
		balance, err = client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			return
		}
		value = new(big.Int).Quo(balance, big.NewInt(2))
	}
	if value.Uint64() < uint64(100*len(swapPathBatch)) {
		err = fmt.Errorf("value too low")
		panic(err)
	}

	resultBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{From: account, Value: value, Data: append(common.FromHex(abi3.SwapVerifierMetaData.Bin), packed...)}, nil)
	if err != nil {
		err = fmt.Errorf("CallContract:%v args:%s value:%x From:%s", err, hex.EncodeToString(packed), value, defi.FastAddrHex(account))
		return
	}

	err = defi.UnpackVerifierResult(resultBytes, &result)
	if err != nil {
		err = fmt.Errorf("UnpackVerifierResult:%v", err)
		return
	}

	for _, suc := range result.Success {
		if !suc {
			fmt.Printf("swapVerifyPathImpl !result.Success: args:%s value:%x From:%s\n", hex.EncodeToString(packed), value, defi.FastAddrHex(account))
			break
		}
	}

	return
}

func swapProfit(ctx *cli.Context) (err error) {
	var arbTokens []common.Address
	if ctx.String(flag.ArbFlag.Name) != "" {
		arbTokens = toEthAddrList(strings.Split(ctx.String(flag.ArbFlag.Name), ","))
	}

	err = swapProfitImpl(ctx.String(flag.InFlag.Name), arbTokens)
	return
}

func swapProfitImpl(in string, arbTokens []common.Address) (err error) {

	arbMap := make(map[common.Address]bool)
	for _, token := range arbTokens {
		arbMap[token] = true
	}

	pools, tokens, err := swapReadPoolsAndTokens(in)
	if err != nil {
		return
	}
	// prepare some basic data structures
	poolMap := make(map[common.Address]*defi.Pool)
	tokenMap := make(map[common.Address]*defi.Token)

	for _, pool := range pools {
		poolMap[pool.Address] = pool
	}

	for _, token := range tokens {
		tokenMap[token.Address] = token
	}

	var tokenWithPrice, tokenWithoutPrice int
	for _, token := range tokenMap {
		if token.Price == nil {
			tokenWithoutPrice++
		} else {
			tokenWithPrice++
		}
	}
	fmt.Println("tokenWithPrice", tokenWithPrice, "tokenWithoutPrice", tokenWithoutPrice)

	_, token2way, token3way := defi.SwapArbEssentials(poolMap, tokenMap, arbMap)
	{
		ok, wrongToken, wrongSwaps := defi.SwapArbValidateCircularSwaps(token2way)
		if !ok {
			panic(fmt.Sprintf("bug, token:%v #swappaths:%d, wrong swappath:%v", tokenMap[wrongToken].Name, len(token2way[wrongToken]), swapPathNames(wrongSwaps, tokenMap)))
		}
		ok, wrongToken, wrongSwaps = defi.SwapArbValidateCircularSwaps(token3way)
		if !ok {
			panic(fmt.Sprintf("bug, token:%v #swappaths:%d, wrong swappath:%v", tokenMap[wrongToken].Name, len(token3way[wrongToken]), swapPathNames(wrongSwaps, tokenMap)))
		}
	}

	totalProfitUSD := big.NewFloat(0)
	fmt.Println("#token2way", len(token2way), "#flat", swapShowSwapPaths(token2way, tokenMap, poolMap))
	start := time.Now()
	for tokenAddr, swapPaths := range token2way {
		tokenMaxProfit := big.NewFloat(0)
		for _, swapPath := range swapPaths {
			if tokenAddr != swapPath.From {
				panic("bug")
			}

			amountIn := defi.SwapArbAmountIn(swapPath, poolMap, nil)
			if amountIn == nil || amountIn.Sign() <= 0 {
				fmt.Println("negative amount", amountIn)
				continue
			}
			profit := swapPath.Profit(amountIn, poolMap, nil)
			if profit.Sign() <= 0 {
				fmt.Println("negative profit", profit)
				continue
			}

			profitUSD := showProfitSwapPath(amountIn, profit, swapPath, tokenMap, poolMap)
			if profitUSD != nil {
				if profitUSD.Cmp(tokenMaxProfit) > 0 {
					tokenMaxProfit = profitUSD
				}
			}
		}
		fmt.Println(tokenMap[tokenAddr].Name, "tokenMaxProfit", tokenMaxProfit)
		totalProfitUSD.Add(totalProfitUSD, tokenMaxProfit)
	}

	fmt.Println("token2way took", time.Since(start))

	fmt.Println("#token3way", len(token3way), "#flat", swapShowSwapPaths(token3way, tokenMap, poolMap))
	start = time.Now()
	for tokenAddr, swapPaths := range token3way {
		tokenMaxProfit := big.NewFloat(0)
		for _, swapPath := range swapPaths {
			if tokenAddr != swapPath.From {
				panic("bug")
			}

			amountIn := defi.SwapArbAmountIn(swapPath, poolMap, nil)
			if amountIn == nil || amountIn.Sign() <= 0 {
				fmt.Println("negative amount", amountIn)
				continue
			}
			profit := swapPath.Profit(amountIn, poolMap, nil)
			if profit.Sign() <= 0 {
				fmt.Println("negative profit", profit)
				continue
			}

			profitUSD := showProfitSwapPath(amountIn, profit, swapPath, tokenMap, poolMap)
			if profitUSD != nil {
				if profitUSD.Cmp(tokenMaxProfit) > 0 {
					tokenMaxProfit = profitUSD
				}
			}
		}
		totalProfitUSD.Add(totalProfitUSD, tokenMaxProfit)
	}
	fmt.Println("token3way took", time.Since(start))
	fmt.Println("totalProfit", totalProfitUSD)
	return
}

func swapPathNames(swapPaths *defi.Swaps, tokenMap map[common.Address]*defi.Token) (names []string) {

	names = make([]string, 0, len(swapPaths.Swaps)+1)
	names = append(names, tokenMap[swapPaths.From].Name)
	for _, swap := range swapPaths.Swaps {
		toToken := tokenMap[swap.To]
		names = append(names, toToken.Name)
	}

	return
}

var chances uint64

const bigChanceDir = "cmd/data/big_chance"

func showProfitSwapPath(amount, profit *big.Int, swapPaths *defi.Swaps, tokenMap map[common.Address]*defi.Token, poolMap map[common.Address]*defi.Pool) (profitUsd *big.Float) {
	token := tokenMap[swapPaths.From]
	profitUsd = token.Value(profit, nil)

	names := swapPathNames(swapPaths, tokenMap)
	mark := ""
	if profitUsd.Cmp(big.NewFloat(10)) > 0 {
		mark = "chance"
	}
	if profitUsd.Cmp(big.NewFloat(100)) > 0 {
		arbResult := defi.ArbResult{Profit: profit, Amount: amount, SwapPaths: swapPaths}
		arbResultBytes, _ := json.Marshal(arbResult)
		chances++
		os.MkdirAll(bigChanceDir, 0600)
		os.WriteFile(fmt.Sprintf("%s/%d.json", bigChanceDir, chances), arbResultBytes, defaultPermission)
	}
	profitPercent := fmt.Sprintf("%v%%", big.NewInt(0).Quo(big.NewInt(0).Mul(big.NewInt(100), profit), amount))
	fmt.Println(mark, "profit token", token.Name, "symbol", token.Symbol, "amountInUnit", token.ToUnit(amount), "profitPercent", profitPercent, "profitUnit", token.ToUnit(profit), "profitUSD", profitUsd, "addr", token.Address, "path", names)

	return
}

func swapTokens(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	var invokes []multicall.Invoke
	err = json.Unmarshal([]byte(ctx.String(flag.InvokesFlag.Name)), &invokes)
	if err != nil {
		return
	}
	result := make([]interface{}, len(invokes))
	ab, _ := abi.JSON(strings.NewReader(homeabi.ICherryERC20ABI))

	_, err = multicall.Do(context.Background(), client, &ab, invokes, result)
	if err != nil {
		for _, invoke := range invokes {
			oneResult := make([]interface{}, 1)
			_, err := multicall.Do(context.Background(), client, &ab, []multicall.Invoke{invoke}, oneResult)
			if err != nil {
				fmt.Printf("invoke failed for contract:%v name:%s, err:%v\n", invoke.Contract, invoke.Name, err)
				continue
			}
		}
		return
	}

	fmt.Println(result)
	return
}

func swapFee(ctx *cli.Context) (err error) {
	numbers := strings.Split(ctx.String(flag.NumbersFlag.Name), ",")

	if len(numbers) != 5 {
		err = fmt.Errorf("need to specify amountIn, amountOut, reverse0, reserve1, dir")
		return
	}
	var dir bool
	if numbers[4] == "1" {
		dir = true
	}
	numbers = numbers[0:4]

	var bigInts []*big.Int
	for _, number := range numbers {
		n, ok := new(big.Int).SetString(number, 10)
		if !ok {
			panic(fmt.Sprintf("invalid number:%v", number))
		}
		bigInts = append(bigInts, n)
	}

	fmt.Println("fee", defi.CalcFee(bigInts[0], bigInts[1], bigInts[2], bigInts[3], dir))
	return
}

func swapBiSwapFee(ctx *cli.Context) (err error) {
	pair := common.HexToAddress(ctx.String(flag.PairFlag.Name))

	client, err := ethclient.Dial(arbConfig["bsc"].Bot.Eth.RPCs[0])
	if err != nil {
		return
	}

	ab, err := homeabi.ParseFunctionAsABI("function swapFee() returns (uint32 fee)")
	if err != nil {
		return
	}
	invokes := []multicall.Invoke{{Contract: pair, Name: "swapFee", Args: []interface{}{}}}

	result := make([]uint32, len(invokes))
	_, err = multicall.Do(context.Background(), client, &ab, invokes, result)
	if err != nil {
		return
	}

	fmt.Println(result[0])
	return
}

func swapMdexFee(ctx *cli.Context) (err error) {
	config, ok := arbConfig[ctx.String(flag.NetworkFlag.Name)]
	if !ok {
		err = fmt.Errorf("invalid network:%s", ctx.String(flag.NetworkFlag.Name))
		return
	}
	pair := common.HexToAddress(ctx.String(flag.PairFlag.Name))

	client, err := ethclient.Dial(config.Bot.Eth.RPCs[0])
	if err != nil {
		return
	}

	_, factorys, err := swapFactoryImpl(config.Bot.Eth.RPCs[0], []common.Address{pair})
	if err != nil {
		return
	}

	factory := factorys[0]

	ab, err := homeabi.ParseFunctionAsABI("function getPairFees(address pair) returns (uint256 fee);")
	if err != nil {
		return
	}

	invokes := []multicall.Invoke{{Contract: factory, Name: "getPairFees", Args: []interface{}{pair}}}

	result := make([]*big.Int, len(invokes))
	_, err = multicall.Do(context.Background(), client, &ab, invokes, result)
	if err != nil {
		return
	}

	fmt.Println(result[0])
	return
}

func swapOut(ctx *cli.Context) (err error) {
	numbers := strings.Split(ctx.String(flag.NumbersFlag.Name), ",")
	if len(numbers) != 4 {
		err = fmt.Errorf("need to specify amountIn, reserveIn, reserveOut, fee")
		return
	}

	var bigInts []*big.Int
	for _, number := range numbers {
		n, ok := new(big.Int).SetString(number, 10)
		if !ok {
			panic(fmt.Sprintf("invalid number:%v", number))
		}
		bigInts = append(bigInts, n)
	}

	fmt.Println(defi.GetAmountOut(bigInts[0], bigInts[1], bigInts[2], bigInts[3].Uint64()))
	return
}
func swapBlocks(ctx *cli.Context) (err error) {
	from := ctx.Uint64(flag.HeightFlag.Name)
	n := uint64(50)
	filter := ethereum.FilterQuery{
		Topics:    [][]common.Hash{{common.HexToHash(eth.SYNC)}},
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(from + n)),
	}

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}
	logs, err := client.FilterLogs(context.Background(), filter)
	if err != nil {
		err = fmt.Errorf("FilterLogs:%v", err)
		return
	}

	swapPerHeight := make(map[uint64]int)
	for _, log := range logs {
		swapPerHeight[log.BlockNumber]++
	}
	for h := from; h <= from+n; h++ {
		if swapPerHeight[h] == 0 {
			fmt.Printf("------------------------no swap @%d\n", h)
		} else {
			fmt.Printf("%d@%d\n", swapPerHeight[h], h)
		}
	}

	return
}

func swapEstimateDeploySwapExecutor(ctx *cli.Context) (err error) {
	client, err := ethclient.Dial("https://bsc-dataseed1.binance.org/")
	if err != nil {
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}

	fmt.Println("gasPrice", gasPrice)

	swapExecutorABI, _ := abi.JSON(strings.NewReader(abi2.SwapExecutorMetaData.ABI))
	method := swapExecutorABI.Constructor

	packed, err := method.Inputs.Pack(common.Address{})
	if err != nil {
		return
	}

	callMsg := ethereum.CallMsg{
		From: common.HexToAddress("0x0000000000000000000000000000000000001004"), Gas: 0, GasPrice: gasPrice,
		Data: append(common.FromHex(abi2.SwapExecutorMetaData.Bin), packed...),
	}

	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return
	}

	fmt.Println("gasLimit", gasLimit)

	product := new(big.Float).SetInt(new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	fmt.Println(new(big.Float).Quo(product, new(big.Float).SetInt(unit)))

	return
}

func swapDeploySwapExecutor(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return
	}

	pkBytes, err := os.ReadFile(ctx.String(flag.PKFileFlag.Name))
	if err != nil {
		return
	}

	pk, err := crypto.HexToECDSA(string(pkBytes))
	if err != nil {
		return
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return
	}

	addr, tx, _, err := swap_executor.DeploySwapExecutor(transactor, client, common.HexToAddress(ctx.String(flag.AccountFlag.Name)))
	if err != nil {
		return
	}

	fmt.Println("contract", addr)

	fmt.Println("tx", tx.Hash())
	return
}

func findBadContracts(client *ethclient.Client, ab *abi.ABI, invokes []multicall.Invoke) (contracts map[common.Address]bool) {
	contracts = make(map[common.Address]bool)
	for _, invoke := range invokes {
		oneResult := make([]interface{}, 1)
		_, err := multicall.Do(context.Background(), client, ab, []multicall.Invoke{invoke}, oneResult)
		if err != nil {
			contracts[invoke.Contract] = true
		}
	}
	return
}

func swapPercent(ctx *cli.Context) (err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}
	stashDir := root + "/cmd/data/stash"
	pairsStash := stashDir + "/pairs.json"
	pairTokensStash := stashDir + "/pairTokens.json"

	var pairs [][]common.Address

	pairsBytes, err := os.ReadFile(pairsStash)
	if err != nil {
		return
	}
	err = json.Unmarshal(pairsBytes, &pairs)
	if err != nil {
		return
	}

	pairTokensBytes, err := os.ReadFile(pairTokensStash)
	if err != nil {
		return
	}
	pairTokens := make([][]struct {
		Token0 common.Address
		Token1 common.Address
	}, 0, len(pairs))
	err = json.Unmarshal(pairTokensBytes, &pairTokens)
	if err != nil {
		return
	}

	wbnb := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	usdt := common.HexToAddress("0x55d398326f99059fF775485246999027B3197955") // usdt
	dai := common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3")  // dai
	usdc := common.HexToAddress("0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d") // usdc
	busd := common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56") // busd

	var bnbPairs, usdtPairs, daiPairs, usdcPairs, busdPairs, total int
	for _, factoryPairTokens := range pairTokens {
		total += len(factoryPairTokens)
		for _, tokens := range factoryPairTokens {
			if tokens.Token0 == wbnb || tokens.Token1 == wbnb {
				bnbPairs++
			}
			if tokens.Token0 == usdt || tokens.Token1 == usdt {
				usdtPairs++
			}
			if tokens.Token0 == dai || tokens.Token1 == dai {
				daiPairs++
			}
			if tokens.Token0 == usdc || tokens.Token1 == usdc {
				usdcPairs++
			}
			if tokens.Token0 == busd || tokens.Token1 == busd {
				busdPairs++
			}
		}

	}

	fmt.Println("total", total)
	fmt.Println("bnbPairs", bnbPairs, "percent", float32(bnbPairs)/float32(total))
	fmt.Println("usdtPairs", usdtPairs, "percent", float32(usdtPairs)/float32(total))
	fmt.Println("daiPairs", daiPairs, "percent", float32(daiPairs)/float32(total))
	fmt.Println("usdcPairs", usdcPairs, "percent", float32(usdcPairs)/float32(total))
	fmt.Println("busdPairs", busdPairs, "percent", float32(busdPairs)/float32(total))
	return
}

func swapShowSwapPaths(tokenSwapPaths map[common.Address][]*defi.Swaps, tokenMap map[common.Address]*defi.Token, poolMap map[common.Address]*defi.Pool) (total int) {
	for tokenAddr, swapPaths := range tokenSwapPaths {
		token, ok := tokenMap[tokenAddr]
		if !ok {
			panic("bug")
		}
		fmt.Println("token", token.Name, "symbol", token.Symbol, "#swapPath", len(swapPaths), "addr", token.Address)
		total += len(swapPaths)
		for _, swapPath := range swapPaths {
			fmt.Printf("\t[%s]\n", swapPath.String(tokenMap, poolMap))
		}
	}

	return
}

const defaultThreshold = 10000

func swapReadPoolsAndTokens(in string) (pools []*defi.Pool, tokens []*defi.Token, err error) {
	inBytes, err := os.ReadFile(in)
	if err != nil {
		return
	}
	data := []interface{}{&pools, &tokens}
	err = json.Unmarshal(inBytes, &data)
	if err != nil {
		return
	}

	return
}

func swapOutputData(poolList []*defi.Pool, tokenList []*defi.Token, out string) (err error) {
	var outBytes []byte
	outBytes, err = json.Marshal([]interface{}{poolList, tokenList})
	if err != nil {
		return
	}

	fmt.Println("output #pools", len(poolList), "#tokens", len(tokenList))
	err = os.WriteFile(out, outBytes, defaultPermission)
	if err != nil {
		return
	}
	return
}

func swapTVLImpl(in string, stables []common.Address, threshold uint, out string) (err error) {
	if threshold <= 0 {
		threshold = defaultThreshold
	}
	if out == "" {
		out = in
	}

	pools, tokens, err := swapReadPoolsAndTokens(in)
	if err != nil {
		return
	}

	fmt.Println("#pools", len(pools), "#tokens", len(tokens), "#stables", len(stables))

	// TVL start
	{
		// prepare some maps

		tokenMap := make(map[common.Address]*defi.Token)
		for _, token := range tokens {
			tokenMap[token.Address] = token
		}

		var validPools []*defi.Pool
		for _, pool := range pools {
			if tokenMap[pool.Tokens[0]] == nil || tokenMap[pool.Tokens[1]] == nil {
				continue
			}

			validPools = append(validPools, pool)
		}
		if len(pools) > len(validPools) {
			fmt.Printf("filtered %d pools with no token info\n", len(pools)-len(validPools))
		}
		pools = validPools

		token2pool := make(map[common.Address]map[common.Address]bool)
		poolMap := make(map[common.Address]*defi.Pool)
		todo := make(map[common.Address]bool)
		for _, pool := range pools {
			todo[pool.Address] = true
			poolMap[pool.Address] = pool
			for _, tokenAddr := range pool.Tokens {
				tokenPools, ok := token2pool[tokenAddr]
				if !ok {
					tokenPools = make(map[common.Address]bool)
					token2pool[tokenAddr] = tokenPools
				}
				tokenPools[pool.Address] = true
			}
		}

		// setup
		currentPricedTokens := make(map[common.Address]*big.Float)
		for _, stable := range stables {
			currentPricedTokens[stable] = big.NewFloat(1) // unit price
		}
		updateTokenPrice := func(pricedTokens map[common.Address]*big.Float) {
			for tokenAddr, price := range pricedTokens {
				tokenMap[tokenAddr].Price = price
			}
		}
		updateTokenPrice(currentPricedTokens)

		// core logic
		for {
			// loop until no new price discovery
			if len(currentPricedTokens) == 0 {
				break
			}
			nextPricedTokens := make(map[common.Address]*big.Float)
			for tokenAddr, unitPrice := range currentPricedTokens {
				tokenPools := token2pool[tokenAddr]
				token, ok := tokenMap[tokenAddr]
				if !ok {
					panic(fmt.Sprintf("token %v not found", tokenAddr))
				}
				for poolAddr := range tokenPools {
					if todo[poolAddr] {
						pool := poolMap[poolAddr]
						tokenReserve := pool.Reserve(tokenAddr)

						value := token.Value(tokenReserve, unitPrice)
						floatValue, _ := value.Float32()
						pool.TVL = uint(2 * floatValue)
						if pool.TVL > threshold {
							otherPrice := pool.OtherPrice(tokenAddr, unitPrice, tokenMap)
							otherAddr := pool.OtherAddr(tokenAddr)
							if _, ok := currentPricedTokens[otherAddr]; !ok {
								nextPricedTokens[otherAddr] = otherPrice
							}
						}
						delete(todo, poolAddr)
					}
				}
			}

			currentPricedTokens = nextPricedTokens
			updateTokenPrice(currentPricedTokens)
		}

		// filter useless pools and tokens
		var finalPools []*defi.Pool
		finalTokenMap := map[common.Address]bool{}
		for _, pool := range pools {
			if pool.TVL > 0 {
				finalPools = append(finalPools, pool)
				finalTokenMap[pool.Tokens[0]] = true
				finalTokenMap[pool.Tokens[1]] = true
			}
		}
		pools = finalPools
		var finalTokens []*defi.Token
		for token := range finalTokenMap {
			finalTokens = append(finalTokens, tokenMap[token])
		}
		tokens = finalTokens
	}
	// TVL end

	fmt.Println("final #pools", len(pools), "#tokens", len(tokens))
	err = swapOutputData(pools, tokens, out)

	return
}

func unitReserve(decimals uint8) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}

func swapFactory(ctx *cli.Context) (err error) {

	config, ok := arbConfig[ctx.String(flag.NetworkFlag.Name)]
	if !ok {
		err = fmt.Errorf("invalid network:%s", ctx.String(flag.NetworkFlag.Name))
		return
	}

	routerOrPair := strings.Split(ctx.String(flag.ContractFlag.Name), ",")

	var routerAddrs []common.Address
	for _, router := range routerOrPair {
		routerAddrs = append(routerAddrs, common.HexToAddress(router))
	}
	height, factorys, err := swapFactoryImpl(config.Bot.Eth.RPCs[0], routerAddrs)
	if err != nil {
		return
	}

	fmt.Println("height", height, "factory", factorys)

	return
}

func swapFactoryImpl(rpc string, routers []common.Address) (height uint64, factorys []common.Address, err error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}

	var invokes []multicall.Invoke
	for _, router := range routers {
		invokes = append(invokes, multicall.Invoke{
			Contract: router,
			Name:     "factory",
			Args:     []interface{}{},
		})
	}

	factorys = make([]common.Address, len(invokes))

	ab, _ := abi.JSON(strings.NewReader(homeabi.CherryRouterABI))
	height, err = multicall.Do(context.Background(), client, &ab, invokes, factorys)
	if err != nil {
		return
	}

	return
}

func swapExactTokensForTokensTxData(ctx *cli.Context) (err error) {
	ab, _ := abi.JSON(strings.NewReader(homeabi.CherryRouterABI))
	amountIn, ok := big.NewInt(0).SetString(ctx.String(flag.AmountInFlag.Name), 10)
	if !ok {
		err = fmt.Errorf("invalid amountIn:%s", ctx.String(flag.AmountInFlag.Name))
		return
	}
	amountOutMin, ok := big.NewInt(0).SetString(ctx.String(flag.AmountOutFlag.Name), 10)
	if !ok {
		err = fmt.Errorf("invalid amountOut:%s", ctx.String(flag.AmountOutFlag.Name))
		return
	}
	paths := toEthAddrList(strings.Split(ctx.String(flag.PathFlag.Name), ","))
	// 2分钟内有效
	deadline := big.NewInt(time.Now().Unix() + 60*5)
	to := common.HexToAddress(ctx.String(flag.ToFlag.Name))

	txData, err := ab.Methods["swapExactTokensForTokens"].Inputs.Pack(amountIn, amountOutMin, paths, to, deadline)
	if err != nil {
		return
	}

	fmt.Println(hex.EncodeToString(txData))
	return
}

type PairTokens struct {
	TokenA common.Address
	TokenB common.Address
}

func swapGetPairs(client *ethclient.Client, factory common.Address, pairs []PairTokens) (pools []common.Address, err error) {
	invokes := make([]multicall.Invoke, 0, len(pairs))
	for _, pairTokens := range pairs {
		invokes = append(invokes, multicall.Invoke{
			Contract: factory,
			Name:     "getPair",
			Args:     []interface{}{pairTokens.TokenA, pairTokens.TokenB},
		})
	}

	factoryAB, _ := abi.JSON(strings.NewReader(homeabi.CherryFactoryABI))
	pools = make([]common.Address, len(invokes))
	_, err = multicall.Do(context.Background(), client, &factoryAB, invokes, pools)
	return
}

func fixedFeeSingleDexPathToSwapPath(client *ethclient.Client, router *defi.Router, paths []common.Address, fee uint64) (swaps *defi.Swaps, err error) {
	pairs := make([]PairTokens, 0, len(paths)-1)
	for i := 0; i < len(paths)-1; i++ {
		pairs = append(pairs, PairTokens{TokenA: paths[i], TokenB: paths[i+1]})
	}
	pools, err := swapGetPairs(client, router.Factory, pairs)
	if err != nil {
		return
	}
	swaps = &defi.Swaps{From: paths[0]}

	for i := 0; i < len(pools); i++ {
		dir := bytes.Compare(paths[i][:], paths[i+1][:]) < 0
		swaps.Swaps = append(swaps.Swaps, &defi.Swap{Exchange: router.Exchange, Pool: pools[i], Fee: fee, From: paths[i], To: paths[i+1], Dir: dir})
	}
	return
}

func swapExactTokensForTokensDryRun(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	config, ok := arbConfig[ctx.String(flag.ForkFlag.Name)]
	if !ok {
		err = fmt.Errorf("unknown fork:%s", ctx.String(flag.ForkFlag.Name))
		return
	}
	account := config.Bot.Eth.Verify.Account
	router := config.Bot.Eth.Verify.Router
	wtoken := config.Bot.Eth.Verify.Wtoken

	var targetRouter *defi.Router
	for i := range config.Bot.Eth.Pool.Routers {
		routerInfo := config.Bot.Eth.Pool.Routers[i]
		if routerInfo.Router == router {
			targetRouter = routerInfo
		}
	}
	if targetRouter == nil {
		err = fmt.Errorf("router not found in config.Bot.Eth.Pool.Routers")
		return
	}

	paths := toEthAddrList(strings.Split(ctx.String(flag.PathFlag.Name), ","))
	swapPath, err := fixedFeeSingleDexPathToSwapPath(client, targetRouter, paths, uint64(ctx.Int(flag.FeeFlag.Name)))
	if err != nil {
		return
	}

	amountIn, ok := big.NewInt(0).SetString(ctx.String(flag.AmountInFlag.Name), 10)
	if !ok {
		err = fmt.Errorf("invalid amountIn:%s", ctx.String(flag.AmountInFlag.Name))
		return
	}

	var wrapPath []common.Address
	if ctx.String(flag.WrapPathFlag.Name) != "" {
		wrapPath = toEthAddrList(strings.Split(ctx.String(flag.WrapPathFlag.Name), ","))
	}
	result, err := swapVerifyPathImpl(client, nil, []*defi.Swaps{swapPath}, []*big.Int{amountIn}, []*big.Int{big.NewInt(0)}, account, router, wtoken, nil, wrapPath)
	if err != nil {
		return
	}

	out := new(big.Int).SetBytes(result.Reason[0])

	fmt.Println(out)
	return
}

func swapBnbPrice(ctx *cli.Context) (err error) {
	config, _ := arbConfig["bsc"]

	client, err := ethclient.Dial(config.Bot.Eth.RPCs[0])
	if err != nil {
		return
	}

	pairAB, _ := abi.JSON(strings.NewReader(homeabi.CherryPairABI))

	reserves, err := swapGetReserves(client, &pairAB, []common.Address{common.HexToAddress("0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16")})
	if err != nil {
		return
	}

	reserve := reserves[0]

	arbValue := new(big.Int).Quo(reserve.Reserve1, reserve.Reserve0)
	fmt.Println("price", arbValue)
	return
}
