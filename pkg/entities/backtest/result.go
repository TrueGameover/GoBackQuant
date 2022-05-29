package backtest

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/metrics"
)

type Result struct {
	Metrics metrics.TradeMetrics
	Error   error
}
