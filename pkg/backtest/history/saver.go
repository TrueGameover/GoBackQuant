package history

import (
	"bufio"
	_ "embed"
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

func (saver *Saver) GenerateReport(graph *graph.Graph, template string, path string) error {
	reportHtml := strings.Clone(template)
	json, err := saver.prepareData(graph.GetTicks())

	reportHtml = strings.ReplaceAll(reportHtml, "{{ JSON_OHLC_DATA }}", json)

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

func (saver *Saver) prepareData(ticks []*graph.Tick) (string, error) {
	builder := strings.Builder{}
	_, err := builder.WriteRune('[')
	if err != nil {
		return "", err
	}

	for _, tick := range ticks {
		tickStr := saver.convertTick(tick)

		if builder.Len() > 1 {
			_, err = builder.WriteRune(',')
			if err != nil {
				return "", err
			}
		}

		_, err = builder.WriteString(tickStr)
		if err != nil {
			return "", err
		}
	}

	_, err = builder.WriteRune(']')
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (saver Saver) convertTick(tick *graph.Tick) string {
	return fmt.Sprintf(
		"[%d, %s, %s, %s, %s]",
		tick.Date.UnixMilli(),
		tick.Open.String(),
		tick.High.String(),
		tick.Low.String(),
		tick.Close.String(),
	)
}
