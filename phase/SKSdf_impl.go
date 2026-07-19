package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type SKSdf_Impl struct {
	sksdfMap resource.TravelTimeMap
}

func (p *SKSdf_Impl) Name() string {
	return "SKIKS"
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
	// TauP's computed SKIKS branch begins just beyond 105 degrees for these
	// depths, despite the rounded published table containing 105-degree cells.
	removeEntries(mapObj, 105, 200, 300, 500, 700)

	return &SKSdf_Impl{sksdfMap: mapObj}, nil
}

// SKIKS_Impl is TauP's canonical name for the phase labelled SKSdf in the
// published AK135 tables.
type SKIKS_Impl = SKSdf_Impl

func NewSKIKS(skiksData string) (*SKIKS_Impl, error) {
	return NewSKSdf(skiksData)
}
