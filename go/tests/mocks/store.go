/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
package mocks

// MockDataStore holds IDs for business-layer entities only.
// Machine inventory data comes from the Nayax simulator via the collection pipeline.
type MockDataStore struct {
	LocationIDs  []string
	GroupIDs     []string
	FacilityIDs []string
	SupplierIDs  []string
	TruckIDs     []string
	DriverIDs    []string
}
