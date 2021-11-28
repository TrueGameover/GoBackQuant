package provider

import "github.com/TrueGameover/GoBackQuant/graph"

type TickProvider interface {
	GetNextTick() graph.Tick
	HasTicks() bool
}
