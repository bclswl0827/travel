package travel

import "math"

func GetDeltaByCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	dLonRad := (lon2 - lon1) * math.Pi / 180.0
	cosDelta := math.Sin(lat1Rad)*math.Sin(lat2Rad) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Cos(dLonRad)
	// Floating-point rounding can put the cosine a few ulps outside [-1, 1].
	cosDelta = math.Max(-1, math.Min(1, cosDelta))
	return math.Acos(cosDelta) * 180.0 / math.Pi
}

func GetDistanceByCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	delta := GetDeltaByCoordinates(lat1, lon1, lat2, lon2)
	return GetDistanceByDelta(delta)
}

func GetDistanceByDelta(delta float64) float64 {
	return RADIUS_EARTH_KM * (delta * math.Pi / 180.0)
}
