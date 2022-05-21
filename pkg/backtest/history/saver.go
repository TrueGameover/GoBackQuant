package history

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/trade"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"os"
	"strings"
)

type Trade struct {
	Id        uint64
	Success   bool
	MoneyDiff decimal.Decimal
	Position  *trade.Position
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

type Saver struct {
	counter uint64
	deals   []*Trade
}

func (saver *Saver) AddToHistory(positions []*trade.Position) {
	for _, position := range positions {
		saver.saveToHistory(position)
	}
}

func (saver *Saver) saveToHistory(position *trade.Position) {
	diff := position.GetPipsAfterClose()

	t := Trade{
		Id:        saver.counter,
		Success:   diff.IsPositive(),
		MoneyDiff: diff,
		Position:  position,
	}

	saver.counter++
	saver.deals = append(saver.deals, &t)
}

func (saver *Saver) GetDealsCount() int {
	return len(saver.deals)
}

func (saver *Saver) GetProfitDealsCount() int {
	trades := funk.Filter(saver.deals, func(trade *Trade) bool {
		return trade.Success
	}).([]*Trade)

	return len(trades)
}

func (saver *Saver) GetLossDealsCount() int {
	trades := funk.Filter(saver.deals, func(trade *Trade) bool {
		return !trade.Success
	}).([]*Trade)

	return len(trades)
}

func (saver *Saver) GenerateReport(graph *graph.Graph, template string, path string, title string) error {
	reportHtml := strings.Clone(template)
	reportHtml = strings.ReplaceAll(reportHtml, "{{ TITLE }}", title)

	candlesJson, err := saver.prepareCandles(graph.GetBars())
	if err != nil {
		return err
	}

	reportHtml = strings.ReplaceAll(reportHtml, "{{ JSON_OHLC_DATA }}", candlesJson)

	openTradesJson, err := saver.prepareOpenPositions(saver.deals)
	if err != nil {
		return err
	}

	reportHtml = strings.ReplaceAll(reportHtml, "{{ JSON_TRADES_OPEN }}", openTradesJson)

	closedTradesJson, err := saver.prepareClosedPositions(saver.deals)
	if err != nil {
		return err
	}

	reportHtml = strings.ReplaceAll(reportHtml, "{{ JSON_TRADES_CLOSED }}", closedTradesJson)

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
	_, err = writer.WriteString(reportHtml)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (saver *Saver) prepareCandles(bars []*graph.Bar) (string, error) {
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

func (saver *Saver) prepareOpenPositions(trades []*Trade) (string, error) {
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

func (saver *Saver) prepareClosedPositions(trades []*Trade) (string, error) {
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

func (saver *Saver) convertTick(tick *graph.Tick) *highchartsCandle {
	candle := highchartsCandle{
		X:     tick.Date.UnixMilli(),
		Open:  tick.Open.InexactFloat64(),
		High:  tick.High.InexactFloat64(),
		Low:   tick.Low.InexactFloat64(),
		Close: tick.Close.InexactFloat64(),
	}

	return &candle
}

func (saver Saver) convertOpenTrade(trade2 *Trade) *highchartsTrade {
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

func (saver Saver) convertCloseTrade(trade2 *Trade) *highchartsTrade {
	t := highchartsTrade{
		X:     trade2.Position.Closed.Tick.Date.UnixMilli(),
		Title: "",
		Text:  "",
	}

	if trade2.Position.PositionType == trade.TypeLong {
		t.Title = "Long closed"
		t.Text = fmt.Sprintf("Position #%d<br>Profit: %s", trade2.Position.Id, trade2.MoneyDiff.String())

	} else {
		t.Title = "Short closed"
		t.Text = fmt.Sprintf("Position #%d", trade2.Position.Id)
	}

	return &t
}

func (saver *Saver) findOpenDeal(bar *graph.Bar) *Trade {
	for _, deal := range saver.deals {
		if deal.Position.Open.Id == bar.Id {
			return deal
		}
	}

	return nil
}

func (saver *Saver) findCloseDeal(bar *graph.Bar) *Trade {
	for _, deal := range saver.deals {
		if deal.Position.Closed.Id == bar.Id {
			return deal
		}
	}

	return nil
}
