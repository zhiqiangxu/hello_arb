package eth

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
)

type Config struct {
	RPCs      []string
	RPCs4Send []string
	FullSub   string
	Gas       GasConfig
	Pool      PoolConfig
	Arb       ArbConfig
	Verify    VerifyConfig
	CPKFile   string
	PKs       []*ecdsa.PrivateKey // read from CPKFile
}

type GasConfig struct {
	MaxGas   uint64
	MinPrice *big.Int
	MaxPrice *big.Int
}

type ArbTokenInfo struct {
	Token     common.Address
	UnitPrice *big.Float
	detail    *defi.Token
}

type ArbConfig struct {
	Tokens           []*ArbTokenInfo
	MinProfitUSD     int
	MaxTxInOneBlock  int32
	SwapExecutorAddr common.Address

	tokenMap map[common.Address]*ArbTokenInfo
}

type PoolConfig struct {
	File    string
	Routers []*defi.Router

	routerMap map[common.Address]*defi.Router
}

type VerifyConfig struct {
	Account common.Address
	Wtoken  common.Address
	Router  common.Address

	// set automatically
	value *big.Int
}
