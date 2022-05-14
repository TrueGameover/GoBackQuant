package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/backtest/graph"
	"github.com/shopspring/decimal"
)

type Strategy interface {
	BeforeTick(graph *graph.Graph)
	Tick(price decimal.Decimal)
	AfterTick(graph *graph.Graph)
	GetTradeFee() decimal.Decimal
	ShouldContinue() bool
	IsOpenPosition() bool
	GetStopLoss(price decimal.Decimal) decimal.Decimal
	GetTakeProfit(price decimal.Decimal) decimal.Decimal
	GetPositionType() uint
	GetLotSize() decimal.Decimal
	GetSinglePipPrice() decimal.Decimal
	GetSingleLotPrice() decimal.Decimal
	GetSinglePipValue() decimal.Decimal
	GetPositionsLimit() uint
}
