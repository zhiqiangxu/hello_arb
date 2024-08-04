package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	zlog "github.com/rs/zerolog/log"
	homeabi "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/multicall"
)

func (b *Bot) initPools() (err error) {
	// load pools and tokens
	inBytes, err := os.ReadFile(b.Pool.File)
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

	poolMap := make(map[common.Address]*defi.Pool)
	for _, pool := range pools {
		poolMap[pool.Address] = pool
	}
	tokenMap := make(map[common.Address]*defi.Token)
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}

	b.pools = poolMap
	b.tokens = tokenMap

	// router map
	if len(b.Pool.Routers) == 0 {
		err = fmt.Errorf("empty routers")
		return
	}
	routerMap := make(map[common.Address]*defi.Router)
	for _, router := range b.Pool.Routers {
		routerMap[router.Router] = router
	}
	b.Pool.routerMap = routerMap

	return
}

func (b *Bot) initReserves() (err error) {

	zlog.Info().Msg("initReserves start")

	start := time.Now()
	defer func() {
		zlog.Info().Dur("took", time.Since(start)).Err(err).Msg("initReserves end")
	}()

	poolIDs := make([]common.Address, 0, len(b.pools))
	for pid := range b.pools {
		poolIDs = append(poolIDs, pid)
	}

	height, reserves, err := b.fetchReserves(poolIDs)
	if err != nil {
		return
	}
	if len(reserves) != len(poolIDs) {
		zlog.Fatal().Int("#reserves", len(reserves)).Int("#poolIDs", len(poolIDs)).Msg("#reserves != #poolIDs")
	}

	for i, pid := range poolIDs {
		if reserves[i].Reserve0 == nil || reserves[i].Reserve1 == nil {
			zlog.Fatal().Str("pool", defi.FastAddrHex(pid)).Msg("nil reserve")
		}
		b.pools[pid].Reserves.Store(&[]*big.Int{reserves[i].Reserve0, reserves[i].Reserve1})
	}

	b.height = height

	return
}

type reserveResult struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func (b *Bot) fetchReserves(poolIDs []common.Address) (height uint64, reserves []reserveResult, err error) {
	if b.pairABI == nil {
		var pairABI abi.ABI
		pairABI, err = abi.JSON(strings.NewReader(homeabi.CherryPairABI))
		if err != nil {
			return
		}
		b.pairABI = &pairABI
	}

	reserves = make([]reserveResult, len(poolIDs))

	batch := 100
	start := time.Now()
	height, err = multicall.DoSliceConcurrent(context.Background(), b.fulls, b.pairABI, len(poolIDs), batch, func(i int) []multicall.Invoke {
		return []multicall.Invoke{{
			Contract: poolIDs[i],
			Name:     "getReserves",
			Args:     []interface{}{},
		}}
	}, func(from, to int) {

		if from == 0 {
			fmt.Println("[getReserves] from", from, "to", to)
			return
		}

		var nilCount int
		for i := 0; i < batch; i++ {
			lastIdx := from - batch + i
			if reserves[lastIdx].Reserve0 == nil || reserves[lastIdx].Reserve1 == nil {
				nilCount++
			}
		}
		if nilCount > 0 {
			fmt.Println("nilCount", nilCount)
		}

		took := time.Since(start)
		fmt.Println("[getReserves] from", from, "to", to, "took", took, "eta", time.Duration(float64(took)*float64(len(poolIDs)-from)/float64((from))))
	}, func(subInvokes []multicall.Invoke, err error, client *ethclient.Client) {
		subInvokesBytes, _ := json.Marshal(subInvokes)

		fmt.Println("invoke err", err)
		fmt.Println("subInvokes", string(subInvokesBytes))
	}, reserves)
	if err != nil {
		return
	}

	var nilReservePool int
	for i, reserve := range reserves {
		if reserve.Reserve0 == nil || reserve.Reserve1 == nil {
			nilReservePool++
			fmt.Printf("nil reserve pool:%s reserve0:%v reserve1:%v ts:%d nilIdx:%d", poolIDs[i], reserve.Reserve0, reserve.Reserve1, reserve.BlockTimestampLast, nilReservePool)
		}
	}
	if nilReservePool > 0 {
		fmt.Println("nilReservePool", nilReservePool)
		panic("nilReservePool")
	}

	return
}
