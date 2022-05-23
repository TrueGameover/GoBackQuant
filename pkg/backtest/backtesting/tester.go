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
	tradeHistories   []*history.TradeHistory
}

func (tester *StrategyTester) Init(balanceManager *money.BalanceManager, tickProviders []tick.Provider) {
	tester.balanceManager = balanceManager
	tester.tickProviders = tickProviders
	tester.graphs = make([]*graph.Graph, len(tickProviders), len(tickProviders))
	tester.positionManagers = make([]*trade.PositionManager, len(tickProviders), len(tickProviders))
	tester.tradeHistories = make([]*history.TradeHistory, len(tickProviders), len(tickProviders))

	for i, tickProvider := range tickProviders {
		g := graph.Graph{
			Timeframe: tickProvider.GetTimeFrame(),
			Title:     tickProvider.GetTitle(),
		}

		tester.graphs[i] = &g
		tester.positionManagers[i] = &trade.PositionManager{}
		tester.tradeHistories[i] = &history.TradeHistory{Graph: &g}
	}
}

func (tester *StrategyTester) Run(s *Strategy) error {
	strategy := *s

	var totalMax uint64 = 0
	for _, tickProvider := range tester.tickProviders {
		if totalMax < tickProvider.GetTotal() {
			totalMax = tickProvider.GetTotal()
		}
	}

	strategy.BeforeStart()

	for totalStep := uint64(0); totalStep < totalMax; totalStep++ {
		for i, tickProvider := range tester.tickProviders {
			g := tester.graphs[i]
			positionManager := tester.positionManagers[i]
			historySaver := tester.tradeHistories[i]

			nextTick, err := tickProvider.GetNextTick()
			if err != nil {
				return err
			}

			if nextTick == nil {
				// end
				continue
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

				historySaver.AddToHistory(closedPositions, tester.balanceManager.GetTotalBalance(), tester.balanceManager.GetFree())
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
							strategy.GetStopLoss(nextTick, g),
							strategy.GetTakeProfit(nextTick, g),
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
func (tester *StrategyTester) GetTradeHistories() []*history.TradeHistory {
	return tester.tradeHistories
}
