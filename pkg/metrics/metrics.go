package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"

	"github.com/zhiqiangxu/util/metrics"
)

var (
	LatencyMetric kitmetrics.Histogram

	HeightMetric                                      gometrics.Gauge
	SubNewBlockMeter, ArbitrageTxMeter, DispatchMeter gometrics.Meter // count / seconds
)

func init() {
	LatencyMetric = metrics.RegisterHist("latency", []string{"method"})

	HeightMetric = gometrics.GetOrRegisterGauge("height", nil)
	SubNewBlockMeter = gometrics.GetOrRegisterMeter("subNewBlock", nil)
	ArbitrageTxMeter = gometrics.GetOrRegisterMeter("arbitrageTx", nil)
	DispatchMeter = gometrics.GetOrRegisterMeter("dispatch", nil)

	exp.Exp(gometrics.DefaultRegistry)
}
