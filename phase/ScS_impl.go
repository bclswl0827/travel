package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type ScS_Impl struct {
	scsMap resource.TravelTimeMap
}

func (p *ScS_Impl) Name() string {
	return "ScS"
}

func (p *ScS_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.scsMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *ScS_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.scsMap.Query(delta, depth, interplate)
}

func NewScS(scsData string) (*ScS_Impl, error) {
	if len(scsData) == 0 {
		return nil, errors.New("empty ScS data")
	}

	mapObj, err := resource.NewTravelTimeMap(scsData)
	if err != nil {
		return nil, err
	}

	return &ScS_Impl{scsMap: mapObj}, nil
}
