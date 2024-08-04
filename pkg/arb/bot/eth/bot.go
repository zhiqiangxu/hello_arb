package eth

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/tx"
	"github.com/zhiqiangxu/litenode"
	"github.com/zhiqiangxu/lru"
)

type Bot struct {
	Config
	lite       *litenode.Node
	fulls      []*ethclient.Client
	fulls4Send []*ethclient.Client
	fullSub    *ethclient.Client

	ctx      context.Context
	cancelFn context.CancelFunc
	wg       sync.WaitGroup

	pairABI         *abi.ABI
	swapExecutorABI *abi.ABI
	swapVerifierABI *abi.ABI
	syncArguments   abi.Arguments

	txDispather *tx.Dispatcher
	// mutable

	txCache lru.Cache

	sentInBlock int32

	height uint64

	tokens          map[common.Address]*defi.Token
	pools           map[common.Address]*defi.Pool
	pair2Pool       map[string]common.Address
	pairToken2Swaps map[string][]*defi.Swaps
}

func NewBot(config *Config, lite *litenode.Node) *Bot {
	if lite.Eth == nil {
		panic("lite.Eth is nil for eth bot")
	}
	txCache := lru.NewCache(10000, 0, nil)
	bot := &Bot{Config: *config, lite: lite, txCache: txCache}
	return bot
}

func (b *Bot) Start() (err error) {

	b.ctx, b.cancelFn = context.WithCancel(context.Background())

	err = b.initClients()
	if err != nil {
		return
	}

	err = b.initGasPrice()
	if err != nil {
		return
	}

	err = b.initDispatcher()
	if err != nil {
		return
	}

	err = b.initVerify()
	if err != nil {
		return
	}

	err = b.initPools()
	if err != nil {
		return
	}

	err = b.initReserves()
	if err != nil {
		return
	}

	err = b.initSub()
	if err != nil {
		return
	}

	err = b.initArb()
	if err != nil {
		return
	}

	var arbSwapPaths int
	for _, swapPaths := range b.pairToken2Swaps {
		arbSwapPaths += len(swapPaths)
	}
	fmt.Println("#pairToken2Swaps", len(b.pairToken2Swaps), "arbSwapPaths", arbSwapPaths, "#pool", len(b.pools), "#token", len(b.tokens))

	return
}

func (b *Bot) initClients() (err error) {

	{
		clients := make([]*ethclient.Client, 0, len(b.RPCs))
		for _, rpc := range b.RPCs {
			var client *ethclient.Client
			client, err = ethclient.Dial(rpc)
			if err != nil {
				return
			}
			clients = append(clients, client)
		}
		b.fulls = clients
	}

	if len(b.RPCs4Send) == 0 {
		b.fulls4Send = b.fulls
	} else {
		clients := make([]*ethclient.Client, 0, len(b.RPCs4Send))
		for _, rpc := range b.RPCs4Send {
			var client *ethclient.Client
			client, err = ethclient.Dial(rpc)
			if err != nil {
				return
			}
			clients = append(clients, client)
		}
		b.fulls4Send = clients
	}

	if b.FullSub != "" {
		b.fullSub, err = ethclient.Dial(b.FullSub)
		if err != nil {
			return
		}
	}
	return

}

func (b *Bot) Stop() {
	b.txDispather.Stop()
	b.cancelFn()
	b.wg.Wait()
}
