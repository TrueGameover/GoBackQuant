package graph

import (
	"github.com/shopspring/decimal"
	"time"
)

type Tick struct {
	Date   time.Time
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}

const (
	TIMEFRAME_M1  = 1
	TIMEFRAME_M5  = 5
	TIMEFRAME_M10 = 10
	TIMEFRAME_M15 = 15
	TIMEFRAME_M30 = 30
	TIMEFRAME_H1  = 60
	TIMEFRAME_H4  = 4 * 60
	TIMEFRAME_D1  = 24 * 60
)

type Bar struct {
	Date  time.Time
	Ticks []*Tick
}

type Graph struct {
	Bars []Bar
}
