package cmd

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/eth/protocols/snap"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	zlog "github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/abi"
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/defi"
	txPkg "github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth/tx"
	snap2 "github.com/zhiqiangxu/arbbot/pkg/protocol/eth/snap"
	"github.com/zhiqiangxu/litenode"
	common2 "github.com/zhiqiangxu/litenode/eth/common"
)

var EthCmd = cli.Command{
	Name:  "eth",
	Usage: "eth actions",
	Subcommands: []cli.Command{
		ethTxCmd,
		ethChainIDCmd,
		ethTxParamCmd,
		ethTxLogCmd,
		ethTxSendCmd,
		ethTopicCmd,
		ethHeadCmd,
		ethEciesEncryptCmd,
		ethEciesDecryptCmd,
		ethSlotCmd,
		ethBalanceCmd,
		ethManualTxCmd,
		ethManualEstimateCmd,
		ethGenWalletCmd,
		ethGenPKCmd,
		ethGenMnemonicCmd,
		ethMnemonicCmd,
		ethEmptyTrieHashCmd,
	},
}

var ethTxCmd = cli.Command{
	Name:   "tx",
	Usage:  "query tx info",
	Action: ethTx,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.TxFlag,
	},
}

var ethGenWalletCmd = cli.Command{
	Name:   "genw",
	Usage:  "generate a wallet",
	Action: ethGenWallet,
	Flags: []cli.Flag{
		flag.AuthFlag,
		flag.WalletFlag,
	},
}

var ethGenPKCmd = cli.Command{
	Name:   "genpk",
	Usage:  "generate a private key",
	Action: ethGenPK,
}

var ethGenMnemonicCmd = cli.Command{
	Name:   "genmn",
	Usage:  "generate a mnemonic and pk",
	Action: ethGenMnemonic,
}

var ethMnemonicCmd = cli.Command{
	Name:   "mn",
	Usage:  "show mnemonic info",
	Action: ethMnemonic,
}

var ethEmptyTrieHashCmd = cli.Command{
	Name:   "etrie",
	Usage:  "print hash of empty trie",
	Action: ethEmptyTrieHash,
}

var ethChainIDCmd = cli.Command{
	Name:   "chainid",
	Usage:  "query chainid",
	Action: ethChainID,
	Flags: []cli.Flag{
		flag.NetworkFlag,
	},
}

var ethHeadCmd = cli.Command{
	Name:   "head",
	Usage:  "query head",
	Action: ethHead,
	Flags: []cli.Flag{
		flag.NetworkFlag,
	},
}

var ethEciesEncryptCmd = cli.Command{
	Name:   "ecies_enc",
	Usage:  "encrypt with ecies",
	Action: ethEciesEncrypt,
	Flags: []cli.Flag{
		flag.PKFlag,
	},
}

var ethEciesDecryptCmd = cli.Command{
	Name:   "ecies_dec",
	Usage:  "decrypt with ecies",
	Action: ethEciesDecrypt,
	Flags: []cli.Flag{
		flag.PKFlag,
	},
}

var ethSlotCmd = cli.Command{
	Name:   "slot",
	Usage:  "inspect storage at slot",
	Action: ethSlot,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.AccountFlag,
		flag.SlotFlag,
		flag.OptionalHeightFlag,
	},
}

var ethBalanceCmd = cli.Command{
	Name:   "balance",
	Usage:  "get account balance",
	Action: ethBalance,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.AccountFlag,
		flag.OptionalHeightFlag,
	},
}

var ethTxParamCmd = cli.Command{
	Name:   "tx_param",
	Usage:  "parse tx param",
	Action: ethTxParam,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.TxFlag,
		flag.OptionalMethodFlag,
	},
}

var ethTxLogCmd = cli.Command{
	Name:   "tx_log",
	Usage:  "parse tx log",
	Action: ethTxLog,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.TxFlag,
		flag.EventFlag,
	},
}

var ethTxSendCmd = cli.Command{
	Name:   "tx_send",
	Usage:  "send tx",
	Action: ethTxSend,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.TxFlag,
		flag.PKFlag,
	},
}

var ethTopicCmd = cli.Command{
	Name:   "topic",
	Usage:  "calculate topic by event signature",
	Action: ethTopic,
	Flags: []cli.Flag{
		flag.EventFlag,
	},
}

var ethManualTxCmd = cli.Command{
	Name:   "manual_tx",
	Usage:  "send tx manually",
	Action: ethManualTx,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.TxDataFlag,
		flag.PKFileFlag,
		flag.GasLimitFlag,
		flag.GasPriceFlag,
	},
}

var ethManualEstimateCmd = cli.Command{
	Name:   "manual_est",
	Usage:  "estimate manually",
	Action: ethManualEstimat,
	Flags: []cli.Flag{
		flag.NetworkFlag,
		flag.ContractFlag,
		flag.TxDataFlag,
		flag.AccountFlag,
	},
}

