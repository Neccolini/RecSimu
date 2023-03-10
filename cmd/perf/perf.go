package perf

const FlitLen = 6

type Perf struct {
	TotalLatency      int
	TotalPacketNum    int
	TotalFlitNum      int
	InjectionId2Cycle map[string]int
	recPerf           RecPerf
}

func NewPerf() *Perf {
	return &Perf{0, 0, 0, map[string]int{}, RecPerf{}}
}

func (p *Perf) Contains(injectionId string) bool {
	if _, ok := p.InjectionId2Cycle[injectionId]; ok {
		return true
	}
	return false
}

func (p *Perf) Start(injectionId string, cycle int) {
	p.InjectionId2Cycle[injectionId] = cycle
	p.TotalPacketNum += 1
}

func (p *Perf) End(injectionId string, cycle int) {
	p.TotalLatency += cycle - p.InjectionId2Cycle[injectionId] + 1 + FlitLen
	delete(p.InjectionId2Cycle, injectionId)
	p.TotalFlitNum += FlitLen
}

func (p *Perf) FailedPacketNum() int {
	return len(p.InjectionId2Cycle)
}
