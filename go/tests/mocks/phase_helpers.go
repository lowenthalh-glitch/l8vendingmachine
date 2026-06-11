package mocks

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func runOp(client *VendClient, label, endpoint string, list interface{}, ids []string, storeIDs *[]string) error {
	fmt.Printf("  Creating %s...", label)
	resp, err := client.Post(endpoint, list)
	if err != nil {
		fmt.Printf(" FAILED\n")
		fmt.Printf("    Error: %v\n", err)
		if resp != "" {
			fmt.Printf("    Response: %s\n", resp)
		}
		return fmt.Errorf("%s: %w", label, err)
	}
	if storeIDs != nil && ids != nil {
		*storeIDs = append(*storeIDs, ids...)
	}
	fmt.Printf(" %d created\n", len(ids))
	return nil
}

func extractIDs(items interface{}, getter func(interface{}) string) []string {
	v := reflect.ValueOf(items)
	ids := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		ids[i] = getter(v.Index(i).Interface())
	}
	return ids
}

func runPhase(label string, fn func() error) {
	fmt.Printf("\n%s\n", label)
	fmt.Printf("%s\n", strings.Repeat("-", len(label)))
	if err := fn(); err != nil {
		fmt.Printf("%s failed: %v\n", label, err)
		os.Exit(1)
	}
}
