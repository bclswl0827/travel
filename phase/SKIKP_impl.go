package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type SKIKP_Impl struct {
	skikpMap resource.TravelTimeMap
}

func (p *SKIKP_Impl) Name() string {
	return "SKIKP"
}

func (p *SKIKP_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	return p.skikpMap.Boundary()
}

func (p *SKIKP_Impl) Estimate(delta, depth float64, interpolate bool) (resource.TravelTimeEntry, error) {
	return p.skikpMap.Query(delta, depth, interpolate)
}

func NewSKIKP(skikpData string) (*SKIKP_Impl, error) {
	if len(skikpData) == 0 {
		return nil, errors.New("empty SKIKP data")
	}

	mapObj, err := resource.NewTravelTimeMap(skikpData)
	if err != nil {
		return nil, err
	}

	return &SKIKP_Impl{skikpMap: mapObj}, nil
}
