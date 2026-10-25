package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type PKPbc_Impl struct {
	pkpbcMap resource.TravelTimeMap
}

func (p *PKPbc_Impl) Name() string {
	return "PKPbc"
}

func (p *PKPbc_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.pkpbcMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *PKPbc_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.pkpbcMap.Query(delta, depth, interplate)
}

func NewPKPbc(pkpbcData string) (*PKPbc_Impl, error) {
	if len(pkpbcData) == 0 {
		return nil, errors.New("empty PKPbc data")
	}

	mapObj, err := resource.NewTravelTimeMap(pkpbcData)
	if err != nil {
		return nil, err
	}

	return &PKPbc_Impl{pkpbcMap: mapObj}, nil
}
