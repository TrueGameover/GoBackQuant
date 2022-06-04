package commission

import "github.com/shopspring/decimal"

type Calculator interface {
	// CalculateDayMarginalCommission
	// Marginal commission per day
	CalculateDayMarginalCommission(marginalAmount decimal.Decimal) decimal.Decimal

	CalculateStockCommission(dealAmount decimal.Decimal) decimal.Decimal

	CalculateFutureCommission(dealAmount decimal.Decimal) decimal.Decimal
}
