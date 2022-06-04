package backtesting

import "C"
import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/history"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/money"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/commission"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/metadata"
	strategy2 "github.com/TrueGameover/GoBackQuant/pkg/entities/strategy"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	"github.com/shopspring/decimal"
	"time"
)

type StrategyTester struct {
	positionManagers     []*trade.PositionManager
	balanceManager       *money.BalanceManager
	commissionCalculator commission.Calculator
	instrumentMetaData   metadata.InstrumentMetaData
	tickProviders        []tick.Provider
	graphs               []*graph.Graph
	tradeHistories       []*history.TradeHistory
}

func (tester *StrategyTester) Init(balanceManager *money.BalanceManager, commissionCalculator commission.Calculator, instrumentMetaData metadata.InstrumentMetaData, tickProviders []tick.Provider) {
	tester.balanceManager = balanceManager
	tester.tickProviders = tickProviders
	tester.commissionCalculator = commissionCalculator
	tester.instrumentMetaData = instrumentMetaData
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

func (tester *StrategyTester) Run(s *strategy2.Strategy) error {
	strategy := *s

	var totalMax uint64 = 0
	for _, tickProvider := range tester.tickProviders {
		if totalMax < tickProvider.GetTotal() {
			totalMax = tickProvider.GetTotal()
		}
	}

	strategy.BeforeStart()

	var lastDay *time.Time

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
					lots := decimal.NewFromInt(closedPosition.Lot)
					newTotal := tester.instrumentMetaData.GetSingleLotPrice(g, nextTick).Mul(lots)
					oldTotal := closedPosition.OneLotPrice.Mul(lots)
					balanceDiff := newTotal.Sub(oldTotal)

					if tester.balanceManager.FreeMoney(oldTotal) {
						tester.balanceManager.AddDiff(balanceDiff)
					}

					commissionAmount := tester.GetCommission(closedPosition.InstrumentType, newTotal)
					tester.balanceManager.Commission(commissionAmount)

					historySaver.SaveToHistory(closedPosition, tester.balanceManager.GetTotalBalance(), tester.balanceManager.GetFree(), balanceDiff)
				}
			}

			strategy.Tick(nextTick, g)

			strategy.AfterTick(tester.graphs)

			if strategy.IsOpenPosition(g) {
				if strategy.GetPositionsLimit(g) == 0 || (positionManager.GetOpenedPositionsCount() < strategy.GetPositionsLimit(g) && strategy.GetPositionsLimit(g) > 0) {
					instrumentType := tester.instrumentMetaData.GetInstrumentType(g)
					lot := strategy.GetLotsAmount(g)
					lotSize := tester.instrumentMetaData.GetLotSize(g)
					singleLotPrice := tester.instrumentMetaData.GetSingleLotPrice(g, nextTick)
					holdMoney := singleLotPrice.Mul(decimal.NewFromInt(lot))

					if tester.balanceManager.HoldMoney(holdMoney) {
						commissionAmount := tester.GetCommission(instrumentType, holdMoney)
						tester.balanceManager.Commission(commissionAmount)

						positionManager.OpenPosition(
							strategy.GetPositionType(g),
							tester.instrumentMetaData.GetInstrumentType(g),
							nextTick,
							g.GetFreshBar(),
							lot,
							lotSize,
							singleLotPrice,
							strategy.GetStopLoss(nextTick, g),
							strategy.GetTakeProfit(nextTick, g),
						)
					}
				}
			}

			if lastDay == nil {
				lastDay = &nextTick.Date

			} else {
				diff := nextTick.Date.Sub(*lastDay)

				if diff.Hours() >= 24 {
					// next day
					lastDay = &nextTick.Date
					shortPositions := positionManager.GetOpenedShortPositions()

					for _, shortPosition := range shortPositions {
						total := shortPosition.Price.Mul(decimal.NewFromInt(shortPosition.Lot))
						marginalCommission := tester.commissionCalculator.CalculateDayMarginalCommission(total)

						tester.balanceManager.Commission(marginalCommission)
					}
				}
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

func (tester *StrategyTester) GetCommission(instrumentType metadata.InstrumentType, dealAmount decimal.Decimal) decimal.Decimal {
	commissionAmount := decimal.NewFromInt(0)

	switch instrumentType {
	case metadata.Future:
		commissionAmount = tester.commissionCalculator.CalculateFutureCommission(dealAmount)
		break
	case metadata.Stock:
		commissionAmount = tester.commissionCalculator.CalculateStockCommission(dealAmount)
		break
	}

	return commissionAmount
}
