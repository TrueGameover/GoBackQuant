package main

import (
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	strategy2 "github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/use_cases/multibacktest"
	"github.com/shopspring/decimal"
	"time"
)

func main() {
	var strategy strategy2.Strategy = &strategy1.TemaAndRStrategy{}
	parameters := []strategy2.Parameter{
		{
			Name: "ind1",
			Min:  1,
			Max:  10,
		},
	}
	err := parameters[0].SetValue(2)
	if err != nil {
		panic(err)
	}
	strategy.UpdateParameters(parameters)

	providers, err := getProviders()
	if err != nil {
		panic(err)
	}

	multiTest := multibacktest.MultiBackTest{}
	results := multiTest.Run(&strategy, providers, decimal.NewFromInt(10000), 4)

	for i, result := range results {
		fmt.Printf("#%d PROM = %s\n", i, result.Metrics.Prom.StringFixed(2))
	}
}

func getProviders() ([][]tick.Provider, error) {
	providers := make([][]tick.Provider, 5)

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
		return nil, err
	}
	providers[0] = make([]tick.Provider, 1)
	providers[0][0] = &csvProvider1

	csvProvider2 := provider.CsvProvider{
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
	err = csvProvider2.Load("SBER_m1.csv", "SBER", graph.TimeFrameM15)
	if err != nil {
		return nil, err
	}
	providers[1] = make([]tick.Provider, 1)
	providers[1][0] = &csvProvider2

	csvProvider3 := provider.CsvProvider{
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
	err = csvProvider3.Load("SBER_m1.csv", "SBER", graph.TimeFrameM15)
	if err != nil {
		return nil, err
	}
	providers[2] = make([]tick.Provider, 1)
	providers[2][0] = &csvProvider3

	csvProvider4 := provider.CsvProvider{
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
	err = csvProvider4.Load("SBER_m1.csv", "SBER", graph.TimeFrameM15)
	if err != nil {
		return nil, err
	}
	providers[3] = make([]tick.Provider, 1)
	providers[3][0] = &csvProvider4

	csvProvider5 := provider.CsvProvider{
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
	err = csvProvider5.Load("SBER_m1.csv", "SBER", graph.TimeFrameM15)
	if err != nil {
		return nil, err
	}
	providers[4] = make([]tick.Provider, 1)
	providers[4][0] = &csvProvider5

	return providers, nil
}
