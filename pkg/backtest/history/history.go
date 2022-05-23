package history

import (
	_ "embed"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
)

type Trade struct {
	Id           uint64
	Success      bool
	PipsDiff     decimal.Decimal
	Position     *trade.Position
	TotalBalance decimal.Decimal
	FreeBalance  decimal.Decimal
	BalanceDiff  decimal.Decimal
}

type TradeHistory struct {
	Graph   *graph.Graph
	counter uint64
	deals   []*Trade
}

func (saver *TradeHistory) SaveToHistory(position *trade.Position, total decimal.Decimal, free decimal.Decimal, balanceDiff decimal.Decimal) {
	diff := position.GetPipsAfterClose()

	t := Trade{
		Id:           saver.counter,
		Success:      diff.IsPositive(),
		PipsDiff:     diff,
		Position:     position,
		TotalBalance: total,
		FreeBalance:  free,
		BalanceDiff:  balanceDiff,
	}

	saver.counter++
	saver.deals = append(saver.deals, &t)
}

func (saver *TradeHistory) GetDealsCount() int64 {
	return int64(len(saver.deals))
}

func (saver *TradeHistory) GetDeals() []*Trade {
	return saver.deals
}
