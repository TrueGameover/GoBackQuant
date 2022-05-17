package tick

import (
	"github.com/TrueGameover/GoBackQuant/pkg/backtest/provider"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/tick"
	tick2 "github.com/TrueGameover/GoBackQuant/pkg/save/tick"
	"github.com/schollz/progressbar/v3"
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	"time"
)

type CsvTicksTransformer struct {
	parseInputFilePath   string
	parseOutputFilePath  string
	parseOutputTimeframe time.Duration
	parseInputTimeframe  time.Duration
}

func (c *CsvTicksTransformer) Help() cmdy.Help {
	return cmdy.Synopsis("timeframe converter")
}

func (c *CsvTicksTransformer) Configure(flags *cmdy.FlagSet, _ *arg.ArgSet) {
	flags.StringVar(&c.parseInputFilePath, "input", "", "--input=/path/to/file.csv")
	flags.StringVar(&c.parseOutputFilePath, "output", "", "--output=/path/to/file.csv")
	flags.DurationVar(&c.parseInputTimeframe, "input-timeframe", time.Minute*5, "--input-timeframe=5m")
	flags.DurationVar(&c.parseOutputTimeframe, "output-timeframe", time.Minute*15, "--output-timeframe=15m")
}

func (c *CsvTicksTransformer) Run(_ cmdy.Context) error {
	inputTimeFrame, err := graph.ParseTimeFrame(c.parseInputTimeframe)
	if err != nil {
		return err
	}

	outputTimeFrame, err := graph.ParseTimeFrame(c.parseOutputTimeframe)
	if err != nil {
		return err
	}

	// <DATE>;<OPEN>;<HIGH>;<LOW>;<CLOSE>;<VOL>
	// 2022-01-03 18:00:00 +0000 UTC;303.45;303.36;303.5;303.12;51795
	tickProvider, err := c.getCsvProvider()
	if err != nil {
		return err
	}

	saver := tick2.CsvSaver{}
	err = saver.WriteTo(c.parseOutputFilePath)
	if err != nil {
		return err
	}

	tickTransformer := tick.Transformer{}
	total := int64(tickProvider.GetTotal())
	step := total / int64(outputTimeFrame)
	progressBar := progressbar.New64(total)

	for i := int64(0); i < total; i += step {
		result, err := tickTransformer.Transform(&tickProvider, inputTimeFrame, outputTimeFrame, uint(step))
		if err != nil {
			return err
		}

		err = saver.Write(result.GetTicks())
		if err != nil {
			return err
		}

		err = progressBar.Set64(i)
		if err != nil {
			return err
		}
	}

	err = progressBar.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *CsvTicksTransformer) getCsvProvider() (tick.Provider, error) {
	csvProvider := provider.CsvProvider{
		DateParseTemplate: time.RFC3339,
		Delimiter:         ';',
		Positions: provider.Positions{
			Date:   0,
			Open:   1,
			High:   2,
			Low:    3,
			Close:  4,
			Volume: 5,
		},
		FieldsPerRecord: 6,
	}
	err := csvProvider.Load(c.parseInputFilePath)
	if err != nil {
		return nil, err
	}

	return &csvProvider, nil
}
