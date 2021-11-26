package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/shopspring/decimal"
)

type Strategy interface {
	BeforeTick(graph *graph.Graph)
	Tick(price decimal.Decimal)
	AfterTick()
	GetTradeFee() decimal.Decimal
}
