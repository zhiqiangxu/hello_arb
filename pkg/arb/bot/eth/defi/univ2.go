package defi

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"

	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	abi "github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
)

func GetAmountOut(amountIn, reserveIn, reserveOut *big.Int, fee uint64) *big.Int {
	a := new(big.Int).Mul(
		new(big.Int).Mul(reserveOut, amountIn),
		big.NewInt(int64(10000-fee)),
	)
	b := new(big.Int).Add(
		new(big.Int).Mul(reserveIn, big.NewInt(10000)),
		new(big.Int).Mul(amountIn, big.NewInt(int64(10000-fee))),
	)

	return new(big.Int).Quo(a, b)
}

func GetAmountIn(amountOut, reserveIn, reserveOut *big.Int, fee uint64) *big.Int {
	a := new(big.Int).Mul(
		new(big.Int).Mul(reserveIn, amountOut),
		big.NewInt(10000),
	)

	b := new(big.Int).Mul(
		new(big.Int).Sub(reserveOut, amountOut),
		big.NewInt(int64(10000-fee)),
	)

	return new(big.Int).Quo(a, b)
}

func GetAmountsOut(amountIn *big.Int, reserves [][2]*big.Int, fees []uint64) (amountsOut []*big.Int) {
	amountsOut = make([]*big.Int, 0, len(fees)+1)
	amountsOut = append(amountsOut, amountIn)
	for i := 0; i < len(fees); i++ {
		amountIn = GetAmountOut(amountIn, reserves[i][0], reserves[i][1], fees[i])
		amountsOut = append(amountsOut, amountIn)
	}
	return
}

func SwapExactIn(amountIn *big.Int, reserves [][2]*big.Int, fees []uint64) *big.Int {
	for i := 0; i < len(fees); i++ {
		amountIn = GetAmountOut(amountIn, reserves[i][0], reserves[i][1], fees[i])
	}
	return amountIn
}

func GetAmountsIn(amoutOut *big.Int, reserves [][2]*big.Int, fees []uint64) (amountsIn []*big.Int) {
	amountsIn = make([]*big.Int, len(fees)+1)
	amountsIn[len(fees)] = amoutOut
	for i := len(fees) - 1; i >= 0; i-- {
		amoutOut = GetAmountIn(amoutOut, reserves[i][0], reserves[i][1], fees[i])
		amountsIn[i] = amoutOut
	}
	return
}

func SwapExactOut(amoutOut *big.Int, reserves [][2]*big.Int, fees []uint64) *big.Int {
	for i := len(fees) - 1; i >= 0; i-- {
		amoutOut = GetAmountIn(amoutOut, reserves[i][0], reserves[i][1], fees[i])
	}
	return amoutOut
}

// dir true means swap from token0 to token1
func CalcFee(amountIn, amountOut, reserve0, reserve1 *big.Int, dir bool) (fee uint64) {
	if dir {
		numerator := new(big.Int).Sub(new(big.Int).Mul(amountOut, reserve0), new(big.Int).Mul(amountIn, amountOut))
		denominator := new(big.Int).Mul(amountIn, reserve1)
		fee = 10000 - new(big.Int).Quo(new(big.Int).Mul(big.NewInt(10000), numerator), denominator).Uint64()
		return
	} else {
		numerator := new(big.Int).Sub(new(big.Int).Mul(amountOut, reserve1), new(big.Int).Mul(amountIn, amountOut))
		denominator := new(big.Int).Mul(amountIn, reserve0)
		fee = 10000 - new(big.Int).Quo(new(big.Int).Mul(big.NewInt(10000), numerator), denominator).Uint64()
		return
	}
}

// -----------------model definitions-------------------
type Router struct {
	Factory  common.Address `json:"factory"`
	Router   common.Address `json:"router"`
	Exchange uint8          `json:"exchange"`
}

type Pool struct {
	Exchange  uint8                      `json:"exchange"`
	Address   common.Address             `json:"address"`
	Factory   common.Address             `json:"factory"`
	Router    common.Address             `json:"router"`
	SwapFee   uint64                     `json:"swapFee"` // 万分
	Tokens    []common.Address           `json:"tokens"`
	Reserves  atomic.Pointer[[]*big.Int] `json:"reserves"` // need to update before use
	TVL       uint                       `json:"tvl"`
	Timestamp int64                      `json:"timestamp"`
}

type pool struct {
	Exchange  uint8            `json:"exchange"`
	Address   common.Address   `json:"address"`
	Factory   common.Address   `json:"factory"`
	SwapFee   uint64           `json:"swapFee"` // 万分
	Tokens    []common.Address `json:"tokens"`
	Reserves  []*big.Int       `json:"reserves"` // need to update before use
	TVL       uint             `json:"tvl"`
	Timestamp int64            `json:"timestamp"`
}

