package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type S_Impl struct {
	shallowMap resource.TravelTimeMap
	deepMap    resource.TravelTimeMap
}

func (p *S_Impl) Name() string {
	return "P"
}

func (p *S_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, _ = p.shallowMap.Boundary()
	_, _, _, maxDepth = p.deepMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *S_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	_, _, _, shallowMaxDepth := p.shallowMap.Boundary()

	mapObj := p.shallowMap
	if depth > shallowMaxDepth {
		mapObj = p.deepMap
	}

	return mapObj.Query(delta, depth, interplate)
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
		shallowMap: shallowMap,
		deepMap:    deepMap,
	}, nil
}
