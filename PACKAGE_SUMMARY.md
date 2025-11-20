# Palletizer Client Package

## What We Built

We've created a production-ready Go client package for the Palletizer API that can be imported and used by any Go application.

## Package Location

```
pkg/client/
├── client.go       # Main client implementation
├── client_test.go  # Comprehensive tests
└── README.md       # Full documentation
```

## Key Features

### 1. **Easy to Use**
```go
import "raidframe.com/palletizer/pkg/client"

c := client.New("https://palletizer.app")
response, err := c.Pack(context.Background(), request)
```

### 2. **No API Client Needed**
Instead of writing your own HTTP client code, just import the package:
```bash
go get raidframe.com/palletizer/pkg/client
```

### 3. **Built-in Unit Conversions**
```go
Length: client.InchesToMM(24)      // 24 inches → 609.6 mm
Weight: client.PoundsToGrams(40)   // 40 lbs → 18143.68 g
```

### 4. **Standard Pallets**
```go
client.StandardPallet()      // 40x72x48 inch pallet (1500 lbs)
client.StandardPallet4048()  // 40x48x48 inch pallet (1500 lbs)
```

### 5. **Full Type Safety**
All request and response types are strongly typed:
- `PackingRequest`
- `PackingResponse`
- `Carton`
- `PalletConstraints`
- `PackingOptions`
- `Pallet`
- `PlacedCarton`
- `PackingSummary`

### 6. **Context Support**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
response, err := c.Pack(ctx, request)
```

### 7. **Health & Metrics**
```go
health, err := c.Health(ctx)
metrics, err := c.Metrics(ctx)
```

### 8. **Custom HTTP Client**
```go
httpClient := &http.Client{Timeout: 60 * time.Second}
c := client.NewWithHTTPClient("https://palletizer.app", httpClient)
```

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "raidframe.com/palletizer/pkg/client"
)

func main() {
    // Create client
    c := client.New("https://palletizer.app")
    
    // Create request
    request := &client.PackingRequest{
        Cartons: []client.Carton{
            {
                ID:            "BOX001",
                Length:        client.InchesToMM(24),
                Width:         client.InchesToMM(18),
                Height:        client.InchesToMM(16),
                Weight:        client.PoundsToGrams(40),
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
    
    // Process results
    fmt.Printf("Packed %d cartons onto %d pallets\n",
        response.Summary.TotalCartonsPacked,
        response.Summary.TotalPallets)
    
    for _, pallet := range response.Pallets {
        fmt.Printf("Pallet %d: %.2f%% utilization, %.2f lbs\n",
            pallet.PalletID,
            pallet.UtilizationPercentage,
            client.GramsToPounds(pallet.TotalWeight))
    }
}
```

## Testing

The client package includes comprehensive unit tests:

```bash
go test -v ./pkg/client/...
```

All tests pass:
- ✅ TestNew - Client creation
- ✅ TestPack - Packing requests
- ✅ TestHealth - Health checks
- ✅ TestMetrics - Metrics retrieval
- ✅ TestStandardPallet - Pallet presets
- ✅ TestConversions - Unit conversions

## Integration with Existing Code

### Before (Manual API Calls)
```go
// Define all types yourself
type PackingRequest struct { ... }
type PackingResponse struct { ... }

// Write HTTP client code
jsonData, _ := json.Marshal(request)
resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
var response PackingResponse
json.NewDecoder(resp.Body).Decode(&response)
```

### After (Using Client Package)
```go
import "raidframe.com/palletizer/pkg/client"

c := client.New("https://palletizer.app")
response, err := c.Pack(context.Background(), request)
```

## Documentation

- **Client Package**: [pkg/client/README.md](../pkg/client/README.md)
- **API Examples**: [examples/api_client.go](../examples/api_client.go)
- **Test Examples**: [examples/README.md](../examples/README.md)

## Deployment

When deploying to https://palletizer.app/, users can simply:

```bash
go get raidframe.com/palletizer/pkg/client
```

And start using the client immediately without writing any HTTP code.

## Advantages

1. **No boilerplate** - No need to write HTTP client code
2. **Type safety** - All types are pre-defined and validated
3. **Error handling** - Built-in error handling and context support
4. **Unit conversions** - Helper functions for imperial ↔ metric
5. **Tested** - Comprehensive unit tests included
6. **Documented** - Full README with examples
7. **Idiomatic Go** - Follows Go best practices (context, error handling)
8. **Production ready** - Ready to use with https://palletizer.app/
