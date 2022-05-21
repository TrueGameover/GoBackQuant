package main

import (
	_ "embed"
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/backtesting"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/shopspring/decimal"
	"time"
)

//go:embed resources/report.template.html
var reportTemplate string

func main() {
	csvProvider := provider.CsvProvider{
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
	err := csvProvider.Load("SBER_m1.csv")
	var tickProvider tick.Provider = &csvProvider

	if err != nil {
		panic(err)
	}

	balanceManager := money.BalanceManager{}
	balanceManager.SetInitialBalance(decimal.NewFromInt(10000))
	balanceManager.Reset()

	positionManager := trade.PositionManager{}

	tester := backtesting.StrategyTester{}
	tester.Init(&positionManager, &balanceManager, &tickProvider, graph.TimeFrameM15)

	var strategy backtesting.Strategy = &strategy1.TemaAndRStrategy{}

	err = tester.Run(&strategy)
	if err != nil {
		panic(err)
	}

	history := tester.GetHistorySaver()
	total := history.GetDealsCount()
	profitDealsCount := history.GetProfitDealsCount()

	err = history.GenerateReport(tester.GetGraph(), reportTemplate, "report.html", "SBER")
	if err != nil {
		panic(err)
	}

	if total > 0 {
		fmt.Printf("Финальный баланс: %s\n", balanceManager.GetBalance().String())
		fmt.Printf("Всего сделок: %d\n", total)
		fmt.Printf("Успешных сделок: %d\n", history.GetProfitDealsCount())
		fmt.Printf("Убыточных сделок: %d\n", history.GetLossDealsCount())

		profitPercent := float32(profitDealsCount) / float32(total) * 100
		fmt.Printf("Процент успешных сделок: %.2f%% \n", profitPercent)
	}
}
