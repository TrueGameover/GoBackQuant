package trade

import "github.com/shopspring/decimal"

type Trade struct {
	Id        uint64
	Success   bool
	MoneyDiff decimal.Decimal
	Position  *Position
}
