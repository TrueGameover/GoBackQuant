package report

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/history"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
	template "html/template"
	"math"
	"os"
)

type GraphReport struct {
}

type highchartsCandle struct {
	X     int64   `json:"x"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

type highchartsTrade struct {
	X     int64  `json:"x"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (saver *GraphReport) GenerateReport(goHtmlTemplate string, path string, title string, tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal, finalBalance decimal.Decimal) error {
	reportTemplate, err := template.New("report").Parse(goHtmlTemplate)
	if err != nil {
		return err
	}

	type graphData struct {
		CandlesJson      template.JS
		OpenTradesJson   template.JS
		ClosedTradesJson template.JS
		Title            string
	}

	data := struct {
		Title                  string
		GraphData              []graphData
		DealsProfitPercent     string
		StandardErrorPercent   string
		MaxBalanceDropPercent  string
		MaxAbsoluteBalanceDrop string
		InitialBalance         string
		MaxBalance             string
		MinBalance             string
		ProfitAmount           string
		ProfitPercentAmount    string
		FinalBalance           string
		Prom                   string
		PromWithoutWinStreak   string
	}{
		Title:                  title,
		DealsProfitPercent:     saver.GetProfitDealsPercent(tradeHistories).StringFixed(2),
		StandardErrorPercent:   saver.GetStandardError(tradeHistories).StringFixed(2),
		MaxBalanceDropPercent:  saver.GetMaxPercentBalanceDrop(tradeHistories).StringFixed(2),
		MaxAbsoluteBalanceDrop: saver.GetMaxAbsoluteBalanceDrop(tradeHistories).StringFixed(2),
		InitialBalance:         initialBalance.StringFixed(2),
		MaxBalance:             saver.findMaxBalance(tradeHistories).StringFixed(2),
		MinBalance:             saver.findMinBalance(tradeHistories).StringFixed(2),
		FinalBalance:           finalBalance.StringFixed(2),
		Prom:                   saver.CalculatePROMPercent(tradeHistories, initialBalance).StringFixed(2),
		PromWithoutWinStreak:   saver.CalculatePROMPercentWithoutWinStreak(tradeHistories, initialBalance).StringFixed(2),
		ProfitAmount:           finalBalance.Sub(initialBalance).StringFixed(2),
		ProfitPercentAmount:    finalBalance.Sub(initialBalance).Div(initialBalance).Mul(decimal.NewFromInt(100)).StringFixed(2),
	}

	for _, tradeHistory := range tradeHistories {
		g := tradeHistory.Graph

		candlesJson, err := saver.prepareCandles(g.GetBars())
		if err != nil {
			return err
		}

		openTradesJson, err := saver.prepareOpenPositions(tradeHistory.GetDeals())
		if err != nil {
			return err
		}

		closedTradesJson, err := saver.prepareClosedPositions(tradeHistory.GetDeals())
		if err != nil {
			return err
		}

		data.GraphData = append(data.GraphData, graphData{
			CandlesJson:      template.JS(candlesJson),
			OpenTradesJson:   template.JS(openTradesJson),
			ClosedTradesJson: template.JS(closedTradesJson),
			Title:            tradeHistory.Graph.Title,
		})
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()

	writer := bufio.NewWriter(file)
	err = reportTemplate.Execute(writer, data)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (saver *GraphReport) prepareCandles(bars []*graph.Bar) (string, error) {
	candles := make([]*highchartsCandle, len(bars), len(bars))

	for i, tick := range bars {
		candle := saver.convertTick(&tick.Tick)

		candles[i] = candle
	}

	candlesJson, err := json.Marshal(&candles)
	if err != nil {
		return "", err
	}

	return string(candlesJson), nil
}

func (saver *GraphReport) prepareOpenPositions(trades []*history.Trade) (string, error) {
	hTrades := make([]*highchartsTrade, len(trades), len(trades))

	for i, t := range trades {
		hTrade := saver.convertOpenTrade(t)

		hTrades[i] = hTrade
	}

	tradesJson, err := json.Marshal(hTrades)
	if err != nil {
		return "", err
	}

	return string(tradesJson), nil
}

func (saver *GraphReport) prepareClosedPositions(trades []*history.Trade) (string, error) {
	hTrades := make([]*highchartsTrade, len(trades), len(trades))

	for i, t := range trades {
		hTrade := saver.convertCloseTrade(t)

		hTrades[i] = hTrade
	}

	tradesJson, err := json.Marshal(hTrades)
	if err != nil {
		return "", err
	}

	return string(tradesJson), nil
}

func (saver *GraphReport) convertTick(tick *graph.Tick) *highchartsCandle {
	candle := highchartsCandle{
		X:     tick.Date.UnixMilli(),
		Open:  tick.Open.InexactFloat64(),
		High:  tick.High.InexactFloat64(),
		Low:   tick.Low.InexactFloat64(),
		Close: tick.Close.InexactFloat64(),
	}

	return &candle
}

func (saver *GraphReport) convertOpenTrade(trade2 *history.Trade) *highchartsTrade {
	t := highchartsTrade{
		X:     trade2.Position.Open.Tick.Date.UnixMilli(),
		Title: "",
		Text:  "",
	}

	if trade2.Position.PositionType == trade.TypeLong {
		t.Title = "Long"
		t.Text = fmt.Sprintf("Position #%d", trade2.Position.Id)

	} else {
		t.Title = "Short"
		t.Text = fmt.Sprintf("Position #%d", trade2.Position.Id)
	}

	return &t
}

func (saver *GraphReport) convertCloseTrade(trade2 *history.Trade) *highchartsTrade {
	t := highchartsTrade{
		X:     trade2.Position.Closed.Tick.Date.UnixMilli(),
		Title: "",
		Text:  "",
	}

	if trade2.Position.PositionType == trade.TypeLong {
		t.Title = "Long closed"
		t.Text = fmt.Sprintf("Position #%d<br>Pips: %s", trade2.Position.Id, trade2.PipsDiff.String())

	} else {
		t.Title = "Short closed"
		t.Text = fmt.Sprintf("Position #%d", trade2.Position.Id)
	}

	return &t
}

func (saver *GraphReport) GetProfitDealsCount(tradeHistories []*history.TradeHistory) int64 {
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

func (saver *GraphReport) GetAverageProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
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

func (saver *GraphReport) GetAverageLoss(tradeHistories []*history.TradeHistory) decimal.Decimal {
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

func (saver *GraphReport) GetProfitDealsPercent(tradeHistories []*history.TradeHistory) decimal.Decimal {
	total := decimal.NewFromInt(saver.GetDealsCount(tradeHistories))
	profit := decimal.NewFromInt(saver.GetProfitDealsCount(tradeHistories))

	return profit.Div(total).Mul(decimal.NewFromInt(100))
}

func (saver *GraphReport) GetLossDealsCount(tradeHistories []*history.TradeHistory) int64 {
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

func (saver *GraphReport) GetDealsCount(tradeHistories []*history.TradeHistory) int64 {
	count := int64(0)

	for _, tradeHistory := range tradeHistories {
		count += tradeHistory.GetDealsCount()
	}

	return count
}

// GetStandardError
// стандартная статистическая ошибка
func (saver *GraphReport) GetStandardError(tradeHistories []*history.TradeHistory) decimal.Decimal {
	one := decimal.New(1, 0)
	count := decimal.NewFromInt(saver.GetDealsCount(tradeHistories))

	return one.Div(count.Add(one))
}

func (saver *GraphReport) findMaxBalance(tradeHistories []*history.TradeHistory) decimal.Decimal {
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

func (saver *GraphReport) findMinBalance(tradeHistories []*history.TradeHistory) decimal.Decimal {
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

func (saver *GraphReport) findMaxAbsoluteBalanceDrop(tradeHistory *history.TradeHistory) (maxDiff decimal.Decimal, max decimal.Decimal) {
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

func (saver *GraphReport) GetMaxAbsoluteBalanceDrop(tradeHistories []*history.TradeHistory) decimal.Decimal {
	max := decimal.NewFromInt(0)

	for _, tradeHistory := range tradeHistories {
		diff, _ := saver.findMaxAbsoluteBalanceDrop(tradeHistory)

		if diff.GreaterThan(max) {
			max = diff
		}
	}

	return max
}

func (saver *GraphReport) GetMaxPercentBalanceDrop(tradeHistories []*history.TradeHistory) decimal.Decimal {
	maxDiff := decimal.NewFromInt(0)
	max := decimal.NewFromInt(0)

	for _, tradeHistory := range tradeHistories {
		diff, m := saver.findMaxAbsoluteBalanceDrop(tradeHistory)

		if diff.GreaterThan(maxDiff) {
			maxDiff = diff
			max = m
		}
	}

	return maxDiff.Div(max).Mul(decimal.NewFromInt(100))
}

func (saver *GraphReport) calculateCorrectedProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
	profitTradesCount := saver.GetProfitDealsCount(tradeHistories)
	averageProfit := saver.GetAverageProfit(tradeHistories)
	correctedProfitTradesCount := float64(profitTradesCount) - math.Sqrt(float64(profitTradesCount))

	return averageProfit.Mul(decimal.NewFromFloat(correctedProfitTradesCount))
}

func (saver *GraphReport) calculateCorrectedLoss(tradeHistories []*history.TradeHistory) decimal.Decimal {
	lossTradesCount := saver.GetLossDealsCount(tradeHistories)
	averageLoss := saver.GetAverageLoss(tradeHistories)
	correctedLossTradesCount := float64(lossTradesCount) - math.Sqrt(float64(lossTradesCount))

	return averageLoss.Abs().Mul(decimal.NewFromFloat(correctedLossTradesCount))
}

func (saver GraphReport) GetMaxStreakProfit(tradeHistories []*history.TradeHistory) decimal.Decimal {
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

func (saver *GraphReport) CalculatePROMPercent(tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal) decimal.Decimal {
	correctedProfit := saver.calculateCorrectedProfit(tradeHistories)
	correctedLoss := saver.calculateCorrectedLoss(tradeHistories)

	return correctedProfit.Sub(correctedLoss).Div(initialBalance).Mul(decimal.NewFromInt(100))
}

func (saver GraphReport) CalculatePROMPercentWithoutWinStreak(tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal) decimal.Decimal {
	correctedProfit := saver.calculateCorrectedProfit(tradeHistories)
	correctedLoss := saver.calculateCorrectedLoss(tradeHistories)

	correctedProfit = correctedProfit.Sub(saver.GetMaxStreakProfit(tradeHistories))

	return correctedProfit.Sub(correctedLoss).Div(initialBalance).Mul(decimal.NewFromInt(100))
}
