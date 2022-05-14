package main

import (
	"context"
	"github.com/TrueGameover/GoBackQuant/pkg/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/save/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/tinkoff/loader"
	"github.com/TrueGameover/GoBackQuant/pkg/tinkoff/token"
	"time"
)

func main() {
	tinkoffToken := &token.TinkoffToken{Token: "t.rBVNYp1_Hla6mH7QredSMHHJJ0ip0bRE5DNf3Q-Cpdd1ipfRMVlcCRHVk28rjXIdyaLi83Dc9M9JwZx1w8Y8uA", Sandbox: true}
	fromDate := time.Date(2021, 3, 5, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2021, 11, 24, 0, 0, 0, 0, time.UTC)

	ctx, cancelTimeout := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancelTimeout()

	tickLoader := &loader.PartialTickLoader{Token: tinkoffToken}
	err := tickLoader.Init("SBER", fromDate, toDate, graph.TimeFrameD1, ctx)

	if err != nil {
		panic(err)
	}

	ctx, loaderTimeout := context.WithTimeout(context.TODO(), 30*time.Second)
	defer loaderTimeout()

	err = tickLoader.LoadNext(ctx)
	if err != nil {
		panic(err)
	}

	saver := tick.CsvSaver{}
	err = saver.WriteTo("test.csv")
	if err != nil {
		panic(err)
	}

	err = saver.Write(tickLoader.GetTicks())
	if err != nil {
		panic(err)
	}

	err = saver.Close()
	if err != nil {

	}
}
