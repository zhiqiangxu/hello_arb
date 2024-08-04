package tx

import (
	"context"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	gometrics "github.com/rcrowley/go-metrics"
	zlog "github.com/rs/zerolog/log"
	"github.com/zhiqiangxu/arbbot/pkg/metrics"
	"github.com/zhiqiangxu/util/parallel"
)

type request struct {
	cancel   chan struct{}
	resultCh chan common.Hash
	contract *common.Address
	gasLimit uint64
	gasPrice *big.Int
	txData   []byte
}

type Dispatcher struct {
	wg             sync.WaitGroup
	ctx            context.Context
	cancelFunc     func()
	requestCh      chan *request
	cancelChs      []chan struct{}
	clients4Send   []*ethclient.Client
	clients4Query  []*ethclient.Client
	cancelPointers []atomic.Pointer[chan struct{}]
	opts           []*bind.TransactOpts // only Signer and From is used

	dispatchMeter gometrics.Meter
}

func NewDispatcher(clients4Send []*ethclient.Client, clients4Query []*ethclient.Client, opts []*bind.TransactOpts) (d *Dispatcher) {

	cancelChs := make([]chan struct{}, len(opts))
	for i := 0; i < len(opts); i++ {
		cancelChs[i] = make(chan struct{})
	}
	if len(clients4Query) == 0 {
		clients4Query = clients4Send
	}
	d = &Dispatcher{
		requestCh:      make(chan *request),
		cancelChs:      cancelChs,
		clients4Send:   clients4Send,
		clients4Query:  clients4Query,
		cancelPointers: make([]atomic.Pointer[chan struct{}], len(opts)),
		opts:           opts,
	}

	return
}

