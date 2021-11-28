package strategy1

import (
	"github.com/TrueGameover/GoBackQuant/backtesting"
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/TrueGameover/GoBackQuant/trade"
	"github.com/shopspring/decimal"
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
	val := decimal.Decimal{}

	if err := val.Scan(10); err != nil {
		panic("Trade fee needed")
	}

	return val
}

func (strategy *TemaAndRStrategy) ShouldContinue() bool {
	return true
}

func (strategy *TemaAndRStrategy) IsOpenPosition() bool {
	return false
}

func (strategy *TemaAndRStrategy) GetStopLoss() decimal.Decimal {
	return decimal.New(0, 0)
}

func (strategy *TemaAndRStrategy) GetTakeProfit() decimal.Decimal {
	return decimal.New(0, 0)
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
	return 0
}
