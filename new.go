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
	PKIKP IPhase
	S     IPhase
	ScP   IPhase
	ScS   IPhase
	SKP   IPhase
	SKIKP IPhase
	SKS   IPhase
	SKIKS IPhase
}

type Estimation struct {
	P     *resource.TravelTimeEntry
	PcP   *resource.TravelTimeEntry
	PKPab *resource.TravelTimeEntry
	PKPbc *resource.TravelTimeEntry
	PKIKP *resource.TravelTimeEntry
	S     *resource.TravelTimeEntry
	ScP   *resource.TravelTimeEntry
	ScS   *resource.TravelTimeEntry
	SKP   *resource.TravelTimeEntry
	SKIKP *resource.TravelTimeEntry
	SKS   *resource.TravelTimeEntry
	SKIKS *resource.TravelTimeEntry
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
	PKIKP, err := phase.NewPKIKP(resource.GetAK135Table(resource.AK135_PKIKP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewPKIKP: %w", err)
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
	SKIKP, err := phase.NewSKIKP(resource.GetAK135Table(resource.AK135_SKIKP))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKIKP: %w", err)
	}
	SKS, err := phase.NewSKS(resource.GetAK135Table(resource.AK135_SKS))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKS: %w", err)
	}
	SKIKS, err := phase.NewSKIKS(resource.GetAK135Table(resource.AK135_SKIKS))
	if err != nil {
		return nil, fmt.Errorf("phase.NewSKIKS: %w", err)
	}
	return &AK135{
		P:     P,
		PcP:   PcP,
		PKPab: PKPab,
		PKPbc: PKPbc,
		PKIKP: PKIKP,
		S:     S,
		ScP:   ScP,
		ScS:   ScS,
		SKP:   SKP,
		SKIKP: SKIKP,
		SKS:   SKS,
		SKIKS: SKIKS,
	}, nil
}
