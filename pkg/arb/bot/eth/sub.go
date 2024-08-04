package eth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	zlog "github.com/rs/zerolog/log"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/arbbot/pkg/metrics"
	common2 "github.com/zhiqiangxu/litenode/eth/common"
	"github.com/zhiqiangxu/util"
	"github.com/zhiqiangxu/util/parallel"
)

func (b *Bot) initSub() (err error) {

	_, syncArguments, err := homeabi.ParseEvent("event Sync(uint112 reserve0, uint112 reserve1)")
	if err != nil {
		return
	}
	b.syncArguments = syncArguments

	// cagtch up first

	height, err := b.fulls[0].BlockNumber(context.Background())
	if err != nil {
		return
	}
	zlog.Info().Int("b.height", int(b.height)).Int("height diff", int(height-b.height)).Msg("catching up...")
	for b.height+maxContinuousEmptyBlocks < height {
		b.handleNewBlock(height)
	}
	zlog.Info().Int("b.height", int(b.height)).Int("height", int(height)).Msg("catch up ok...")

	if b.fullSub != nil {
		util.GoFunc(&b.wg, b.subNewBlockFull)
	} else {
		util.GoFunc(&b.wg, b.subNewBlockLite)
	}

	util.GoFunc(&b.wg, b.subNewPendingTx)
	return
}

const maxContinuousEmptyBlocks = 10

