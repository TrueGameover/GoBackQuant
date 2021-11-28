package trade

import "github.com/shopspring/decimal"

type Trade struct {
	Id        uint64
	Success   bool
	MoneyDiff decimal.Decimal
	Position  *Position
}

type HistorySaver struct {
	counter uint64
	deals   []*Trade
}

func (saver *HistorySaver) AddToHistory(positions []*Position) {
	for _, position := range positions {
		saver.saveToHistory(position)
	}
}

func (saver *HistorySaver) saveToHistory(position *Position) {
	diff := position.GetPipsAfterClose()

	trade := Trade{
		Id:        saver.counter,
		Success:   diff.IsPositive(),
		MoneyDiff: diff,
		Position:  position,
	}

	saver.counter++
	saver.deals = append(saver.deals, &trade)
}

func (saver *HistorySaver) GetDealsCount() int {
	return len(saver.deals)
}
