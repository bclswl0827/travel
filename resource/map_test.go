package resource

import (
	"math"
	"testing"
	"time"
)

func testEntry(delta, depth, seconds, slowness float64) TravelTimeEntry {
	return TravelTimeEntry{
		Delta:    delta,
		Depth:    depth,
		Duration: time.Duration(seconds * float64(time.Second)),
		Slowness: slowness,
	}
}

func TestBoundaryIgnoresRowsWithoutArrivals(t *testing.T) {
	m := TravelTimeMap{
		0: {},
		1: {10: testEntry(1, 10, 5, 1)},
		2: {20: testEntry(2, 20, 10, 1)},
		3: {},
	}
	minDelta, maxDelta, minDepth, maxDepth := m.Boundary()
	if minDelta != 1 || maxDelta != 2 || minDepth != 10 || maxDepth != 20 {
		t.Fatalf("boundary = (%g, %g, %g, %g), want (1, 2, 10, 20)", minDelta, maxDelta, minDepth, maxDepth)
	}
	for _, empty := range []TravelTimeMap{nil, {}, {0: {}}} {
		minDelta, maxDelta, minDepth, maxDepth := empty.Boundary()
		if !math.IsNaN(minDelta) || !math.IsNaN(maxDelta) || !math.IsNaN(minDepth) || !math.IsNaN(maxDepth) {
			t.Fatalf("empty table boundary = (%g, %g, %g, %g), want NaN", minDelta, maxDelta, minDepth, maxDepth)
		}
	}
}

func TestQueryRequiresExactEntryWhenInterpolationDisabled(t *testing.T) {
	m := TravelTimeMap{
		0: {
			0:   testEntry(0, 0, 10, 1),
			100: testEntry(0, 100, 20, 1),
		},
		2: {
			0:   testEntry(2, 0, 12, 1),
			100: testEntry(2, 100, 22, 1),
		},
	}

	if _, err := m.Query(0, 0, false); err != nil {
		t.Fatalf("exact query failed: %v", err)
	}
	if _, err := m.Query(0, 50, false); err == nil {
		t.Fatal("non-exact depth succeeded with interpolation disabled")
	}
	if _, err := m.Query(1, 0, false); err == nil {
		t.Fatal("non-exact distance succeeded with interpolation disabled")
	}
}

func TestQueryInterpolatesDepthAtMaximumDistance(t *testing.T) {
	m := TravelTimeMap{
		0: {
			0:   testEntry(0, 0, 10, 1),
			100: testEntry(0, 100, 20, 3),
		},
		2: {
			0:   testEntry(2, 0, 20, 2),
			100: testEntry(2, 100, 30, 4),
		},
	}

	got, err := m.Query(2, 50, true)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if got.Duration != 25*time.Second || got.Slowness != 3 {
		t.Fatalf("got duration=%v slowness=%v, want 25s and 3", got.Duration, got.Slowness)
	}
}

func TestQueryUsesTauPDistanceInterpolation(t *testing.T) {
	m := TravelTimeMap{
		0: {0: testEntry(0, 0, 100, 10)},
		2: {0: testEntry(2, 0, 119, 8)},
	}

	got, err := m.Query(1, 0, true)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	// TauP interpolates p to 9 s/degree and integrates from the closest
	// endpoint: 119 + 9*(1-2) = 110 seconds. A duration lerp would be 109.5.
	if got.Duration != 110*time.Second || got.Slowness != 9 {
		t.Fatalf("got duration=%v slowness=%v, want 110s and 9", got.Duration, got.Slowness)
	}
}

func TestQueryDoesNotCrossMissingPhaseCell(t *testing.T) {
	m := TravelTimeMap{
		0: {
			0:   testEntry(0, 0, 10, 1),
			100: testEntry(0, 100, 20, 1),
		},
		2: {
			0: testEntry(2, 0, 12, 1),
		},
	}

	if _, err := m.Query(1, 50, true); err == nil {
		t.Fatal("query interpolated across a missing phase cell")
	}
}

func TestMergeTravelTimeMapsJoinsDepthRanges(t *testing.T) {
	shallow := TravelTimeMap{1: {250: testEntry(1, 250, 10, 1)}}
	deep := TravelTimeMap{1: {300: testEntry(1, 300, 20, 1)}}

	got, err := MergeTravelTimeMaps(shallow, deep).Query(1, 275, true)
	if err != nil {
		t.Fatalf("query in joined depth range failed: %v", err)
	}
	if got.Duration != 15*time.Second {
		t.Fatalf("got %v, want 15s", got.Duration)
	}
}
