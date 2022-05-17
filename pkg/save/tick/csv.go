package tick

import (
	"encoding/csv"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"os"
	"time"
)

type CsvSaver struct {
	file   *os.File
	writer *csv.Writer
}

func (saver *CsvSaver) WriteTo(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	w.Comma = ';'
	w.UseCRLF = true

	err = w.Write([]string{
		"Date",
		"Close",
		"Open",
		"High",
		"Low",
		"Volume",
	})
	if err != nil {
		return err
	}

	saver.writer = w
	saver.file = file
	return nil
}

func (saver *CsvSaver) Close() error {
	saver.writer.Flush()

	err := saver.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (saver *CsvSaver) Write(ticks []*graph.Tick) error {
	for _, tick := range ticks {
		// 'Close', 'Open', 'High', 'Low', 'Volume'
		record := []string{
			tick.Date.Format(time.RFC3339),
			tick.Close.String(),
			tick.Open.String(),
			tick.High.String(),
			tick.Low.String(),
			tick.Volume.String(),
		}

		err := saver.writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
