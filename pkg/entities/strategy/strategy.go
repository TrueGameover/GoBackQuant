package strategy

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
)

type Strategy interface {
	BeforeStart()
	BeforeTick(graphs []*graph.Graph)
	Tick(tick *graph.Tick, currentGraph *graph.Graph)
	AfterTick(graph []*graph.Graph)
	GetTradeFee(currentGraph *graph.Graph) decimal.Decimal
	ShouldContinue() bool
	IsOpenPosition(currentGraph *graph.Graph) bool
	GetStopLoss(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal
	GetTakeProfit(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal
	GetPositionType(currentGraph *graph.Graph) trade.PositionType
	GetLotSize(currentGraph *graph.Graph) decimal.Decimal
	GetSinglePipPrice(currentGraph *graph.Graph) decimal.Decimal
	GetSingleLotPrice(currentGraph *graph.Graph) decimal.Decimal
	GetSinglePipValue(currentGraph *graph.Graph) decimal.Decimal
	GetPositionDayTransferCommission() decimal.Decimal
	GetPositionsLimit(currentGraph *graph.Graph) uint
	UpdateParameters(parameters []Parameter)
	GetParameters() []Parameter
}
