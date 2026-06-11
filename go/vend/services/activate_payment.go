/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/payment/cashpositions"
	"github.com/saichler/l8vendingmachine/go/vend/payment/collections"
)

func collectPaymentActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { cashpositions.Activate(creds, dbname, nic) },
		func() { collections.Activate(creds, dbname, nic) },
	}
}
