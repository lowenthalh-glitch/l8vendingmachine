/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/route/drivers"
	"github.com/saichler/l8vendingmachine/go/vend/route/routes"
	"github.com/saichler/l8vendingmachine/go/vend/route/trucks"
)

func collectRouteActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { routes.Activate(creds, dbname, nic) },
		func() { drivers.Activate(creds, dbname, nic) },
		func() { trucks.Activate(creds, dbname, nic) },
	}
}
