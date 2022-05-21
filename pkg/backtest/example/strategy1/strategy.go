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
	rand.Seed(time.Now().Unix())
	return rand.Intn(10)%2 == 0
}

func (strategy *TemaAndRStrategy) GetStopLoss(price decimal.Decimal, currentGraph *graph.Graph) decimal.Decimal {
	return price.Sub(decimal.NewFromFloat(15))
}

func (strategy *TemaAndRStrategy) GetTakeProfit(price decimal.Decimal, currentGraph *graph.Graph) decimal.Decimal {
	return price.Add(decimal.NewFromFloat(30))
}

func (strategy *TemaAndRStrategy) GetPositionType(currentGraph *graph.Graph) trade.PositionType {
	return trade.TypeLong
}

func (strategy *TemaAndRStrategy) GetLotSize(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromInt(2)
}

func (strategy *TemaAndRStrategy) GetSinglePipPrice(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(7)
}

func (strategy *TemaAndRStrategy) GetSingleLotPrice(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(3000)
}

func (strategy *TemaAndRStrategy) GetSinglePipValue(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(0.01)
}

func (strategy *TemaAndRStrategy) GetPositionsLimit(currentGraph *graph.Graph) uint {
	return 1
}
