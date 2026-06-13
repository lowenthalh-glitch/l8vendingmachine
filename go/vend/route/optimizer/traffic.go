/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Google Maps Directions API integration for traffic-aware route refinement.
 * Uses raw net/http — no external SDK dependency (same pattern as l8collector).
 */
package optimizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
)

const (
	directionsURL = "https://maps.googleapis.com/maps/api/directions/json"
	credentialKey = "google-maps"
	credentialDb  = "google-maps"
)

// directionsResponse is the subset of Google Maps Directions API response we need.
type directionsResponse struct {
	Status string           `json:"status"`
	Routes []directionsRoute `json:"routes"`
}

type directionsRoute struct {
	Legs []directionsLeg `json:"legs"`
}

type directionsLeg struct {
	Distance          directionsValue `json:"distance"`
	Duration          directionsValue `json:"duration"`
	DurationInTraffic directionsValue `json:"duration_in_traffic"`
}

type directionsValue struct {
	Value int64  `json:"value"`
	Text  string `json:"text"`
}

// RefineWithTraffic calls Google Maps Directions API to replace haversine
// estimates with real road distances and traffic-aware durations.
// Falls back silently if no API key or API failure.
func RefineWithTraffic(route *BuiltRoute, startTime int64, config *RouteConfig, nic ifs.IVNic) error {
	apiKey := getAPIKey(nic)
	if apiKey == "" {
		nic.Resources().Logger().Info("Route optimizer: no Google Maps API key, using haversine estimates")
		return nil
	}

	if len(route.Stops) < 2 {
		return nil
	}

	// Build waypoints: origin → waypoints → destination
	origin := fmt.Sprintf("%f,%f", route.Stops[0].Lat, route.Stops[0].Lng)
	destination := fmt.Sprintf("%f,%f", route.Stops[len(route.Stops)-1].Lat, route.Stops[len(route.Stops)-1].Lng)

	var waypoints []string
	for i := 1; i < len(route.Stops)-1; i++ {
		waypoints = append(waypoints, fmt.Sprintf("%f,%f", route.Stops[i].Lat, route.Stops[i].Lng))
	}

	url := fmt.Sprintf("%s?origin=%s&destination=%s&departure_time=%d&traffic_model=best_guess&key=%s",
		directionsURL, origin, destination, startTime, apiKey)
	if len(waypoints) > 0 {
		url += "&waypoints=" + strings.Join(waypoints, "|")
	}

	resp, err := http.Get(url)
	if err != nil {
		nic.Resources().Logger().Info("Route optimizer: Google Maps API error, keeping haversine: ", err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var dirResp directionsResponse
	if err := json.Unmarshal(body, &dirResp); err != nil {
		nic.Resources().Logger().Info("Route optimizer: failed to parse Google Maps response")
		return nil
	}

	if dirResp.Status != "OK" || len(dirResp.Routes) == 0 {
		nic.Resources().Logger().Info("Route optimizer: Google Maps returned status: ", dirResp.Status)
		return nil
	}

	legs := dirResp.Routes[0].Legs
	if len(legs) != len(route.Stops) {
		nic.Resources().Logger().Info("Route optimizer: leg count mismatch, keeping haversine")
		return nil
	}

	// Convert Google legs to RouteLeg
	newLegs := make([]RouteLeg, len(legs))
	for i, leg := range legs {
		distMiles := float64(leg.Distance.Value) / 1609.34
		durSecs := leg.DurationInTraffic.Value
		if durSecs == 0 {
			durSecs = leg.Duration.Value
		}
		newLegs[i] = RouteLeg{
			DistanceMiles:   distMiles,
			DurationSeconds: durSecs,
			IsReload:        route.Stops[i].IsReload,
		}
	}

	// Recalculate metrics using the same ComputeRouteMetrics
	route.Legs = newLegs
	route.Metrics = ComputeRouteMetrics(newLegs, startTime, route.TruckMPG,
		config.FuelPriceGal, config.ServiceMinutes, config.ReloadMinutes)

	return nil
}

func getAPIKey(nic ifs.IVNic) string {
	// Read from the Credentials service (not ShallowSecurityProvider)
	results, err := vendcommon.GetEntities("Creds", byte(75), &l8api.L8Credentials{Id: credentialKey}, nic)
	if err != nil || len(results) == 0 {
		return ""
	}
	cred, ok := results[0].(*l8api.L8Credentials)
	if !ok || cred.Creds == nil {
		return ""
	}
	entry, exists := cred.Creds[credentialDb]
	if !exists || entry == nil {
		return ""
	}
	return entry.Zside
}
