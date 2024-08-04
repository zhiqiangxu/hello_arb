package eth

import (
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	zlog "github.com/rs/zerolog/log"
	abi2 "github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
)

func (b *Bot) initArb() (err error) {
	// arbTokenMap
	if len(b.Arb.Tokens) == 0 {
		err = fmt.Errorf("empty arb tokens")
		return
	}
	arbTokenMap := make(map[common.Address]*ArbTokenInfo)
	for _, arbToken := range b.Arb.Tokens {
		detail := b.tokens[arbToken.Token]
		if detail == nil {
			err = fmt.Errorf("%v missing from tokenMap", arbToken)
			return
		}
		arbToken.detail = detail
		arbTokenMap[arbToken.Token] = arbToken
	}
	b.Arb.tokenMap = arbTokenMap

	// populate pairToken2Swaps
	pairToken2Swaps := make(map[string][]*defi.Swaps)
	_, token2way, token3way := defi.SwapArbEssentials(b.pools, b.tokens, arbTokenMap)
	for token, swapPaths := range token2way {
		for _, swapPath := range swapPaths {
			key := defi.PairTokenKey(swapPath.Swaps[0].Pool, token)
			pairToken2Swaps[key] = append(pairToken2Swaps[key], swapPath)
		}
	}
	for token, swapPaths := range token3way {
		for _, swapPath := range swapPaths {
			key := defi.PairTokenKey(swapPath.Swaps[0].Pool, token)
			pairToken2Swaps[key] = append(pairToken2Swaps[key], swapPath)
		}
	}
	{
		for _, swapPaths := range pairToken2Swaps {
			dedup := make(map[*defi.Swaps]bool)
			for _, swapPath := range swapPaths {
				if dedup[swapPath] {
					err = fmt.Errorf("dup swapPath found")
					return
				}
				dedup[swapPath] = true
			}
		}
	}
	b.pairToken2Swaps = pairToken2Swaps
	b.trimPoolAndTokens()

	if b.Arb.SwapExecutorAddr == (common.Address{}) {
		err = fmt.Errorf("SwapExecutorAddr empty")
		return
	}

	swapExecutorABI, err := abi.JSON(strings.NewReader(abi2.SwapExecutorMetaData.ABI))
	if err != nil {
		return
	}
	b.swapExecutorABI = &swapExecutorABI

	return

}

func (b *Bot) trimPoolAndTokens() {
	usedPools := make(map[common.Address]bool)
	usedTokens := make(map[common.Address]bool)
	for _, swapPaths := range b.pairToken2Swaps {
		for _, swapPath := range swapPaths {
			for _, swap := range swapPath.Swaps {
				usedPools[swap.Pool] = true
				usedTokens[swap.From] = true
			}
		}
	}

	pools := make(map[common.Address]*defi.Pool)
	tokens := make(map[common.Address]*defi.Token)
	for addr, pool := range b.pools {
		if usedPools[addr] {
			pools[addr] = pool
		}
	}
	for addr, token := range b.tokens {
		if usedTokens[addr] {
			tokens[addr] = token
		}
	}
	b.pools = pools
	b.tokens = tokens

	// simple verification and populate pair2Pool
	pair2Pool := make(map[string]common.Address)
	for _, pool := range pools {
		if len(pool.Tokens) != 2 {
			panic("invalid #pool.Tokens")
		}
		if tokens[pool.Tokens[0]] == nil {
			panic(fmt.Sprintf("token %v missing", pool.Tokens[0]))
		}
		if tokens[pool.Tokens[1]] == nil {
			panic(fmt.Sprintf("token %v missing", pool.Tokens[1]))
		}

		reserves := *pool.Reserves.Load()
		if reserves[0] == nil || reserves[1] == nil {
			panic(fmt.Sprintf("nil reserve, pool:%s", defi.FastAddrHex(pool.Address)))
		}

		key := defi.PairKey(pool.Factory, pool.Tokens[0], pool.Tokens[1])
		pair2Pool[key] = pool.Address

	}
	b.pair2Pool = pair2Pool
}

var zero = big.NewInt(0)

func (b *Bot) execArb(arbCtx *defi.ArbContext, result *defi.ArbResult) (ok bool) {
	arbToken := b.Arb.tokenMap[result.SwapPaths.From]
	profitUSD, _ := arbToken.detail.Value(result.Amount, arbToken.UnitPrice).Float32()
	profitUSDInt := int(profitUSD)
	if profitUSDInt < b.Arb.MinProfitUSD {
		zlog.Info().Int("profitUSD", profitUSDInt).Msg("low profit")
		return
	}
	err := b.verifySwapPaths(result, arbCtx.UpdatedReserves)
	if err != nil {
		zlog.Info().Err(err).Msg("execArb.verifySwapPaths")
		return
	}
	saveBigChance(result)

	if atomic.LoadInt32(&b.sentInBlock) >= b.Arb.MaxTxInOneBlock {
		zlog.Info().Msg("tx limit reached")
		return
	}

	sentInBlock := atomic.AddInt32(&b.sentInBlock, 1)
	if sentInBlock > b.Arb.MaxTxInOneBlock {
		atomic.CompareAndSwapInt32(&b.sentInBlock, sentInBlock, sentInBlock-1)
		zlog.Info().Msg("tx limit reached")
		return
	}

	txData, err := defi.SwapArbConstructTxData(b.swapExecutorABI, result, zero /*zero for tx*/)
	if err != nil {
		zlog.Warn().Str("tx_hash", arbCtx.Tx.Hash().String()).Err(err).Msg("SwapArbConstructTxData")
		return
	}

	beforeDispatch := time.Now()
	b.txDispather.Dispatch(&b.Arb.SwapExecutorAddr, b.Gas.MaxGas, arbCtx.Tx.GasPrice(), txData, false)
	zlog.Info().Str("tx_hash", arbCtx.Tx.Hash().String()).Dur("took", time.Since(beforeDispatch)).Msg("execArb dispatch")

	ok = true
	return
}

var arbAmountIn, amountOutMin []*big.Int

func init() {
	arbAmountIn = append(arbAmountIn, big.NewInt(100))
	amountOutMin = append(amountOutMin, big.NewInt(0))
}
