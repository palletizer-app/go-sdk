# Palletizer Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/palletizer-app/sdk.svg)](https://pkg.go.dev/github.com/palletizer-app/sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/palletizer-app/sdk)](https://goreportcard.com/report/github.com/palletizer-app/sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go client library for the [Palletizer.app](https://palletizer.app) API - a high-performance 3D bin packing service for warehouse and logistics operations.

## Installation

```bash
go get github.com/palletizer-app/sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/palletizer-app/sdk"
)

func main() {
    // Create client
    client := palletizer.New("https://api.palletizer.app")

    // Create packing request
    request := &palletizer.PackingRequest{
        Cartons: []palletizer.Carton{
            {
                ID:            "BOX001",
                Length:        palletizer.InchesToMM(24),    // 24 inches
                Width:         palletizer.InchesToMM(18),    // 18 inches
                Height:        palletizer.InchesToMM(16),    // 16 inches
                Weight:        palletizer.PoundsToGrams(40), // 40 lbs
                Quantity:      30,
                AllowRotation: true,
            },
        },
        PalletConstraints: palletizer.StandardPallet(),
        PackingOptions: palletizer.PackingOptions{
            SupportPercentage: 80.0,
        },
    }

    // Pack the cartons
    response, err := client.Pack(context.Background(), request)
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

The API uses metric units (millimeters and grams), but the SDK provides helper functions:

```go
carton := palletizer.Carton{
    ID:            "BOX001",
    Length:        palletizer.InchesToMM(24),    // 24 inches → 609.6 mm
    Width:         palletizer.InchesToMM(18),    // 18 inches → 457.2 mm
    Height:        palletizer.InchesToMM(16),    // 16 inches → 406.4 mm
    Weight:        palletizer.PoundsToGrams(40), // 40 lbs → 18143.68 g
    Quantity:      30,
    AllowRotation: true,
}
```

Convert results back to imperial:

```go
heightInches := palletizer.MMToInches(pallet.TotalHeight)
weightPounds := palletizer.GramsToPounds(pallet.TotalWeight)
```

## Standard Pallet Sizes

```go
// 40×72×48 inch pallet (1500 lbs) - most common
constraints := palletizer.StandardPallet()

// 40×48×48 inch pallet (1500 lbs)
constraints := palletizer.StandardPallet4048()
```

## Multiple Carton Types

```go
request := &palletizer.PackingRequest{
    Cartons: []palletizer.Carton{
        {
            ID:            "LARGE_BOX",
            Length:        palletizer.InchesToMM(24),
            Width:         palletizer.InchesToMM(18),
            Height:        palletizer.InchesToMM(16),
            Weight:        palletizer.PoundsToGrams(40),
            Quantity:      20,
            AllowRotation: true,
        },
        {
            ID:            "MEDIUM_BOX",
            Length:        palletizer.InchesToMM(18),
            Width:         palletizer.InchesToMM(12),
            Height:        palletizer.InchesToMM(12),
            Weight:        palletizer.PoundsToGrams(20),
            Quantity:      30,
            AllowRotation: true,
        },
        {
            ID:            "SMALL_BOX",
            Length:        palletizer.InchesToMM(12),
            Width:         palletizer.InchesToMM(8),
            Height:        palletizer.InchesToMM(8),
            Weight:        palletizer.PoundsToGrams(10),
            Quantity:      50,
            AllowRotation: true,
        },
    },
    PalletConstraints: palletizer.StandardPallet(),
    PackingOptions: palletizer.PackingOptions{
        SupportPercentage: 80.0,
    },
}
```

## Processing Results

```go
response, err := client.Pack(context.Background(), request)
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
    fmt.Printf("  Weight: %.2f lbs\n", palletizer.GramsToPounds(pallet.TotalWeight))
    fmt.Printf("  Height: %.2f inches\n", palletizer.MMToInches(pallet.TotalHeight))
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
health, err := client.Health(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", health.Status)
```

## Metrics

```go
metrics, err := client.Metrics(context.Background())
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

client := palletizer.NewWithHTTPClient("https://api.palletizer.app", httpClient)
```

## Context and Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.Pack(ctx, request)
```

## Error Handling

```go
response, err := client.Pack(context.Background(), request)
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

## Performance

The Palletizer service provides excellent performance:

- 30 cartons: < 10ms
- 279 cartons: ~267ms
- 1,000 cartons: ~1.6 seconds
- 10,000 cartons: ~86 seconds

Typical space utilization: 85-95%

## Testing

Run the SDK tests:

```bash
go test -v
```

## Support

- **Website**: https://palletizer.app
- **Documentation**: https://docs.palletizer.app
- **Issues**: https://github.com/palletizer-app/sdk/issues
- **Email**: info@palletizer.app

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

Made with ❤️ by [Palletizer.app](https://palletizer.app)
