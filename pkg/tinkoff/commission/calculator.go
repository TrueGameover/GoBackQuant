package commission

import (
	"github.com/shopspring/decimal"
)

type TinkoffCommissionCalculator struct {
}

func (t *TinkoffCommissionCalculator) CalculateDayMarginalCommission(marginalAmount decimal.Decimal) decimal.Decimal {
	switch true {
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(5000)):
		return decimal.NewFromInt(0)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(50000)):
		return decimal.NewFromInt(35)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(100000)):
		return decimal.NewFromInt(70)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(200000)):
		return decimal.NewFromInt(135)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(300000)):
		return decimal.NewFromInt(200)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(500000)):
		return decimal.NewFromInt(320)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(1000000)):
		return decimal.NewFromInt(620)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(2000000)):
		return decimal.NewFromInt(1200)
	case marginalAmount.LessThanOrEqual(decimal.NewFromInt(5000000)):
		return decimal.NewFromInt(3000)
	}

	percent := decimal.NewFromFloat(0.058).Div(decimal.NewFromInt(100))

	return marginalAmount.Mul(percent)
}

func (t *TinkoffCommissionCalculator) CalculateStockCommission(dealAmount decimal.Decimal) decimal.Decimal {
	percent := decimal.NewFromFloat(0.04).Div(decimal.NewFromInt(100))

	return dealAmount.Mul(percent)
}

func (t *TinkoffCommissionCalculator) CalculateFutureCommission(dealAmount decimal.Decimal) decimal.Decimal {
	percent := decimal.NewFromFloat(0.04).Div(decimal.NewFromInt(100))

	return dealAmount.Mul(percent)
}
