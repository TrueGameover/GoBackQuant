package money

import "github.com/shopspring/decimal"

type BalanceManager struct {
	free  decimal.Decimal
	hold  decimal.Decimal
	total decimal.Decimal
}

func (manager *BalanceManager) holdMoney(amount decimal.Decimal) bool {
	hold := manager.free.Sub(amount)

	if hold.IsPositive() {
		manager.free = manager.free.Sub(hold)
		manager.hold = manager.hold.Add(hold)
		manager.total = manager.free.Add(manager.hold)
	}

	return false
}

func (manager *BalanceManager) freeMoney(amount decimal.Decimal) bool {
	hold := manager.hold.Sub(amount)

	if hold.IsPositive() {
		manager.free = manager.free.Add(hold)
		manager.hold = manager.hold.Sub(hold)
		manager.total = manager.free.Add(manager.hold)
	}

	return false
}
