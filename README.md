# travel

travel is a Go library for estimating seismic-wave travel times from the published AK135 tables. It provides a compact, table-driven API for common phases such as P, PcP, PKPab, PKPbc, PKIKP, S, ScP, ScS, SKP, SKS, and SKIKS.

## Features

- Create an AK135 model with `NewAK135`
- Estimate travel times for all supported phases in one call
- Optionally interpolate travel times in distance and source depth
- Convert geographic coordinates to angular distance with helper functions

## Installation

```bash
go get github.com/bclswl0827/travel
```

## Quick start

```go
package main

import (
    "fmt"

    "github.com/bclswl0827/travel"
)

func main() {
    model, err := travel.NewAK135()
    if err != nil {
        panic(err)
    }

    estimation := model.Estimate(
        travel.GetDeltaByCoordinates(28.25, 104.87, 29.81, 106.40),
        10,
        true,
    )

    fmt.Printf("P travel time: %v\n", estimation.P.Duration)
}
```

## API notes

- `Estimate(delta, depth, interpolate)` returns a struct containing one result per supported phase.
- Each phase result is a pointer to a `TravelTimeEntry` when the phase exists for the requested distance and depth; otherwise it is `nil`.
- The `Duration` field is a Go `time.Duration` (nanoseconds), and `Slowness` is reported in seconds per degree.
- Interpolation follows TauP's ray-parameter rule for distance and uses linear interpolation in depth between adjacent published table depths.
- The package does not trace rays through a velocity model or return multiple arrivals for a single phase; it is intentionally limited to the published AK135 tables.

## License

The original software source code is licensed under Apache-2.0.

The AK135 numerical data under `resource/data` are derived from B. L. N. Kennett's _Seismological Tables: ak135_ and are not claimed to be licensed under Apache-2.0.

See `resource/data/README.md`.
