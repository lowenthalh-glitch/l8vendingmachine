/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Router provides road distance/duration via OSRM (local Docker container).
 * Falls back to haversine if OSRM is unavailable.
 * Single source for all distance calculations in the optimizer.
 */
package optimizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DefaultOSRMUrl  = "http://localhost:5000"
	osrmRetryAfter  = 5 * time.Minute
	defaultAvgSpeed = 25.0 // mph fallback for haversine duration
)

// Router provides distance and duration calculations.
// Tries OSRM first, falls back to haversine.
type Router struct {
	osrmURL    string
	httpClient *http.Client
	available  bool
	lastCheck  time.Time
	mtx        sync.Mutex
}

// NewRouter creates a Router with a persistent http.Client.
func NewRouter(osrmURL string) *Router {
	if osrmURL == "" {
		osrmURL = DefaultOSRMUrl
	}
	r := &Router{
		osrmURL: osrmURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     60 * time.Second,
				MaxConnsPerHost:     10,
			},
		},
		available: true,
	}
	return r
}

// Distance returns road distance (miles) and duration (seconds) between two points.
func (r *Router) Distance(lat1, lng1, lat2, lng2 float64) (float64, int64) {
	if r.isAvailable() {
		dist, dur, err := r.osrmRoute(lat1, lng1, lat2, lng2)
		if err == nil {
			return dist, dur
		}
		r.markUnavailable()
	}
	return haversineFallback(lat1, lng1, lat2, lng2)
}

// Matrix returns N×N distance and duration matrices for a set of points.
// points[i] = [lat, lng]
func (r *Router) Matrix(points [][2]float64) ([][]float64, [][]int64) {
	n := len(points)
	if n == 0 {
		return nil, nil
	}
	if r.isAvailable() {
		distM, durM, err := r.osrmTable(points)
		if err == nil {
			return distM, durM
		}
		r.markUnavailable()
	}
	return haversineMatrix(points)
}

// Leg creates a RouteLeg between two points using OSRM or haversine.
func (r *Router) Leg(lat1, lng1, lat2, lng2 float64, isReload bool) RouteLeg {
	dist, dur := r.Distance(lat1, lng1, lat2, lng2)
	return RouteLeg{DistanceMiles: dist, DurationSeconds: dur, IsReload: isReload}
}

func (r *Router) isAvailable() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if !r.available && time.Since(r.lastCheck) > osrmRetryAfter {
		r.available = true // retry
	}
	return r.available
}

func (r *Router) markUnavailable() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.available = false
	r.lastCheck = time.Now()
}

// --- OSRM API calls ---

type osrmRouteResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Distance float64 `json:"distance"` // meters
		Duration float64 `json:"duration"` // seconds
	} `json:"routes"`
}

func (r *Router) osrmRoute(lat1, lng1, lat2, lng2 float64) (float64, int64, error) {
	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=false",
		r.osrmURL, lng1, lat1, lng2, lat2)
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	var result osrmRouteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}
	if result.Code != "Ok" || len(result.Routes) == 0 {
		return 0, 0, fmt.Errorf("OSRM: %s", result.Code)
	}
	distMiles := result.Routes[0].Distance / 1609.34
	durSecs := int64(result.Routes[0].Duration)
	return distMiles, durSecs, nil
}

type osrmTableResponse struct {
	Code      string      `json:"code"`
	Distances [][]float64 `json:"distances"` // meters
	Durations [][]float64 `json:"durations"` // seconds
}

func (r *Router) osrmTable(points [][2]float64) ([][]float64, [][]int64, error) {
	// Build coordinates string: lng,lat;lng,lat;...
	coords := make([]string, len(points))
	for i, p := range points {
		coords[i] = fmt.Sprintf("%f,%f", p[1], p[0]) // OSRM expects lng,lat
	}
	url := fmt.Sprintf("%s/table/v1/driving/%s?annotations=distance,duration",
		r.osrmURL, strings.Join(coords, ";"))
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var result osrmTableResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}
	if result.Code != "Ok" || len(result.Distances) == 0 {
		return nil, nil, fmt.Errorf("OSRM table: %s", result.Code)
	}

	n := len(points)
	distM := make([][]float64, n)
	durM := make([][]int64, n)
	for i := 0; i < n; i++ {
		distM[i] = make([]float64, n)
		durM[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			distM[i][j] = result.Distances[i][j] / 1609.34 // meters → miles
			durM[i][j] = int64(result.Durations[i][j])
		}
	}
	return distM, durM, nil
}

// --- Haversine fallback ---

func haversineFallback(lat1, lng1, lat2, lng2 float64) (float64, int64) {
	dist := Haversine(lat1, lng1, lat2, lng2)
	dur := int64(dist / defaultAvgSpeed * 3600)
	return dist, dur
}

func haversineMatrix(points [][2]float64) ([][]float64, [][]int64) {
	n := len(points)
	distM := make([][]float64, n)
	durM := make([][]int64, n)
	for i := 0; i < n; i++ {
		distM[i] = make([]float64, n)
		durM[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			d := Haversine(points[i][0], points[i][1], points[j][0], points[j][1])
			distM[i][j] = d
			durM[i][j] = int64(d / defaultAvgSpeed * 3600)
		}
	}
	return distM, durM
}
