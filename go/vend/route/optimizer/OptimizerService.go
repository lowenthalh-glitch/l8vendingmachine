/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Command service for on-demand route optimization.
 * Follows the l8collector ExecuteService pattern (NOT standard ActivateService).
 */
package optimizer

import (
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/web"
	"github.com/saichler/l8vendingmachine/go/types/vend"
)

const (
	ServiceName = "OptRoute"
	ServiceArea = byte(10)
)

type OptimizerService struct {
	serviceArea byte
}

func (this *OptimizerService) Activate(sla *ifs.ServiceLevelAgreement, vnic ifs.IVNic) error {
	this.serviceArea = sla.ServiceArea()
	vnic.Resources().Registry().Register(&vend.VendRouteOptRequest{})
	vnic.Resources().Logger().Info("Route Optimizer service activated")
	return nil
}

func (this *OptimizerService) DeActivate() error {
	return nil
}

func (this *OptimizerService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	req := pb.Element().(*vend.VendRouteOptRequest)

	_, err := GenerateRoutes(vnic, req)
	if err != nil {
		req.Error = err.Error()
	}

	return object.New(nil, req)
}

func (this *OptimizerService) Put(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}

func (this *OptimizerService) Patch(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}

func (this *OptimizerService) Delete(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}

func (this *OptimizerService) Get(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}

func (this *OptimizerService) GetCopy(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}

func (this *OptimizerService) Failed(pb ifs.IElements, vnic ifs.IVNic, msg *ifs.Message) ifs.IElements {
	return nil
}

func (this *OptimizerService) TransactionConfig() ifs.ITransactionConfig {
	return nil
}

func (this *OptimizerService) WebService() ifs.IWebService {
	ws := web.New(ServiceName, this.serviceArea, 0)
	ws.AddEndpoint(&vend.VendRouteOptRequest{}, ifs.POST, &vend.VendRouteOptRequest{})
	return ws
}

// ActivateOptimizer registers the optimizer command service with the vnic.
func ActivateOptimizer(nic ifs.IVNic) {
	sla := ifs.NewServiceLevelAgreement(&OptimizerService{}, ServiceName, ServiceArea, false, nil)
	nic.Resources().Services().Activate(sla, nic)
}
