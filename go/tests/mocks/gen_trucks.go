package mocks

import (
	"math/rand"

	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// Last known GPS positions for trucks (scattered around Austin TX)
var truckPositions = [][2]float64{
	{30.2710, -97.7510}, {30.3050, -97.7380}, {30.2430, -97.7620},
	{30.3200, -97.7100}, {30.2560, -97.7490}, {30.2900, -97.7250},
	{30.3500, -97.7050}, {30.2100, -97.7550}, {30.2780, -97.7400},
	{30.2350, -97.7300},
}

func generateTrucks(store *MockDataStore) []*vend.VendDeliveryTruck {
	items := make([]*vend.VendDeliveryTruck, len(TruckNames))
	for i, name := range TruckNames {
		makeIdx := i % len(TruckMakes)
		year := int32(2018 + rand.Intn(8))
		mpg := 8.0 + rand.Float64()*14.0

		// Status distribution: 60% Active, 20% En-Route, 10% Maintenance, 10% Decommissioned
		var status vend.VendTruckStatus
		switch {
		case i < 6:
			status = vend.VendTruckStatus_VEND_TRUCK_STATUS_ACTIVE
		case i < 8:
			status = vend.VendTruckStatus_VEND_TRUCK_STATUS_EN_ROUTE
		case i < 9:
			status = vend.VendTruckStatus_VEND_TRUCK_STATUS_MAINTENANCE
		default:
			status = vend.VendTruckStatus_VEND_TRUCK_STATUS_DECOMMISSIONED
		}

		// Fuel type distribution
		fuelTypes := []vend.VendFuelType{
			vend.VendFuelType_VEND_FUEL_TYPE_DIESEL,
			vend.VendFuelType_VEND_FUEL_TYPE_GASOLINE,
			vend.VendFuelType_VEND_FUEL_TYPE_DIESEL,
			vend.VendFuelType_VEND_FUEL_TYPE_DIESEL,
			vend.VendFuelType_VEND_FUEL_TYPE_GASOLINE,
			vend.VendFuelType_VEND_FUEL_TYPE_HYBRID,
			vend.VendFuelType_VEND_FUEL_TYPE_DIESEL,
			vend.VendFuelType_VEND_FUEL_TYPE_ELECTRIC,
			vend.VendFuelType_VEND_FUEL_TYPE_DIESEL,
			vend.VendFuelType_VEND_FUEL_TYPE_GASOLINE,
		}

		truckTypes := []vend.VendTruckType{
			vend.VendTruckType_VEND_TRUCK_TYPE_BOX_TRUCK,
			vend.VendTruckType_VEND_TRUCK_TYPE_CARGO_VAN,
			vend.VendTruckType_VEND_TRUCK_TYPE_REFRIGERATED,
			vend.VendTruckType_VEND_TRUCK_TYPE_SPRINTER,
			vend.VendTruckType_VEND_TRUCK_TYPE_BOX_TRUCK,
			vend.VendTruckType_VEND_TRUCK_TYPE_CARGO_VAN,
			vend.VendTruckType_VEND_TRUCK_TYPE_BOX_TRUCK,
			vend.VendTruckType_VEND_TRUCK_TYPE_REFRIGERATED,
			vend.VendTruckType_VEND_TRUCK_TYPE_SPRINTER,
			vend.VendTruckType_VEND_TRUCK_TYPE_PICKUP,
		}

		items[i] = &vend.VendDeliveryTruck{
			TruckId:                genID("trk", i),
			PlateNumber:            randomPlate(),
			Vin:                    randomVin(i),
			Name:                   name,
			Make:                   TruckMakes[makeIdx],
			Model:                  TruckModels[makeIdx],
			Year:                   year,
			Type:                   truckTypes[i],
			CargoCapacityCuFt:      int32(400 + rand.Intn(401)),
			MaxPayloadLbs:          int32(3000 + rand.Intn(5001)),
			FuelType:               fuelTypes[i],
			Mileage:                int32(5000 + rand.Intn(95001)),
			MilesPerGallon:         mpg,
			Status:                 status,
			HomeDepotId:            pickRef(store.FacilityIDs, i),
			LastLatitude:           truckPositions[i][0],
			LastLongitude:          truckPositions[i][1],
			LastLocationUpdate:     randomPastDate(1, 1),
			LastMaintenanceDate:    randomPastDate(3, 30),
			NextMaintenanceDate:    randomFutureDate(3, 30),
			NextMaintenanceMileage: int32(50000 + rand.Intn(50001)),
			InsuranceExpiry:        randomFutureDate(12, 30),
			RegistrationExpiry:     randomFutureDate(12, 30),
			RefrigerationEquipped:  truckTypes[i] == vend.VendTruckType_VEND_TRUCK_TYPE_REFRIGERATED,
			CashCollectionEquipped: i < 7,
			CoinChangerEquipped:    i < 5,
			Stock:                  generateTruckStock(i),
			AuditInfo:              createAuditInfo(),
		}
	}
	return items
}

// Simulator product catalog — matches the Nayax simulator's inventory slots exactly.
var simulatorProducts = []struct {
	name     string
	sku      string
	price    int64
	maxQty   int32
}{
	{"Coca-Cola 12oz", "SKU-CC12", 175, 48},
	{"Pepsi 12oz", "SKU-PP12", 175, 48},
	{"Snickers Bar", "SKU-SN01", 200, 36},
	{"Lay's Classic Chips", "SKU-LC01", 225, 36},
	{"Monster Energy 16oz", "SKU-ME16", 350, 24},
	{"Bottled Water 16oz", "SKU-BW16", 150, 60},
	{"KitKat Bar", "SKU-KK01", 200, 36},
	{"Doritos Nacho", "SKU-DN01", 225, 36},
	{"Red Bull 8.4oz", "SKU-RB08", 325, 24},
	{"Gatorade Fruit Punch", "SKU-GF20", 250, 36},
	{"M&M's Peanut", "SKU-MM01", 200, 48},
	{"Pringles Original", "SKU-PR01", 275, 24},
	{"Tropicana OJ 10oz", "SKU-TJ10", 225, 36},
	{"Reese's Cups", "SKU-RC01", 200, 48},
	{"Smartwater 20oz", "SKU-SW20", 250, 36},
	{"Nature Valley Granola", "SKU-NV01", 175, 60},
}

func generateTruckStock(truckIdx int) []*vend.VendTruckStockItem {
	stock := make([]*vend.VendTruckStockItem, len(simulatorProducts))
	for j, p := range simulatorProducts {
		// Vary stock levels: some trucks fully loaded, some partially
		fillPct := 0.4 + rand.Float64()*0.6 // 40-100% fill
		qty := int32(float64(p.maxQty) * fillPct)
		stock[j] = &vend.VendTruckStockItem{
			ProductName: p.name,
			Sku:         p.sku,
			Price:       p.price,
			Quantity:    qty,
			MaxQuantity: p.maxQty,
		}
	}
	return stock
}

func randomPlate() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	plate := make([]byte, 7)
	for j := 0; j < 3; j++ {
		plate[j] = letters[rand.Intn(len(letters))]
	}
	plate[3] = '-'
	for j := 4; j < 7; j++ {
		plate[j] = digits[rand.Intn(len(digits))]
	}
	return string(plate)
}

func randomVin(seed int) string {
	chars := "0123456789ABCDEFGHJKLMNPRSTUVWXYZ"
	vin := make([]byte, 17)
	for j := range vin {
		vin[j] = chars[(seed*17+j*13)%len(chars)]
	}
	return string(vin)
}