func (d *Dispatcher) Start() {
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	d.wg.Add(len(d.opts))
	for i := 0; i < len(d.opts); i++ {
		go func(i int) {
			defer d.wg.Done()

			opt := d.opts[i]
			var (
				nonce uint64
				err   error
			)
			updateNonceFunc := func() {
				for {
					idx := rand.Intn(len(d.clients4Query))
					client := d.clients4Query[idx]
					nonce, err = client.NonceAt(context.Background(), opt.From, nil)
					if err == nil {
						break
					}
					zlog.Warn().Int("client_idx", idx).Err(err).Msg("NonceAt")
					time.Sleep(time.Second)
				}
			}
			updateNonceFunc()
			fetchTxStatusFunc := func(hash common.Hash) {
				for {
					idx := rand.Intn(len(d.clients4Query))
					client := d.clients4Query[idx]
					receipt, err := client.TransactionReceipt(context.Background(), hash)
					if err != nil {
						zlog.Warn().Int("client_idx", idx).Err(err).Msg("TransactionReceipt")
						time.Sleep(time.Second)
						continue
					}

					zlog.Info().Str("tx_hash", hash.String()).Uint64("status", receipt.Status).Msg("TransactionReceipt")
					break
				}
			}

			for {
				select {
				case <-d.ctx.Done():
					return
				case req := <-d.requestCh:
					func() {
						d.cancelPointers[i].Store(&req.cancel)
						// fetch gas price
						if req.gasPrice == nil {
							for {
								idx := rand.Intn(len(d.clients4Query))
								client := d.clients4Query[idx]
								req.gasPrice, err = client.SuggestGasPrice(context.Background())
								if err == nil {
									break
								}
								zlog.Warn().Int("client_idx", idx).Err(err).Msg("SuggestGasPrice")
								time.Sleep(time.Second)
								select {
								case <-d.ctx.Done():
									return
								case <-req.cancel:
									return
								default:
								}
							}
						}

						// fetch gas limit
						if req.gasLimit == 0 {
							maxTimes := 3
							var j int
							for {
								idx := rand.Intn(len(d.clients4Query))
								client := d.clients4Query[idx]
								callMsg := ethereum.CallMsg{
									From: opt.From, To: req.contract, GasPrice: req.gasPrice, Value: big.NewInt(0), Data: req.txData,
								}
								req.gasLimit, err = client.EstimateGas(context.Background(), callMsg)
								if err == nil {
									break
								}
								zlog.Warn().Int("client_idx", idx).Err(err).Msg("EstimateGas")
								time.Sleep(time.Second)
								select {
								case <-d.ctx.Done():
									return
								case <-req.cancel:
									return
								default:
								}
								j++
								if j >= maxTimes {
									old := d.cancelPointers[i].Swap(nil)
									if old != nil {
										close(*old)
									}
									return
								}
							}
						}

						tx := types.NewTransaction(nonce, *req.contract, big.NewInt(0), req.gasLimit, req.gasPrice, req.txData)
						tx, err := opt.Signer(opt.From, tx)
						if err != nil {
							zlog.Fatal().Int("opt_i", i).Str("tx_hash", tx.Hash().String()).Err(err).Msg("opt.Signer")
						}

						// send tx
						var sent uint32
						parallel.All(d.ctx, len(d.clients4Send), 1, len(d.clients4Send), func(ctx context.Context, workerIdx, from, to int) error {
							err := d.clients4Send[workerIdx].SendTransaction(ctx, tx)
							if err == nil {
								atomic.AddUint32(&sent, 1)
								return nil
							}
							if _, ok := err.(rpc.Error); ok {
								zlog.Warn().Int("client_idx", workerIdx).Err(err).Msg("SendTransaction rpc error treated as success")
								atomic.AddUint32(&sent, 1)
								return nil
							}
							select {
							case <-d.ctx.Done():
								zlog.Warn().Int("client_idx", workerIdx).Err(err).Msg("SendTransaction error treated as success because context is done")
								return nil
							case <-req.cancel:
								zlog.Warn().Int("client_idx", workerIdx).Err(err).Msg("SendTransaction error treated as success because request is canceled")
								return nil
							default:
								return err
							}
						}, nil, 99, time.Second)
						if sent > 0 {
							nonce++
						}

						// wait tx
						var ispending bool
						for {

							idx := rand.Intn(len(d.clients4Query))
							client := d.clients4Query[idx]
							_, ispending, err = client.TransactionByHash(context.Background(), tx.Hash())
							if err == nil && !ispending {
								fetchTxStatusFunc(tx.Hash())
								break
							}
							if err != nil {
								zlog.Warn().Int("client_idx", idx).Err(err).Msg("TransactionByHash")
							}
							time.Sleep(time.Second)

							select {
							case <-d.ctx.Done():
								return
							case <-req.cancel:
								return
							default:
							}

						}

						if req.resultCh != nil {
							req.resultCh <- tx.Hash()
						}

						updateNonceFunc()
					}()
				case <-time.NewTimer(time.Minute).C:
					// refresh nonce every minute
					updateNonceFunc()
				}
			}
		}(i)
	}
}

func (d *Dispatcher) Stop() {
	d.cancelFunc()
	d.wg.Wait()
}

func (d *Dispatcher) Dispatch(contract *common.Address, gasLimit uint64, gasPrice *big.Int, txData []byte, sync bool) {

	metrics.DispatchMeter.Mark(1)

	req := &request{cancel: make(chan struct{}), contract: contract, gasLimit: gasLimit, gasPrice: gasPrice, txData: txData}
	if sync {
		req.resultCh = make(chan common.Hash, 1)
	}
	select {
	case <-d.ctx.Done():
		return
	case d.requestCh <- req:
	default:
		i := rand.Intn(len(d.opts))
		old := d.cancelPointers[i].Swap(&req.cancel)
		if old != nil {
			close(*old)
		}

		select {
		case <-d.ctx.Done():
			return
		case d.requestCh <- req:
		case <-req.cancel:
			zlog.Warn().Msg("Dispatch canceled")
			return
		}
	}

	if sync {
		select {
		case <-d.ctx.Done():
			return
		case <-req.cancel:
			zlog.Warn().Msg("Dispatch canceled")
			return
		case hash := <-req.resultCh:
			zlog.Info().Str("tx_hash", hash.String()).Msg("Dispatch OK")

		}
	}

}
