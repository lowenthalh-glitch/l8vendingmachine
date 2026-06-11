package mocks

import (
	"fmt"

	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8vendingmachine/go/vend/common/seeddata"
)

func seedAlarmDefinitions(client *VendClient) error {
	defs := seeddata.GetAlarmDefinitions()
	defList := &alm.AlarmDefinitionList{List: defs}
	fmt.Printf("  Seeding %d alarm definitions...", len(defs))
	_, err := client.Post("/vend/10/AlmDef", defList)
	if err != nil {
		fmt.Printf(" FAILED: %v\n", err)
		return err
	}
	fmt.Printf(" done\n")

	rules := seeddata.GetCorrelationRules()
	ruleList := &alm.CorrelationRuleList{List: rules}
	fmt.Printf("  Seeding %d correlation rules...", len(rules))
	_, err = client.Post("/vend/10/CorrRule", ruleList)
	if err != nil {
		fmt.Printf(" FAILED: %v\n", err)
		return err
	}
	fmt.Printf(" done\n")

	return nil
}
