package provider

import "github.com/TrueGameover/GoBackQuant/pkg/backtest/graph"

type TickProvider interface {
	GetNextTick() *graph.Tick
}
