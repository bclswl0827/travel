package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type P_Impl struct {
	travelTimeMap resource.TravelTimeMap
}

func (p *P_Impl) Name() string {
	return "P"
}

func (p *P_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	return p.travelTimeMap.Boundary()
}

func (p *P_Impl) Estimate(delta, depth float64, interpolate bool) (result resource.TravelTimeEntry, err error) {
	return p.travelTimeMap.Query(delta, depth, interpolate)
}

func NewP(pShallowData, pDeepData string) (*P_Impl, error) {
	if len(pShallowData) == 0 {
		return nil, errors.New("empty P shallow data")
	}
	if len(pDeepData) == 0 {
		return nil, errors.New("empty P deep data")
	}

	shallowMap, err := resource.NewTravelTimeMap(pShallowData)
	if err != nil {
		return nil, err
	}
	deepMap, err := resource.NewTravelTimeMap(pDeepData)
	if err != nil {
		return nil, err
	}

	return &P_Impl{
		travelTimeMap: resource.MergeTravelTimeMaps(shallowMap, deepMap),
	}, nil
}
