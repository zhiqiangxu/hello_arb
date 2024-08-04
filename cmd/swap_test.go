package cmd

import (
	"math"
	"math/big"
	"testing"

	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
)

func TestSwapArbAmountRaw(t *testing.T) {

	reserves := [][2]float64{{100, 200}, {300, 400}}
	rations := []float64{0.997, 0.997}

	amountIn := defi.SwapArbAmountInRaw(reserves, rations)

	amountInf := arbTest(reserves, rations, len(rations))
	amountIn2, _ := big.NewFloat(amountInf).Int(nil)

	if amountIn.Cmp(amountIn2) != 0 {
		t.Fatal("bug")
	}
}

func arbTest(rs [][2]float64, fs []float64, ln int) float64 {
	if ln == 2 {
		a := math.Sqrt(fs[0]*fs[1]*rs[0][0]*rs[0][1]*rs[1][0]*rs[1][1]) - (rs[0][0] * rs[1][0])
		b := fs[0]*rs[1][0] + fs[0]*fs[1]*rs[0][1] // RyF1F0 + Ry'*F0
		return a / b
	} else {
		a := math.Sqrt(fs[0]*fs[1]*fs[2]*rs[0][0]*rs[0][1]*rs[1][0]*rs[1][1]*rs[2][0]*rs[2][1]) - (rs[0][0] * rs[1][0] * rs[2][0])
		b := fs[0]*rs[1][0]*rs[2][0] + fs[0]*fs[1]*rs[0][1]*rs[2][0] + fs[0]*fs[1]*fs[2]*rs[0][1]*rs[1][1]
		return a / b
	}
}