var ethSnapAccountCmd = cli.Command{
	Name:   "snap_account",
	Usage:  "query all accounts via p2p",
	Action: ethSnapAccount,
	Flags: []cli.Flag{
		flag.NetworkFlag,
	},
}

func ethTx(ctx *cli.Context) (err error) {

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	hash := common.HexToHash(ctx.String(flag.TxFlag.Name))

	tx, _, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return
	}

	fmt.Printf(
		"ChainId:\t%v\nGasPrice:\t%v\nGas:\t\t%v\nCost:\t\t%v\nNonce:\t\t%v\nType:\t\t%v\nSize:\t\t%v\nGasFeeCap:\t%v\nGasTipCap:\t%v\nTo:\t\t%v\nProtected:\t%v\nData:\t\t%v\n",
		tx.ChainId(),
		tx.GasPrice(),
		tx.Gas(),
		tx.Cost(),
		tx.Nonce(),
		tx.Type(),
		tx.Size(),
		tx.GasFeeCap(),
		tx.GasTipCap(),
		tx.To(),
		tx.Protected(),
		hex.EncodeToString(tx.Data()),
	)

	var receipt *types.Receipt
	receipt, err = client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return
	}

	fmt.Println("\nreceipt.ContractAddress", receipt.ContractAddress)

	var logBytes []byte
	logBytes, err = json.MarshalIndent(receipt.Logs, "", "  ")
	if err != nil {
		return
	}
	fmt.Println("\nlogs", string(logBytes))
	return
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *keystore.Key {
	id, _ := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

func ethGenWallet(ctx *cli.Context) (err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	storeKey := newKeyFromECDSA(key)

	resultBytes, err := keystore.EncryptKey(storeKey, ctx.String(flag.AuthFlag.Name), keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(ctx.String(flag.WalletFlag.Name), resultBytes, 0777)

	fmt.Println("pk", hex.EncodeToString(crypto.FromECDSA(key)))
	fmt.Println("addr", crypto.PubkeyToAddress(key.PublicKey).Hex())

	return
}

func ethMnemonic(ctx *cli.Context) (err error) {
	mnemonic := ctx.Args()[0]
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return
	}

	// m / purpose' / coin_type' / account' / change / address_index
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		err = fmt.Errorf("Derive:%v", err)
		return
	}
	priKey, err := wallet.PrivateKey(account)
	if err != nil {
		return
	}
	address := crypto.PubkeyToAddress(priKey.PublicKey)
	fmt.Printf("mnemonic:\t%s\npk:\t%s\naddress: %s\n\n", mnemonic, hex.EncodeToString(crypto.FromECDSA(priKey)), address)
	return
	return
}

func ethGenMnemonic(ctx *cli.Context) (err error) {

	mnemonic, err := hdwallet.NewMnemonic(256 / 2)
	if err != nil {
		return
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return
	}

	// m / purpose' / coin_type' / account' / change / address_index
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		err = fmt.Errorf("Derive:%v", err)
		return
	}
	priKey, err := wallet.PrivateKey(account)
	if err != nil {
		return
	}
	address := crypto.PubkeyToAddress(priKey.PublicKey)
	fmt.Printf("mnemonic:\t%s\npk:\t%s\naddress: %s\n\n", mnemonic, hex.EncodeToString(crypto.FromECDSA(priKey)), address)
	return
}

func ethGenPK(ctx *cli.Context) (err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	fmt.Println("bitlen", key.D.BitLen())
	fmt.Println("pk", hex.EncodeToString(crypto.FromECDSA(key)))
	fmt.Println("addr", crypto.PubkeyToAddress(key.PublicKey).Hex())

	t1 := time.Now()
	data := crypto.Keccak256([]byte("abc"))
	sig, err := crypto.Sign(data, key)
	if err != nil {
		return
	}
	t2 := time.Now()
	pubkey, err := crypto.Ecrecover(data, sig)
	if err != nil {
		return
	}

	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	fmt.Println("signer", signer, "#sig", len(sig), "sign took", t2.Sub(t1), "verify took", time.Now().Sub(t2))
	return
}

func ethEmptyTrieHash(ctx *cli.Context) (err error) {
	hasher := trie.NewStackTrie(nil)
	fmt.Println(hasher.Hash())
	return
}

func ethChainID(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return
	}

	fmt.Println(chainID)
	return
}
func ethHead(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return
	}
	headerBytes, _ := json.Marshal(header)

	fmt.Println(string(headerBytes))
	return
}

