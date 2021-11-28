package money

import "github.com/shopspring/decimal"

type BalanceManager struct {
	free           decimal.Decimal
	hold           decimal.Decimal
	total          decimal.Decimal
	initialBalance decimal.Decimal
}

func (manager *BalanceManager) HoldMoney(amount decimal.Decimal) bool {
	hold := manager.free.Sub(amount)

	if hold.IsPositive() {
		manager.free = manager.free.Sub(hold)
		manager.hold = manager.hold.Add(hold)
		manager.total = manager.free.Add(manager.hold)

		return true
	}

	return false
}

func (manager *BalanceManager) FreeMoney(amount decimal.Decimal) bool {
	hold := manager.hold.Sub(amount)

	if hold.IsPositive() {
		manager.free = manager.free.Add(hold)
		manager.hold = manager.hold.Sub(hold)
		manager.total = manager.free.Add(manager.hold)

		return true
	}

	return false
}

func (manager *BalanceManager) SetInitialBalance(amount decimal.Decimal) {
	manager.initialBalance = amount
}

func (manager *BalanceManager) Reset() {
	manager.free = manager.initialBalance
}
