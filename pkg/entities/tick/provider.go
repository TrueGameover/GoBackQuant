package tick

import "github.com/TrueGameover/GoBackQuant/pkg/entities/graph"

type Provider interface {
	GetNextTick() (*graph.Tick, error)
	GetTotal() uint64
	GetTitle() string
	GetTimeFrame() graph.TimeFrame
}
