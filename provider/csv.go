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
		Date:   time.Now(),
		Open:   decimal.New(0, 0),
		High:   decimal.New(0, 0),
		Low:    decimal.New(0, 0),
		Close:  decimal.New(0, 0),
		Volume: decimal.New(0, 0),
	}
}

func (provider *CsvProvider) HasTicks() bool {
	return true
}
