package multibacktest

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/backtesting"
	metrics2 "github.com/TrueGameover/GoBackQuant/pkg/backtest/metrics"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/backtest"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/shopspring/decimal"
	"sync"
)

type MultiBackTest struct {
}

func (t *MultiBackTest) Run(strategy *strategy.Strategy, providers [][]tick.Provider, initialBalance decimal.Decimal, parallelTests int) []backtest.Result {
	results := make([]backtest.Result, len(providers))

	if len(providers) < parallelTests {
		parallelTests = len(providers)
	}

	for i := 0; i < len(providers); i += parallelTests {
		p := make([][]tick.Provider, parallelTests)
		waitGroup := sync.WaitGroup{}

		for j := 0; j < parallelTests; j++ {
			p[j] = providers[j]
		}

		mutex := sync.Mutex{}
		for k, testProviders := range p {
			waitGroup.Add(1)

			go func(index int, testProviders []tick.Provider, mutex *sync.Mutex, results *[]backtest.Result) {
				balanceManager := money.BalanceManager{}
				balanceManager.SetInitialBalance(initialBalance)
				balanceManager.Reset()

				tester := backtesting.StrategyTester{}
				tester.Init(&balanceManager, testProviders)

				clonedStrategy := (*strategy).Clone()

				err := tester.Run(&clonedStrategy)
				result := backtest.Result{}

				if err != nil {
					result.Error = err
				}

				historySavers := tester.GetTradeHistories()
				metrics := metrics2.CalculateMetrics(historySavers, balanceManager.GetInitialBalance(), balanceManager.GetTotalBalance())

				result.Metrics = metrics

				mutex.Lock()
				(*results)[index] = result
				mutex.Unlock()

				waitGroup.Done()
			}(k, testProviders, &mutex, &results)
		}

		waitGroup.Wait()
	}

	return results
}
