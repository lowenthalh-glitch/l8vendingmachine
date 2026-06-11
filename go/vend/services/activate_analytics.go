/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/fleetinventory"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/forecasts"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/performance"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/snapshots"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/topperformers"
	analyticsprofiles "github.com/saichler/l8vendingmachine/go/vend/analytics/profiles"
	"github.com/saichler/l8vendingmachine/go/vend/analytics/restock"
)

func collectAnalyticsActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { forecasts.Activate(creds, dbname, nic) },
		func() { performance.Activate(creds, dbname, nic) },
		func() { fleetinventory.Activate(creds, dbname, nic) },
		func() { snapshots.Activate(creds, dbname, nic) },
		func() { topperformers.Activate(creds, dbname, nic) },
		func() { analyticsprofiles.Activate(creds, dbname, nic) },
		func() { restock.Activate(creds, dbname, nic) },
	}
}
