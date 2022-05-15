package loader

import (
	"context"
	"errors"
	sdk "github.com/Tinkoff/invest-openapi-go-sdk"
	"github.com/TrueGameover/GoBackQuant/pkg/communication/token"
	"github.com/TrueGameover/GoBackQuant/pkg/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/utils"
	"github.com/shopspring/decimal"
	"time"
)

type PartialTickLoader struct {
	Token         token.Token
	sandboxClient utils.Nullable[*sdk.SandboxRestClient]
	client        utils.Nullable[*sdk.RestClient]
	figi          string
	timeFrame     sdk.CandleInterval
	ticks         []graph.Tick
	ticksCount    uint64
}

func (loader *PartialTickLoader) Init(ticker string, timeframe graph.TimeFrame, ctx context.Context) error {
	client := loader.getActualClient()

	loader.timeFrame = loader.getTimeFrame(timeframe)
	loader.ticksCount = 0

	instruments, err := client.InstrumentByTicker(ctx, ticker)
	if err != nil {
		return err
	}

	if len(instruments) == 0 {
		return errors.New("instrument not found")
	}

	loader.figi = instruments[0].FIGI

	return nil
}

func (loader *PartialTickLoader) LoadNext(ctx context.Context, startDate time.Time, endDate time.Time) error {
	client := loader.getActualClient()

	candles, err := client.Candles(ctx, startDate, endDate, loader.timeFrame, loader.figi)
	if err != nil {
		return err
	}

	if len(candles) > 0 {
		loader.ticks = make([]graph.Tick, 0)

		for _, candle := range candles {
			tick := loader.convertCandle(candle)

			loader.ticks = append(loader.ticks, tick)
			loader.ticksCount++
		}
	}

	return nil
}

func (loader *PartialTickLoader) getActualClient() *sdk.RestClient {
	if loader.Token.IsSandbox() {
		return loader.getSandboxClient().RestClient
	}

	return loader.getClient()
}

func (loader *PartialTickLoader) getSandboxClient() *sdk.SandboxRestClient {
	if loader.sandboxClient.HasValue() {
		return loader.sandboxClient.GetValue()
	}

	if loader.Token.IsSandbox() {
		client := sdk.NewSandboxRestClient(loader.Token.GetToken())
		loader.sandboxClient.SetValue(client)

	} else {
		panic("token should be sandbox type")
	}

	return loader.sandboxClient.GetValue()
}

func (loader *PartialTickLoader) getClient() *sdk.RestClient {
	if loader.client.HasValue() {
		return loader.client.GetValue()
	}

	if !loader.Token.IsSandbox() {
		client := sdk.NewRestClient(loader.Token.GetToken())
		loader.client.SetValue(client)

	} else {
		panic("token should not be sandbox type")
	}

	return loader.client.GetValue()
}

func (loader PartialTickLoader) getTimeFrame(frame graph.TimeFrame) sdk.CandleInterval {
	switch frame {
	case graph.TimeFrameM1:
		return sdk.CandleInterval1Min
	case graph.TimeFrameM5:
		return sdk.CandleInterval5Min
	case graph.TimeFrameM10:
		return sdk.CandleInterval10Min
	case graph.TimeFrameM15:
		return sdk.CandleInterval15Min
	case graph.TimeFrameM30:
		return sdk.CandleInterval30Min
	case graph.TimeFrameH1:
		return sdk.CandleInterval1Hour
	case graph.TimeFrameH4:
		return sdk.CandleInterval4Hour
	case graph.TimeFrameD1:
		return sdk.CandleInterval1Day
	case graph.TimeFrameD4:
		break
	case graph.TimeFrameW1:
		return sdk.CandleInterval1Week
	}

	panic("unsupported timeframe")
}

func (loader *PartialTickLoader) GetTicks() []graph.Tick {
	return loader.ticks
}

func (loader *PartialTickLoader) convertCandle(candle sdk.Candle) graph.Tick {
	tick := graph.Tick{
		Id:     loader.ticksCount,
		Date:   candle.TS,
		Open:   decimal.NewFromFloat(candle.OpenPrice),
		High:   decimal.NewFromFloat(candle.HighPrice),
		Low:    decimal.NewFromFloat(candle.LowPrice),
		Close:  decimal.NewFromFloat(candle.ClosePrice),
		Volume: decimal.NewFromFloat(candle.Volume),
	}

	return tick
}
