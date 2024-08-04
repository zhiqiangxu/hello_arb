package defi

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Address  common.Address `json:"address"`
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Decimals uint8          `json:"decimals"`
	Price    *big.Float     `json:"price,omitempty"` // unit price, only for temporary use
}

// calculate amount by value and unitPrice
func (t *Token) Amount(value int, unitPrice *big.Float) (amountInt *big.Int) {
	if unitPrice == nil {
		unitPrice = t.Price
	}
	if unitPrice == nil {
		panic("no unitPrice")
	}

	amount := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(unitReserve(t.Decimals), big.NewInt(int64(value)))), unitPrice)
	amountInt, _ = amount.Int(nil)
	return
}

func (t *Token) Value(reserve *big.Int, unitPrice *big.Float) *big.Float {
	if t == nil {
		panic("nil token")
	}
	if reserve == nil {
		panic("nil reserve")
	}
	if unitPrice == nil {
		unitPrice = t.Price
	}
	if unitPrice == nil {
		panic("no unitPrice")
	}
	return new(big.Float).Mul(t.ToUnit(reserve), unitPrice)
}

func (t *Token) ToUnit(reserve *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(reserve),
		new(big.Float).SetInt(unitReserve(t.Decimals)))
}

func unitReserve(decimals uint8) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}
