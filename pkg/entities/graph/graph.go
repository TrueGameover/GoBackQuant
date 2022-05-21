package graph

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
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

func (t *Tick) IsGrowing() bool {
	return t.Close.GreaterThanOrEqual(t.Open)
}

func (t *Tick) IsFalling() bool {
	return t.Close.LessThan(t.Open)
}

type TimeFrame uint

const (
	TimeFrameM1  TimeFrame = 1
	TimeFrameM5  TimeFrame = 5
	TimeFrameM10 TimeFrame = 10
	TimeFrameM15 TimeFrame = 15
	TimeFrameM30 TimeFrame = 30
	TimeFrameH1  TimeFrame = 60
	TimeFrameH4  TimeFrame = 4 * 60
	TimeFrameD1  TimeFrame = 24 * 60
	TimeFrameD4  TimeFrame = 4 * 24 * 60
	TimeFrameW1  TimeFrame = 7 * 24 * 60
)

func ParseTimeFrame(d time.Duration) (TimeFrame, error) {
	minutes := uint(d.Minutes())

	switch minutes {
	case uint(TimeFrameM1):
		return TimeFrameM1, nil
	case uint(TimeFrameM5):
		return TimeFrameM5, nil
	case uint(TimeFrameM10):
		return TimeFrameM10, nil
	case uint(TimeFrameM15):
		return TimeFrameM15, nil
	case uint(TimeFrameM30):
		return TimeFrameM30, nil
	case uint(TimeFrameH1):
		return TimeFrameH1, nil
	case uint(TimeFrameH4):
		return TimeFrameH4, nil
	case uint(TimeFrameD1):
		return TimeFrameD1, nil
	case uint(TimeFrameD4):
		return TimeFrameD4, nil
	case uint(TimeFrameW1):
		return TimeFrameW1, nil
	}

	return 0, errors.New("not supported")
}

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
	Timeframe  TimeFrame
	Title      string
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

func (graph *Graph) AddTick(tick *Tick) {
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

func (graph *Graph) GetTicks() []*Tick {
	result := funk.Map(graph.bars, func(t *Bar) *Tick {
		return &t.Tick
	})

	ticks, ok := result.([]*Tick)

	if ok {
		return ticks
	}

	return nil
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
