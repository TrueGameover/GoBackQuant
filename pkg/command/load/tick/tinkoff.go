package tick

import (
	"context"
	"github.com/TrueGameover/GoBackQuant/pkg/communication/token"
	"github.com/TrueGameover/GoBackQuant/pkg/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/save/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/tinkoff/loader"
	"github.com/schollz/progressbar/v3"
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	"log"
	"time"
)

type TinkoffTicksLoader struct {
	token.Token
	graph.TimeFrame
	LoadStepDate time.Duration
	FromDate     time.Time
	ToDate       time.Time
	Ticker       string
	FilePath     string

	parseStartDate string
	parseEndDate   string
	parseTimeFrame time.Duration
}

func (t *TinkoffTicksLoader) loadTicks(ctx context.Context, tickLoader *loader.PartialTickLoader, start time.Time, end time.Time, logger *log.Logger) error {
	ctx, loaderTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer loaderTimeout()

	err := tickLoader.LoadNext(ctx, start, end)
	if err != nil {
		logger.Panicln(err)
	}

	return nil
}

func (t *TinkoffTicksLoader) Help() cmdy.Help {
	return cmdy.Synopsis(" - ticks loader from tinkoff")
}

func (t *TinkoffTicksLoader) Configure(flags *cmdy.FlagSet, _ *arg.ArgSet) {
	flags.StringVar(&t.FilePath, "output", "", "--output=path_to_file.csv")
	flags.StringVar(&t.Ticker, "ticker", "", "--ticker=SBER")
	flags.StringVar(&t.parseStartDate, "start", "", "--start=2022-01-01")
	flags.StringVar(&t.parseEndDate, "end", "", "--end=2022-12-01")
	flags.DurationVar(&t.LoadStepDate, "step", time.Hour*24, "--step=24h")
	flags.DurationVar(&t.parseTimeFrame, "timeframe", time.Minute*1, "--timeframe=1m")
}

func (t *TinkoffTicksLoader) Run(ctx cmdy.Context) error {
	err := t.parseArgs()
	if err != nil {
		return err
	}

	logger := log.New(ctx.Stdout(), "TinkoffTicksLoader", log.Ltime)
	tickLoader := &loader.PartialTickLoader{Token: t.Token}
	err = tickLoader.Init(t.Ticker, t.TimeFrame, ctx)

	if err != nil {
		logger.Panicln(err)
	}

	saver := tick.CsvSaver{}
	err = saver.WriteTo(t.FilePath)
	if err != nil {
		logger.Panicln(err)
	}
	defer func() {
		err = saver.Close()
		if err != nil {
			logger.Panicln(err)
		}
	}()

	// interval
	start := t.FromDate
	end := t.FromDate
	end = end.Add(t.LoadStepDate)
	bar := progressbar.New64(t.ToDate.Unix() - t.FromDate.Unix())

	for start.Before(t.ToDate) {
		err = t.loadTicks(ctx, tickLoader, start, end, logger)
		if err != nil {
			return err
		}

		err = saver.Write(tickLoader.GetTicks())
		if err != nil {
			logger.Panicln(err)
		}

		start = start.Add(t.LoadStepDate)
		end = end.Add(t.LoadStepDate)

		err = bar.Add64(int64(t.LoadStepDate.Seconds()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TinkoffTicksLoader) parseArgs() error {
	fromDate, err := time.Parse("2006-01-02", t.parseStartDate)
	if err != nil {
		return err
	}
	t.FromDate = fromDate

	toDate, err := time.Parse("2006-01-02", t.parseEndDate)
	if err != nil {
		return err
	}
	t.ToDate = toDate

	frame, err := graph.ParseTimeFrame(t.parseTimeFrame)
	if err != nil {
		return err
	}
	t.TimeFrame = frame

	return nil
}
