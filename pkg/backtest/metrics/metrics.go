package metrics

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/history"
	"github.com/shopspring/decimal"
	"math"
)

type TradeMetrics struct {
	DealsProfitPercent     decimal.Decimal
	StandardErrorPercent   decimal.Decimal
	MaxBalanceDropPercent  decimal.Decimal
	MaxAbsoluteBalanceDrop decimal.Decimal
	InitialBalance         decimal.Decimal
	MaxBalance             decimal.Decimal
	MinBalance             decimal.Decimal
	ProfitAmount           decimal.Decimal
	ProfitPercentAmount    decimal.Decimal
	FinalBalance           decimal.Decimal
	Prom                   decimal.Decimal
	PromWithoutWinStreak   decimal.Decimal
	DealsCount             int64
}

func GetProfitDealsCount(tradeHistories []*history.TradeHistory) int64 {
	count := int64(0)

	for _, tradeHistory := range tradeHistories {
		for _, deal := range tradeHistory.GetDeals() {
			if deal.Success {
				count++
			}
		}
	}

	return count
}

func GetAverageProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
	profit := decimal.NewFromInt(0)
	two := decimal.NewFromInt(2)

	for _, tradeHistory := range tradeHistories {
		for _, deal := range tradeHistory.GetDeals() {
			if deal.Success {
				profit = profit.Add(deal.BalanceDiff).Div(two)
			}
		}
	}

	return profit
}

func GetAverageLoss(tradeHistories []*history.TradeHistory) decimal.Decimal {
	profit := decimal.NewFromInt(0)
	two := decimal.NewFromInt(2)

	for _, tradeHistory := range tradeHistories {
		for _, deal := range tradeHistory.GetDeals() {
			if !deal.Success {
				profit = profit.Add(deal.BalanceDiff).Div(two)
			}
		}
	}

	return profit
}

func GetProfitDealsPercent(tradeHistories []*history.TradeHistory) decimal.Decimal {
	total := decimal.NewFromInt(GetDealsCount(tradeHistories))
	profit := decimal.NewFromInt(GetProfitDealsCount(tradeHistories))

	return profit.Div(total).Mul(decimal.NewFromInt(100))
}

func GetLossDealsCount(tradeHistories []*history.TradeHistory) int64 {
	count := int64(0)

	for _, tradeHistory := range tradeHistories {
		for _, deal := range tradeHistory.GetDeals() {
			if !deal.Success {
				count++
			}
		}
	}

	return count
}

func GetDealsCount(tradeHistories []*history.TradeHistory) int64 {
	count := int64(0)

	for _, tradeHistory := range tradeHistories {
		count += tradeHistory.GetDealsCount()
	}

	return count
}

// GetStandardError
// стандартная статистическая ошибка
func GetStandardError(tradeHistories []*history.TradeHistory) decimal.Decimal {
	one := decimal.New(1, 0)
	count := decimal.NewFromInt(GetDealsCount(tradeHistories))

	return one.Div(count.Add(one))
}

func findMaxBalance(tradeHistories []*history.TradeHistory) decimal.Decimal {
	max := decimal.NewFromInt(0)
	for _, tradeHistory := range tradeHistories {
		deals := tradeHistory.GetDeals()

		for _, deal := range deals {
			if deal.TotalBalance.GreaterThan(max) {
				max = deal.TotalBalance
			}
		}
	}

	return max
}

func findMinBalance(tradeHistories []*history.TradeHistory) decimal.Decimal {
	min := decimal.NewFromInt(0)
	for _, tradeHistory := range tradeHistories {
		deals := tradeHistory.GetDeals()

		for _, deal := range deals {
			if min.Equals(decimal.Zero) {
				min = deal.TotalBalance
			}

			if deal.TotalBalance.LessThan(min) {
				min = deal.TotalBalance
			}
		}
	}

	return min
}

func findMaxAbsoluteBalanceDrop(tradeHistory *history.TradeHistory) (maxDiff decimal.Decimal, max decimal.Decimal) {
	deals := tradeHistory.GetDeals()
	max = decimal.NewFromInt(0)
	maxDiff = decimal.NewFromInt(0)

	for _, deal := range deals {
		if deal.TotalBalance.GreaterThan(max) {
			max = deal.TotalBalance
		}

		if deal.TotalBalance.LessThan(max) {
			diff := max.Sub(deal.TotalBalance)

			if diff.GreaterThan(maxDiff) {
				maxDiff = diff
			}
		}
	}

	return
}

