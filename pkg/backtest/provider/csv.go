package provider

import (
	"bufio"
	"encoding/csv"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
	"io"
	"os"
	"time"
)

type CsvProvider struct {
	reader            *csv.Reader
	DateParseTemplate string
	Delimiter         rune
	Positions         Positions
	FieldsPerRecord   uint
	totalLines        uint64
}

type Positions struct {
	Date   uint
	Open   uint
	High   uint
	Low    uint
	Close  uint
	Volume uint
}

func (provider *CsvProvider) GetNextTick() (*graph.Tick, error) {
	row, err := provider.reader.Read()

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// <DATE>;<TIME>;<OPEN>;<HIGH>;<LOW>;<CLOSE>;<VOL>
	// 20210305;070100;26.5100000;26.5100000;26.5100000;26.5100000;0
	date, err := time.Parse(provider.DateParseTemplate, row[provider.Positions.Date])

	if err != nil {
		return nil, err
	}

	open, err := decimal.NewFromString(row[provider.Positions.Open])

	if err != nil {
		return nil, err
	}

	high, err := decimal.NewFromString(row[provider.Positions.High])

	if err != nil {
		return nil, err
	}

	low, err := decimal.NewFromString(row[provider.Positions.Low])

	if err != nil {
		return nil, err
	}

	closePrice, err := decimal.NewFromString(row[provider.Positions.Close])

	if err != nil {
		return nil, err
	}

	volume, err := decimal.NewFromString(row[provider.Positions.Volume])

	if err != nil {
		return nil, err
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

	return &tick, nil
}

func (provider *CsvProvider) GetTotal() uint64 {
	return provider.totalLines
}

func (provider *CsvProvider) HasTicks() bool {
	return true
}

func (provider *CsvProvider) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := bufio.NewScanner(file)
	provider.totalLines = 0
	for reader.Scan() {
		provider.totalLines++
	}
	// skip title
	provider.totalLines--

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	provider.reader = csv.NewReader(file)
	provider.reader.Comma = ';'
	provider.reader.FieldsPerRecord = int(provider.FieldsPerRecord)

	_, err = provider.reader.Read()
	if err != nil {
		return err
	}

	return nil
}
