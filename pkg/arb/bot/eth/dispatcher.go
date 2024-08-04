package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/tx"
)

func (b *Bot) initDispatcher() (err error) {
	if len(b.PKs) == 0 {
		err = fmt.Errorf("empty PKs")
		return
	}

	chainID, err := b.fulls[0].ChainID(context.Background())
	if err != nil {
		err = fmt.Errorf("initDispatcher ChainID:%v", err)
		return
	}

	var opts []*bind.TransactOpts
	for _, pk := range b.PKs {
		var transactor *bind.TransactOpts
		transactor, err = bind.NewKeyedTransactorWithChainID(pk, chainID)
		if err != nil {
			err = fmt.Errorf("initDispatcher NewKeyedTransactorWithChainID:%v", err)
			return
		}
		opts = append(opts, transactor)
	}

	txDispatcher := tx.NewDispatcher(b.fulls4Send, b.fulls, opts)
	txDispatcher.Start()
	b.txDispather = txDispatcher

	return
}
