package provider

import (
	"encoding/csv"
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/graph"
	"github.com/shopspring/decimal"
	"os"
	"time"
)

type CsvProvider struct {
	TickProvider
	reader *csv.Reader
}

func (provider *CsvProvider) GetNextTick() *graph.Tick {
	row, err := provider.reader.Read()

	if err != nil {
		return nil
	}

	// <DATE>;<TIME>;<OPEN>;<HIGH>;<LOW>;<CLOSE>;<VOL>
	// 20210305;070100;26.5100000;26.5100000;26.5100000;26.5100000;0
	date, err := time.Parse("20060102 150405", row[0]+" "+row[1])

	if err != nil {
		return nil
	}

	open, err := decimal.NewFromString(row[2])

	if err != nil {
		return nil
	}

	high, err := decimal.NewFromString(row[3])

	if err != nil {
		return nil
	}

	low, err := decimal.NewFromString(row[4])

	if err != nil {
		return nil
	}

	closePrice, err := decimal.NewFromString(row[5])

	if err != nil {
		return nil
	}

	volume, err := decimal.NewFromString(row[6])

	if err != nil {
		return nil
	}

	tick := graph.Tick{
		Id:     0,
		Date:   date,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  closePrice,
		Volume: volume,
	}

	return &tick
}

func (provider *CsvProvider) HasTicks() bool {
	return true
}

func (provider *CsvProvider) Load(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	provider.reader = csv.NewReader(file)
	provider.reader.Comma = ';'
	provider.reader.FieldsPerRecord = 7

	_, err = provider.reader.Read()

	if err != nil {
		return err
	}

	return nil
}
