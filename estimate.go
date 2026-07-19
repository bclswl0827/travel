package travel

func (p *AK135) Estimate(delta, depth float64, interpolate bool) Estimation {
	var estimation Estimation

	if result, err := p.P.Estimate(delta, depth, interpolate); err == nil {
		estimation.P = &result
	}
	if result, err := p.PcP.Estimate(delta, depth, interpolate); err == nil {
		estimation.PcP = &result
	}
	if result, err := p.PKPab.Estimate(delta, depth, interpolate); err == nil {
		estimation.PKPab = &result
	}
	if result, err := p.PKPbc.Estimate(delta, depth, interpolate); err == nil {
		estimation.PKPbc = &result
	}
	if result, err := p.PKIKP.Estimate(delta, depth, interpolate); err == nil {
		estimation.PKIKP = &result
	}
	if result, err := p.S.Estimate(delta, depth, interpolate); err == nil {
		estimation.S = &result
	}
	if result, err := p.ScP.Estimate(delta, depth, interpolate); err == nil {
		estimation.ScP = &result
	}
	if result, err := p.ScS.Estimate(delta, depth, interpolate); err == nil {
		estimation.ScS = &result
	}
	if result, err := p.SKP.Estimate(delta, depth, interpolate); err == nil {
		estimation.SKP = &result
	}
	if result, err := p.SKS.Estimate(delta, depth, interpolate); err == nil {
		estimation.SKS = &result
	}
	if result, err := p.SKIKS.Estimate(delta, depth, interpolate); err == nil {
		estimation.SKIKS = &result
	}

	return estimation
}
