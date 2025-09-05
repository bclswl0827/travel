package travel

import (
	"fmt"

	"github.com/bclswl0827/travel/phase"
	"github.com/bclswl0827/travel/resource"
)

type AK135 struct {
	P     IPhase
	PcP   IPhase
	PKPab IPhase
	PKPbc IPhase
	PKPdf IPhase
	S     IPhase
	ScP   IPhase
	ScS   IPhase
	SKP   IPhase
	SKSac IPhase
	SKSdf IPhase
}

type Estimation struct {
	P     *resource.TravelTimeEntry
	PcP   *resource.TravelTimeEntry
	PKPab *resource.TravelTimeEntry
	PKPbc *resource.TravelTimeEntry
	PKPdf *resource.TravelTimeEntry
	S     *resource.TravelTimeEntry
	ScP   *resource.TravelTimeEntry
	ScS   *resource.TravelTimeEntry
	SKP   *resource.TravelTimeEntry
	SKSac *resource.TravelTimeEntry
	SKSdf *resource.TravelTimeEntry
}

func NewAK135() (*AK135, error) {
	P, err := phase.NewP(resource.GetAK135Table(resource.AK135_P_SHALLOW), resource.GetAK135Table(resource.AK135_P_DEEP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewP: %w", err)
	}

	PcP, err := phase.NewPcP(resource.GetAK135Table(resource.AK135_PcP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewPcP: %w", err)
	}
	PKPab, err := phase.NewPKPab(resource.GetAK135Table(resource.AK135_PKPab))
	if err != nil {
		return nil, fmt.Errorf("phase.NewPKPab: %w", err)
	}
	PKPbc, err := phase.NewPKPbc(resource.GetAK135Table(resource.AK135_PKPbc))
	if err != nil {
		return nil, fmt.Errorf("phase.NewPKPbc: %w", err)
	}
	PKPdf, err := phase.NewPKPdf(resource.GetAK135Table(resource.AK135_PKPdf))
	if err != nil {
		return nil, fmt.Errorf("phase.NewPKPdf: %w", err)
	}
	S, err := phase.NewS(resource.GetAK135Table(resource.AK135_S_SHALLOW), resource.GetAK135Table(resource.AK135_S_DEEP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewS: %w", err)
	}
	ScP, err := phase.NewScP(resource.GetAK135Table(resource.AK135_ScP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewScP: %w", err)
	}
	ScS, err := phase.NewScS(resource.GetAK135Table(resource.AK135_ScS))
	if err != nil {
		return nil, fmt.Errorf("phase.NewScS: %w", err)
	}
	SKP, err := phase.NewSKP(resource.GetAK135Table(resource.AK135_SKP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKP: %w", err)
	}
	SKSac, err := phase.NewSKSac(resource.GetAK135Table(resource.AK135_SKSac))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKSac: %w", err)
	}
	SKSdf, err := phase.NewSKSdf(resource.GetAK135Table(resource.AK135_SKSdf))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKSdf: %w", err)
	}
	return &AK135{
		P:     P,
		PcP:   PcP,
		PKPab: PKPab,
		PKPbc: PKPbc,
		PKPdf: PKPdf,
		S:     S,
		ScP:   ScP,
		ScS:   ScS,
		SKP:   SKP,
		SKSac: SKSac,
		SKSdf: SKSdf,
	}, nil
}
