package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
)

type StrategyTester struct {
	positionManager *trade.PositionManager
	balanceManager  *money.BalanceManager
	tickProvider    *tick.Provider
	graph           *graph.Graph
	historySaver    *trade.HistorySaver
}

func (tester *StrategyTester) Init(positionManager *trade.PositionManager, balanceManager *money.BalanceManager, tickProvider *tick.Provider, timeframe graph.TimeFrame) {
	tester.positionManager = positionManager
	tester.balanceManager = balanceManager
	tester.tickProvider = tickProvider
	tester.graph = &graph.Graph{Timeframe: timeframe}
	tester.historySaver = &trade.HistorySaver{}
}

func (tester *StrategyTester) Run(target *Strategy) error {
	tickProvider := *tester.tickProvider
	strategy := *target

	nextTick, err := tickProvider.GetNextTick()
	if err != nil {
		return err
	}

	for nextTick != nil {
		strategy.BeforeTick(tester.graph)

		tester.graph.AddTick(nextTick)
		closedPositions := tester.positionManager.UpdateForClosePositions(nextTick, tester.graph.GetFreshBar())
		if len(closedPositions) > 0 {
			for _, closedPosition := range closedPositions {
				usedMoney := strategy.GetSingleLotPrice().Mul(strategy.GetLotSize())

				if tester.balanceManager.FreeMoney(usedMoney) {
					balanceDiff := strategy.GetSinglePipPrice().Mul(closedPosition.GetPipsAfterClose()).Div(strategy.GetSinglePipValue())
					tester.balanceManager.AddDiff(balanceDiff)
				}
			}

			tester.historySaver.AddToHistory(closedPositions)
		}

		strategy.Tick(nextTick.Close)

		strategy.AfterTick(tester.graph)

		if strategy.IsOpenPosition() {
			if strategy.GetPositionsLimit() == 0 || (tester.positionManager.GetOpenedPositionsCount() < strategy.GetPositionsLimit() && strategy.GetPositionsLimit() > 0) {
				holdMoney := strategy.GetSingleLotPrice().Mul(strategy.GetLotSize())

				if tester.balanceManager.HoldMoney(holdMoney) {
					tester.balanceManager.Commission(strategy.GetTradeFee())

					tester.positionManager.OpenPosition(
						strategy.GetPositionType(),
						nextTick,
						tester.graph.GetFreshBar(),
						strategy.GetLotSize(),
						strategy.GetStopLoss(nextTick.Close),
						strategy.GetTakeProfit(nextTick.Close),
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

	tester.positionManager.CloseAll(tester.graph.GetFreshBar())
	return nil
}

func (tester *StrategyTester) GetGraph() *graph.Graph {
	return tester.graph
}
func (tester *StrategyTester) GetHistorySaver() *trade.HistorySaver {
	return tester.historySaver
}