func (p *Pool) MarshalJSON() ([]byte, error) {
	reserves := *p.Reserves.Load()
	poolData := pool{
		Exchange:  p.Exchange,
		Address:   p.Address,
		Factory:   p.Factory,
		SwapFee:   p.SwapFee,
		Tokens:    p.Tokens,
		Reserves:  reserves,
		TVL:       p.TVL,
		Timestamp: p.Timestamp,
	}

	return json.Marshal(poolData)
}

func (p *Pool) UnmarshalJSON(data []byte) (err error) {
	var poolData pool
	err = json.Unmarshal(data, &poolData)
	if err != nil {
		return
	}

	reserves := poolData.Reserves

	p.Exchange = poolData.Exchange
	p.Address = poolData.Address
	p.Factory = poolData.Factory
	p.SwapFee = poolData.SwapFee
	p.Tokens = poolData.Tokens
	p.Reserves.Store(&reserves)
	p.TVL = poolData.TVL
	p.Timestamp = poolData.Timestamp

	return
}

func (p *Pool) Reserve(token common.Address) *big.Int {
	reserves := *p.Reserves.Load()
	switch token {
	case p.Tokens[0]:
		return reserves[0]
	case p.Tokens[1]:
		return reserves[1]
	default:
		panic(fmt.Sprintf("token %v not in pool", token))
	}
}

func (p *Pool) Token0() common.Address {
	return p.Tokens[0]
}

func (p *Pool) OtherAddr(token common.Address) common.Address {
	switch token {
	case p.Tokens[0]:
		return p.Tokens[1]
	case p.Tokens[1]:
		return p.Tokens[0]
	default:
		panic(fmt.Sprintf("token %v not in pool", token))
	}
}

func (p *Pool) OtherPrice(tokenAddr common.Address, unitPrice *big.Float, tokenMap map[common.Address]*Token) *big.Float {
	token0, ok := tokenMap[p.Tokens[0]]
	if !ok {
		panic(fmt.Sprintf("token %v not in tokenMap", p.Tokens[0]))
	}
	token1, ok := tokenMap[p.Tokens[1]]
	if !ok {
		panic(fmt.Sprintf("token %v not in tokenMap", p.Tokens[1]))
	}

	reserves := *p.Reserves.Load()
	switch tokenAddr {
	case p.Tokens[0]:
		return new(big.Float).Quo(token0.Value(reserves[0], unitPrice), token1.ToUnit(reserves[1]))
	case p.Tokens[1]:
		return new(big.Float).Quo(token1.Value(reserves[1], unitPrice), token0.ToUnit(reserves[0]))
	default:
		panic(fmt.Sprintf("token %v not in pool", tokenAddr))
	}
}

type Swaps struct {
	From     common.Address `json:"from"`
	TotalTVL uint           `json:"total_tvl"`
	MinTVL   uint           `json:"min_tvl"`
	Swaps    []*Swap        `json:"swaps"`
}

func (s *Swaps) init(poolMap map[common.Address]*Pool) {
	for _, swap := range s.Swaps {
		pool := poolMap[swap.Pool]
		s.TotalTVL += pool.TVL
		if s.MinTVL == 0 || s.MinTVL < pool.TVL {
			s.MinTVL = pool.TVL
		}
	}
}

func (s *Swaps) AbiPath() (paths []abi.Path) {
	paths = make([]abi.Path, 0, len(s.Swaps))
	for _, swap := range s.Swaps {
		paths = append(paths, abi.Path{Exchange: swap.Exchange, Pool: swap.Pool, To: swap.To, Fee: big.NewInt(int64(swap.Fee))})
	}
	return
}

func (s *Swaps) String(tokenMap map[common.Address]*Token, poolMap map[common.Address]*Pool) (r string) {
	tokenPaths := s.TokenPath()
	pairPaths := s.PairPath()
	if len(tokenPaths) != len(pairPaths)+1 {
		panic("bug")
	}
	tokenNames := make([]string, len(tokenPaths))
	for i := 0; i < len(tokenPaths); i++ {
		tokenNames[i] = tokenMap[tokenPaths[i]].Name
	}
	tokenPathStr := strings.Join(tokenNames, " > ")

	pairNames := make([]string, len(pairPaths))
	for i := 0; i < len(pairPaths); i++ {
		pool := poolMap[pairPaths[i]]
		var suffix string
		if tokenPaths[i] == pool.Token0() {
			suffix = "+"
		} else {
			suffix = "-"
		}
		pairNames[i] = "0x" + FastAddrHex(pool.Address) + suffix
	}
	pairPathStr := strings.Join(pairNames, " > ")
	r = tokenPathStr + " / " + pairPathStr

	return
}
func (s *Swaps) TokenPath() []common.Address {
	paths := []common.Address{s.From}
	for _, swap := range s.Swaps {
		paths = append(paths, swap.To)
	}
	return paths
}

