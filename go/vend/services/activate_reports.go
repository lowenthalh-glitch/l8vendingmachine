/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/reports/reports"
)

func collectReportsActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { reports.Activate(creds, dbname, nic) },
	}
}
