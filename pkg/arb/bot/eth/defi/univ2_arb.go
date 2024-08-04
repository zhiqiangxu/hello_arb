package defi

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	zlog "github.com/rs/zerolog/log"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
)

func SwapArbEssentials[T any](
	pools map[common.Address]*Pool,
	tokens map[common.Address]*Token,
	arbMap map[common.Address]T) (
	ttps map[common.Address]map[common.Address]map[common.Address]bool,
	token2way,
	token3way map[common.Address][]*Swaps,
) {
	// prepare some basic data structures
	// token -> token -> [pools]
	ttps = make(map[common.Address]map[common.Address]map[common.Address]bool)
	token2way = make(map[common.Address][]*Swaps)
	token3way = make(map[common.Address][]*Swaps)

	for _, pool := range pools {
		token0 := pool.Tokens[0]
		token1 := pool.Tokens[1]
		if ttps[token0] == nil {
			ttps[token0] = make(map[common.Address]map[common.Address]bool)
		}
		if ttps[token1] == nil {
			ttps[token1] = make(map[common.Address]map[common.Address]bool)
		}
		if ttps[token0][token1] == nil {
			ttps[token0][token1] = make(map[common.Address]bool)
		}
		if ttps[token1][token0] == nil {
			ttps[token1][token0] = make(map[common.Address]bool)
		}
		ttps[token0][token1][pool.Address] = true
		ttps[token1][token0][pool.Address] = true
	}

	// find all 2-way and 3-way circular pairs
	for token0 := range ttps {
		if len(arbMap) > 0 {
			if _, ok := arbMap[token0]; !ok {
				continue
			}
		}
		for token1, pools1 := range ttps[token0] {
			for poolID1 := range pools1 {
				pool1, ok := pools[poolID1]
				if !ok {
					panic("bug")
				}
				swap1 := &Swap{Exchange: pool1.Exchange, Pool: poolID1, From: token0, To: token1, Fee: pool1.SwapFee, Dir: pool1.Token0() == token0}
				for token2, pools2 := range ttps[token1] {
					for poolID2 := range pools2 {
						if poolID2 == poolID1 {
							continue
						}

						pool2, ok := pools[poolID2]
						if !ok {
							panic("bug")
						}
						swap2 := &Swap{Exchange: pool2.Exchange, Pool: poolID2, From: token1, To: token2, Fee: pool2.SwapFee, Dir: pool2.Token0() == token1}

						if token0 == token2 {
							token2way[token0] = append(token2way[token0], NewSwapPath(token0, []*Swap{swap1, swap2}, pools))
						} else {
							for poolID3 := range ttps[token2][token0] {
								pool3, ok := pools[poolID3]
								if !ok {
									panic("bug")
								}
								swap3 := &Swap{Exchange: pool3.Exchange, Pool: poolID3, From: token2, To: token0, Fee: pool3.SwapFee, Dir: pool3.Token0() == token2}
								token3way[token0] = append(token3way[token0], NewSwapPath(token0, []*Swap{swap1, swap2, swap3}, pools))
							}
						}
					}
				}
			}
		}
	}
	return
}

func SwapArbAmountIn(swapPath *Swaps, poolMap map[common.Address]*Pool, updatedPoolReserves map[common.Address][]*big.Int) *big.Int {
	reserves, ratios := swapPath.ReservesAndRatios(poolMap, updatedPoolReserves)

	return SwapArbAmountInRaw(reserves, ratios)
}

func SwapArbAmountInRaw(reserves [][2]float64, ratios []float64) *big.Int {
	switch len(ratios) {
	case 2:
		a := math.Sqrt(reserves[0][1]*reserves[0][0]*reserves[1][1]*reserves[1][0]*ratios[0]*ratios[1]) - reserves[0][0]*reserves[1][0]
		if a <= 0 {
			return nil
		}

		b := ratios[0] * (reserves[1][0] + reserves[0][1]*ratios[1])

		result, _ := big.NewFloat(a / b).Int(nil)
		return result
	case 3:
		a := math.Sqrt(reserves[0][0]*reserves[0][1]*reserves[1][0]*reserves[1][1]*reserves[2][0]*reserves[2][1]*ratios[0]*ratios[1]*ratios[2]) - reserves[0][0]*reserves[1][0]*reserves[2][0]
		if a <= 0 {
			return nil
		}

		b := ratios[0] * (reserves[1][0]*reserves[2][0] + reserves[0][1]*ratios[1]*(reserves[2][0]+reserves[1][1]*ratios[2]))
		result, _ := big.NewFloat(a / b).Int(nil)
		return result
	default:
		panic(fmt.Sprintf("unexpected swapPath length:%v", len(ratios)))
	}
}

