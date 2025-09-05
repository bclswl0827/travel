package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type ScP_Impl struct {
	scpMap resource.TravelTimeMap
}

func (p *ScP_Impl) Name() string {
	return "ScP"
}

func (p *ScP_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.scpMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *ScP_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.scpMap.Query(delta, depth, interplate)
}

func NewScP(scpData string) (*ScP_Impl, error) {
	if len(scpData) == 0 {
		return nil, errors.New("empty ScP data")
	}

	mapObj, err := resource.NewTravelTimeMap(scpData)
	if err != nil {
		return nil, err
	}

	return &ScP_Impl{scpMap: mapObj}, nil
}
