package trade

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	TypeLong  uint = 1
	TypeShort      = 2
)

type Position struct {
	Id              uint64
	Price           decimal.Decimal
	StopLossPrice   decimal.Decimal
	TakeProfitPrice decimal.Decimal
	Lot             uint
	Date            time.Time
	PositionType    uint
}

func (p *Position) IsClosed(price decimal.Decimal) bool {
	switch p.PositionType {
	case TypeLong:
		return price.LessThanOrEqual(p.StopLossPrice)
	case TypeShort:
		return price.GreaterThanOrEqual(p.StopLossPrice)
	}

	return false
}
