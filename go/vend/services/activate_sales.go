/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/sales/settlements"
	"github.com/saichler/l8vendingmachine/go/vend/sales/transactions"
)

func collectSalesActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { transactions.Activate(creds, dbname, nic) },
		func() { settlements.Activate(creds, dbname, nic) },
	}
}
