package money

import "github.com/shopspring/decimal"

type BalanceManager struct {
	free           decimal.Decimal
	hold           decimal.Decimal
	total          decimal.Decimal
	initialBalance decimal.Decimal
}

func (manager *BalanceManager) HoldMoney(amount decimal.Decimal) bool {
	free := manager.free.Sub(amount)

	if free.IsPositive() {
		manager.free = free
		manager.hold = manager.hold.Add(amount)
		manager.total = manager.free.Add(manager.hold)

		return true
	}

	return false
}

func (manager *BalanceManager) FreeMoney(amount decimal.Decimal) bool {
	hold := manager.hold.Sub(amount)

	if hold.IsPositive() {
		manager.free = manager.free.Add(amount)
		manager.hold = hold
		manager.total = manager.free.Add(manager.hold)

		return true
	}

	return false
}

func (manager *BalanceManager) Commission(amount decimal.Decimal) {
	manager.free = manager.free.Sub(amount)
}

func (manager *BalanceManager) SetInitialBalance(amount decimal.Decimal) {
	manager.initialBalance = amount
}

func (manager *BalanceManager) Reset() {
	manager.free = manager.initialBalance
}
