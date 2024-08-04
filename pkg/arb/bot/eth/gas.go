package eth

import (
	"context"
	"fmt"
	"math/big"
)

func (b *Bot) initGasPrice() (err error) {
	chainID, err := b.fulls[0].ChainID(context.Background())
	if err != nil {
		return
	}

	switch chainID.Uint64() {
	case 56: //bsc
		b.Gas.MaxGas = 200000
		b.Gas.MinPrice = big.NewInt(5000000000) //5 Gwei
		b.Gas.MaxPrice = new(big.Int).Mul(b.Gas.MinPrice, big.NewInt(4))
	default:
		err = fmt.Errorf("unknown chainID:%v", chainID)
	}

	return
}
