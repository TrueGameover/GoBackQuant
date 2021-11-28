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
	GetStopLoss(price decimal.Decimal) decimal.Decimal
	GetTakeProfit(price decimal.Decimal) decimal.Decimal
	GetPositionType() uint
	GetLotSize() decimal.Decimal
	GetSingleLotPrice() decimal.Decimal
	GetPositionsLimit() uint
}
