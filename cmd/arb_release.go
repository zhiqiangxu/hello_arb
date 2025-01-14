//go:build release
// +build release

package cmd

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/zhiqiangxu/arbbot/pkg/arb"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot"
	ethbot "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	"github.com/zhiqiangxu/litenode"
	"github.com/zhiqiangxu/litenode/eth"
	common2 "github.com/zhiqiangxu/litenode/eth/common"
)

var arbConfig map[string]arb.Config = map[string]arb.Config{
	"bsc": {
		Lite: litenode.Config{
			Eth: &common2.NodeConfig{
				P2P: p2p.Config{
					MaxPeers: 50,
					BootstrapNodes: eth.Nodes{
						"enode://1cc4534b14cfe351ab740a1418ab944a234ca2f702915eadb7e558a02010cb7c5a8c295a3b56bcefa7701c07752acd5539cb13df2aab8ae2d98934d712611443@52.71.43.172:30311", "enode://28b1d16562dac280dacaaf45d54516b85bc6c994252a9825c5cc4e080d3e53446d05f63ba495ea7d44d6c316b54cd92b245c5c328c37da24605c4a93a0d099c4@34.246.65.14:30311", "enode://5a7b996048d1b0a07683a949662c87c09b55247ce774aeee10bb886892e586e3c604564393292e38ef43c023ee9981e1f8b335766ec4f0f256e57f8640b079d5@35.73.137.11:30311",
					}.Convert(),
					StaticNodes: eth.Nodes{
						"enode://9f90d69c5fef1ca0b1417a1423038aa493a7f12d8e3d27e10a5a8fd3da216e485cf6c15f48ee310a14729bc3a4b05038479476c0aa82eed3c5d9d2e64ba3a2b3@52.69.42.169:30311", "enode://78ef719ebb2f4fc222aa988a356274dcd3624fb808936ca2ea77388ca229773d4351f795abf505e86db1a30ed1523ded9f9674d916b295bfb98516b78d2844be@13.231.200.147:30311", "enode://a8ff9670029785a644fb709ec7cd7e7e2d2b93761872bfe1b011a1ed1c601b23ffa69ead0901b759d780ed65aa81444261905b6964bdf8647bf5b061a4796d2d@54.168.191.244:30311", "enode://0f0abad52d6e3099776f70fda913611ad33c9f4b7cafad6595691ea1dd57a37804738be65315fc417d41ab52632c55a5f5f1e5ed3123ed64a312341a8c3f9e3c@52.193.230.222:30311", "enode://ecc277f466f35b249b62de8ca567dfe759162ffecc79f40339655537ee58132aec892bc0c4ad3dfb0ba5441bb7a68301c0c09e3f66454110c2c03ccca084c6b5@54.238.240.9:30311", "enode://dd3fb5f4da631067d0a9206bb0ac4400d3a076102194257911b632c5aa56f6a3289a855cc0960ad7f2cda3ba5162e0d879448775b07fa73ccd2e4e0477290d9a@54.199.96.72:30311", "enode://74481dd5079320755588b5243f82ddec7364ad36108ac77272b8e003194bb3f5e6386fcd5e50a0604db1032ac8cb9b58bb813f8e57125ad84ec6ceec65d29b4b@52.192.99.35:30311", "enode://190df80c16509d9d205145495f169a605d1459e270558f9684fcd7376934e43c65a38999d5e49d2ad118f49abfb6ff62068051ce49acc029da7d2be9910fe9fd@13.113.113.139:30311", "enode://368fc439d8f86f459822f67d9e8d1984bab32098096dc13d4d361f8a4eaf8362caae3af86e6b31524bda9e46910ac61b075728b14af163eca45413421386b7e2@52.68.165.102:30311", "enode://2038dac8d835db7c4c1f9d2647e37e6f5c5dc5474853899adb9b61700e575d237156539a720ff53cdb182ee56ac381698f357c7811f8eadc56858e0d141dcce0@18.182.11.67:30311", "enode://fc0bb7f6fc79ad7d867332073218753cb9fe5687764925f8405459a98b30f8e39d4da3a10f87fe06aa10df426c2c24c3907a4d81df4e3c88e890f7de8f8980de@54.65.239.152:30311", "enode://3aaaa0e0c7961ef3a9bf05f879f84308ca59651327cf94b64252f67448e582dcd6a6dbe996264367c8aa27fc302736db0283a3516c7406d48f268c5e317b9d49@34.250.1.192:30311", "enode://62c516645635f0389b4c851bfc4545720fac0607de74942e4ea7e923f4fa2ac0c438c146e2f0721c8ce06dca4e7f30f5c0136569d9f4b6a827c62b980fd53272@52.215.57.20:30311", "enode://5df2f71ae6b2e3bb92f92badbce1f601feabd2d6ce899cf8265c39c38ff446136d74f5bfa089532c7074bb7606a509a54a2ac66397aaaab2363dad3f43c687a8@79.125.103.83:30311", "enode://760b5fde9bc14155fa2a87e56cf610701ad6c1adcf44555a7b839baf71f86f11cdadcaf925e50b17c98cc28e20e0df3c3463caad7c6658a76ab68389af639f33@34.243.1.225:30311",
					}.Convert(),
					TrustedNodes: eth.Nodes{}.Convert(),
					// EnableMsgEvents: true,
				},
				Handler: common2.HandlerConfig{
					NetworkID:   56,
					GenesisHash: common.HexToHash("0x0d21840abff46b96c84b2ac9e10e4f5cdaeb5693cb665db62a2f3b02d2d57b5b"),
					Upgrade:     true,
				},
				TxPool:   common2.TxPoolConfig{HashCap: 10000, TxCap: 1000},
				LogLevel: log.LvlInfo,
				EthProtocolVersions: common2.ProtocolVersions{
					Versions: []uint{common2.ETH67, common2.ETH66, common2.ETH65},
					Lengths:  map[uint]uint64{common2.ETH67: 18, common2.ETH66: 17, common2.ETH65: 17},
				},
			},
		},
		Bot: bot.Config{
			Eth: &ethbot.Config{
				RPCs: []string{
					"https://solemn-quiet-voice.bsc.discover.quiknode.pro/855f0733f2342556485cd743d267fe9ca889bc85/",
					"https://bsc-dataseed1.ninicoin.io",
					"https://bsc-dataseed2.ninicoin.io",
					"https://bsc-dataseed3.ninicoin.io",
					"https://bsc-dataseed4.ninicoin.io",
					"https://bsc-dataseed1.binance.org/",
					"https://bsc-dataseed2.binance.org/",
					"https://bsc-dataseed3.binance.org/",
					"https://bsc-dataseed4.binance.org/",
					"https://bsc-dataseed1.defibit.io/",
					"https://bsc-dataseed2.defibit.io/",
					"https://bsc-dataseed3.defibit.io/",
					"https://bsc-dataseed4.defibit.io/",
					"https://bscrpc.com",
					"https://rpc-bsc.bnb48.club",
				},
				Pool: ethbot.PoolConfig{
					File: "cmd/artifact/bsc.verified.json",
					Routers: []*defi.Router{
						{
							// pancake
							Router:  common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e"),
							Factory: common.HexToAddress("0xca143ce32fe78f1f7019d7d551a6402fc5350c73"),
						},
						{
							// biswap
							Router:  common.HexToAddress("0x3a6d8ca21d1cf76f653a67577fa0d27453350dd8"),
							Factory: common.HexToAddress("0x858E3312ed3A876947EA49d572A7C42DE08af7EE"),
						},
						{
							// nomiswap
							Router:  common.HexToAddress("0xd654953d746f0b114d1f85332dc43446ac79413d"),
							Factory: common.HexToAddress("0xd6715a8be3944ec72738f0bfdc739d48c3c29349"),
						},
						{
							// mdex
							Router:  common.HexToAddress("0x7dae51bd3e3376b8c7c4900e9107f12be3af1ba8"),
							Factory: common.HexToAddress("0x3CD1C46068dAEa5Ebb0d3f55F6915B10648062B8"),
						},
					},
				},
				Arb: ethbot.ArbConfig{
					Tokens: []*ethbot.ArbTokenInfo{
						{
							Token:     common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"), // wbnb
							UnitPrice: big.NewFloat(300),                                                 // TODO fetch from pcs
						},
						{
							Token:     common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"), // usdt
							UnitPrice: big.NewFloat(1),
						},
						{
							Token:     common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"), // dai
							UnitPrice: big.NewFloat(1),
						},
						{
							Token:     common.HexToAddress("0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"), // usdc
							UnitPrice: big.NewFloat(1),
						},
						{
							Token:     common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"), // busd
							UnitPrice: big.NewFloat(1),
						},
					},
					MinProfitUSD:     50,
					MaxTxInOneBlock:  1,
					SwapExecutorAddr: common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"), // wrong
				},
				Verify: ethbot.VerifyConfig{
					Account: common.HexToAddress("0x0000000000000000000000000000000000001004"), //bsc top account
					Router:  common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e"), //pancake
					Wtoken:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"), //wbnb
				},
				CPKFile: "cmd/artifact/cpk.json",
			},
		},
	},
	"eth": {
		Lite: litenode.Config{
			Eth: &common2.NodeConfig{
				P2P: p2p.Config{
					MaxPeers: 999,
					BootstrapNodes: eth.Nodes{
						"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
					}.Convert(),
					StaticNodes:  eth.Nodes{}.Convert(),
					TrustedNodes: eth.Nodes{}.Convert(),
					// EnableMsgEvents: true,
				},
				Handler: common2.HandlerConfig{
					NetworkID:   1,
					GenesisHash: common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"),
				},
				SyncChallengeHeaderPool: &common2.SyncChallengeHeaderPoolConfig{Cap: 10000, Expire: 5 * 60},
				TxPool:                  common2.TxPoolConfig{HashCap: 10000, TxCap: 1000},
				LogLevel:                log.LvlInfo,
				EthProtocolVersions: common2.ProtocolVersions{
					Versions: []uint{common2.ETH67, common2.ETH66, common2.ETH65},
					Lengths:  map[uint]uint64{common2.ETH67: 17, common2.ETH66: 17, common2.ETH65: 17},
				},
			},
		},
		Bot: bot.Config{
			Eth: &ethbot.Config{
				RPCs: []string{
					"https://eth-mainnet.public.blastapi.io",
					"https://eth-mainnet.nodereal.io/v1/1659dfb40aa24bbb8153a677b98064d7",
					"https://rpc.ankr.com/eth",
				},
				Pool: ethbot.PoolConfig{
					File: "cmd/artifact/eth.verified.json",
					Routers: []*defi.Router{
						{
							// uniswap v2
							Router:  common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
							Factory: common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"),
						},
					},
				},
				Verify: ethbot.VerifyConfig{
					Account: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), //weth reused as eth top account
					Router:  common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), //uniswap v2
					Wtoken:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), //weth
				},
			},
		},
	},
	"polygon": {
		Lite: litenode.Config{
			Eth: &common2.NodeConfig{
				P2P: p2p.Config{
					MaxPeers: 999,
					BootstrapNodes: eth.Nodes{
						"enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303",
					}.Convert(),
					StaticNodes:  eth.Nodes{}.Convert(),
					TrustedNodes: eth.Nodes{}.Convert(),
					// EnableMsgEvents: true,
				},
				Handler: common2.HandlerConfig{
					NetworkID:   137,
					GenesisHash: common.HexToHash("0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b"),
				},
				TxPool:   common2.TxPoolConfig{HashCap: 10000, TxCap: 1000},
				LogLevel: log.LvlInfo,
				EthProtocolVersions: common2.ProtocolVersions{
					Versions: []uint{common2.ETH66},
					Lengths:  map[uint]uint64{common2.ETH66: 17},
				},
			},
		},
		Bot: bot.Config{
			Eth: &ethbot.Config{
				RPCs: []string{
					"https://polygon-rpc.com",
				},
			},
		},
	},
}
