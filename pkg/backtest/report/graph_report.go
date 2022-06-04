package report

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/history"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/metrics"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
	template "html/template"
	"os"
	"strconv"
)

//go:embed resources/report.template.gohtml
var goHtmlReportTemplate string

type GraphReport struct {
}

type Metrics struct {
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
	DealsCount             string
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

func (saver *GraphReport) GenerateReport(path string, title string, tradeHistories []*history.TradeHistory, initialBalance decimal.Decimal, finalBalance decimal.Decimal) error {
	reportTemplate, err := template.New("report").Parse(goHtmlReportTemplate)
	if err != nil {
		return err
	}

	type graphData struct {
		CandlesJson      template.JS
		OpenTradesJson   template.JS
		ClosedTradesJson template.JS
		Title            string
	}

	tradeMetrics := metrics.CalculateMetrics(tradeHistories, initialBalance, finalBalance)

	data := struct {
		Metrics   Metrics
		GraphData []graphData
		Title     string
	}{
		Title: title,
		Metrics: Metrics{
			DealsProfitPercent:     tradeMetrics.DealsProfitPercent.StringFixed(2),
			StandardErrorPercent:   tradeMetrics.StandardErrorPercent.StringFixed(2),
			MaxBalanceDropPercent:  tradeMetrics.MaxBalanceDropPercent.StringFixed(2),
			MaxAbsoluteBalanceDrop: tradeMetrics.MaxAbsoluteBalanceDrop.StringFixed(2),
			InitialBalance:         tradeMetrics.InitialBalance.StringFixed(2),
			MaxBalance:             tradeMetrics.MaxBalance.StringFixed(2),
			MinBalance:             tradeMetrics.MinBalance.StringFixed(2),
			ProfitAmount:           tradeMetrics.ProfitAmount.StringFixed(2),
			ProfitPercentAmount:    tradeMetrics.ProfitPercentAmount.StringFixed(2),
			FinalBalance:           tradeMetrics.FinalBalance.StringFixed(2),
			Prom:                   tradeMetrics.Prom.StringFixed(2),
			PromWithoutWinStreak:   tradeMetrics.PromWithoutWinStreak.StringFixed(2),
			DealsCount:             strconv.FormatInt(tradeMetrics.DealsCount, 10),
		},
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
