/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/locations"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machinegroups"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
)

func collectFleetActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { locations.Activate(creds, dbname, nic) },
		func() { machinegroups.Activate(creds, dbname, nic) },
		func() { machines.Activate(creds, dbname, nic) },
	}
}
