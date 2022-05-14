package graph

import (
	"github.com/shopspring/decimal"
	"math"
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

func (bar *Bar) GetLastTick() *Tick {
	count := len(bar.Ticks)

	if count > 0 {
		return bar.Ticks[count-1]
	}

	return nil
}

func (bar *Bar) Recalculate() {
	if len(bar.Ticks) > 0 {
		bar.Open = bar.Ticks[0].Open
		bar.Close = bar.Ticks[len(bar.Ticks)-1].Close

		for _, tick := range bar.Ticks {
			if tick.Low.LessThan(bar.Low) {
				bar.Low = tick.Low
			}

			if tick.High.GreaterThan(bar.High) {
				bar.High = tick.High
			}

			bar.Volume = bar.Volume.Add(tick.Volume)
		}
	}
}

type Graph struct {
	Timeframe  uint
	bars       []*Bar
	lastTickId uint64
	currentBar *Bar
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

func (graph *Graph) Tick(tick *Tick) {
	bar := graph.currentBar

	if bar == nil {
		bar = &Bar{
			Tick: Tick{
				Id:     0,
				Date:   tick.Date,
				Open:   decimal.New(0, 0),
				High:   decimal.New(math.MinInt64, 0),
				Low:    decimal.New(math.MaxInt64, 0),
				Close:  decimal.New(0, 0),
				Volume: decimal.New(0, 0),
			},
			Ticks: []*Tick{},
		}
	}

	bar.Ticks = append(bar.Ticks, tick)
	bar.Recalculate()

	if len(bar.Ticks) >= int(graph.Timeframe) {
		graph.AddBar(bar)
		bar = nil
	}

	graph.currentBar = bar
}

func (graph *Graph) GetBars() []*Bar {
	return graph.bars
}

func (graph *Graph) Reset() {
	graph.lastTickId = 1
	graph.bars = []*Bar{}
}

func (graph *Graph) GetTicksCount() uint64 {
	return graph.lastTickId - 1
}

func (graph *Graph) GetFreshBar() *Bar {
	if graph.currentBar == nil {
		return graph.bars[len(graph.bars)-1]
	}

	return graph.currentBar
}
