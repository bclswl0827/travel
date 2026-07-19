package travel_test

import (
	"math"
	"testing"

	"github.com/bclswl0827/travel"
)

func TestTauPCanonicalPhaseNames(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}

	checks := []struct {
		name string
		got  travel.IPhase
	}{
		{"PKIKP", model.PKIKP},
		{"SKS", model.SKS},
		{"SKIKS", model.SKIKS},
	}
	for _, check := range checks {
		if check.got.Name() != check.name {
			t.Errorf("phase name = %q, want %q", check.got.Name(), check.name)
		}
	}
	if model.S.Name() != "S" {
		t.Fatalf("S phase name = %q, want S", model.S.Name())
	}
}

func TestPAndSInterpolateAcrossPublishedTableSplit(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	for _, phase := range []travel.IPhase{model.P, model.S} {
		entry, err := phase.Estimate(10, 275, true)
		if err != nil {
			t.Errorf("%s failed at 275 km: %v", phase.Name(), err)
			continue
		}
		if entry.Depth != 275 || entry.Delta != 10 {
			t.Errorf("%s returned wrong coordinates: %+v", phase.Name(), entry)
		}
	}
}

func TestTauPAK135PublishedTableExceptionsAreExcluded(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	tests := []struct {
		phase travel.IPhase
		delta float64
		depth float64
	}{
		{model.PKPbc, 155, 500},
		{model.SKP, 110, 300},
		{model.SKS, 144, 500},
		{model.SKIKS, 105, 200},
	}
	for _, test := range tests {
		if _, err := test.phase.Estimate(test.delta, test.depth, false); err == nil {
			t.Errorf("%s returned TauP-excluded entry at %.0f degrees, %.0f km", test.phase.Name(), test.delta, test.depth)
		}
	}
}

func TestAK135RadiusAndSphericalDistance(t *testing.T) {
	if travel.RADIUS_EARTH_KM != 6371.0 {
		t.Fatalf("AK135 radius = %v, want 6371.0", travel.RADIUS_EARTH_KM)
	}
	if got := travel.GetDeltaByCoordinates(0, 0, 0, 90); math.Abs(got-90) > 1e-12 {
		t.Fatalf("quarter-circle delta = %.15g, want 90", got)
	}
	if got := travel.GetDeltaByCoordinates(12.3, 45.6, 12.3, 45.6); got != 0 {
		t.Fatalf("identical-coordinate delta = %.15g, want 0", got)
	}
	if got, want := travel.GetDistanceByDelta(180), math.Pi*6371.0; math.Abs(got-want) > 1e-9 {
		t.Fatalf("antipodal distance = %.15g, want %.15g", got, want)
	}
}
