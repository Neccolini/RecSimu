package perf

type RecPerf struct {
	recCycles    []int
	recInitCycle int
}

func (p *Perf) RecStart(cycle int) {
	if p.recPerf.recInitCycle != 0 {
		return
	}
	p.recPerf.recInitCycle = cycle
}

func (p *Perf) RecEnd(cycle int) {
	if p.recPerf.recInitCycle == 0 {
		return
	}
	p.recPerf.recCycles = append(p.recPerf.recCycles, cycle-p.recPerf.recInitCycle)
	p.recPerf.recInitCycle = 0
}

func (p *Perf) RecResult() []int {
	return p.recPerf.recCycles
}
