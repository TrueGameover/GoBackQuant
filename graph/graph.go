package graph

import (
	"github.com/shopspring/decimal"
	"time"
)

type Tick struct {
	Id     uint64
	Date   time.Time
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}

const (
	TIMEFRAME_M1  uint = 1
	TIMEFRAME_M5       = 5
	TIMEFRAME_M10      = 10
	TIMEFRAME_M15      = 15
	TIMEFRAME_M30      = 30
	TIMEFRAME_H1       = 60
	TIMEFRAME_H4       = 4 * 60
	TIMEFRAME_D1       = 24 * 60
)

type Bar struct {
	Tick
	Ticks []*Tick
}

type Graph struct {
	Timeframe  uint
	bars       []*Bar
	lastTickId uint64
}

func (graph *Graph) AddBar(bar *Bar) {
	var index = graph.lastTickId

	for _, tick := range bar.Ticks {
		tick.Id = index
		index++
	}

	bar.Id = index
	index++

	graph.bars = append(graph.bars, bar)
	graph.lastTickId = index
}

func (graph *Graph) GetBars() []*Bar {
	return graph.bars
}

func (graph *Graph) Reset() {
	graph.lastTickId = 1
}

func (graph *Graph) GetTicksCount() uint64 {
	return graph.lastTickId - 1
}
