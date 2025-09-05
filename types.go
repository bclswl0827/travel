package travel

import (
	"github.com/bclswl0827/travel/resource"
)

const RADIUS_EARTH_KM = 6371.0008

type IPhase interface {
	Name() string
	Boundary() (minDelta, maxDelta, minDepth, maxDepth float64)
	Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error)
}
