package tick

import (
	"errors"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
)

type Transformer struct {
}

func (t *Transformer) Transform(tickProvider *Provider, sourceTimeFrame graph.TimeFrame, targetTimeframe graph.TimeFrame, step uint) (*graph.Graph, error) {
	if sourceTimeFrame > targetTimeframe {
		return nil, errors.New("source should be lower lower timeframe than target")
	}

	var err error = nil
	var tick *graph.Tick
	targetGraph := &graph.Graph{Timeframe: targetTimeframe}
	i := uint(0)

	for i < step {
		tick, err = (*tickProvider).GetNextTick()

		if err != nil {
			return nil, err
		}

		if tick != nil {
			targetGraph.AddTick(tick)
		}

		i++
	}

	if err != nil {
		return nil, err
	}

	return targetGraph, nil
}