func (s *Swaps) PairPath() []common.Address {
	paths := []common.Address{}
	for _, swap := range s.Swaps {
		paths = append(paths, swap.Pool)
	}
	return paths
}

func (s *Swaps) Profit(amountIn *big.Int, poolMap map[common.Address]*Pool, updatedPoolReserves map[common.Address][]*big.Int) *big.Int {
	reserves := make([][2]*big.Int, 0, len(s.Swaps))
	fees := make([]uint64, 0, len(s.Swaps))

	for _, swap := range s.Swaps {
		fees = append(fees, swap.Fee)

		updatedReserve := updatedPoolReserves[swap.Pool]
		if updatedReserve != nil {
			if swap.Dir {
				reserves = append(reserves, [2]*big.Int{updatedReserve[0], updatedReserve[1]})
			} else {
				reserves = append(reserves, [2]*big.Int{updatedReserve[1], updatedReserve[0]})
			}
			continue
		}

		pool := poolMap[swap.Pool]
		poolReserves := *pool.Reserves.Load()
		if swap.Dir {
			reserves = append(reserves, [2]*big.Int{poolReserves[0], poolReserves[1]})
		} else {
			reserves = append(reserves, [2]*big.Int{poolReserves[1], poolReserves[0]})
		}
	}

	amountOut := SwapExactIn(amountIn, reserves, fees)
	return new(big.Int).Sub(amountOut, amountIn)
}

func (s *Swaps) ReservesAndRatios(poolMap map[common.Address]*Pool, updatedPoolReserves map[common.Address][]*big.Int) (reserves [][2]float64, ratios []float64) {

	for _, swap := range s.Swaps {
		ratios = append(ratios, float64(10000-swap.Fee)/10000.)
		updatedReserve := updatedPoolReserves[swap.Pool]
		if updatedReserve != nil {
			if swap.Dir {
				reserves = append(reserves, [2]float64{bigInt2Float64(updatedReserve[0]), bigInt2Float64(updatedReserve[1])})
			} else {
				reserves = append(reserves, [2]float64{bigInt2Float64(updatedReserve[1]), bigInt2Float64(updatedReserve[0])})
			}
			continue
		}
		pool := poolMap[swap.Pool]
		poolReserves := *pool.Reserves.Load()
		if swap.Dir {
			reserves = append(reserves, [2]float64{bigInt2Float64(poolReserves[0]), bigInt2Float64(poolReserves[1])})
		} else {
			reserves = append(reserves, [2]float64{bigInt2Float64(poolReserves[1]), bigInt2Float64(poolReserves[0])})
		}
	}
	return

}

func (s *Swaps) IsReverse(s2 *Swaps) (result bool) {
	if s.From != s2.From {
		return
	}

	if len(s.Swaps) != len(s2.Swaps) {
		return
	}
	n := len(s.Swaps)

	for i := 0; i < n; i++ {
		if !s.Swaps[i].IsReverse(s2.Swaps[n-1-i]) {
			return
		}
	}

	result = true
	return
}

func bigInt2Float64(v *big.Int) float64 {
	bf := new(big.Float).SetInt(v)
	f, _ := bf.Float64()
	return f
}

type Swap struct {
	Exchange uint8          `json:"exchange"`
	Pool     common.Address `json:"pool"`
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Fee      uint64         `json:"fee"`
	Dir      bool           `json:"dir"` // true means token0->token1
}

func NewSwapPath(from common.Address, swaps []*Swap, poolMap map[common.Address]*Pool) *Swaps {
	swapPath := &Swaps{From: from, Swaps: swaps}
	swapPath.init(poolMap)
	return swapPath
}

func (s *Swap) IsReverse(s2 *Swap) (result bool) {

	if s.From != s2.To || s.To != s2.From {
		return
	}

	if s.Dir == s2.Dir {
		panic("dir must be wrong")
	}

	result = true
	return
}

func PairKey(factory, tokenA, tokenB common.Address) string {
	if bytes.Compare(tokenA[:], tokenB[:]) < 0 {
		return fmt.Sprintf("%s|%s|%s", FastAddrHex(factory), FastAddrHex(tokenA), FastAddrHex(tokenB))
	} else {
		return fmt.Sprintf("%s|%s|%s", FastAddrHex(factory), FastAddrHex(tokenB), FastAddrHex(tokenA))
	}
}

func PairTokenKey(pool, token common.Address) string {
	return fmt.Sprintf("%s|%s", FastAddrHex(pool), FastAddrHex(token))
}

func FastAddrHex(addr common.Address) string {
	return hex.EncodeToString(addr.Bytes())
}

func FastAddrsHex(addrs []common.Address) (result []string) {
	result = make([]string, 0, len(addrs))
	for _, addr := range addrs {
		result = append(result, FastAddrHex(addr))
	}
	return
}
