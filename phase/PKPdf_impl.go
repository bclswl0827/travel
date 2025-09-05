package phase

import (
	"errors"

	"github.com/bclswl0827/travel/resource"
)

type PKPdf_Impl struct {
	pkpdfMap resource.TravelTimeMap
}

func (p *PKPdf_Impl) Name() string {
	return "PKPdf"
}

func (p *PKPdf_Impl) Boundary() (minDelta, maxDelta, minDepth, maxDepth float64) {
	minDelta, maxDelta, minDepth, maxDepth = p.pkpdfMap.Boundary()
	return minDelta, maxDelta, minDepth, maxDepth
}

func (p *PKPdf_Impl) Estimate(delta, depth float64, interplate bool) (result resource.TravelTimeEntry, err error) {
	return p.pkpdfMap.Query(delta, depth, interplate)
}

func NewPKPdf(pkpdfData string) (*PKPdf_Impl, error) {
	if len(pkpdfData) == 0 {
		return nil, errors.New("empty PKPdf data")
	}

	mapObj, err := resource.NewTravelTimeMap(pkpdfData)
	if err != nil {
		return nil, err
	}

	return &PKPdf_Impl{pkpdfMap: mapObj}, nil
}
