package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type SKSac_Impl struct {
	sksacMap resource.TravelTimeMap
}

func (p *SKSac_Impl) Name() string {
	return "SKS"
}

func (p *SKSac_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.sksacMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *SKSac_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.sksacMap.Query(delta, depth, interplate)
}

func NewSKSac(sksacData string) (*SKSac_Impl, error) {
	if len(sksacData) == 0 {
		return nil, errors.New("empty SKSac data")
	}

	mapObj, err := resource.NewTravelTimeMap(sksacData)
	if err != nil {
		return nil, err
	}

	return &SKSac_Impl{sksacMap: mapObj}, nil
}

// SKS_Impl is TauP's canonical name for the branch labelled SKSac in the
// published AK135 tables.
type SKS_Impl = SKSac_Impl

func NewSKS(sksData string) (*SKS_Impl, error) {
	return NewSKSac(sksData)
}
