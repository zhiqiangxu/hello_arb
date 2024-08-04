package clients

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zlog "github.com/rs/zerolog/log"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/arbbot/pkg/sequential"
	"github.com/zhiqiangxu/multicall"
)

type Sequential []*ethclient.Client

var cd = time.Second * 2

func (c Sequential) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {
	lastAccess := make([]time.Time, len(c))
	block, _, err = sequential.MustDo(ctx, len(c), func(ctx context.Context, i int) (block *types.Block, err error) {
		if duration := time.Since(lastAccess[i]); duration <= cd {
			time.Sleep(duration)
		}
		lastAccess[i] = time.Now()
		block, err = c[i].BlockByNumber(ctx, number)
		if err != nil {
			zlog.Warn().Uint64("height", number.Uint64()).Int("idx", i).Err(err).Msg("BlockByNumber")
		}
		return
	})
	return
}

func (c Sequential) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	lastAccess := make([]time.Time, len(c))
	receipt, _, err = sequential.MustDo(ctx, len(c), func(ctx context.Context, i int) (receipt *types.Receipt, err error) {
		if duration := time.Since(lastAccess[i]); duration <= cd {
			time.Sleep(duration)
		}
		lastAccess[i] = time.Now()
		receipt, err = c[i].TransactionReceipt(ctx, txHash)
		if err != nil {
			zlog.Warn().Str("tx_hash", txHash.String()).Int("idx", i).Err(err).Msg("TransactionReceipt")
		}
		return
	})
	return
}

type PoolTokens struct {
	Token0 common.Address
	Token1 common.Address
}

func (c Sequential) PoolTokens(ctx context.Context, pool common.Address) (tokens PoolTokens, err error) {
	lastAccess := make([]time.Time, len(c))

	tokens, _, err = sequential.MustDo(ctx, len(c), func(ctx context.Context, i int) (tokens PoolTokens, err error) {
		if duration := time.Since(lastAccess[i]); duration <= cd {
			time.Sleep(duration)
		}
		lastAccess[i] = time.Now()

		invokes := []multicall.Invoke{
			{
				Contract: pool,
				Name:     "token0",
				Args:     []interface{}{},
			},
			{
				Contract: pool,
				Name:     "token1",
				Args:     []interface{}{},
			},
		}
		result := make([]common.Address, len(invokes))
		_, err = multicall.Do(ctx, c[i], &homeabi.PairABI, invokes, result)
		if err != nil {
			zlog.Warn().Str("pool", defi.FastAddrHex(pool)).Int("idx", i).Err(err).Msg("PoolTokens")
			return
		}

		tokens = PoolTokens{Token0: result[0], Token1: result[1]}
		return
	})

	return
}

func (c Sequential) FilterLogs(ctx context.Context, filter ethereum.FilterQuery) (logs []types.Log, err error) {
	lastAccess := make([]time.Time, len(c))

	logs, _, err = sequential.MustDo(ctx, len(c), func(ctx context.Context, i int) (logs []types.Log, err error) {
		if duration := time.Since(lastAccess[i]); duration <= cd {
			time.Sleep(duration)
		}
		lastAccess[i] = time.Now()

		logs, err = c[i].FilterLogs(ctx, filter)
		if err != nil {
			zlog.Warn().Int("idx", i).Err(err).Msg("FilterLogs")
		}
		return
	})
	return
}
