// Package client provides a Go client for the Palletizer API.
//
// The client can be used to interact with a Palletizer service hosted at
// https://palletizer.app/ or any other endpoint.
//
// Example usage:
//
//	client := client.New("https://palletizer.app")
//
//	request := &client.PackingRequest{
//	    Cartons: []client.Carton{
//	        {
//	            ID:            "BOX001",
//	            Length:        609.6,
//	            Width:         457.2,
//	            Height:        406.4,
//	            Weight:        18143.68,
//	            Quantity:      30,
//	            AllowRotation: true,
//	        },
//	    },
//	    PalletConstraints: client.StandardPallet(),
//	    PackingOptions: client.PackingOptions{
//	        SupportPercentage: 80.0,
//	    },
//	}
//
//	response, err := client.Pack(context.Background(), request)
//	if err != nil {
//	    log.Fatal(err)
//	}
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Palletizer API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new Palletizer API client
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// NewWithHTTPClient creates a client with a custom HTTP client
func NewWithHTTPClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// Carton represents a carton to be packed
type Carton struct {
	ID            string  `json:"id"`
	Length        float64 `json:"length"`         // millimeters
	Width         float64 `json:"width"`          // millimeters
	Height        float64 `json:"height"`         // millimeters
	Weight        float64 `json:"weight"`         // grams
	Quantity      int     `json:"quantity"`       // number of identical cartons
	Fragile       bool    `json:"fragile"`        // whether carton is fragile
	AllowRotation bool    `json:"allow_rotation"` // whether carton can be rotated
}

// PalletConstraints defines the maximum dimensions and weight for a pallet
type PalletConstraints struct {
	MaxLength float64 `json:"max_length"` // millimeters
	MaxWidth  float64 `json:"max_width"`  // millimeters
	MaxHeight float64 `json:"max_height"` // millimeters
	MaxWeight float64 `json:"max_weight"` // grams
}

// PackingOptions configures the packing algorithm
type PackingOptions struct {
	SupportPercentage float64 `json:"support_percentage"` // minimum support area percentage (0-100)
}

// PackingRequest is the request sent to the Pack API
type PackingRequest struct {
	Cartons           []Carton          `json:"cartons"`
	PalletConstraints PalletConstraints `json:"pallet_constraints"`
	PackingOptions    PackingOptions    `json:"packing_options"`
}

// Point3D represents a 3D coordinate
type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Dimensions represents 3D dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// PlacedCarton represents a carton placed on a pallet
type PlacedCarton struct {
	CartonID    string     `json:"carton_id"`
	Position    Point3D    `json:"position"`
	Dimensions  Dimensions `json:"dimensions"`
	Orientation string     `json:"orientation"`
	Weight      float64    `json:"weight"`
}

// Pallet represents a packed pallet
type Pallet struct {
	PalletID              int            `json:"pallet_id"`
	TotalWeight           float64        `json:"total_weight"`
	TotalHeight           float64        `json:"total_height"`
	UtilizationPercentage float64        `json:"utilization_percentage"`
	Cartons               []PlacedCarton `json:"cartons"`
	CenterOfGravity       Point3D        `json:"center_of_gravity"`
}

// PackingSummary provides overall statistics
type PackingSummary struct {
	TotalPallets       int     `json:"total_pallets"`
	TotalCartonsPacked int     `json:"total_cartons_packed"`
	AverageUtilization float64 `json:"average_utilization"`
	ComputationTimeMs  int     `json:"computation_time_ms"`
}

// PackingResponse is the response from the Pack API
type PackingResponse struct {
	Pallets []Pallet       `json:"pallets"`
	Summary PackingSummary `json:"summary"`
	Error   string         `json:"error,omitempty"`
}

// HealthResponse is the response from the Health API
type HealthResponse struct {
	Status string `json:"status"`
}

// MetricsResponse is the response from the Metrics API
type MetricsResponse struct {
	TotalRequests  int     `json:"total_requests"`
	TotalCartons   int     `json:"total_cartons"`
	TotalPallets   int     `json:"total_pallets"`
	AverageTimeMs  float64 `json:"average_time_ms"`
	AverageUtilPct float64 `json:"average_util_pct"`
	SuccessRate    float64 `json:"success_rate"`
	UptimeSeconds  int     `json:"uptime_seconds"`
	MemoryAllocMB  float64 `json:"memory_alloc_mb"`
	MemorySysMB    float64 `json:"memory_sys_mb"`
	NumGoroutines  int     `json:"num_goroutines"`
	NumGC          uint32  `json:"num_gc"`
	LastGCPauseMs  float64 `json:"last_gc_pause_ms"`
	GoVersion      string  `json:"go_version"`
	BuildVersion   string  `json:"build_version"`
	BuildTime      string  `json:"build_time"`
}

// Pack sends a packing request and returns the packed pallets
func (c *Client) Pack(ctx context.Context, request *PackingRequest) (*PackingResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/pack", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response PackingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error != "" {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, response.Error)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return &response, nil
}

// Health checks if the API is healthy
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &health, nil
}

// Metrics retrieves API metrics
func (c *Client) Metrics(ctx context.Context) (*MetricsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/metrics", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics request failed with status %d", resp.StatusCode)
	}

	var metrics MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &metrics, nil
}

// StandardPallet returns constraints for a standard 40x72x48 inch pallet (1500 lbs)
func StandardPallet() PalletConstraints {
	return PalletConstraints{
		MaxLength: 1016.0,   // 40 inches
		MaxWidth:  1828.8,   // 72 inches
		MaxHeight: 1219.2,   // 48 inches
		MaxWeight: 680388.0, // 1500 lbs
	}
}

// StandardPallet4048 returns constraints for a 40x48x48 inch pallet (1500 lbs)
func StandardPallet4048() PalletConstraints {
	return PalletConstraints{
		MaxLength: 1016.0,   // 40 inches
		MaxWidth:  1219.2,   // 48 inches
		MaxHeight: 1219.2,   // 48 inches
		MaxWeight: 680388.0, // 1500 lbs
	}
}

// InchesToMM converts inches to millimeters
func InchesToMM(inches float64) float64 {
	return inches * 25.4
}

// PoundsToGrams converts pounds to grams
func PoundsToGrams(pounds float64) float64 {
	return pounds * 453.592
}

// MMToInches converts millimeters to inches
func MMToInches(mm float64) float64 {
	return mm / 25.4
}

// GramsToPounds converts grams to pounds
func GramsToPounds(grams float64) float64 {
	return grams / 453.592
}
