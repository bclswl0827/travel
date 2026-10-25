package resource

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

type TravelTimeEntry struct {
	Depth    float64
	Delta    float64
	Slowness float64
	Duration time.Duration
}

// map[delta][depth]TravelTimeEntry
type TravelTimeMap map[float64]map[float64]TravelTimeEntry

func (m TravelTimeMap) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	if len(m) == 0 {
		return math.NaN(), math.NaN(), math.NaN(), math.NaN()
	}

	minDepth = math.Inf(1)
	maxDepth = math.Inf(-1)
	minDelta = math.Inf(1)
	maxDelta = math.Inf(-1)

	for delta, deltaMap := range m {
		if len(deltaMap) == 0 {
			continue
		}
		if delta < minDelta {
			minDelta = delta
		}
		if delta > maxDelta {
			maxDelta = delta
		}

		for depth := range deltaMap {
			if depth < minDepth {
				minDepth = depth
			}
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	if math.IsInf(minDelta, 1) {
		return math.NaN(), math.NaN(), math.NaN(), math.NaN()
	}
	return minDelta, maxDelta, minDepth, maxDepth
}

func (m TravelTimeMap) Query(delta, depth float64, interpolate bool) (TravelTimeEntry, error) {
	if len(m) == 0 {
		return TravelTimeEntry{}, errors.New("empty travel time map")
	}
	if math.IsNaN(delta) || math.IsInf(delta, 0) || math.IsNaN(depth) || math.IsInf(depth, 0) {
		return TravelTimeEntry{}, errors.New("distance and depth must be finite")
	}

	// Return an exact table value without applying any interpolation.
	if row, ok := m[delta]; ok {
		if entry, found := row[depth]; found {
			return entry, nil
		}
	}
	if !interpolate {
		return TravelTimeEntry{}, errors.New("interpolation disabled and exact entry is not found")
	}

	deltas := make([]float64, 0, len(m))
	depthSet := make(map[float64]struct{})
	for d := range m {
		deltas = append(deltas, d)
		for z := range m[d] {
			depthSet[z] = struct{}{}
		}
	}
	sort.Float64s(deltas)
	depths := make([]float64, 0, len(depthSet))
	for z := range depthSet {
		depths = append(depths, z)
	}
	sort.Float64s(depths)

	d0, d1, ok := bracket(deltas, delta)
	if !ok {
		return TravelTimeEntry{}, errors.New("delta out of range")
	}
	z0, z1, ok := bracket(depths, depth)
	if !ok {
		return TravelTimeEntry{}, errors.New("depth out of range")
	}

	// Interpolate only between adjacent cells in the original table grid. A
	// missing corner marks a distance/depth where TauP has no arrival; crossing
	// it would incorrectly fill a shadow or phase-existence gap.
	entry0, ok := m.interpolateDistance(d0, d1, z0, delta)
	if !ok {
		return TravelTimeEntry{}, errors.New("phase does not exist in interpolation interval")
	}
	if z0 == z1 {
		entry0.Depth = depth
		entry0.Delta = delta
		return entry0, nil
	}
	entry1, ok := m.interpolateDistance(d0, d1, z1, delta)
	if !ok {
		return TravelTimeEntry{}, errors.New("phase does not exist in interpolation interval")
	}

	t := (depth - z0) / (z1 - z0)
	return TravelTimeEntry{
		Depth:    depth,
		Delta:    delta,
		Duration: lerpDuration(entry0.Duration, entry1.Duration, t),
		Slowness: lerp(entry0.Slowness, entry1.Slowness, t),
	}, nil
}

func bracket(values []float64, target float64) (lower, upper float64, ok bool) {
	if len(values) == 0 || target < values[0] || target > values[len(values)-1] {
		return 0, 0, false
	}
	i := sort.SearchFloat64s(values, target)
	if i < len(values) && values[i] == target {
		return target, target, true
	}
	if i == 0 || i == len(values) {
		return 0, 0, false
	}
	return values[i-1], values[i], true
}

func (m TravelTimeMap) interpolateDistance(d0, d1, depth, delta float64) (TravelTimeEntry, bool) {
	e0, ok := m[d0][depth]
	if !ok {
		return TravelTimeEntry{}, false
	}
	if d0 == d1 {
		return e0, true
	}
	e1, ok := m[d1][depth]
	if !ok {
		return TravelTimeEntry{}, false
	}

	t := (delta - d0) / (d1 - d0)
	slowness := lerp(e0.Slowness, e1.Slowness, t)

	// TauP linearly interpolates ray parameter, then uses that derivative and
	// the closest bracketing arrival to interpolate travel time. The table ray
	// parameter is seconds/degree, so no radian conversion is needed here.
	var seconds float64
	if math.Abs(delta-d0) < math.Abs(delta-d1) {
		seconds = e0.Duration.Seconds() + slowness*(delta-d0)
	} else {
		seconds = e1.Duration.Seconds() + slowness*(delta-d1)
	}
	return TravelTimeEntry{
		Depth:    depth,
		Delta:    delta,
		Duration: durationFromSeconds(seconds),
		Slowness: slowness,
	}, true
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func lerpDuration(a, b time.Duration, t float64) time.Duration {
	return time.Duration(math.Round(lerp(float64(a), float64(b), t)))
}

func durationFromSeconds(seconds float64) time.Duration {
	return time.Duration(math.Round(seconds * float64(time.Second)))
}

// MergeTravelTimeMaps combines tables that contain different source-depth
// ranges for the same phase, such as the shallow and deep AK135 P/S tables.
func MergeTravelTimeMaps(maps ...TravelTimeMap) TravelTimeMap {
	merged := make(TravelTimeMap)
	for _, m := range maps {
		for delta, row := range m {
			if merged[delta] == nil {
				merged[delta] = make(map[float64]TravelTimeEntry)
			}
			for depth, entry := range row {
				merged[delta][depth] = entry
			}
		}
	}
	return merged
}

func (ttm TravelTimeMap) PrettyPrint(w io.Writer) {
	var deltas []float64
	for delta := range ttm {
		deltas = append(deltas, delta)
	}
	sort.Float64s(deltas)

	depthSet := make(map[float64]struct{})
	for _, m := range ttm {
		for depth := range m {
			depthSet[depth] = struct{}{}
		}
	}
	var depths []float64
	for depth := range depthSet {
		depths = append(depths, depth)
	}
	sort.Float64s(depths)

	fmt.Fprintln(w, "Travel Times (s):")
	fmt.Fprintf(w, "%8s", "Δ\\Depth")
	for _, depth := range depths {
		fmt.Fprintf(w, "%12.1f", depth)
	}
	fmt.Fprintln(w)

	for _, delta := range deltas {
		fmt.Fprintf(w, "%8.1f", delta)
		for _, depth := range depths {
			if entry, ok := ttm[delta][depth]; ok {
				fmt.Fprintf(w, "%12.1f", entry.Duration.Seconds())
			} else {
				fmt.Fprintf(w, "%12s", "-")
			}
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "\nSlowness:")
	fmt.Fprintf(w, "%8s", "Δ\\Depth")
	for _, depth := range depths {
		fmt.Fprintf(w, "%12.1f", depth)
	}
	fmt.Fprintln(w)

	for _, delta := range deltas {
		fmt.Fprintf(w, "%8.1f", delta)
		for _, depth := range depths {
			if entry, ok := ttm[delta][depth]; ok {
				fmt.Fprintf(w, "%12.2f", entry.Slowness)
			} else {
				fmt.Fprintf(w, "%12s", "-")
			}
		}
		fmt.Fprintln(w)
	}
}
