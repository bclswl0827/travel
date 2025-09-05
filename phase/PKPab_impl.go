package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type PKPab_Impl struct {
	pkpabMap resource.TravelTimeMap
}

func (p *PKPab_Impl) Name() string {
	return "PKPab"
}

func (p *PKPab_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.pkpabMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *PKPab_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.pkpabMap.Query(delta, depth, interplate)
}

func NewPKPab(pkpabData string) (*PKPab_Impl, error) {
	if len(pkpabData) == 0 {
		return nil, errors.New("empty PKPab data")
	}

	mapObj, err := resource.NewTravelTimeMap(pkpabData)
	if err != nil {
		return nil, err
	}

	return &PKPab_Impl{pkpabMap: mapObj}, nil
}
