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
		{"SKIKP", model.SKIKP},
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

func TestPAndSHaveCompleteTableGrids(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	depths := []float64{0, 15, 35, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700}
	deltas := make([]float64, 0, 144)
	for delta := 0.0; delta <= 125; delta++ {
		deltas = append(deltas, delta)
	}
	for tenth := 1; tenth < 20; tenth++ {
		if tenth != 10 {
			deltas = append(deltas, float64(tenth)/10)
		}
	}
	for _, phase := range []travel.IPhase{model.P, model.S} {
		for _, delta := range deltas {
			for _, depth := range depths {
				if _, err := phase.Estimate(delta, depth, false); err != nil {
					t.Errorf("%s table is missing %.1f degrees, %.0f km: %v", phase.Name(), delta, depth, err)
				}
			}
		}
	}
}

func TestDeepSInterpolationAgainstAK135Reference(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	// Independent TauP 3.2.1 AK135 calculations. The tolerance accounts for
	// interpolation between table cells rounded to 0.01 seconds and s/degree.
	tests := []struct {
		delta   float64
		depth   float64
		seconds float64
	}{
		{25.5, 275, 549.9626},
		{29.4319, 275, 611.70996},
		{49.5, 275, 908.8707},
		{25.5, 480, 521.5524},
		{29.4319, 480, 582.95654}, // 3276.3 km epicentral distance in Observer.
		{49.5, 480, 876.69324},
	}
	for _, test := range tests {
		result := model.Estimate(test.delta, test.depth, true)
		if result.S == nil {
			t.Errorf("S arrival is missing at %.4f degrees, %.0f km", test.delta, test.depth)
			continue
		}
		if got := result.S.Duration.Seconds(); math.Abs(got-test.seconds) > 0.25 {
			t.Errorf("S at %.4f degrees, %.0f km = %.5f s, want %.5f s within 0.25 s", test.delta, test.depth, got, test.seconds)
		}
		if result.P == nil || result.S.Duration <= result.P.Duration {
			t.Errorf("expected a P arrival followed by an S arrival at %.4f degrees, %.0f km", test.delta, test.depth)
		}
	}
}

func TestTauPUnavailableArrivals(t *testing.T) {
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
		{model.PKPbc, 155, 700},
		{model.SKP, 110, 300},
		{model.SKIKP, 110, 300},
		{model.SKS, 144, 700},
		{model.SKIKS, 105, 200},
		{model.PcP, 97, 700},
		{model.PcP, 98, 500},
		{model.ScS, 98, 700},
	}
	for _, test := range tests {
		if _, err := test.phase.Estimate(test.delta, test.depth, false); err == nil {
			t.Errorf("%s returned an unavailable arrival at %.0f degrees, %.0f km", test.phase.Name(), test.delta, test.depth)
		}
	}
}

func TestCorrectedAK135ReferenceArrivals(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	tests := []struct {
		phase   travel.IPhase
		delta   float64
		depth   float64
		seconds float64
	}{
		{model.PcP, 97, 500, 757.5992},
		{model.PcP, 99, 0, 822.5759},
		{model.ScP, 63, 300, 843.03406},
		{model.ScS, 100, 0, 1522.8375},
		{model.SKS, 144, 500, 1507.6077},
		{model.SKIKS, 105, 300, 1469.0819},
		{model.SKIKS, 105, 500, 1428.8163},
		{model.SKIKS, 105, 700, 1392.9469},
		{model.SKIKP, 110, 500, 1221.7473},
		{model.SKIKP, 110, 700, 1185.8779},
		{model.SKP, 130, 200, 1309.4484},
		{model.SKIKP, 130, 200, 1321.8884},
		{model.SKP, 149, 0, 1408.6007},
		{model.SKIKP, 149, 0, 1400.8081},
	}
	for _, test := range tests {
		entry, err := test.phase.Estimate(test.delta, test.depth, false)
		if err != nil {
			t.Errorf("%s is missing at %.0f degrees, %.0f km: %v", test.phase.Name(), test.delta, test.depth, err)
			continue
		}
		if got := entry.Duration.Seconds(); math.Abs(got-test.seconds) > 0.1 {
			t.Errorf("%s at %.0f degrees, %.0f km = %.5f s, want %.5f s within 0.1 s", test.phase.Name(), test.delta, test.depth, got, test.seconds)
		}
	}
}

func TestNearSourceInterpolationAgainstAK135Reference(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	tests := []struct {
		phase   travel.IPhase
		delta   float64
		depth   float64
		seconds float64
	}{
		{model.P, 0.5, 25, 10.19191},
		{model.S, 0.5, 25, 17.123648},
		{model.P, 1.5, 7.5, 27.24975},
		{model.S, 1.5, 7.5, 47.022785},
	}
	for _, test := range tests {
		entry, err := test.phase.Estimate(test.delta, test.depth, true)
		if err != nil {
			t.Errorf("%s failed at %.1f degrees, %.1f km: %v", test.phase.Name(), test.delta, test.depth, err)
			continue
		}
		if got := entry.Duration.Seconds(); math.Abs(got-test.seconds) > 0.25 {
			t.Errorf("%s at %.1f degrees, %.1f km = %.5f s, want %.5f s within 0.25 s", test.phase.Name(), test.delta, test.depth, got, test.seconds)
		}
	}
}

func TestSKPAndSKIKPAreIndependent(t *testing.T) {
	model, err := travel.NewAK135()
	if err != nil {
		t.Fatalf("NewAK135 failed: %v", err)
	}
	result := model.Estimate(129.5, 225, true)
	if result.SKP != nil {
		t.Fatal("SKP was interpolated across the SKIKP/SKP branch transition")
	}
	if result.SKIKP == nil {
		t.Fatal("SKIKP arrival is missing")
	}
	if got := result.SKIKP.Duration.Seconds(); math.Abs(got-1315.4497) > 0.25 {
		t.Fatalf("SKIKP = %.5f s, want 1315.4497 s within 0.25 s", got)
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

func TestSphericalDistanceIsStableNearCoincidentAndAntipodalPoints(t *testing.T) {
	tests := []struct {
		lat1, lon1, lat2, lon2 float64
		delta                  float64
	}{
		{29.81, 106.4, 29.81, 106.4, 0},
		{0, 0, 0, 1e-9, 1e-9},
		{0, 179.999999, 0, -179.999999, 0.000002},
		{0, 0, 0, 179.999999, 179.999999},
		{30, 0, -30, 180, 180},
	}
	for _, test := range tests {
		if got := travel.GetDeltaByCoordinates(test.lat1, test.lon1, test.lat2, test.lon2); math.Abs(got-test.delta) > 1e-12 {
			t.Errorf("delta for (%g, %g), (%g, %g) = %.15g, want %.15g", test.lat1, test.lon1, test.lat2, test.lon2, got, test.delta)
		}
	}
}
