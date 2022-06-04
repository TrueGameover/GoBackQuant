package metadata

import (
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/shopspring/decimal"
)

type InstrumentType uint

const (
	Future InstrumentType = 1
	Stock  InstrumentType = 2
)

type InstrumentMetaData interface {
	GetInstrumentType(currentGraph *graph.Graph) InstrumentType
	GetLotSize(currentGraph *graph.Graph) int64
	GetSingleLotPrice(currentGraph *graph.Graph, tick *graph.Tick) decimal.Decimal
}
