# Palletizer Go Client

The official Go client library for the Palletizer API hosted at https://palletizer.app/

## Installation

```bash
go get github.com/yourusername/palletizer/pkg/client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yourusername/palletizer/pkg/client"
)

func main() {
    // Create client
    c := client.New("https://palletizer.app")

    // Create packing request
    request := &client.PackingRequest{
        Cartons: []client.Carton{
            {
                ID:            "BOX001",
                Length:        609.6,    // 24 inches in mm
                Width:         457.2,    // 18 inches in mm
                Height:        406.4,    // 16 inches in mm
                Weight:        18143.68, // 40 lbs in grams
                Quantity:      30,
                AllowRotation: true,
            },
        },
        PalletConstraints: client.StandardPallet(),
        PackingOptions: client.PackingOptions{
            SupportPercentage: 80.0,
        },
    }

    // Pack the cartons
    response, err := c.Pack(context.Background(), request)
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    fmt.Printf("Packed %d cartons onto %d pallets\n",
        response.Summary.TotalCartonsPacked,
        response.Summary.TotalPallets)
    fmt.Printf("Average utilization: %.2f%%\n",
        response.Summary.AverageUtilization)
}
```

## Using Imperial Units

The API uses metric units (millimeters and grams), but the client provides helper functions for imperial units:

```go
carton := client.Carton{
    ID:            "BOX001",
    Length:        client.InchesToMM(24),    // 24 inches
    Width:         client.InchesToMM(18),    // 18 inches
    Height:        client.InchesToMM(16),    // 16 inches
    Weight:        client.PoundsToGrams(40), // 40 pounds
    Quantity:      30,
    AllowRotation: true,
}
```

## Standard Pallet Sizes

```go
// 40x72x48 inch pallet (1500 lbs)
constraints := client.StandardPallet()

// 40x48x48 inch pallet (1500 lbs)
constraints := client.StandardPallet4048()
```

## Multiple Carton Types

```go
request := &client.PackingRequest{
    Cartons: []client.Carton{
        {
            ID:            "LARGE_BOX",
            Length:        client.InchesToMM(24),
            Width:         client.InchesToMM(18),
            Height:        client.InchesToMM(16),
            Weight:        client.PoundsToGrams(40),
            Quantity:      20,
            AllowRotation: true,
        },
        {
            ID:            "MEDIUM_BOX",
            Length:        client.InchesToMM(18),
            Width:         client.InchesToMM(12),
            Height:        client.InchesToMM(12),
            Weight:        client.PoundsToGrams(20),
            Quantity:      30,
            AllowRotation: true,
        },
        {
            ID:            "SMALL_BOX",
            Length:        client.InchesToMM(12),
            Width:         client.InchesToMM(8),
            Height:        client.InchesToMM(8),
            Weight:        client.PoundsToGrams(10),
            Quantity:      50,
            AllowRotation: true,
        },
    },
    PalletConstraints: client.StandardPallet(),
    PackingOptions: client.PackingOptions{
        SupportPercentage: 80.0,
    },
}
```

## Processing Results

```go
response, err := c.Pack(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

// Summary
fmt.Printf("Total Pallets: %d\n", response.Summary.TotalPallets)
fmt.Printf("Total Cartons Packed: %d\n", response.Summary.TotalCartonsPacked)
fmt.Printf("Average Utilization: %.2f%%\n", response.Summary.AverageUtilization)
fmt.Printf("Computation Time: %d ms\n", response.Summary.ComputationTimeMs)

// Iterate through pallets
for _, pallet := range response.Pallets {
    fmt.Printf("\nPallet %d:\n", pallet.PalletID)
    fmt.Printf("  Weight: %.2f lbs\n", client.GramsToPounds(pallet.TotalWeight))
    fmt.Printf("  Height: %.2f inches\n", client.MMToInches(pallet.TotalHeight))
    fmt.Printf("  Utilization: %.2f%%\n", pallet.UtilizationPercentage)
    fmt.Printf("  Cartons: %d\n", len(pallet.Cartons))
    
    // Center of gravity
    fmt.Printf("  Center of Gravity: (%.1f, %.1f, %.1f) mm\n",
        pallet.CenterOfGravity.X,
        pallet.CenterOfGravity.Y,
        pallet.CenterOfGravity.Z)
    
    // List cartons
    for _, carton := range pallet.Cartons {
        fmt.Printf("    - %s at (%.1f, %.1f, %.1f) [%s]\n",
            carton.CartonID,
            carton.Position.X,
            carton.Position.Y,
            carton.Position.Z,
            carton.Orientation)
    }
}
```

## Health Check

```go
health, err := c.Health(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", health.Status)
```

## Metrics

```go
metrics, err := c.Metrics(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
fmt.Printf("Average Time: %.2f ms\n", metrics.AverageTimeMs)
fmt.Printf("Success Rate: %.2f%%\n", metrics.SuccessRate*100)
```

## Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:    10,
        IdleConnTimeout: 30 * time.Second,
    },
}

c := client.NewWithHTTPClient("https://palletizer.app", httpClient)
```

## Context and Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := c.Pack(ctx, request)
```

## Error Handling

```go
response, err := c.Pack(context.Background(), request)
if err != nil {
    // Network error or timeout
    log.Printf("Request failed: %v", err)
    return
}

if response.Error != "" {
    // API returned an error (e.g., oversized cartons)
    log.Printf("Packing error: %s", response.Error)
    return
}
```

## Units Reference

| Measurement | Unit | Conversion |
|------------|------|------------|
| Length, Width, Height | millimeters (mm) | 1 inch = 25.4 mm |
| Weight | grams (g) | 1 pound = 453.592 g |

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/pack` | POST | Pack cartons onto pallets |
| `/api/v1/health` | GET | Health check |
| `/api/v1/metrics` | GET | Service metrics |

## Support

- Documentation: https://palletizer.app/docs
- Issues: https://github.com/yourusername/palletizer/issues
- Email: support@palletizer.app
