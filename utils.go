package travel

import "math"

func GetDeltaByCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	dLonRad := (lon2 - lon1) * math.Pi / 180.0
	// atan2 avoids the precision loss of acos for coincident or nearby points,
	// and remains stable close to antipodes.
	sinLat1, cosLat1 := math.Sincos(lat1Rad)
	sinLat2, cosLat2 := math.Sincos(lat2Rad)
	sinLon, cosLon := math.Sincos(dLonRad)
	sinDelta := math.Hypot(cosLat2*sinLon, cosLat1*sinLat2-sinLat1*cosLat2*cosLon)
	cosDelta := sinLat1*sinLat2 + cosLat1*cosLat2*cosLon
	return math.Atan2(sinDelta, cosDelta) * 180.0 / math.Pi
}

func GetDistanceByCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	delta := GetDeltaByCoordinates(lat1, lon1, lat2, lon2)
	return GetDistanceByDelta(delta)
}

func GetDistanceByDelta(delta float64) float64 {
	return RADIUS_EARTH_KM * (delta * math.Pi / 180.0)
}
