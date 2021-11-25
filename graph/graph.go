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
	TimeFrameM1  uint = 1
	TimeFrameM5       = 5
	TimeFrameM10      = 10
	TimeFrameM15      = 15
	TimeFrameM30      = 30
	TimeFrameH1       = 60
	TimeFrameH4       = 4 * 60
	TimeFrameD1       = 24 * 60
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
