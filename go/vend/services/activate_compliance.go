/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/compliance/certifications"
	"github.com/saichler/l8vendingmachine/go/vend/compliance/findings"
	"github.com/saichler/l8vendingmachine/go/vend/compliance/inspections"
)

func collectComplianceActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { inspections.Activate(creds, dbname, nic) },
		func() { findings.Activate(creds, dbname, nic) },
		func() { certifications.Activate(creds, dbname, nic) },
	}
}
