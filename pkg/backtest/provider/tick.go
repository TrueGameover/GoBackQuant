package provider

import "github.com/TrueGameover/GoBackQuant/pkg/graph"

type TickProvider interface {
	GetNextTick() *graph.Tick
}