func (b *Bot) subNewBlockFull() {
	chainHeadCh := make(chan *types.Header)
	sub, err := b.fullSub.SubscribeNewHead(b.ctx, chainHeadCh)
	if err != nil {
		zlog.Fatal().Err(err).Msg("SubscribeNewHead")
	}
	defer sub.Unsubscribe()

	gate := make(chan struct{}, 1)

	for {
		select {
		case nextHeader, ok := <-chainHeadCh:
			if !ok {
				select {
				case <-b.ctx.Done():
				default:
					zlog.Warn().Msg("chainHeadCh closed, re-subNewBlockFull...")
					time.Sleep(time.Second)
					util.GoFunc(&b.wg, b.subNewBlockFull)
				}

				return
			}

			metrics.SubNewBlockMeter.Mark(1)

			start := time.Now()
			oldHeight := b.height
			b.handleNewBlock(nextHeader.Number.Uint64())
			for b.height+maxContinuousEmptyBlocks < nextHeader.Number.Uint64() {
				zlog.Warn().Msgf("FilterLogs is lagging behind SubscribeChainHeadEvent, %d vs %d, catching up...", b.height, nextHeader.Number)
				b.handleNewBlock(nextHeader.Number.Uint64())
			}
			if b.height > oldHeight {
				atomic.StoreInt32(&b.sentInBlock, 0)
			}

			metrics.LatencyMetric.With("method", "handle_nextEvent").Observe(float64(time.Since(start).Nanoseconds()))
			metrics.HeightMetric.Update(int64(b.height))

			select {
			case gate <- struct{}{}:
				go b.checkpoint(gate)
			default:
			}

		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Bot) subNewBlockLite() {

	chainHeadCh := make(chan common2.ChainHeadEvent)
	sub := b.lite.Eth.SubscribeChainHeadEvent(chainHeadCh)
	defer sub.Unsubscribe()

	gate := make(chan struct{}, 1)

	for {
		select {
		case nextEvent, ok := <-chainHeadCh:
			if !ok {
				return
			}

			metrics.SubNewBlockMeter.Mark(1)

			start := time.Now()
			oldHeight := b.height
			b.handleNewBlock(nextEvent.Number)
			for b.height+maxContinuousEmptyBlocks < nextEvent.Number {
				zlog.Warn().Msgf("FilterLogs is lagging behind SubscribeChainHeadEvent, %d vs %d, catching up...", b.height, nextEvent.Number)
				b.handleNewBlock(nextEvent.Number)
			}
			if b.height > oldHeight {
				atomic.StoreInt32(&b.sentInBlock, 0)
			}

			metrics.LatencyMetric.With("method", "handle_nextEvent").Observe(float64(time.Since(start).Nanoseconds()))
			metrics.HeightMetric.Update(int64(b.height))

			select {
			case gate <- struct{}{}:
				go b.checkpoint(gate)
			default:
			}

		case <-b.ctx.Done():
			return
		}
	}

}

func (b *Bot) subNewPendingTx() {
	txCh := make(chan core.NewTxsEvent)
	sub := b.lite.Eth.SubscribeNewTxsEvent(txCh)
	defer sub.Unsubscribe()

	for {
		select {
		case event, ok := <-txCh:
			if !ok {
				return
			}
			util.GoFunc(&b.wg, func() {

				metrics.ArbitrageTxMeter.Mark(1)

				b.arbitrageTx(event.Txs)
			})
		case <-b.ctx.Done():
			return
		}
	}
}

const SYNC = "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"

var backOffDelay = time.Second

func (b *Bot) handleNewBlock(height uint64) {

	if height <= b.height {
		return
	}

	start := time.Now()
	toBlock := b.height + 100
	if toBlock > height {
		toBlock = height
	}
	filter := ethereum.FilterQuery{
		Topics:    [][]common.Hash{{common.HexToHash(SYNC)}},
		FromBlock: big.NewInt(int64(b.height + 1)),
		ToBlock:   big.NewInt(int64(toBlock)),
	}

	logs, idx, err := parallel.First(b.ctx, len(b.fulls), func(ctx context.Context, i int) (logs []types.Log, err error) {
		logs, err = b.fulls[i].FilterLogs(ctx, filter)
		if err != nil {
			zlog.Warn().Uint64("height", height).Int("i", i).Err(err).Msg("FilterLogs")
			return
		}
		return
	}, backOffDelay)
	if err != nil {
		zlog.Warn().Uint64("height", height).Err(err).Msg("handleNewBlock err")
		return
	}

	metrics.LatencyMetric.With("method", "FilterLogs").Observe(float64(time.Since(start).Nanoseconds()))

	var poolUpdate = make(map[common.Address][]*big.Int)
	if len(logs) == 0 {
		zlog.Warn().Uint64("height", height).Int("idx", idx).Msg("no logs")
		return
	}

	zlog.Info().Int("#logs", len(logs)).Uint64("height", height).Uint64("start_height", logs[0].BlockNumber).Uint64("end_height", logs[len(logs)-1].BlockNumber).Int("idx", idx).Msg("logs available")

	for _, log := range logs {
		// only subscribe to interested pools
		if b.pools[log.Address] == nil {
			continue
		}

		values, err := b.syncArguments.UnpackValues(log.Data)
		if err != nil {
			zlog.Warn().Uint64("height", height).Str("data", hex.EncodeToString(log.Data)).Err(err).Msg("syncArguments.UnpackValues")
			time.Sleep(backOffDelay)
			return
		}

		poolUpdate[log.Address] = []*big.Int{values[0].(*big.Int), values[1].(*big.Int)}
		height := log.BlockNumber
		if height > b.height {
			b.height = height
		}
	}
	if len(poolUpdate) > 0 {
		for pool := range poolUpdate {
			reserves := poolUpdate[pool]
			if reserves[0] == nil || reserves[1] == nil {
				panic(fmt.Sprintf("nil reserve, pool:%s", defi.FastAddrHex(pool)))
			}
			b.pools[pool].Reserves.Store(&reserves)
		}
	}

}

func (b *Bot) arbitrageTx(txs []*types.Transaction) {
	fmt.Println("#txs", len(txs))

	for i := range txs {
		tx := txs[i]
		go func() {
			if tx.To() == nil {
				return
			}
			router := b.Pool.routerMap[*tx.To()]
			if router == nil {
				return
			}

			start := time.Now()
			ok := false
			defer func() {
				if ok {
					metrics.LatencyMetric.With("method", "arbitrageTxOK").Observe(float64(time.Since(start).Nanoseconds()))
				} else {
					metrics.LatencyMetric.With("method", "arbitrageTxNG").Observe(float64(time.Since(start).Nanoseconds()))
				}
			}()
			if !b.txCache.Add(tx.Hash(), true, 0) {
				return
			}
			cacheAdd := time.Now()

			if b.Gas.MaxPrice != nil {
				if tx.GasPrice().Cmp(b.Gas.MaxPrice) > 0 {
					zlog.Info().Uint64("GasPrice", tx.GasPrice().Uint64()).Msg("HighGas")
					return
				}
			}
			if b.Gas.MinPrice != nil {
				if tx.GasPrice().Cmp(b.Gas.MinPrice) < 0 {
					zlog.Info().Uint64("GasPrice", tx.GasPrice().Uint64()).Msg("LowGas")
					return
				}
			}

			methodID, params, err := defi.SwapArbGetParams(tx)
			if err != nil {
				zlog.Warn().Str("hash", tx.Hash().String()).Err(err).Msg("SwapArbGetParams")
				return
			}

			arbCtx := defi.SwapArgGetContext(router, tx, methodID, params)

			swaps := arbCtx.ChooseBranch(b.pairToken2Swaps)
			if len(swaps) == 0 {
				zlog.Info().Str("tx_hash", tx.Hash().String()).Msg("NoSwapPaths")
				return
			}

			updatedReserves := arbCtx.CalcUpdatedReserves(b.pools, b.pair2Pool)
			arbCtx.UpdatedReserves = updatedReserves

			resultCh := make(chan *defi.ArbResult, len(swaps))
			var wg sync.WaitGroup
			wg.Add(len(swaps))
			for i := range swaps {
				swapPath := swaps[i]
				go func() {
					defer wg.Done()

					amountIn := defi.SwapArbAmountIn(swapPath, b.pools, updatedReserves)
					profit := swapPath.Profit(amountIn, b.pools, updatedReserves)
					if profit.Sign() < 0 {
						resultCh <- nil
						return
					}

					resultCh <- &defi.ArbResult{Profit: profit, SwapPaths: swapPath, Amount: amountIn}

				}()
			}

			wg.Wait()

			maxProfit := big.NewInt(0)
			var result *defi.ArbResult
			for arbResult := range resultCh {
				if arbResult == nil {
					continue
				}
				if arbResult.Profit.Cmp(maxProfit) > 0 {
					maxProfit = arbResult.Profit
					result = arbResult
				}
			}

			newReserves := updatedReserves[result.SwapPaths.Swaps[0].Pool]
			oldReserves := *b.pools[result.SwapPaths.Swaps[0].Pool].Reserves.Load()
			zlog.Info().Str("token", defi.FastAddrHex(result.SwapPaths.From)).Str("old", oldReserves[0].String()).Str("new", newReserves[0].String()).Msg("ArbTokenDelta")
			ok = b.execArb(&arbCtx, result)

			zlog.Info().Str("tx_hash", tx.Hash().String()).Bool("ok", ok).Dur("took", time.Since(start)).Dur("cacheAdd took", cacheAdd.Sub(start)).Msg("arbitrageTx")
		}()

	}
}

var chances int64

const bigChanceDir = "cmd/data/big_chance"

func (b *Bot) checkpoint(gate <-chan struct{}) {
	start := time.Now()
	empty := true
	var (
		pathChecked        int
		dispatchWaiting    time.Duration
		sendInBlockLimited bool
	)
	defer func() {

		metrics.LatencyMetric.With("method", "checkpoint").Observe(float64(time.Since(start).Nanoseconds()))

		zlog.Info().Bool("empty", empty).Int("pathChecked", pathChecked).Bool("sendInBlockLimited", sendInBlockLimited).Dur("took", time.Since(start)).Int64("dispatchWaitingNS", dispatchWaiting.Nanoseconds()).Msg("checkpoint done")
		<-gate
	}()
	for _, swapPaths := range b.pairToken2Swaps {
		if len(swapPaths) == 0 {
			continue
		}
		arbToken := b.Arb.tokenMap[swapPaths[0].From]

		for _, swapPath := range swapPaths {
			pathChecked++
			amountIn := defi.SwapArbAmountIn(swapPath, b.pools, nil)
			if amountIn == nil || amountIn.Sign() <= 0 {
				continue
			}
			profit := swapPath.Profit(amountIn, b.pools, nil)
			profitUint64, _ := arbToken.detail.Value(profit, arbToken.UnitPrice).Uint64()
			profitInt := int(profitUint64)
			if profitInt > b.Arb.MinProfitUSD {
				arbResult := defi.ArbResult{Profit: profit, SwapPaths: swapPath, Amount: amountIn}
				err := b.verifySwapPaths(&arbResult, nil)
				if err != nil {
					zlog.Info().Err(err).Msg("checkpoint.verifySwapPaths")
					continue
				}

				saveBigChance(&arbResult)

				sentInBlock := atomic.AddInt32(&b.sentInBlock, 1)
				if sentInBlock > b.Arb.MaxTxInOneBlock {
					atomic.CompareAndSwapInt32(&b.sentInBlock, sentInBlock, sentInBlock-1)
					zlog.Info().Msg("checkpoint: tx limit reached")
					sendInBlockLimited = true
					return
				}

				txData, err := defi.SwapArbConstructTxData(b.swapExecutorABI, &arbResult, zero /* zero for tx */)
				if err != nil {
					zlog.Warn().Err(err).Msg("SwapArbConstructTxData")
					return
				}

				empty = false
				zlog.Info().Int("profit", profitInt).Str("arbtoken", defi.FastAddrHex(arbToken.Token)).Uint64("amountIn", amountIn.Uint64()).Msg("checkpoint: start dispatch")
				beforeDispatch := time.Now()
				b.txDispather.Dispatch(&b.Arb.SwapExecutorAddr, b.Gas.MaxGas, b.Gas.MaxPrice, txData, false)
				dispatchWaiting = time.Since(beforeDispatch)
				return
			}
		}
	}

}

func saveBigChance(arbResult *defi.ArbResult) {

	// save big chance
	idx := atomic.AddInt64(&chances, 1)
	arbResultBytes, _ := json.Marshal(arbResult)
	err := os.WriteFile(fmt.Sprintf("%s/%d.json", bigChanceDir, idx), arbResultBytes, 0777)
	if err != nil {
		panic(err)
	}

}

func init() {
	os.MkdirAll(bigChanceDir, 0777)
}
