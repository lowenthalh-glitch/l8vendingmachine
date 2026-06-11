/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
package mocks

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"
)

func RunMockGenerator(address, user, password string, insecure bool, simulatorIP string, simulatorPort int32) {
	fmt.Printf("L8 Vending Machine Mock Data Generator\n")
	fmt.Printf("=======================================\n")
	fmt.Printf("Server: %s\n", address)
	fmt.Printf("User: %s\n", user)
	if insecure {
		fmt.Printf("TLS: Insecure (certificate verification disabled)\n")
	}
	if simulatorIP != "" {
		fmt.Printf("Simulator: %s:%d\n", simulatorIP, simulatorPort)
	}
	fmt.Println()

	httpClient := &http.Client{Timeout: 30 * time.Second}
	if insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := NewVendClient(address, httpClient)
	err := client.Authenticate(user, password)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Authentication successful\n\n")

	store := &MockDataStore{}
	RunAllPhases(client, store)

	if simulatorIP != "" {
		runPhase("Setup: Pollaris & Simulator Target", func() error {
			return setupPollarisAndTarget(client, simulatorIP, simulatorPort)
		})
	}

	PrintSummary(store)
}
