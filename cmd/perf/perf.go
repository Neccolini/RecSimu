package perf

const FlitLen = 6

type Perf struct {
	TotalLatency      int
	TotalPacketNum    int
	TotalFlitNum      int
	InjectionId2Cycle map[string]int
}

func NewPerf() *Perf {
	return &Perf{0, 0, 0, map[string]int{}}
}

func (p *Perf) Contains(injectionId string) bool {
	if _, ok := p.InjectionId2Cycle[injectionId]; ok {
		return true
	}
	return false
}

func (p *Perf) Start(injectionId string, cycle int) {
	p.InjectionId2Cycle[injectionId] = cycle
}

func (p *Perf) End(injectionId string, cycle int) {
	p.TotalLatency += cycle - p.InjectionId2Cycle[injectionId] + 1 + FlitLen
	p.TotalPacketNum += 1
	p.TotalFlitNum += FlitLen
}
