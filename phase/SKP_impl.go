package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type SKP_Impl struct {
	skpMap resource.TravelTimeMap
}

func (p *SKP_Impl) Name() string {
	return "SKP"
}

func (p *SKP_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.skpMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *SKP_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.skpMap.Query(delta, depth, interplate)
}

func NewSKP(skpData string) (*SKP_Impl, error) {
	if len(skpData) == 0 {
		return nil, errors.New("empty SKP data")
	}

	mapObj, err := resource.NewTravelTimeMap(skpData)
	if err != nil {
		return nil, err
	}

	return &SKP_Impl{skpMap: mapObj}, nil
}
