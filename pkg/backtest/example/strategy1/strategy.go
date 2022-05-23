package strategy1

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/backtesting"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

type TemaAndRStrategy struct {
	backtesting.Strategy
}

func (strategy *TemaAndRStrategy) BeforeTick(graphs []*graph.Graph) {

}

func (strategy *TemaAndRStrategy) Tick(tick *graph.Tick, currentGraph *graph.Graph) {
}

func (strategy *TemaAndRStrategy) AfterTick(graph []*graph.Graph) {

}

func (strategy *TemaAndRStrategy) GetTradeFee(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.New(10, 0)
}

func (strategy *TemaAndRStrategy) ShouldContinue() bool {
	return true
}

func (strategy *TemaAndRStrategy) IsOpenPosition(currentGraph *graph.Graph) bool {
	return rand.Intn(10)%2 == 0
}

func (strategy *TemaAndRStrategy) GetStopLoss(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal {
	return tick.Close.Sub(decimal.NewFromFloat(15))
}

func (strategy *TemaAndRStrategy) GetTakeProfit(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal {
	return tick.Close.Add(decimal.NewFromFloat(30))
}

func (strategy *TemaAndRStrategy) GetPositionType(currentGraph *graph.Graph) trade.PositionType {
	return trade.TypeLong
}

func (strategy *TemaAndRStrategy) GetLotSize(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromInt(2)
}

func (strategy *TemaAndRStrategy) GetSinglePipPrice(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(1)
}

func (strategy *TemaAndRStrategy) GetSingleLotPrice(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(100)
}

func (strategy *TemaAndRStrategy) GetSinglePipValue(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(0.01)
}

func (strategy *TemaAndRStrategy) GetPositionsLimit(currentGraph *graph.Graph) uint {
	return 1
}

func (strategy TemaAndRStrategy) BeforeStart() {
	rand.Seed(time.Now().UnixMilli())
}
