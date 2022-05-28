package main

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/metrics"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	strategy2 "github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/optimization"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	optimizer := optimization.StrategyOptimizer{}

	var strategy strategy2.Strategy = &strategy1.TemaAndRStrategy{}
	strategy.UpdateParameters([]strategy2.Parameter{
		{
			Name: "random_factor1",
			Min:  1,
			Max:  100,
		},
		{
			Name: "random_factor2",
			Min:  1,
			Max:  100,
		},
	})

	optimizer.Run(strategy, getProviders, fitness, 80, 100, 10000, 4)
}

func getProviders() []tick.Provider {
	providers := make([]tick.Provider, 1)

	csvProvider1 := provider.CsvProvider{
		DateParseTemplate: time.RFC3339,
		Delimiter:         ';',
		Positions: provider.Positions{
			Date:   0,
			Open:   2,
			High:   3,
			Low:    4,
			Close:  1,
			Volume: 5,
		},
		FieldsPerRecord: 6,
	}
	err := csvProvider1.Load("SBER_m1.csv", "SBER", graph.TimeFrameM15)
	if err != nil {
		panic(err)
	}
	providers[0] = &csvProvider1

	return providers
}

func fitness(metrics metrics.TradeMetrics) int {
	return int(metrics.ProfitPercentAmount.IntPart())
}
