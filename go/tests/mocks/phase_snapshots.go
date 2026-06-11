package mocks

import (
	"fmt"

	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

func seedHistoricalSnapshots(client *VendClient) error {
	// Generate 30 days of snapshots for 20 machines (enough for meaningful charts)
	snapshots := generateHistoricalSnapshots(20)

	// Post in batches of 100 to avoid oversized requests
	batchSize := 100
	total := len(snapshots)
	posted := 0

	fmt.Printf("  Seeding %d historical inventory snapshots...", total)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := snapshots[i:end]
		list := &vend.VendInventorySnapshotList{List: batch}
		_, err := client.Post("/vend/10/InvSnap", list)
		if err != nil {
			fmt.Printf(" FAILED at batch %d: %v\n", i/batchSize, err)
			return err
		}
		posted += len(batch)
	}
	fmt.Printf(" %d posted\n", posted)
	return nil
}
