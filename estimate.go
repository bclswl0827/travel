package travel

func (p *AK135) Estimate(delta, depth float64, interplate bool) Estimation {
	var estimation Estimation

	if result, err := p.P.Estimate(delta, depth, interplate); err == nil {
		estimation.P = &result
	}
	if result, err := p.PcP.Estimate(delta, depth, interplate); err == nil {
		estimation.PcP = &result
	}
	if result, err := p.PKPab.Estimate(delta, depth, interplate); err == nil {
		estimation.PKPab = &result
	}
	if result, err := p.PKPbc.Estimate(delta, depth, interplate); err == nil {
		estimation.PKPbc = &result
	}
	if result, err := p.PKPdf.Estimate(delta, depth, interplate); err == nil {
		estimation.PKPdf = &result
	}
	if result, err := p.S.Estimate(delta, depth, interplate); err == nil {
		estimation.S = &result
	}
	if result, err := p.ScP.Estimate(delta, depth, interplate); err == nil {
		estimation.ScP = &result
	}
	if result, err := p.ScS.Estimate(delta, depth, interplate); err == nil {
		estimation.ScS = &result
	}
	if result, err := p.SKP.Estimate(delta, depth, interplate); err == nil {
		estimation.SKP = &result
	}
	if result, err := p.SKSac.Estimate(delta, depth, interplate); err == nil {
		estimation.SKSac = &result
	}
	if result, err := p.SKSdf.Estimate(delta, depth, interplate); err == nil {
		estimation.SKSdf = &result
	}

	return estimation
}
