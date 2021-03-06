package trade

import (
	"github.com/TrueGameover/GoBackQuant/pkg/entities/graph"
	"github.com/TrueGameover/GoBackQuant/pkg/entities/metadata"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
)

type PositionType uint

const (
	TypeLong  PositionType = 1
	TypeShort PositionType = 2
)

type Position struct {
	Id              uint64
	Price           decimal.Decimal
	ClosePrice      decimal.Decimal
	StopLossPrice   decimal.Decimal
	TakeProfitPrice decimal.Decimal
	Lot             int64
	LotSize         int64
	OneLotPrice     decimal.Decimal
	InstrumentType  metadata.InstrumentType
	PositionType    PositionType
	Open            *graph.Bar
	Closed          *graph.Bar
}

func (p *Position) IsShouldClose(tick *graph.Tick) bool {
	switch p.PositionType {
	case TypeLong:
		return tick.Low.LessThanOrEqual(p.StopLossPrice) || tick.High.GreaterThanOrEqual(p.TakeProfitPrice)
	case TypeShort:
		return tick.High.GreaterThanOrEqual(p.StopLossPrice) || tick.Low.LessThanOrEqual(p.TakeProfitPrice)
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

func (manager *PositionManager) OpenPosition(positionType PositionType, instrumentType metadata.InstrumentType, tick *graph.Tick, bar *graph.Bar, lot int64, lotSize int64, oneLotPrice decimal.Decimal, stopLoss decimal.Decimal, takeProfit decimal.Decimal) uint64 {
	position := Position{
		Id:              manager.counter,
		Price:           tick.Close,
		StopLossPrice:   stopLoss,
		TakeProfitPrice: takeProfit,
		Lot:             lot,
		LotSize:         lotSize,
		InstrumentType:  instrumentType,
		PositionType:    positionType,
		Open:            bar,
		Closed:          nil,
		OneLotPrice:     oneLotPrice,
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

func (manager *PositionManager) GetOpenedLongPositions() []*Position {
	var positions []*Position

	for _, position := range manager.openPositions {
		if position.PositionType == TypeLong {
			positions = append(positions, position)
		}
	}

	return positions
}

func (manager *PositionManager) GetOpenedShortPositions() []*Position {
	var positions []*Position

	for _, position := range manager.openPositions {
		if position.PositionType == TypeShort {
			positions = append(positions, position)
		}
	}

	return positions
}

func (manager *PositionManager) GetOpenedLongPositionsCount() int {
	count := 0

	for _, position := range manager.openPositions {
		if position.PositionType == TypeLong {
			count++
		}
	}

	return count
}

func (manager *PositionManager) GetOpenedShortPositionsCount() int {
	count := 0

	for _, position := range manager.openPositions {
		if position.PositionType == TypeShort {
			count++
		}
	}

	return count
}
