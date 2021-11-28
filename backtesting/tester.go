package backtesting

import (
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/TrueGameover/GoBackQuant/money"
	"github.com/TrueGameover/GoBackQuant/provider"
	"github.com/TrueGameover/GoBackQuant/trade"
)

type StrategyTester struct {
	positionManager *trade.PositionManager
	balanceManager  *money.BalanceManager
	tickProvider    *provider.TickProvider
	graph           *graph.Graph
	historySaver    *trade.HistorySaver
}

func (tester *StrategyTester) Init(positionManager *trade.PositionManager, balanceManager *money.BalanceManager, tickProvider *provider.TickProvider, timeframe uint) {
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

		closedPositions := tester.positionManager.UpdateForClosePositions(tick, tester.graph.GetCurrentBar())
		if len(closedPositions) > 0 {
			tester.historySaver.AddToHistory(closedPositions)
		}

		strategy.Tick(tick.Close)

		tester.graph.Tick(tick)

		strategy.AfterTick(tester.graph)

		if strategy.IsOpenPosition() {
			if strategy.GetPositionsLimit() == 0 || tester.positionManager.GetOpenedPositionsCount() < strategy.GetPositionsLimit() {
				holdMoney := strategy.GetSingleLotPrice().Mul(strategy.GetLotSize()).Add(strategy.GetTradeFee())

				if tester.balanceManager.HoldMoney(holdMoney) {
					tester.positionManager.OpenPosition(
						strategy.GetPositionType(),
						tick,
						tester.graph.GetCurrentBar(),
						strategy.GetLotSize(),
						strategy.GetStopLoss(),
						strategy.GetTakeProfit(),
					)
				}
			}
		}

		tick = tickProvider.GetNextTick()
	}
}

func (tester *StrategyTester) GetGraph() *graph.Graph {
	return tester.graph
}
func (tester *StrategyTester) GetHistorySaver() *trade.HistorySaver {
	return tester.historySaver
}
