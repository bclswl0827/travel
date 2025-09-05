package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type PcP_Impl struct {
	pcpMap resource.TravelTimeMap
}

func (p *PcP_Impl) Name() string {
	return "PcP"
}

func (p *PcP_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.pcpMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *PcP_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.pcpMap.Query(delta, depth, interplate)
}

func NewPcP(pcpData string) (*PcP_Impl, error) {
	if len(pcpData) == 0 {
		return nil, errors.New("empty PcP data")
	}

	mapObj, err := resource.NewTravelTimeMap(pcpData)
	if err != nil {
		return nil, err
	}

	return &PcP_Impl{pcpMap: mapObj}, nil
}
