/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/dashboard/dashboards"
	"github.com/saichler/l8vendingmachine/go/vend/dashboard/kpis"
)

func collectDashboardActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { kpis.Activate(creds, dbname, nic) },
		func() { dashboards.Activate(creds, dbname, nic) },
	}
}
