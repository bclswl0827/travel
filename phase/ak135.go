package phase

import "github.com/bclswl0827/travel/resource"

func removeEntries(m resource.TravelTimeMap, delta float64, depths ...float64) {
	for _, depth := range depths {
		delete(m[delta], depth)
	}
}
