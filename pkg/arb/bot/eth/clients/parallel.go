package clients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zlog "github.com/rs/zerolog/log"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/multicall"
	"github.com/zhiqiangxu/util/parallel"
)

type Parallel []*ethclient.Client

func (c Parallel) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {

	block, _, err = parallel.First(ctx, len(c), func(ctx context.Context, i int) (block *types.Block, err error) {
		block, err = c[i].BlockByNumber(ctx, number)
		if err != nil {
			zlog.Warn().Uint64("height", number.Uint64()).Int("idx", i).Err(err).Msg("BlockByNumber")
		}
		return
	}, cd)
	return
}

func (c Parallel) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {

	receipt, _, err = parallel.First(ctx, len(c), func(ctx context.Context, i int) (receipt *types.Receipt, err error) {
		receipt, err = c[i].TransactionReceipt(ctx, txHash)
		if err != nil {
			zlog.Warn().Str("tx_hash", txHash.String()).Int("idx", i).Err(err).Msg("TransactionReceipt")
		}
		return
	}, cd)
	return
}

func (c Parallel) PoolTokens(ctx context.Context, pool common.Address) (tokens PoolTokens, err error) {
	tokens, _, err = parallel.First(ctx, len(c), func(ctx context.Context, i int) (tokens PoolTokens, err error) {

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
	}, cd)

	return
}

func (c Parallel) FilterLogs(ctx context.Context, filter ethereum.FilterQuery) (logs []types.Log, err error) {

	logs, _, err = parallel.First(ctx, len(c), func(ctx context.Context, i int) (logs []types.Log, err error) {

		logs, err = c[i].FilterLogs(ctx, filter)
		if err != nil {
			zlog.Warn().Int("idx", i).Err(err).Msg("FilterLogs")
		}
		return
	}, cd)
	return
}
