package provider

import (
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/shopspring/decimal"
	"time"
)

type CsvProvider struct {
	TickProvider
}

func (provider *CsvProvider) GetNextTick() graph.Tick {

	return graph.Tick{
		Id:     0,
		Date:   time.Time{},
		Open:   decimal.Decimal{},
		High:   decimal.Decimal{},
		Low:    decimal.Decimal{},
		Close:  decimal.Decimal{},
		Volume: decimal.Decimal{},
	}
}
