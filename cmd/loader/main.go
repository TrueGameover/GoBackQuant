package main

import (
	"context"
	"encoding/json"
	tick2 "github.com/TrueGameover/GoBackQuant/pkg/command/load/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/command/transform/tick"
	"github.com/TrueGameover/GoBackQuant/pkg/tinkoff/token"
	"github.com/shabbyrobe/cmdy"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Tinkoff token.TinkoffToken
}

func main() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := Configuration{}
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	err = cmdy.Run(context.Background(), os.Args[1:], func() cmdy.Command {
		return cmdy.NewGroup("", cmdy.Builders{
			"download": func() cmdy.Command {
				return &tick2.TinkoffTicksLoader{
					Token: &config.Tinkoff,
				}
			},
			"transform": func() cmdy.Command {
				return &tick.CsvTicksTransformer{}
			},
		})
	})

	if err != nil {
		cmdy.Fatal(err)
	}
}