func SwapArbValidateCircularSwaps(tokenNway map[common.Address][]*Swaps) (result bool, wrongToken common.Address, wrongSwaps *Swaps) {
	for token, swaps := range tokenNway {
		for len(swaps) > 0 {
			target := -1
			for i := 1; i < len(swaps); i++ {
				if swaps[0].IsReverse(swaps[i]) {
					target = i
					break
				}
			}
			if target == -1 {
				wrongToken = token
				wrongSwaps = swaps[0]
				return
			}
			if len(swaps) == 2 {
				break
			}
			swaps[target] = swaps[1]
			swaps = swaps[2:]
		}
	}

	result = true
	return
}

func SwapArbGetParams(tx *types.Transaction) (methodID string, params map[string]interface{}, err error) {
	// ab.Methods
	data := tx.Data()
	if len(data) <= 4 {
		err = fmt.Errorf("no method")
		return
	}
	methodID = hex.EncodeToString(data[:4])
	args, ok := SwapMethods[methodID]
	if !ok {
		err = fmt.Errorf("not target method")
		return
	}
	params = make(map[string]interface{})
	err = args.UnpackIntoMap(params, data[4:])
	if err != nil {
		return
	}

	return
}

func SwapArgGetContext(router *Router, tx *types.Transaction, methodID string, params map[string]interface{}) (ctx ArbContext) {
	ctx = ArbContext{Tx: tx, Router: router, MethodID: methodID}

	switch methodID {
	case swapExactETHForTokensID, swapExactBNBForTokensID,
		swapExactETHForTokensSupportingFeeOnTransferTokensID, swapExactBNBForTokensSupportingFeeOnTransferTokensID:
		ctx.AmountIn = tx.Value()
		ctx.AmountOut = params["amountOutMin"].(*big.Int)
		ctx.ExactIn = true
	case swapExactTokensForTokensID, swapExactTokensForETHID, swapExactTokensForBNBID,
		swapExactTokensForTokensSupportingFeeOnTransferTokensID, swapExactTokensForETHSupportingFeeOnTransferTokensID, swapExactTokensForBNBSupportingFeeOnTransferTokensID:
		ctx.AmountIn = params["amountIn"].(*big.Int)
		ctx.AmountOut = params["amountOutMin"].(*big.Int)
		ctx.ExactIn = true
	case swapETHForExactTokensID, swapBNBForExactTokensID:
		ctx.AmountIn = tx.Value()
		ctx.AmountOut = params["amountOut"].(*big.Int)
		ctx.ExactIn = false
	case swapTokensForExactTokensID, swapTokensForExactETHID, swapTokensForExactBNBID:
		ctx.AmountIn = params["amountInMax"].(*big.Int)
		ctx.AmountOut = params["amountOut"].(*big.Int)
		ctx.ExactIn = false
	default:
		panic(fmt.Sprintf("unkown methodID:%s", methodID))
	}

	return
}

type ArbContext struct {
	Tx              *types.Transaction
	Router          *Router
	MethodID        string
	AmountIn        *big.Int
	AmountOut       *big.Int
	ExactIn         bool
	Paths           []common.Address
	Pools           []*Pool
	UpdatedReserves map[common.Address][]*big.Int
}

