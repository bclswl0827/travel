# travel

travel is a Go library for estimating seismic-wave travel times from AK135 tables. It provides a compact, table-driven API for P, PcP, PKPab, PKPbc, PKIKP, S, ScP, ScS, SKP, SKIKP, SKS, and SKIKS.

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
- Each phase result is a pointer to a `TravelTimeEntry` when the bundled table can estimate it; otherwise it is `nil`.
- The `Duration` field is a Go `time.Duration` (nanoseconds), and `Slowness` is reported in seconds per degree.
- Distance interpolation uses ray parameters; source-depth interpolation is linear between adjacent table depths.
- P and S represent the first arrivals of `p/P/Pdiff` and `s/S/Sdiff`, respectively. `PKPab` and `PKPbc` select the separate PKP branches.
- `SKP` and `SKIKP` are independent phases. The published SKP table mixes them; callers needing the earliest of the two should compare the available results rather than interpolate between phases.
- P/S distance grids use 0.1° spacing below 2° and 1° spacing elsewhere. See [data provenance](resource/data/README.md) for corrections, supplements, and reproduction instructions.
- Source depths are limited to 0–700 km. Distance coverage depends on the phase. Near a shadow or phase boundary, interpolation may return `nil` when adjacent cells are unavailable.
- The package does not trace rays through a velocity model or return multiple arrivals for a single phase; it uses only the bundled AK135 grids.

## License

The original software source code is licensed under Apache-2.0.

The AK135 numerical data under `resource/data` are derived from B. L. N. Kennett's _Seismological Tables: ak135_ and are not claimed to be licensed under Apache-2.0.

See `resource/data/README.md`.
