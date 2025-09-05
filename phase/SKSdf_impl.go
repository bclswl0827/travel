package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type SKSdf_Impl struct {
	sksdfMap resource.TravelTimeMap
}

func (p *SKSdf_Impl) Name() string {
	return "SKSdf"
}

func (p *SKSdf_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.sksdfMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *SKSdf_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.sksdfMap.Query(delta, depth, interplate)
}

func NewSKSdf(sksdfData string) (*SKSdf_Impl, error) {
	if len(sksdfData) == 0 {
		return nil, errors.New("empty SKSdf data")
	}

	mapObj, err := resource.NewTravelTimeMap(sksdfData)
	if err != nil {
		return nil, err
	}

	return &SKSdf_Impl{sksdfMap: mapObj}, nil
}
