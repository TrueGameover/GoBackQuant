package trade

import (
	"github.com/TrueGameover/GoBackQuant/graph"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
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
	PositionType    uint
	Open            *graph.Bar
	Closed          *graph.Bar
}

func (p *Position) IsShouldClose(price decimal.Decimal) bool {
	switch p.PositionType {
	case TypeLong:
		return price.LessThanOrEqual(p.StopLossPrice)
	case TypeShort:
		return price.GreaterThanOrEqual(p.StopLossPrice)
	}

	return false
}

type PositionManager struct {
	openPositions   []*Position
	closedPositions []*Position
	counter         uint64
}

func (manager *PositionManager) OpenPosition(positionType uint, price decimal.Decimal, bar *graph.Bar, lot uint, stopLoss decimal.Decimal, takeProfit decimal.Decimal) uint64 {
	position := Position{
		Id:              manager.counter,
		Price:           price,
		StopLossPrice:   stopLoss,
		TakeProfitPrice: takeProfit,
		Lot:             lot,
		PositionType:    positionType,
		Open:            bar,
		Closed:          nil,
	}

	manager.openPositions = append(manager.openPositions, &position)
	manager.counter++

	return position.Id
}

func (manager *PositionManager) ClosePosition(id uint64, bar *graph.Bar) (found bool) {
	found = false
	manager.openPositions = funk.Filter(manager.openPositions, func(position *Position) bool {
		if position.Id == id {
			position.Closed = bar
			found = true
			return false
		}

		return true
	}).([]*Position)

	return
}

func (manager *PositionManager) CloseAll(bar *graph.Bar) {
	for _, position := range manager.openPositions {
		manager.ClosePosition(position.Id, bar)
	}
}
