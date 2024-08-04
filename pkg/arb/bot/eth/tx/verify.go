package tx

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var zeroAddr common.Address

func Verify(ctx context.Context, client *ethclient.Client, tx *types.Transaction, from common.Address, signer types.Signer) (err error) {
	if from == zeroAddr {
		from, err = types.Sender(signer, tx)
		if err != nil {
			return
		}
	}

	_, err = client.CallContract(ctx, ethereum.CallMsg{From: from, To: tx.To(), Data: tx.Data()}, nil)

	return
}
