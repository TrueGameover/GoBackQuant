package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/graph"
)

type StrategyTester struct {
	positionManager *trade.PositionManager
	balanceManager  *money.BalanceManager
	tickProvider    *provider.TickProvider
	graph           *graph.Graph
	historySaver    *trade.HistorySaver
}

func (tester *StrategyTester) Init(positionManager *trade.PositionManager, balanceManager *money.BalanceManager, tickProvider *provider.TickProvider, timeframe graph.TimeFrame) {
	tester.positionManager = positionManager
	tester.balanceManager = balanceManager
	tester.tickProvider = tickProvider
	tester.graph = &graph.Graph{Timeframe: timeframe}
	tester.historySaver = &trade.HistorySaver{}
}

func (tester *StrategyTester) Run(target *Strategy) {
	tickProvider := *tester.tickProvider
	strategy := *target

	tick := tickProvider.GetNextTick()

	for tick != nil {
		strategy.BeforeTick(tester.graph)

		tester.graph.Tick(tick)
		closedPositions := tester.positionManager.UpdateForClosePositions(tick, tester.graph.GetFreshBar())
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

		strategy.Tick(tick.Close)

		strategy.AfterTick(tester.graph)

		if strategy.IsOpenPosition() {
			if strategy.GetPositionsLimit() == 0 || (tester.positionManager.GetOpenedPositionsCount() < strategy.GetPositionsLimit() && strategy.GetPositionsLimit() > 0) {
				holdMoney := strategy.GetSingleLotPrice().Mul(strategy.GetLotSize())

				if tester.balanceManager.HoldMoney(holdMoney) {
					tester.balanceManager.Commission(strategy.GetTradeFee())

					tester.positionManager.OpenPosition(
						strategy.GetPositionType(),
						tick,
						tester.graph.GetFreshBar(),
						strategy.GetLotSize(),
						strategy.GetStopLoss(tick.Close),
						strategy.GetTakeProfit(tick.Close),
					)
				}
			}
		}

		tick = tickProvider.GetNextTick()

		if !strategy.ShouldContinue() {
			break
		}
	}

	tester.positionManager.CloseAll(tester.graph.GetFreshBar())
}

func (tester *StrategyTester) GetGraph() *graph.Graph {
	return tester.graph
}
func (tester *StrategyTester) GetHistorySaver() *trade.HistorySaver {
	return tester.historySaver
}