func GetMaxAbsoluteBalanceDrop(tradeHistories []*history.TradeHistory) decimal.Decimal {
	max := decimal.NewFromInt(0)

	for _, tradeHistory := range tradeHistories {
		diff, _ := findMaxAbsoluteBalanceDrop(tradeHistory)

		if diff.GreaterThan(max) {
			max = diff
		}
	}

	return max
}

func GetMaxPercentBalanceDrop(tradeHistories []*history.TradeHistory) decimal.Decimal {
	maxDiff := decimal.NewFromInt(0)
	max := decimal.NewFromInt(0)

	for _, tradeHistory := range tradeHistories {
		diff, m := findMaxAbsoluteBalanceDrop(tradeHistory)

		if diff.GreaterThan(maxDiff) {
			maxDiff = diff
			max = m
		}
	}

	return maxDiff.Div(max).Mul(decimal.NewFromInt(100))
}

func calculateCorrectedProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
	profitTradesCount := GetProfitDealsCount(tradeHistories)
	averageProfit := GetAverageProfit(tradeHistories)
	correctedProfitTradesCount := float64(profitTradesCount) - math.Sqrt(float64(profitTradesCount))

	return averageProfit.Mul(decimal.NewFromFloat(correctedProfitTradesCount))
}

func calculateCorrectedLoss(tradeHistories []*history.TradeHistory) decimal.Decimal {
	lossTradesCount := GetLossDealsCount(tradeHistories)
	averageLoss := GetAverageLoss(tradeHistories)
	correctedLossTradesCount := float64(lossTradesCount) - math.Sqrt(float64(lossTradesCount))

	return averageLoss.Abs().Mul(decimal.NewFromFloat(correctedLossTradesCount))
}

func GetMaxStreakProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
	maxStreakCount := 0
	maxStreakProfit := decimal.NewFromInt(0)
	streakCount := 0
	streakProfit := decimal.NewFromInt(0)

	for _, tradeHistory := range tradeHistories {
		for _, deal := range tradeHistory.GetDeals() {
			if deal.Success {
				streakCount++
				streakProfit = streakProfit.Add(deal.BalanceDiff)

				if streakCount > maxStreakCount {
					maxStreakCount = streakCount
					maxStreakProfit = streakProfit
				}

			} else {
				streakCount = 0
				streakProfit = decimal.NewFromInt(0)
			}
		}
	}

	return maxStreakProfit
}

func CalculatePROMPercent(tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal) decimal.Decimal {
	correctedProfit := calculateCorrectedProfit(tradeHistories)
	correctedLoss := calculateCorrectedLoss(tradeHistories)

	return correctedProfit.Sub(correctedLoss).Div(initialBalance).Mul(decimal.NewFromInt(100))
}

func CalculatePROMPercentWithoutWinStreak(tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal) decimal.Decimal {
	correctedProfit := calculateCorrectedProfit(tradeHistories)
	correctedLoss := calculateCorrectedLoss(tradeHistories)

	correctedProfit = correctedProfit.Sub(GetMaxStreakProfit(tradeHistories))

	return correctedProfit.Sub(correctedLoss).Div(initialBalance).Mul(decimal.NewFromInt(100))
}

func CalculateMetrics(tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal, finalBalance decimal.Decimal) TradeMetrics {
	return TradeMetrics{
		DealsProfitPercent:     GetProfitDealsPercent(tradeHistories),
		StandardErrorPercent:   GetStandardError(tradeHistories),
		MaxBalanceDropPercent:  GetMaxPercentBalanceDrop(tradeHistories),
		MaxAbsoluteBalanceDrop: GetMaxAbsoluteBalanceDrop(tradeHistories),
		InitialBalance:         initialBalance,
		MaxBalance:             findMaxBalance(tradeHistories),
		MinBalance:             findMinBalance(tradeHistories),
		FinalBalance:           finalBalance,
		Prom:                   CalculatePROMPercent(tradeHistories, initialBalance),
		PromWithoutWinStreak:   CalculatePROMPercentWithoutWinStreak(tradeHistories, initialBalance),
		ProfitAmount:           finalBalance.Sub(initialBalance),
		ProfitPercentAmount:    finalBalance.Sub(initialBalance).Div(initialBalance).Mul(decimal.NewFromInt(100)),
		DealsCount:             GetDealsCount(tradeHistories),
	}
}