func (ctx *ArbContext) CalcUpdatedReserves(poolMap map[common.Address]*Pool, pair2PoolMap map[string]common.Address) (updatedReserves map[common.Address][]*big.Int) {

	var (
		reserves [][2]*big.Int
		fees     []uint64
		pools    []*Pool
	)
	for i := 0; i < len(ctx.Paths)-1; i++ {
		poolAddr, ok := pair2PoolMap[PairKey(ctx.Router.Factory, ctx.Paths[i], ctx.Paths[i+1])]
		if !ok {
			return
		}
		pool := poolMap[poolAddr]
		poolReserves := *pool.Reserves.Load()
		pools = append(pools, pool)
		fees = append(fees, pool.SwapFee)
		if ctx.Paths[i] == pool.Tokens[0] {
			reserves = append(reserves, [2]*big.Int{poolReserves[0], poolReserves[1]})
		} else {
			reserves = append(reserves, [2]*big.Int{poolReserves[1], poolReserves[0]})
		}
	}

	var amounts []*big.Int

	if ctx.ExactIn {

		amounts = GetAmountsOut(ctx.AmountIn, reserves, fees)
		amountOut := amounts[len(amounts)-1]
		if ctx.AmountOut.Cmp(amountOut) > 0 {
			zlog.Debug().Str("factory", FastAddrHex(ctx.Router.Factory)).
				Strs("paths", FastAddrsHex(ctx.Paths)).
				Uint64("amountOutMin", ctx.AmountOut.Uint64()).
				Uint64("amountOut", amountOut.Uint64()).
				Msg("CalcUpdatedReserves: amountOut too small")
			return
		}

	} else {

		amounts = GetAmountsIn(ctx.AmountOut, reserves, fees)
		amountIn := amounts[0]
		if ctx.AmountIn.Cmp(amountIn) < 0 {
			zlog.Debug().Str("factory", FastAddrHex(ctx.Router.Factory)).
				Strs("paths", FastAddrsHex(ctx.Paths)).
				Uint64("amountInMax", ctx.AmountIn.Uint64()).
				Uint64("amountIn", amountIn.Uint64()).
				Msg("CalcReserve: amountIn too big")
			return
		}
	}

	updatedReserves = make(map[common.Address][]*big.Int)

	for i, pool := range pools {
		poolReserves := *pool.Reserves.Load()
		tokenIn := ctx.Paths[i]
		reserves := make([]*big.Int, 2)
		if tokenIn == pool.Tokens[0] {
			reserves[0] = new(big.Int).Add(poolReserves[0], amounts[i])
			reserves[1] = new(big.Int).Sub(poolReserves[1], amounts[i+1])
			if reserves[1].Sign() < 0 {
				zlog.Warn().Msgf("CalcReserve: factory:%v paths:%v Neg reserves %v", ctx.Router.Factory, ctx.Paths, reserves)
				return nil
			}
		} else {
			reserves[0] = new(big.Int).Sub(poolReserves[0], amounts[i+1])
			reserves[1] = new(big.Int).Add(poolReserves[1], amounts[i])
			if reserves[0].Sign() < 0 {
				zlog.Warn().Msgf("CalcReserve: factory:%v paths:%v Neg reserves %v", ctx.Router.Factory, ctx.Paths, reserves)
				return nil
			}
		}

		updatedReserves[pool.Address] = reserves
	}

	// for reuse in ChooseBranch
	ctx.Pools = pools
	return
}

func (ctx *ArbContext) ChooseBranch(pairToken2Swaps map[string][]*Swaps) (swaps []*Swaps) {
	for i, pool := range ctx.Pools {
		key := PairTokenKey(pool.Address, ctx.Paths[i+1])
		tmpSwaps := pairToken2Swaps[key]
		if len(tmpSwaps) == 0 {
			continue
		}
		if len(swaps) == 0 {
			swaps = tmpSwaps
			continue
		}
		if len(tmpSwaps) < len(swaps) {
			swaps = tmpSwaps
			continue
		}
	}
	return
}

func SwapArbConstructTxData(swapExecutorABI *ethabi.ABI, result *ArbResult, amountOutMin *big.Int) (txData []byte, err error) {

	paths := result.SwapPaths.AbiPath()

	// txData, err = swapExecutorABI.Pack("swap", result.SwapPaths.From, result.Amount, amountOutMin, paths)

	txData, err = swapExecutorABI.Pack("swapFlashloan", result.SwapPaths.From, result.Amount, amountOutMin, paths)

	return
}

type ArbResult struct {
	Profit    *big.Int
	SwapPaths *Swaps
	Amount    *big.Int
}

