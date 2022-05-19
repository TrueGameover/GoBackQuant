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
	rand.Seed(time.Now().Unix())
	return rand.Intn(10)%2 == 0
}

func (strategy *TemaAndRStrategy) GetStopLoss(price decimal.Decimal) decimal.Decimal {
	return price.Sub(decimal.NewFromFloat(0.05))
}

func (strategy *TemaAndRStrategy) GetTakeProfit(price decimal.Decimal) decimal.Decimal {
	return price.Add(decimal.NewFromFloat(0.15))
}

func (strategy *TemaAndRStrategy) GetPositionType() uint {
	return trade.TypeLong
}

func (strategy *TemaAndRStrategy) GetLotSize() decimal.Decimal {
	return decimal.NewFromInt(2)
}

func (strategy *TemaAndRStrategy) GetSinglePipPrice() decimal.Decimal {
	return decimal.NewFromFloat(7)
}

func (strategy *TemaAndRStrategy) GetSingleLotPrice() decimal.Decimal {
	return decimal.NewFromFloat(3000)
}

func (strategy *TemaAndRStrategy) GetSinglePipValue() decimal.Decimal {
	return decimal.NewFromFloat(0.01)
}

func (strategy *TemaAndRStrategy) GetPositionsLimit() uint {
	return 1
}
