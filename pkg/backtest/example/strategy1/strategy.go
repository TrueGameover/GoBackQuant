package strategy1

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

type TemaAndRStrategy struct {
	parameters []strategy.Parameter
}

func (strategy *TemaAndRStrategy) BeforeTick(graphs []*graph.Graph) {

}

func (strategy *TemaAndRStrategy) Tick(tick *graph.Tick, currentGraph *graph.Graph) {
}

func (strategy *TemaAndRStrategy) AfterTick(graph []*graph.Graph) {

}

func (strategy *TemaAndRStrategy) ShouldContinue() bool {
	return true
}

func (strategy *TemaAndRStrategy) IsOpenPosition(currentGraph *graph.Graph) bool {
	val := strategy.parameters[0].GetValue()
	//div := strategy.parameters[1].GetValue()

	//if div < 1 {
	//	div = 1
	//}

	return rand.Intn(10)%val == 0
}

func (strategy *TemaAndRStrategy) GetStopLoss(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal {
	return tick.Close.Sub(decimal.NewFromFloat(15))
}

func (strategy *TemaAndRStrategy) GetTakeProfit(tick *graph.Tick, currentGraph *graph.Graph) decimal.Decimal {
	return tick.Close.Add(decimal.NewFromFloat(30))
}

func (strategy *TemaAndRStrategy) GetPositionType(currentGraph *graph.Graph) trade.PositionType {
	r := rand.Intn(10)

	if r%2 == 0 {
		return trade.TypeLong
	} else {
		return trade.TypeShort
	}
}

func (strategy *TemaAndRStrategy) GetLotSize(currentGraph *graph.Graph) int64 {
	return 1
}

func (strategy *TemaAndRStrategy) GetSingleLotPrice(currentGraph *graph.Graph) decimal.Decimal {
	return decimal.NewFromFloat(100)
}

func (strategy *TemaAndRStrategy) GetPositionsLimit(currentGraph *graph.Graph) uint {
	return 1
}

func (strategy *TemaAndRStrategy) BeforeStart() {
	rand.Seed(time.Now().UnixMilli())
}

func (strategy *TemaAndRStrategy) UpdateParameters(parameters []strategy.Parameter) {
	strategy.parameters = parameters
}

func (strategy *TemaAndRStrategy) GetParameters() []strategy.Parameter {
	return strategy.parameters
}

func (strategy *TemaAndRStrategy) Clone() strategy.Strategy {
	return &TemaAndRStrategy{parameters: strategy.parameters}
}

func (strategy *TemaAndRStrategy) GetLotsAmount(currentGraph *graph.Graph) int64 {
	return 5
}