var (
	SwapMethods    map[string]ethabi.Arguments
	SwapFunctions  []string
	SwapMethodMids map[string]string
)

var (
	swapExactTokensForTokensID, swapTokensForExactTokensID, swapExactETHForTokensID, swapTokensForExactETHID, swapExactTokensForETHID, swapETHForExactTokensID,
	swapExactBNBForTokensID, swapTokensForExactBNBID, swapExactTokensForBNBID, swapBNBForExactTokensID, swapExactTokensForTokensSupportingFeeOnTransferTokensID, swapExactETHForTokensSupportingFeeOnTransferTokensID, swapExactTokensForETHSupportingFeeOnTransferTokensID,
	swapExactBNBForTokensSupportingFeeOnTransferTokensID, swapExactTokensForBNBSupportingFeeOnTransferTokensID string
)

func init() {
	SwapMethodMids = map[string]string{}
	SwapFunctions = []string{
		"swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapTokensForExactTokens( uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)",
		"swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)",
		"swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline)",

		"swapExactBNBForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapTokensForExactBNB(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)",
		"swapExactTokensForBNB(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapBNBForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline)",
		"swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn,uint256 amountOutMin, address[] path, address to, uint256 deadline)",

		"swapExactBNBForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)",
		"swapExactTokensForBNBSupportingFeeOnTransferTokens(uint256 amountIn,uint256 amountOutMin, address[] path, address to, uint256 deadline)",
	}

	SwapMethods = make(map[string]ethabi.Arguments)
	for _, swapFunction := range SwapFunctions {
		sig, in, _, err := homeabi.ParseFunction(swapFunction)
		if err != nil {
			panic(fmt.Errorf("err:%v func:%s", err, swapFunction))
		}
		mid := hex.EncodeToString(homeabi.SigToMid(sig))
		SwapMethods[mid] = in
		method := homeabi.SigToMethod(sig)
		SwapMethodMids[method] = mid
	}

	var ok bool
	swapExactTokensForTokensID, ok = SwapMethodMids["swapExactTokensForTokens"]
	assert(ok)
	swapTokensForExactTokensID, ok = SwapMethodMids["swapTokensForExactTokens"]
	assert(ok)
	swapExactETHForTokensID, ok = SwapMethodMids["swapExactETHForTokens"]
	assert(ok)
	swapTokensForExactETHID, ok = SwapMethodMids["swapTokensForExactETH"]
	assert(ok)
	swapExactTokensForETHID, ok = SwapMethodMids["swapExactTokensForETH"]
	assert(ok)
	swapETHForExactTokensID, ok = SwapMethodMids["swapETHForExactTokens"]
	assert(ok)

	swapExactBNBForTokensID, ok = SwapMethodMids["swapExactBNBForTokens"]
	assert(ok)
	swapTokensForExactBNBID, ok = SwapMethodMids["swapTokensForExactBNB"]
	assert(ok)
	swapExactTokensForBNBID, ok = SwapMethodMids["swapExactTokensForBNB"]
	assert(ok)
	swapBNBForExactTokensID, ok = SwapMethodMids["swapBNBForExactTokens"]
	assert(ok)
	swapExactTokensForTokensSupportingFeeOnTransferTokensID, ok = SwapMethodMids["swapExactTokensForTokensSupportingFeeOnTransferTokens"]
	assert(ok)
	swapExactETHForTokensSupportingFeeOnTransferTokensID, ok = SwapMethodMids["swapExactETHForTokensSupportingFeeOnTransferTokens"]
	assert(ok)
	swapExactTokensForETHSupportingFeeOnTransferTokensID, ok = SwapMethodMids["swapExactTokensForETHSupportingFeeOnTransferTokens"]
	assert(ok)

	swapExactBNBForTokensSupportingFeeOnTransferTokensID, ok = SwapMethodMids["swapExactBNBForTokensSupportingFeeOnTransferTokens"]
	assert(ok)
	swapExactTokensForBNBSupportingFeeOnTransferTokensID, ok = SwapMethodMids["swapExactTokensForBNBSupportingFeeOnTransferTokens"]
	assert(ok)

	// fmt.Println(SwapMethodMids)
}

func assert(ok bool) {
	if !ok {
		panic("ng")
	}
}
