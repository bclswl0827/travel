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

	return minDelta, maxDelta, minDepth, maxDepth
}

func (m TravelTimeMap) Query(delta, depth float64, interpolate bool) (TravelTimeEntry, error) {
	if len(m) == 0 {
		return TravelTimeEntry{}, errors.New("empty travel time map")
	}

	// return directly if exact hit
	if row, ok := m[delta]; ok {
		if entry, ok2 := row[depth]; ok2 {
			return entry, nil
		}
	} else if !interpolate {
		return TravelTimeEntry{}, errors.New("interpolation disabled, depth is not found")
	}

	var deltas []float64
	for d := range m {
		deltas = append(deltas, d)
	}
	sort.Float64s(deltas)

	if delta < deltas[0] || delta > deltas[len(deltas)-1] {
		return TravelTimeEntry{}, errors.New("delta out of range")
	}

	// find the nearest two deltas
	var d0, d1 float64
	for i := 0; i < len(deltas)-1; i++ {
		if delta >= deltas[i] && delta <= deltas[i+1] {
			d0, d1 = deltas[i], deltas[i+1]
			break
		}
	}

	// for d0 and d1, interpolate along the depth dimension
	entry0, ok0 := m.interpolateDepth(m[d0], depth, d0)
	entry1, ok1 := m.interpolateDepth(m[d1], depth, d1)
	if !ok0 || !ok1 {
		return TravelTimeEntry{}, errors.New("depth out of range")
	}

	t := (delta - d0) / (d1 - d0)
	return TravelTimeEntry{
		Depth:    depth,
		Delta:    delta,
		Duration: m.lerpDuration(entry0.Duration, entry1.Duration, t),
		Slowness: m.lerp(entry0.Slowness, entry1.Slowness, t),
	}, nil
}

func (m TravelTimeMap) interpolateDepth(row map[float64]TravelTimeEntry, depth, delta float64) (TravelTimeEntry, bool) {
	if row == nil {
		return TravelTimeEntry{}, false
	}

	var depths []float64
	for d := range row {
		depths = append(depths, d)
	}
	sort.Float64s(depths)

	if depth < depths[0] || depth > depths[len(depths)-1] {
		return TravelTimeEntry{}, false
	}

	// exact hit
	if entry, ok := row[depth]; ok {
		return entry, true
	}

	// find the nearest two depths
	var d0, d1 float64
	for i := 0; i < len(depths)-1; i++ {
		if depth >= depths[i] && depth <= depths[i+1] {
			d0, d1 = depths[i], depths[i+1]
			break
		}
	}

	e0 := row[d0]
	e1 := row[d1]
	t := (depth - d0) / (d1 - d0)

	return TravelTimeEntry{
		Depth:    depth,
		Delta:    delta,
		Duration: m.lerpDuration(e0.Duration, e1.Duration, t),
		Slowness: m.lerp(e0.Slowness, e1.Slowness, t),
	}, true
}

func (m TravelTimeMap) lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func (m TravelTimeMap) lerpDuration(a, b time.Duration, t float64) time.Duration {
	af := float64(a)
	bf := float64(b)
	return time.Duration(af + (bf-af)*t)
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
