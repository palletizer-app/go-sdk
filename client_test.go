package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New("https://palletizer.app")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.baseURL != "https://palletizer.app" {
		t.Errorf("expected baseURL https://palletizer.app, got %s", client.baseURL)
	}
	if client.httpClient == nil {
		t.Fatal("expected non-nil httpClient")
	}
}

func TestPack(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/pack" {
			t.Errorf("expected path /api/v1/pack, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Parse request
		var req PackingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Send mock response
		response := PackingResponse{
			Pallets: []Pallet{
				{
					PalletID:              1,
					TotalWeight:           18143.68,
					TotalHeight:           406.4,
					UtilizationPercentage: 95.0,
					Cartons: []PlacedCarton{
						{
							CartonID: "BOX001_1",
							Position: Point3D{X: 0, Y: 0, Z: 0},
							Dimensions: Dimensions{
								Length: 609.6,
								Width:  457.2,
								Height: 406.4,
							},
							Orientation: "original",
							Weight:      18143.68,
						},
					},
					CenterOfGravity: Point3D{X: 304.8, Y: 228.6, Z: 203.2},
				},
			},
			Summary: PackingSummary{
				TotalPallets:       1,
				TotalCartonsPacked: 1,
				AverageUtilization: 95.0,
				ComputationTimeMs:  5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Create request
	request := &PackingRequest{
		Cartons: []Carton{
			{
				ID:            "BOX001",
				Length:        609.6,
				Width:         457.2,
				Height:        406.4,
				Weight:        18143.68,
				Quantity:      1,
				AllowRotation: true,
			},
		},
		PalletConstraints: StandardPallet(),
		PackingOptions: PackingOptions{
			SupportPercentage: 80.0,
		},
	}

	// Send request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Pack(ctx, request)
	if err != nil {
		t.Fatalf("Pack failed: %v", err)
	}

	// Verify response
	if response.Summary.TotalPallets != 1 {
		t.Errorf("expected 1 pallet, got %d", response.Summary.TotalPallets)
	}
	if response.Summary.TotalCartonsPacked != 1 {
		t.Errorf("expected 1 carton packed, got %d", response.Summary.TotalCartonsPacked)
	}
	if len(response.Pallets) != 1 {
		t.Errorf("expected 1 pallet in response, got %d", len(response.Pallets))
	}
}

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("expected path /api/v1/health, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{Status: "healthy"})
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	health, err := client.Health(ctx)
	if err != nil {
		t.Fatalf("Health failed: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", health.Status)
	}
}

func TestMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/metrics" {
			t.Errorf("expected path /api/v1/metrics, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MetricsResponse{
			TotalRequests: 100,
			TotalCartons:  500,
			TotalPallets:  50,
			AverageTimeMs: 10.5,
		})
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	metrics, err := client.Metrics(ctx)
	if err != nil {
		t.Fatalf("Metrics failed: %v", err)
	}

	if metrics.TotalRequests != 100 {
		t.Errorf("expected 100 requests, got %d", metrics.TotalRequests)
	}
	if metrics.TotalCartons != 500 {
		t.Errorf("expected 500 cartons, got %d", metrics.TotalCartons)
	}
}

func TestStandardPallet(t *testing.T) {
	pallet := StandardPallet()
	if pallet.MaxLength != 1016.0 {
		t.Errorf("expected length 1016.0, got %f", pallet.MaxLength)
	}
	if pallet.MaxWidth != 1828.8 {
		t.Errorf("expected width 1828.8, got %f", pallet.MaxWidth)
	}
	if pallet.MaxHeight != 1219.2 {
		t.Errorf("expected height 1219.2, got %f", pallet.MaxHeight)
	}
	if pallet.MaxWeight != 680388.0 {
		t.Errorf("expected weight 680388.0, got %f", pallet.MaxWeight)
	}
}

func TestConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
		fn       func(float64) float64
	}{
		{"inches to mm", 10.0, 254.0, InchesToMM},
		{"pounds to grams", 10.0, 4535.92, PoundsToGrams},
		{"mm to inches", 254.0, 10.0, MMToInches},
		{"grams to pounds", 4535.92, 10.0, GramsToPounds},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)
			if result < tt.expected-0.01 || result > tt.expected+0.01 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}