func ethEciesEncrypt(ctx *cli.Context) (err error) {
	pk, err := crypto.HexToECDSA(ctx.String(flag.PKFlag.Name))
	if err != nil {
		return
	}

	eciespk := ecies.ImportECDSA(pk)
	result, err := ethEciesEncryptImpl(eciespk, []byte(ctx.Args()[0]))
	if err != nil {
		return
	}
	fmt.Println(result)
	return
}

func ethEciesEncryptImpl(eciespk *ecies.PrivateKey, data []byte) (result string, err error) {
	cipher, err := ecies.Encrypt(rand.Reader, &eciespk.PublicKey, data, nil, nil)
	if err != nil {
		return
	}
	result = hex.EncodeToString(cipher)
	return
}

func ethEciesDecrypt(ctx *cli.Context) (err error) {
	pk, err := crypto.HexToECDSA(ctx.String(flag.PKFlag.Name))
	if err != nil {
		return
	}

	eciespk := ecies.ImportECDSA(pk)

	plain, err := ethEciesDecryptImpl(eciespk, ctx.Args()[0])
	if err != nil {
		return
	}

	fmt.Println(string(plain))
	return
}

func ethSlot(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	account := common.HexToAddress(ctx.String(flag.AccountFlag.Name))

	key := common.HexToHash(ctx.String(flag.SlotFlag.Name))
	var height *big.Int
	if ctx.IsSet(flag.OptionalHeightFlag.Name) {
		height = big.NewInt(int64(ctx.Uint64(flag.OptionalHeightFlag.Name)))
	}
	result, err := client.StorageAt(context.Background(), account, key, height)
	if err != nil {
		return
	}

	fmt.Println(hex.EncodeToString(result))
	return
}

func ethBalance(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	account := common.HexToAddress(ctx.String(flag.AccountFlag.Name))
	var height *big.Int
	if ctx.IsSet(flag.OptionalHeightFlag.Name) {
		height = big.NewInt(int64(ctx.Uint64(flag.OptionalHeightFlag.Name)))
	}

	balance, err := client.BalanceAt(context.Background(), account, height)
	if err != nil {
		return
	}
	fmt.Println(balance)
	return
}
func ethEciesDecryptImpl(eciespk *ecies.PrivateKey, cipherHex string) (plain []byte, err error) {
	data, err := hex.DecodeString(cipherHex)
	if err != nil {
		return
	}
	plain, err = eciespk.Decrypt(data, nil, nil)
	return
}

func ethTxParam(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	hash := common.HexToHash(ctx.String(flag.TxFlag.Name))

	tx, _, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return
	}

	method := ctx.String(flag.OptionalMethodFlag.Name)
	var tryMethods []string
	if method == "" {
		tryMethods = defi.SwapFunctions
	} else {
		tryMethods = []string{method}
	}
	data := tx.Data()
	if len(data) <= 4 {
		return
	}

	found := false
	for _, m := range tryMethods {
		var (
			sig string
			in  ethabi.Arguments
		)
		sig, in, _, err = abi.ParseFunction(m)
		if err != nil {
			return
		}
		mid := abi.SigToMid(sig)
		if !bytes.Equal(mid, data[0:4]) {
			continue
		}

		found = true
		values := make(map[string]interface{})
		err = in.UnpackIntoMap(values, data[4:])
		if err != nil {
			return
		}

		fmt.Println(m)
		fmt.Println(values)
		for k, v := range values {
			if vbytes, ok := v.([]byte); ok {
				fmt.Printf("%s:\t%s\n", k, hex.EncodeToString(vbytes))
			}
		}
		break
	}

	if !found {
		err = fmt.Errorf("didn't find matching method")
		return
	}

	return
}

func ethTxSend(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return
	}

	txBytes, err := hex.DecodeString(ctx.String(flag.TxFlag.Name))
	if err != nil {
		return
	}
	var tx types.Transaction
	err = tx.UnmarshalBinary(txBytes)
	if err != nil {
		return
	}

	pk, err := crypto.HexToECDSA(flag.PKFlag.Name)
	if err != nil {
		return
	}

	opt, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return
	}
	dispatcher := txPkg.NewDispatcher([]*ethclient.Client{client}, nil, []*bind.TransactOpts{opt})

	dispatcher.Dispatch(tx.To(), tx.Gas(), tx.GasPrice(), tx.Data(), true)
	return
}

