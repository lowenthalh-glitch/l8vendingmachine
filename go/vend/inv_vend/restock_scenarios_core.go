/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"time"

	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

// Scenario 1: Day-of-week demand patterns
// If the next 1-2 days have historically high demand and current fill is borderline, recommend restock.
func evaluateDayOfWeekDemand(m *vend.VendFleetMachine, profile *vend.VendMachineProfile) *vend.VendRestockRecommendation {
	if len(profile.DowDepletionRate) < 7 {
		return nil
	}

	now := time.Now()
	today := int(now.Weekday())
	tomorrow := (today + 1) % 7
	dayAfter := (today + 2) % 7

	// Average depletion across all days
	totalRate := float64(0)
	for _, r := range profile.DowDepletionRate {
		totalRate += r
	}
	avgRate := totalRate / 7

	if avgRate <= 0 {
		return nil
	}

	// Check if upcoming days have above-average demand
	upcomingRate := (profile.DowDepletionRate[tomorrow] + profile.DowDepletionRate[dayAfter]) / 2
	ratio := upcomingRate / avgRate

	if ratio < 1.3 {
		return nil // not significantly higher than average
	}

	// Project fill% through the high-demand period (48 hours)
	fillPct := calcCurrentFillPct(m)
	if fillPct < 0 {
		return nil
	}

	projectedFill := float64(fillPct) - (upcomingRate * 48)
	if projectedFill > 20 {
		return nil // will survive the high-demand period
	}

	priority := vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_MEDIUM
	if projectedFill < 10 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_HIGH
	}

	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	reason := fmt.Sprintf("%s+%s demand %.1fx average. Projected %.0f%% fill by %s.",
		dayNames[tomorrow], dayNames[dayAfter], ratio, projectedFill, dayNames[dayAfter])

	reasonCode := vend.VendRestockReasonCode_VEND_RESTOCK_REASON_WEEKEND_DEMAND
	if tomorrow >= 1 && tomorrow <= 5 {
		reasonCode = vend.VendRestockReasonCode_VEND_RESTOCK_REASON_WEEKDAY_DEMAND
	}

	return &vend.VendRestockRecommendation{
		RecommendationId: fmt.Sprintf("DOW-%s", m.MachineId),
		MachineId:        m.MachineId,
		MachineName:      m.Name,
		LocationClass:    profile.LocationClass,
		Priority:         priority,
		Reason:           reason,
		ReasonCode:       reasonCode,
		CurrentFillPct:   fillPct,
		ProjectedFillPct: int32(projectedFill),
		Confidence:       0.75,
		AvgDailyRevenue:  int64(profile.AvgDailyRevenue),
	}
}

// Scenario 3: Product-specific fast movers
// Flag when any top product will be empty in < 8 hours.
func evaluateFastMovers(m *vend.VendFleetMachine, profile *vend.VendMachineProfile) *vend.VendRestockRecommendation {
	if len(profile.TopProducts) == 0 {
		return nil
	}

	var criticalProducts []*vend.VendRestockItem
	var minTimeToEmpty float64
	minTimeToEmpty = 999

	for _, pp := range profile.TopProducts {
		if pp.DepletionRatePerHour <= 0 || pp.AvgStock <= 0 {
			continue
		}
		tte := float64(pp.AvgStock) / pp.DepletionRatePerHour
		if tte < 8 {
			criticalProducts = append(criticalProducts, &vend.VendRestockItem{
				ProductName:  pp.ProductName,
				CurrentStock: pp.AvgStock,
				Capacity:     pp.Capacity,
				UnitsToAdd:   pp.Capacity - pp.AvgStock,
				DepletionRate: pp.DepletionRatePerHour,
			})
			if tte < minTimeToEmpty {
				minTimeToEmpty = tte
			}
		}
	}

	if len(criticalProducts) == 0 {
		return nil
	}

	priority := vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_MEDIUM
	if minTimeToEmpty < 3 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_HIGH
	}
	if minTimeToEmpty < 1 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_CRITICAL
	}

	reason := fmt.Sprintf("%d fast-moving products will be empty in %.1f hours. Top: %s (%.1f units/hr).",
		len(criticalProducts), minTimeToEmpty, criticalProducts[0].ProductName, criticalProducts[0].DepletionRate)

	// Revenue at risk: sum of (depletionRate * price * hoursToEmpty) for critical products
	var revenueAtRisk int64
	for _, cp := range criticalProducts {
		pp := findProductProfile(profile, cp.ProductName)
		if pp != nil {
			revenueAtRisk += int64(cp.DepletionRate * float64(pp.Price) * minTimeToEmpty)
		}
	}

	now := time.Now().Unix()
	return &vend.VendRestockRecommendation{
		RecommendationId: fmt.Sprintf("FM-%s", m.MachineId),
		MachineId:        m.MachineId,
		MachineName:      m.Name,
		LocationClass:    profile.LocationClass,
		Priority:         priority,
		Reason:           reason,
		ReasonCode:       vend.VendRestockReasonCode_VEND_RESTOCK_REASON_FAST_MOVER_EMPTY,
		PredictedEmptyTime: now + int64(minTimeToEmpty*3600),
		CurrentFillPct:   calcCurrentFillPct(m),
		RevenueAtRisk:    revenueAtRisk,
		SuggestedProducts: criticalProducts,
		Confidence:       0.8,
		AvgDailyRevenue:  int64(profile.AvgDailyRevenue),
	}
}

