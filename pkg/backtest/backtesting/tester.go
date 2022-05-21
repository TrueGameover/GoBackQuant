package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/history"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
)

type StrategyTester struct {
	positionManagers []*trade.PositionManager
	balanceManager   *money.BalanceManager
	tickProviders    []tick.Provider
	graphs           []*graph.Graph
	historySavers    []*history.TradeHistory
}

func (tester *StrategyTester) Init(balanceManager *money.BalanceManager, tickProviders []tick.Provider) {
	tester.balanceManager = balanceManager
	tester.tickProviders = tickProviders
	tester.graphs = make([]*graph.Graph, len(tickProviders), len(tickProviders))
	tester.positionManagers = make([]*trade.PositionManager, len(tickProviders), len(tickProviders))

	for i, tickProvider := range tickProviders {
		g := graph.Graph{
			Timeframe: tickProvider.GetTimeFrame(),
			Title:     tickProvider.GetTitle(),
		}

		tester.graphs[i] = &g
		tester.positionManagers[i] = &trade.PositionManager{}
		tester.historySavers[i] = &history.TradeHistory{Graph: &g}
	}
}

func (tester *StrategyTester) Run(s *Strategy) error {
	strategy := *s

	for i, tickProvider := range tester.tickProviders {
		g := tester.graphs[i]
		positionManager := tester.positionManagers[i]
		historySaver := tester.historySavers[i]

		nextTick, err := tickProvider.GetNextTick()
		if err != nil {
			return err
		}

		strategy.BeforeTick(tester.graphs)

		g.AddTick(nextTick)

		closedPositions := positionManager.UpdateForClosePositions(nextTick, g.GetFreshBar())
		if len(closedPositions) > 0 {
			for _, closedPosition := range closedPositions {
				usedMoney := strategy.GetSingleLotPrice(g).Mul(strategy.GetLotSize(g))

				if tester.balanceManager.FreeMoney(usedMoney) {
					balanceDiff := strategy.GetSinglePipPrice(g).Mul(closedPosition.GetPipsAfterClose()).Div(strategy.GetSinglePipValue(g))
					tester.balanceManager.AddDiff(balanceDiff)
				}
			}

			historySaver.AddToHistory(closedPositions)
		}

		strategy.Tick(nextTick, g)

		strategy.AfterTick(tester.graphs)

		if strategy.IsOpenPosition(g) {
			if strategy.GetPositionsLimit(g) == 0 || (positionManager.GetOpenedPositionsCount() < strategy.GetPositionsLimit(g) && strategy.GetPositionsLimit(g) > 0) {
				holdMoney := strategy.GetSingleLotPrice(g).Mul(strategy.GetLotSize(g))

				if tester.balanceManager.HoldMoney(holdMoney) {
					tester.balanceManager.Commission(strategy.GetTradeFee(g))

					positionManager.OpenPosition(
						strategy.GetPositionType(g),
						nextTick,
						g.GetFreshBar(),
						strategy.GetLotSize(g),
						strategy.GetStopLoss(nextTick.Close, g),
						strategy.GetTakeProfit(nextTick.Close, g),
					)
				}
			}
		}

		nextTick, err = tickProvider.GetNextTick()
		if err != nil {
			return err
		}

		if !strategy.ShouldContinue() {
			break
		}
	}

	for i, manager := range tester.positionManagers {
		g := tester.graphs[i]
		manager.CloseAll(g.GetFreshBar())
	}

	return nil
}

func (tester *StrategyTester) GetGraphs() []*graph.Graph {
	return tester.graphs
}
func (tester *StrategyTester) GetHistorySavers() []*history.TradeHistory {
	return tester.historySavers
}
