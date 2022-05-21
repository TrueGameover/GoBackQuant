package history

import (
	_ "embed"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
)

type Trade struct {
	Id        uint64
	Success   bool
	MoneyDiff decimal.Decimal
	Position  *trade.Position
}

type TradeHistory struct {
	Graph   *graph.Graph
	counter uint64
	deals   []*Trade
}

func (saver *TradeHistory) AddToHistory(positions []*trade.Position) {
	for _, position := range positions {
		saver.saveToHistory(position)
	}
}

func (saver *TradeHistory) saveToHistory(position *trade.Position) {
	diff := position.GetPipsAfterClose()

	t := Trade{
		Id:        saver.counter,
		Success:   diff.IsPositive(),
		MoneyDiff: diff,
		Position:  position,
	}

	saver.counter++
	saver.deals = append(saver.deals, &t)
}

func (saver *TradeHistory) GetDealsCount() int {
	return len(saver.deals)
}