// Scenario 10: Critical threshold prediction
// Project fill% forward using profile's hourly depletion rate.
func evaluateCriticalPrediction(m *vend.VendFleetMachine, profile *vend.VendMachineProfile) *vend.VendRestockRecommendation {
	fillPct := calcCurrentFillPct(m)
	if fillPct < 0 || fillPct > 60 {
		return nil // only evaluate machines already somewhat low
	}
	if profile.AvgHourlyDepletion <= 0 {
		return nil
	}

	// Use cascade threshold from profile (or default 15%)
	criticalThreshold := float64(profile.CascadeThresholdPct)
	if criticalThreshold <= 0 {
		criticalThreshold = 15
	}

	// Project forward using hour-of-day rates if available
	projected := float64(fillPct)
	now := time.Now()
	hoursToEmpty := float64(0)

	for h := 0; h < 48; h++ {
		hour := (now.Hour() + h) % 24
		rate := profile.AvgHourlyDepletion
		if len(profile.HodDepletionRate) == 24 && profile.HodDepletionRate[hour] > 0 {
			rate = profile.HodDepletionRate[hour]
		}
		projected -= rate
		hoursToEmpty = float64(h + 1)
		if projected < criticalThreshold {
			break
		}
	}

	if projected >= criticalThreshold {
		return nil // won't hit critical in 48 hours
	}

	priority := vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_LOW
	if hoursToEmpty < 24 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_MEDIUM
	}
	if hoursToEmpty < 6 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_HIGH
	}
	if hoursToEmpty < 2 {
		priority = vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_CRITICAL
	}

	predictedTime := now.Add(time.Duration(hoursToEmpty) * time.Hour).Unix()
	reason := fmt.Sprintf("Predicted to hit %.0f%% fill in %.1f hours (by %s). Current: %d%%.",
		criticalThreshold, hoursToEmpty,
		time.Unix(predictedTime, 0).Format("3:04 PM"),
		fillPct)

	return &vend.VendRestockRecommendation{
		RecommendationId: fmt.Sprintf("CP-%s", m.MachineId),
		MachineId:        m.MachineId,
		MachineName:      m.Name,
		LocationClass:    profile.LocationClass,
		Priority:         priority,
		Reason:           reason,
		ReasonCode:       vend.VendRestockReasonCode_VEND_RESTOCK_REASON_CRITICAL_PREDICTION,
		PredictedEmptyTime: predictedTime,
		CurrentFillPct:   fillPct,
		ProjectedFillPct: int32(projected),
		Confidence:       0.85,
		AvgDailyRevenue:  int64(profile.AvgDailyRevenue),
	}
}

// Scenario 9: Revenue-based priority adjustment
// Top 20% revenue machines get priority +1, bottom 20% get -1.
func applyRevenuePriority(recs []*vend.VendRestockRecommendation, ranks map[string]int32) {
	for _, rec := range recs {
		rank := ranks[rec.MachineId]
		rec.RevenueRank = rank

		if rank <= 20 { // Top 20%
			if rec.Priority < vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_CRITICAL {
				rec.Priority++
			}
			rec.Reason = fmt.Sprintf("[Top %d%% revenue] %s", rank, rec.Reason)
		} else if rank >= 80 { // Bottom 20%
			if rec.Priority > vend.VendRestockPriority_VEND_RESTOCK_PRIORITY_LOW {
				rec.Priority--
			}
		}
	}
}

func calcCurrentFillPct(m *vend.VendFleetMachine) int32 {
	totalStock, totalCapacity := int32(0), int32(0)
	for _, slot := range m.Inventory {
		totalStock += slot.CurrentStock
		totalCapacity += slot.Capacity
	}
	if totalCapacity == 0 {
		return -1
	}
	return int32(float64(totalStock) / float64(totalCapacity) * 100)
}
