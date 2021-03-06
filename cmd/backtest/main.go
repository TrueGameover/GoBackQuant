package main

import (
	_ "embed"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/backtesting"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/report"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	strategy2 "github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	commission2 "github.com/TrueGameover/GoBackQuant/pkg/tinkoff/commission"
	"github.com/TrueGameover/GoBackQuant/pkg/tinkoff/metadata"
	"github.com/shopspring/decimal"
	"time"
)

func main() {
	providers, err := getProviders()
	if err != nil {
		panic(err)
	}

	balanceManager := money.BalanceManager{}
	balanceManager.SetInitialBalance(decimal.NewFromInt(10000))
	balanceManager.Reset()

	commissionCalculator := &commission2.TinkoffCommissionCalculator{}
	metaData := &metadata.TinkoffInstrumentsMetaData{}

	tester := backtesting.StrategyTester{}
	tester.Init(&balanceManager, commissionCalculator, metaData, providers)

	var strategy strategy2.Strategy = &strategy1.TemaAndRStrategy{}
	parameters := []strategy2.Parameter{
		{
			Name: "ind1",
			Min:  1,
			Max:  10,
		},
	}
	err = parameters[0].SetValue(2)
	if err != nil {
		panic(err)
	}
	strategy.UpdateParameters(parameters)

	err = tester.Run(&strategy)
	if err != nil {
		panic(err)
	}

	graphReport := report.GraphReport{}
	historySavers := tester.GetTradeHistories()
	//total := historySavers.GetDealsCount()
	//profitDealsCount := historySavers.GetProfitDealsCount()

	err = graphReport.GenerateReport("report.html", "Tradebot", historySavers, balanceManager.GetInitialBalance(), balanceManager.GetTotalBalance())
	if err != nil {
		panic(err)
	}

	/*if total > 0 {
		fmt.Printf("Финальный баланс: %s\n", balanceManager.GetTotalBalance().String())
		fmt.Printf("Всего сделок: %d\n", total)
		fmt.Printf("Успешных сделок: %d\n", historySavers.GetProfitDealsCount())
		fmt.Printf("Убыточных сделок: %d\n", historySavers.GetLossDealsCount())

		profitPercent := float32(profitDealsCount) / float32(total) * 100
		fmt.Printf("Процент успешных сделок: %.2f%% \n", profitPercent)
	}*/
}

func getProviders() ([]tick.Provider, error) {
	providers := make([]tick.Provider, 3)

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
	providers[0] = &csvProvider1

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
	providers[1] = &csvProvider2

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
	providers[2] = &csvProvider3

	return providers, nil
}
