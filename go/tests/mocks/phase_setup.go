package mocks

import (
	"fmt"

	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
)

func setupPollarisAndTarget(client *VendClient, simulatorIP string, simulatorPort int32) error {
	// Pollaris config is registered by the parser binary directly (avoids JSON serialization issues)
	fmt.Printf("  Pollaris config registered by parser (skipping HTTP POST)\n")

	// Create target for the Nayax management API (POST as list, same as all other entities)
	target := createSimulatorTarget(simulatorIP, simulatorPort)
	targetList := &l8tpollaris.L8PTargetList{List: []*l8tpollaris.L8PTarget{target}}
	targetEndpoint := fmt.Sprintf("%d/%s", targets.ServiceArea, targets.ServiceName)
	fmt.Printf("  Adding target '%s'...", target.TargetId)
	_, err := client.Post("/vend/"+targetEndpoint, targetList)
	if err != nil {
		fmt.Printf(" FAILED: %v\n", err)
		return err
	}
	fmt.Printf(" done\n")

	return nil
}

func createSimulatorTarget(ip string, port int32) *l8tpollaris.L8PTarget {
	target := &l8tpollaris.L8PTarget{}
	target.TargetId = ip
	target.LinksId = vendcommon.VendMachine_Links_ID
	target.Hosts = make(map[string]*l8tpollaris.L8PHost)
	target.InventoryType = l8tpollaris.L8PTargetType_Vending_Machine
	target.State = l8tpollaris.L8PTargetState_Up

	host := &l8tpollaris.L8PHost{}
	host.HostId = ip
	host.Configs = make(map[int32]*l8tpollaris.L8PHostProtocol)

	restConfig := &l8tpollaris.L8PHostProtocol{}
	restConfig.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
	restConfig.Addr = ip
	restConfig.Port = port
	restConfig.Timeout = 30
	restConfig.Ainfo = &l8tpollaris.AuthInfo{}

	host.Configs[int32(restConfig.Protocol)] = restConfig
	target.Hosts[host.HostId] = host

	return target
}
