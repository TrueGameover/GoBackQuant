package strategy1

import (
	"github.com/TrueGameover/GoBackQuant/backtesting"
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/TrueGameover/GoBackQuant/trade"
	"github.com/shopspring/decimal"
	"math/rand"
)

type TemaAndRStrategy struct {
	backtesting.Strategy
}

func (strategy *TemaAndRStrategy) BeforeTick(graph *graph.Graph) {

}

func (strategy *TemaAndRStrategy) Tick(price decimal.Decimal) {

}

func (strategy *TemaAndRStrategy) AfterTick(graph *graph.Graph) {

}

func (strategy *TemaAndRStrategy) GetTradeFee() decimal.Decimal {
	return decimal.New(10, 0)
}

func (strategy *TemaAndRStrategy) ShouldContinue() bool {
	return true
}

func (strategy *TemaAndRStrategy) IsOpenPosition() bool {
	return rand.Intn(10)%2 == 0
}

func (strategy *TemaAndRStrategy) GetStopLoss(price decimal.Decimal) decimal.Decimal {
	return price.Sub(decimal.New(0, 5))
}

func (strategy *TemaAndRStrategy) GetTakeProfit(price decimal.Decimal) decimal.Decimal {
	return price.Add(decimal.New(0, 15))
}

func (strategy *TemaAndRStrategy) GetPositionType() uint {
	return trade.TypeLong
}

func (strategy *TemaAndRStrategy) GetLotSize() decimal.Decimal {
	return decimal.New(1, 0)
}

func (strategy *TemaAndRStrategy) GetSingleLotPrice() decimal.Decimal {
	return decimal.New(10, 0)
}

func (strategy *TemaAndRStrategy) GetPositionsLimit() uint {
	return 1
}
