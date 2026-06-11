/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
package mocks

var MachineModels = []string{
	"TCN-ZK(22SP)+BLH-64S", "TCN-ZK(22SP)+BLH-40S",
	"AF-60C(22SP)", "AF-D900-54C(22SP)",
}

var MachineManufacturers = []string{
	"TCN Vending Technology", "TCN Vending Technology",
	"Afen (Hunan TCN Vending Technology)", "Afen (Hunan TCN Vending Technology)",
}

var LocationNames = []string{
	"Building A - Lobby", "Building C - Break Room", "Gym - Main Entrance",
	"Hospital - Cafeteria Wing", "Airport - Terminal B", "Hotel - Lobby",
	"Mall - Food Court", "Factory - Canteen", "University - Student Center",
	"Train Station - Platform 3",
}

var LocationTypes = []string{
	"OFFICE_LOBBY", "OFFICE_BREAKROOM", "GYM", "HOSPITAL", "AIRPORT",
	"HOTEL", "MALL", "FACTORY", "UNIVERSITY", "TRAIN_STATION",
}

var ProductNames = []string{
	"Coca-Cola 355ml", "Pepsi 355ml", "Sprite 355ml", "Dr Pepper 355ml",
	"Mountain Dew 355ml", "Fanta Orange 355ml", "Aquafina 500ml",
	"Dasani 500ml", "Gatorade Cool Blue 591ml", "Arizona Iced Tea 500ml",
	"Red Bull 250ml", "Monster Energy 473ml", "Celsius Sparkling 355ml",
	"Bang Energy 473ml", "Snickers Bar 52g", "Doritos Nacho 28g",
	"Lay's Classic 28g", "Cheetos Crunchy 28g", "Kit Kat 42g",
	"M&M's Peanut 49g", "Twix Bar 50g", "Trail Mix 40g",
	"Nature Valley Granola 35g", "Planters Peanuts 40g",
	"Turkey & Swiss Sandwich", "Caesar Salad Bowl",
	"Chicken Caesar Wrap", "Greek Yogurt 170g",
	"Mixed Fruit Cup", "Rold Gold Pretzels 28g",
}

var ProductPrices = []int64{
	175, 175, 175, 175, 175, 175, 150, 150, 225, 175,
	350, 325, 275, 300, 200, 175, 175, 175, 200, 200,
	200, 200, 175, 150, 550, 650, 575, 225, 350, 150,
}

var ProductCategories = []int32{
	1, 1, 1, 1, 1, 1, 6, 6, 1, 7,
	2, 2, 2, 2, 3, 3, 3, 3, 5, 5,
	5, 3, 9, 3, 4, 4, 4, 4, 4, 3,
}

var GroupNames = []string{
	"Downtown Office District", "Fitness Centers", "Healthcare Facilities",
	"Transportation Hubs", "Education Campus",
}

var WarehouseNames = []string{
	"Central Distribution Center", "North Regional Warehouse", "South Regional Warehouse",
}

var SupplierNames = []string{
	"Coca-Cola Bottling Co.", "PepsiCo Distribution",
	"Red Bull Distribution", "Frito-Lay Snacks", "Mars Wrigley Confectionery",
}

var DriverFirstNames = []string{
	"James", "Maria", "Robert", "Sarah", "Michael", "Jennifer", "David", "Emily",
}

var DriverLastNames = []string{
	"Rodriguez", "Chen", "Williams", "Johnson", "Martinez", "Anderson", "Taylor", "Brown",
}

var RouteNames = []string{
	"Downtown Loop", "Fitness Circuit", "Medical Mile", "Airport Run",
	"Mall & Hotel Route", "Campus Route", "Factory District",
	"Station Express", "North Side AM", "South Side PM",
	"Weekend Coverage", "Emergency Response",
}
