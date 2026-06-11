/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
package main

import (
	"flag"
	"github.com/saichler/l8vendingmachine/go/tests/mocks"
)

func main() {
	address := flag.String("address", "http://localhost:8080", "Vending server address")
	user := flag.String("user", "admin", "Username")
	password := flag.String("password", "admin", "Password")
	insecure := flag.Bool("insecure", false, "Skip TLS verification")
	simulator := flag.String("simulator", "", "Nayax simulator IP (e.g., 192.168.200.1)")
	simulatorPort := flag.Int("simulator-port", 8443, "Nayax simulator HTTPS port")
	flag.Parse()
	mocks.RunMockGenerator(*address, *user, *password, *insecure, *simulator, int32(*simulatorPort))
}
