package metadata

import (
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/metadata"
	"github.com/shopspring/decimal"
)

type TinkoffInstrumentsMetaData struct {
}

func (t *TinkoffInstrumentsMetaData) GetInstrumentType(currentGraph *graph.Graph) metadata.InstrumentType {
	return metadata.Stock
}

func (t *TinkoffInstrumentsMetaData) GetLotSize(currentGraph *graph.Graph) int64 {
	return 1
}

func (t *TinkoffInstrumentsMetaData) GetSingleLotPrice(currentGraph *graph.Graph, tick *graph.Tick) decimal.Decimal {
	return tick.Close
}
