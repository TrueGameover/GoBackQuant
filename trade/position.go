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
	ClosePrice      decimal.Decimal
	StopLossPrice   decimal.Decimal
	TakeProfitPrice decimal.Decimal
	Lot             decimal.Decimal
	PositionType    uint
	Open            *graph.Bar
	Closed          *graph.Bar
}

func (p *Position) IsShouldClose(tick *graph.Tick) bool {
	switch p.PositionType {
	case TypeLong:
		return tick.Close.LessThanOrEqual(p.StopLossPrice) || tick.Close.GreaterThanOrEqual(p.TakeProfitPrice)
	case TypeShort:
		return tick.Close.GreaterThanOrEqual(p.StopLossPrice) || tick.Close.LessThanOrEqual(p.TakeProfitPrice)
	}

	return false
}

func (p *Position) GetPipsAfterClose() decimal.Decimal {
	result := decimal.New(0, 0)

	if p.Closed != nil {
		switch p.PositionType {
		case TypeLong:
			result = p.ClosePrice.Sub(p.Price)
			break
		case TypeShort:
			result = p.Price.Sub(p.ClosePrice)
			break
		}
	}

	return result
}

type PositionManager struct {
	openPositions   []*Position
	closedPositions []*Position
	counter         uint64
}

func (manager *PositionManager) OpenPosition(positionType uint, tick *graph.Tick, bar *graph.Bar, lot decimal.Decimal, stopLoss decimal.Decimal, takeProfit decimal.Decimal) uint64 {
	position := Position{
		Id:              manager.counter,
		Price:           tick.Close,
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

func (manager *PositionManager) ClosePosition(id uint64, bar *graph.Bar) *Position {
	var targetPosition *Position

	manager.openPositions = funk.Filter(manager.openPositions, func(position *Position) bool {
		if position.Id == id {
			barLastTick := bar.GetLastTick()
			if barLastTick == nil {
				position.ClosePrice = bar.Open

			} else {
				position.ClosePrice = barLastTick.Close
			}

			position.Closed = bar
			targetPosition = position
			return false
		}

		return true
	}).([]*Position)

	return targetPosition
}

func (manager *PositionManager) CloseAll(bar *graph.Bar) {
	for _, position := range manager.openPositions {
		manager.ClosePosition(position.Id, bar)
	}
}

func (manager *PositionManager) UpdateForClosePositions(tick *graph.Tick, bar *graph.Bar) []*Position {
	var closedPositions []*Position

	for _, openPosition := range manager.openPositions {
		if openPosition.IsShouldClose(tick) {
			if manager.ClosePosition(openPosition.Id, bar) == nil {
				panic("Can't close position.")
			}

			closedPositions = append(closedPositions, openPosition)
		}
	}

	manager.closedPositions = append(manager.closedPositions, closedPositions...)

	return closedPositions
}

func (manager *PositionManager) GetOpenedPositionsCount() uint {
	return uint(len(manager.openPositions))
}
