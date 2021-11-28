package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/shopspring/decimal"
)

type Strategy interface {
	BeforeTick(graph *graph.Graph)
	Tick(price decimal.Decimal)
	AfterTick(graph *graph.Graph)
	GetTradeFee() decimal.Decimal
	ShouldContinue() bool
	IsOpenPosition() bool
	GetStopLoss() decimal.Decimal
	GetTakeProfit() decimal.Decimal
	GetPositionType() uint
	GetLotSize() decimal.Decimal
	GetSingleLotPrice() decimal.Decimal
	GetPositionsLimit() uint
}
