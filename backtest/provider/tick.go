package provider

import "github.com/TrueGameover/GoBackQuant/backtest/graph"

type TickProvider interface {
	GetNextTick() *graph.Tick
}
