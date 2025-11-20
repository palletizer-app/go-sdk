# Palletizer Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/palletizer-app/sdk.svg)](https://pkg.go.dev/github.com/palletizer-app/sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/palletizer-app/sdk)](https://goreportcard.com/report/github.com/palletizer-app/sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Optimize your warehouse operations with AI-powered 3D pallet packing.**

Official Go SDK for [Palletizer.app](https://palletizer.app) - The smartest way to pack cartons onto pallets. Our advanced algorithm maximizes space utilization while ensuring load stability and safety.

## 🚀 Why Palletizer?

- **Save Money** - Reduce shipping costs with 85-95% space utilization
- **Pack Faster** - Get optimal packing solutions in milliseconds
- **Stay Safe** - Automatic weight distribution and stability checks
- **Think Less** - Just send your cartons, we handle the complexity
- **Scale Easily** - From 10 to 10,000 cartons per request

## ✨ Key Features

- ✅ **Multi-Orientation Support** - Automatically tests all possible carton rotations
- ✅ **Weight Distribution** - Ensures balanced center of gravity
- ✅ **Support Rules** - 80% base support requirement for stability
- ✅ **Mixed Carton Sizes** - Handle multiple SKUs in one request
- ✅ **Fragile Item Handling** - Special placement rules for delicate items
- ✅ **Constraint Validation** - Respects weight, height, and dimensional limits
- ✅ **Lightning Fast** - 1,000 cartons packed in ~1.6 seconds

## 📦 Installation

```bash
go get github.com/palletizer-app/sdk
```

## 🎯 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/palletizer-app/sdk"
)

func main() {
    // Create client (connects to api.palletizer.app automatically)
    client := palletizer.New()

    // Create packing request with imperial units
    request := &palletizer.PackingRequest{
        Cartons: []palletizer.Carton{
            {
                ID:            "BOX001",
                Length:        palletizer.InchesToMM(24),
                Width:         palletizer.InchesToMM(18),
                Height:        palletizer.InchesToMM(16),
                Weight:        palletizer.PoundsToGrams(40),
                Quantity:      30,
                AllowRotation: true,
            },
        },
        PalletConstraints: palletizer.StandardPallet(), // 40×72×48" pallet
        PackingOptions: palletizer.PackingOptions{
            SupportPercentage: 80.0,
        },
    }

    // Get optimized packing solution
    response, err := client.Pack(context.Background(), request)
    if err != nil {
        log.Fatal(err)
    }

    // Results
    fmt.Printf("✅ Packed %d cartons onto %d pallets\n",
        response.Summary.TotalCartonsPacked,
        response.Summary.TotalPallets)
    fmt.Printf("📊 Average utilization: %.2f%%\n",
        response.Summary.AverageUtilization)
    fmt.Printf("⚡ Computed in: %d ms\n",
        response.Summary.ComputationTimeMs)
}
```

## 💡 Real-World Examples

### Mixed Carton Sizes (E-commerce Fulfillment)

```go
client := palletizer.New()

request := &palletizer.PackingRequest{
    Cartons: []palletizer.Carton{
        {
            ID:       "LARGE",
            Length:   palletizer.InchesToMM(24),
            Width:    palletizer.InchesToMM(18),
            Height:   palletizer.InchesToMM(16),
            Weight:   palletizer.PoundsToGrams(40),
            Quantity: 20,
            AllowRotation: true,
        },
        {
            ID:       "MEDIUM",
            Length:   palletizer.InchesToMM(18),
            Width:    palletizer.InchesToMM(12),
            Height:   palletizer.InchesToMM(12),
            Weight:   palletizer.PoundsToGrams(20),
            Quantity: 30,
            AllowRotation: true,
        },
        {
            ID:       "SMALL",
            Length:   palletizer.InchesToMM(12),
            Width:    palletizer.InchesToMM(8),
            Height:   palletizer.InchesToMM(8),
            Weight:   palletizer.PoundsToGrams(10),
            Quantity: 50,
            Fragile:  true, // Handle with care
            AllowRotation: false,
        },
    },
    PalletConstraints: palletizer.StandardPallet(),
}

response, _ := client.Pack(context.Background(), request)

