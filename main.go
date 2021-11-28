package main

import (
	"github.com/TrueGameover/GoBackQuant/backtesting"
	"github.com/TrueGameover/GoBackQuant/example/strategy1"
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/TrueGameover/GoBackQuant/money"
	"github.com/TrueGameover/GoBackQuant/provider"
	"github.com/TrueGameover/GoBackQuant/trade"
	"github.com/shopspring/decimal"
)

func main() {
	var tickProvider provider.TickProvider = &provider.CsvProvider{}

	balanceManager := money.BalanceManager{}
	balanceManager.SetInitialBalance(decimal.New(10000, 0))
	balanceManager.Reset()

	positionManager := trade.PositionManager{}

	tester := backtesting.StrategyTester{}
	tester.Init(&positionManager, &balanceManager, &tickProvider, graph.TimeFrameM15)

	var strategy backtesting.Strategy = &strategy1.TemaAndRStrategy{}

	tester.Run(&strategy)
}