func ethTxLog(ctx *cli.Context) (err error) {
	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	hash := common.HexToHash(ctx.String(flag.TxFlag.Name))

	var receipt *types.Receipt
	receipt, err = client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return
	}

	events := strings.Split(ctx.String(flag.EventFlag.Name), ";")
	inMap := make(map[common.Hash]ethabi.Arguments)
	eventMap := make(map[common.Hash]string)
	for _, event := range events {
		var (
			sig string
			in  ethabi.Arguments
		)
		sig, in, err = abi.ParseEvent(event)
		if err != nil {
			return
		}
		topic := abi.SigToTopic(sig)
		inMap[topic] = in
		eventMap[topic] = abi.SigToEvent(sig)
	}
	for _, log := range receipt.Logs {

		in, ok := inMap[log.Topics[0]]
		if ok {
			value := make(map[string]interface{})
			err = in.UnpackIntoMap(value, log.Data)
			if err != nil {
				err = fmt.Errorf("UnpackIntoMap:%v", err)
				return
			}
			indexed := make(map[string]interface{})
			i := 0
			for _, argument := range in {
				if argument.Indexed {
					indexed[argument.Name] = log.Topics[1+i]
					i++
				}
			}

			fmt.Println(eventMap[log.Topics[0]], value)
			if len(indexed) > 0 {
				fmt.Println("\tindexed", indexed)
			}
		}
	}

	return
}

func ethTopic(ctx *cli.Context) (err error) {
	event := ctx.String(flag.EventFlag.Name)
	sig, _, err := abi.ParseEvent(event)
	if err != nil {
		return
	}
	topic := abi.SigToTopic(sig)

	fmt.Println("topic", topic)
	return
}

func ethManualTx(ctx *cli.Context) (err error) {

	txData, err := hex.DecodeString(ctx.String(flag.TxDataFlag.Name))
	if err != nil {
		return
	}

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return
	}

	pkBytes, err := os.ReadFile(ctx.String(flag.PKFileFlag.Name))
	if err != nil {
		return
	}
	realPK, err := crypto.HexToECDSA(string(pkBytes))
	if err != nil {
		return
	}
	transactor, err := bind.NewKeyedTransactorWithChainID(realPK, chainID)
	if err != nil {
		return
	}

	dispatcher := txPkg.NewDispatcher([]*ethclient.Client{client}, nil, []*bind.TransactOpts{transactor})
	dispatcher.Start()

	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	gasPriceStr := ctx.String(flag.GasPriceFlag.Name)
	var gasPrice *big.Int
	if gasPriceStr != "" {
		var ok bool
		gasPrice, ok = new(big.Int).SetString(gasPriceStr, 10)
		if !ok {
			err = fmt.Errorf("invalid gas price:%s", gasPriceStr)
			return
		}
	}
	dispatcher.Dispatch(&contract, ctx.Uint64(flag.GasLimitFlag.Name), gasPrice, txData, true)

	return
}

func ethManualEstimat(ctx *cli.Context) (err error) {
	txData, err := hex.DecodeString(ctx.String(flag.TxDataFlag.Name))
	if err != nil {
		return
	}

	client, err := networkToClient(ctx)
	if err != nil {
		return
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	fmt.Println("gasPrice", gasPrice)

	account := common.HexToAddress(ctx.String(flag.AccountFlag.Name))
	contract := common.HexToAddress(ctx.String(flag.ContractFlag.Name))

	callMsg := ethereum.CallMsg{
		From: account, To: &contract, GasPrice: gasPrice, Value: big.NewInt(0), Data: txData,
	}
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return
	}

	fmt.Println("gasLimit", gasLimit)

	product := new(big.Float).SetInt(new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	fmt.Println("cost", new(big.Float).Quo(product, new(big.Float).SetInt(unit)))

	return
}

func ethSnapAccount(ctx *cli.Context) (err error) {
	network := ctx.String(flag.NetworkFlag.Name)
	config, ok := arbConfig[network]
	if !ok {
		err = fmt.Errorf("network %s not supported yet", network)
		return
	}

	syncer := snap2.NewSyncer("accounts.json")
	config.Lite.Eth.Handler.StatusFeed = true
	config.Lite.Eth.Handler.SnapSyncer = syncer
	config.Lite.Eth.SnapProtocolVersions = &common2.ProtocolVersions{Name: snap.ProtocolName, Versions: snap.ProtocolVersions, Lengths: map[uint]uint64{snap.SNAP1: 8}}
	node := litenode.New(config.Lite)
	err = node.Start()
	if err != nil {
		return
	}
	defer func() {
		node.Stop()
		syncer.Stop()
	}()

	statusCh := make(chan eth.MinStatus, 1)
	statusSub := node.Eth.SubscribeStatusMsg(statusCh)
	defer statusSub.Unsubscribe()

	snapMsgCh := make(chan common2.SnapSyncPacket)
	snapMsgSub := node.Eth.SubscribeSnapSyncMsg(snapMsgCh)
	defer snapMsgSub.Unsubscribe()

	syncer.Start(node, statusCh, snapMsgCh)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		time.Sleep(time.Second * 10)
		zlog.Info().Int("peers", node.Eth.PeerCount()).Msg("stat")

		select {
		case <-signalChan:
			return
		default:
		}

	}

}