// Typical result: 90%+ utilization with 2-3 pallets
```

### Processing Results

```go
response, err := client.Pack(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

// Summary statistics
fmt.Printf("Total Pallets: %d\n", response.Summary.TotalPallets)
fmt.Printf("Total Cartons: %d\n", response.Summary.TotalCartonsPacked)
fmt.Printf("Utilization: %.2f%%\n", response.Summary.AverageUtilization)
fmt.Printf("Time: %d ms\n", response.Summary.ComputationTimeMs)

// Detailed pallet breakdown
for _, pallet := range response.Pallets {
    fmt.Printf("\n📦 Pallet %d:\n", pallet.PalletID)
    fmt.Printf("  Weight: %.1f lbs\n", palletizer.GramsToPounds(pallet.TotalWeight))
    fmt.Printf("  Height: %.1f in\n", palletizer.MMToInches(pallet.TotalHeight))
    fmt.Printf("  Space used: %.1f%%\n", pallet.UtilizationPercentage)
    fmt.Printf("  Cartons: %d boxes\n", len(pallet.Cartons))
    
    // Center of gravity for forklift operators
    fmt.Printf("  Center: (%.0f, %.0f, %.0f) mm\n",
        pallet.CenterOfGravity.X,
        pallet.CenterOfGravity.Y,
        pallet.CenterOfGravity.Z)
}
```

## 🔧 Configuration Options

### Standard Pallet Sizes

```go
// 40×72×48 inch pallet (1500 lbs) - Most common in North America
palletizer.StandardPallet()

// 40×48×48 inch pallet (1500 lbs) - Square pallet
palletizer.StandardPallet4048()

// Custom pallet
palletizer.PalletConstraints{
    MaxLength: palletizer.InchesToMM(48),
    MaxWidth:  palletizer.InchesToMM(40),
    MaxHeight: palletizer.InchesToMM(60),
    MaxWeight: palletizer.PoundsToGrams(2000),
}
```

### Unit Conversions

```go
// Convert TO metric (API requires mm and grams)
length := palletizer.InchesToMM(24)      // 609.6 mm
weight := palletizer.PoundsToGrams(40)   // 18143.68 g

// Convert FROM metric (for display)
inches := palletizer.MMToInches(609.6)   // 24.0 inches
pounds := palletizer.GramsToPounds(18143.68) // 40.0 lbs
```

### Packing Options

```go
options := palletizer.PackingOptions{
    SupportPercentage: 80.0,  // Require 80% base support (recommended)
}
```

### Custom HTTP Client

```go
// With custom timeout
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}
client := palletizer.NewWithHTTPClient(httpClient)

// With context timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
response, err := client.Pack(ctx, request)
```

### Custom Endpoint (for testing)

```go
// Use custom API endpoint
client := palletizer.NewWithEndpoint("http://localhost:8080")
```

## 🎯 Use Cases

### E-commerce Fulfillment
Pack multi-SKU orders efficiently, reducing shipping costs and transit damage.

### Warehouse Operations
Optimize pallet loading for storage density and forklift stability.

### Freight Shipping
Maximize trailer space utilization while meeting carrier weight limits.

### Manufacturing
Plan production runs with optimal packaging layouts.

### 3PL Services
Provide instant packing solutions for diverse client inventories.

## 📊 Performance Benchmarks

| Cartons | Time | Pallets | Utilization |
|---------|------|---------|-------------|
| 30 | < 10ms | 2 | 75% |
| 100 | ~50ms | 5 | 88% |
| 279 | ~267ms | 12 | 78% |
| 1,000 | ~1.6s | 40 | 88% |
| 10,000 | ~86s | 527 | 87% |

*Measured on production API*

## 🛡️ Error Handling

```go
response, err := client.Pack(context.Background(), request)
if err != nil {
    // Network or timeout error
    log.Printf("Request failed: %v", err)
    return
}

if response.Error != "" {
    // API validation error (e.g., carton too large for pallet)
    log.Printf("Packing error: %s", response.Error)
    return
}

// Success!
```

## 📐 API Input Format

```go
type PackingRequest struct {
    Cartons           []Carton
    PalletConstraints PalletConstraints
    PackingOptions    PackingOptions
}

type Carton struct {
    ID            string  // Your SKU or identifier
    Length        float64 // millimeters
    Width         float64 // millimeters
    Height        float64 // millimeters
    Weight        float64 // grams
    Quantity      int     // number of identical cartons
    Fragile       bool    // special handling
    AllowRotation bool    // can be rotated for better fit
}
```

## 🧪 Testing

```bash
# Run SDK tests
go test -v

# With coverage
go test -cover
```

## 📚 Documentation

- **API Documentation**: https://docs.palletizer.app
- **Algorithm Details**: https://palletizer.app/algorithm
- **Pricing**: https://palletizer.app/pricing

## 🆘 Support

- **Email**: info@palletizer.app
- **Issues**: https://github.com/palletizer-app/sdk/issues
- **Website**: https://palletizer.app

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Built by warehouse professionals, for warehouse professionals.**

[Palletizer.app](https://palletizer.app) | Smarter Packing, Better Shipping
