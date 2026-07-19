package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type S_Impl struct {
	travelTimeMap resource.TravelTimeMap
}

func (p *S_Impl) Name() string {
	return "S"
}

func (p *S_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	return p.travelTimeMap.Boundary()
}

func (p *S_Impl) Estimate(delta, depth float64, interpolate bool) (result resource.TravelTimeEntry, err error) {
	return p.travelTimeMap.Query(delta, depth, interpolate)
}

func NewS(sShallowData, sDeepData string) (*S_Impl, error) {
	if len(sShallowData) == 0 {
		return nil, errors.New("empty S shallow data")
	}
	if len(sDeepData) == 0 {
		return nil, errors.New("empty S deep data")
	}

	shallowMap, err := resource.NewTravelTimeMap(sShallowData)
	if err != nil {
		return nil, err
	}

	deepMap, err := resource.NewTravelTimeMap(sDeepData)
	if err != nil {
		return nil, err
	}

	return &S_Impl{
		travelTimeMap: resource.MergeTravelTimeMaps(shallowMap, deepMap),
	}, nil
}
