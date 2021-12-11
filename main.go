package main

import (
	"fmt"
	"github.com/TrueGameover/GoBackQuant/backtesting"
	"github.com/TrueGameover/GoBackQuant/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/TrueGameover/GoBackQuant/money"
	"github.com/TrueGameover/GoBackQuant/provider"
	"github.com/TrueGameover/GoBackQuant/trade"
	"github.com/shopspring/decimal"
)

func main() {
	csvProvider := provider.CsvProvider{}
	err := csvProvider.Load("data/SPFB.SILV-3.22_210305_211124.txt")
	var tickProvider provider.TickProvider = &csvProvider

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

	tester.Run(&strategy)

	history := tester.GetHistorySaver()
	total := history.GetDealsCount()
	profitDealsCount := history.GetProfitDealsCount()

	if total > 0 {
		fmt.Printf("Финальный баланс: %s\n", balanceManager.GetBalance().String())
		fmt.Printf("Всего сделок: %d\n", total)
		fmt.Printf("Успешных сделок: %d\n", history.GetProfitDealsCount())
		fmt.Printf("Убыточных сделок: %d\n", history.GetLossDealsCount())

		profitPercent := float32(profitDealsCount) / float32(total) * 100
		fmt.Printf("Процент успешных сделок: %.2f%% \n", profitPercent)
	}
}
