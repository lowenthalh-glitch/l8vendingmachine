package mocks

import (
	l8m "github.com/saichler/l8common/go/mocks"
	l8common "github.com/saichler/l8common/go/types/l8common"
)

var (
	pickRef          = l8m.PickRef
	randomPastDate   = l8m.RandomPastDate
	randomFutureDate = l8m.RandomFutureDate
	genID            = l8m.GenID
	genCode          = l8m.GenCode
	createAuditInfo  = l8m.CreateAuditInfo
	randomPhone      = l8m.RandomPhone
	sanitizeEmail    = l8m.SanitizeEmail
	minInt           = l8m.MinInt
)

func money(amount int64) *l8common.Money {
	return &l8common.Money{Amount: amount}
}
