/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"sync"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/route/optimizer"
)

const parallelWorkers = 20

func ActivateAllServices(creds, dbname string, nic ifs.IVNic) {
	var all []func()
	all = append(all, collectFleetActivations(creds, dbname, nic)...)
	// Inventory is now served by inv_vend (l8inventory cache), not CRUD services
	all = append(all, collectSalesActivations(creds, dbname, nic)...)
	all = append(all, collectPaymentActivations(creds, dbname, nic)...)
	all = append(all, collectTemperatureActivations(creds, dbname, nic)...)
	all = append(all, collectMaintenanceActivations(creds, dbname, nic)...)
	all = append(all, collectRouteActivations(creds, dbname, nic)...)
	all = append(all, collectAnalyticsActivations(creds, dbname, nic)...)
	all = append(all, collectAccessActivations(creds, dbname, nic)...)
	all = append(all, collectDexActivations(creds, dbname, nic)...)
	all = append(all, collectWarehouseActivations(creds, dbname, nic)...)
	all = append(all, collectDashboardActivations(creds, dbname, nic)...)
	all = append(all, collectComplianceActivations(creds, dbname, nic)...)
	all = append(all, collectReportsActivations(creds, dbname, nic)...)
	all = append(all, collectRetentionActivations(creds, dbname, nic)...)

	sem := make(chan struct{}, parallelWorkers)
	var wg sync.WaitGroup

	for _, fn := range all {
		wg.Add(1)
		sem <- struct{}{}
		go func(f func()) {
			defer wg.Done()
			defer func() { <-sem }()
			f()
		}(fn)
	}
	wg.Wait()

	// Command services (non-CRUD) — activated after all entity services are ready
	optimizer.ActivateOptimizer(nic)
}
