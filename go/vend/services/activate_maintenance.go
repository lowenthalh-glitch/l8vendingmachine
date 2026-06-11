/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/maintenance/alerts"
	"github.com/saichler/l8vendingmachine/go/vend/maintenance/servicevisits"
	"github.com/saichler/l8vendingmachine/go/vend/maintenance/workorders"
)

func collectMaintenanceActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { alerts.Activate(creds, dbname, nic) },
		func() { workorders.Activate(creds, dbname, nic) },
		func() { servicevisits.Activate(creds, dbname, nic) },
	}
}
