package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	zlog "github.com/rs/zerolog/log"
	abi2 "github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
	abi3 "github.com/zhiqiangxu/arbbot/contracts/abi/swap_verifier"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/arbbot/pkg/metrics"
	"github.com/zhiqiangxu/util/parallel"
)

func (b *Bot) initVerify() (err error) {
	swapVerifierABI, err := abi.JSON(strings.NewReader(abi3.SwapVerifierMetaData.ABI))
	if err != nil {
		return
	}
	b.swapVerifierABI = &swapVerifierABI

	balance, err := b.fulls[0].BalanceAt(context.Background(), b.Verify.Account, nil)
	if err != nil {
		return
	}

	if balance.Uint64() < 100 {
		err = fmt.Errorf("verify account balance too low:%d", balance.Uint64())
		return
	}
	b.Verify.value = new(big.Int).Quo(balance, big.NewInt(2))
	return
}

var zeroBig = big.NewInt(0)

func (b *Bot) verifySwapPaths(arbResult *defi.ArbResult, updatedReserves map[common.Address][]*big.Int) (err error) {
	start := time.Now()
	defer func() {
		metrics.LatencyMetric.With("method", "verifySwapPaths").Observe(float64(time.Since(start).Nanoseconds()))
	}()

	var nilWrapPath []common.Address
	packed, err := b.swapVerifierABI.Constructor.Inputs.Pack(arbResult.SwapPaths.From, b.Verify.Wtoken, b.Verify.Router, nilWrapPath, []*big.Int{arbResult.Amount}, []*big.Int{zeroBig}, [][]abi2.Path{arbResult.SwapPaths.AbiPath()})
	if err != nil {
		return
	}

	resultBytes, idx, err := parallel.First(b.ctx, len(b.fulls), func(ctx context.Context, i int) (resultBytes []byte, err error) {
		resultBytes, err = b.fulls[i].CallContract(ctx, ethereum.CallMsg{From: b.Verify.Account, Value: b.Verify.value, Data: append(common.FromHex(abi3.SwapVerifierMetaData.Bin), packed...)}, nil)
		if err != nil {
			zlog.Warn().Int("i", i).Err(err).Msg("verifySwapPaths.CallContract")
			return
		}
		return
	}, backOffDelay)
	if err != nil {
		zlog.Warn().Err(err).Msg("verifySwapPaths err")
		return
	}

	var result defi.SwapVerifierResult
	err = defi.UnpackVerifierResult(resultBytes, &result)
	if err != nil {
		err = fmt.Errorf("UnpackVerifierResult failed:%v idx:%d", err, idx)
		return
	}

	if !result.Success[0] {
		var errStr string
		errStr, err = defi.UnpackVerifyError(result.Reason[0])
		if err != nil {
			err = fmt.Errorf("UnpackVerifyError failed:%v, reason:%s idx:%d", err, string(result.Reason[0]), idx)
			return
		}

		err = fmt.Errorf(
			"SwapVerifier failed:%s, constructor params:%s value:%x arbAmount:%d arbToken:%s path:[%s]",
			errStr,
			hex.EncodeToString(packed),
			b.Verify.value,
			arbResult.Amount,
			defi.FastAddrHex(arbResult.SwapPaths.From),
			arbResult.SwapPaths.String(b.tokens, b.pools),
		)
		return
	}

	out := new(big.Int).SetBytes(result.Reason[0])
	if out.Cmp(arbResult.Amount) < 0 {
		err = fmt.Errorf("OutLtIn, out:%v in:%v diff:%v", out, arbResult.Amount, new(big.Int).Sub(arbResult.Amount, out))
		return
	}

	fmt.Println("verifySwapPaths pass", "in", arbResult.Amount, "out", out, "idx", idx)
	return
}
