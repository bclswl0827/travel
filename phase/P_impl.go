package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type P_Impl struct {
	shallowMap resource.TravelTimeMap
	deepMap    resource.TravelTimeMap
}

func (p *P_Impl) Name() string {
	return "P"
}

func (p *P_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, _ = p.shallowMap.Boundary()
	_, _, _, maxDepth = p.deepMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *P_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	_, _, _, shallowMaxDepth := p.shallowMap.Boundary()

	mapObj := p.shallowMap
	if depth > shallowMaxDepth {
		mapObj = p.deepMap
	}

	return mapObj.Query(delta, depth, interplate)
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
		shallowMap: shallowMap,
		deepMap:    deepMap,
	}, nil
}
