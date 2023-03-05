package read

import "github.com/Neccolini/RecSimu/cmd/injection"

type RecInfo struct {
	Id            string   `json:"id"`
	Cycle         int      `json:"cycle"`
	Operation     string   `json:"operation"`
	AdjacencyList []string `json:"adjacencies"`
}

func MapCycle2RecInfo(recInfos []RecInfo) map[int][]RecInfo {
	m := map[int][]RecInfo{}
	for _, recInfo := range recInfos {
		m[recInfo.Cycle] = append(m[recInfo.Cycle], recInfo)
	}
	return m
}

func shapeInjections(injections []injection.Injection) *injection.InjectionTable {
	// for文を用いてinjectionMapと cycleQueryを生成する
	injectionMap := map[string]injection.Injection{}
	cycleQuery := map[int][]string{}
	for _, inj := range injections {
		injectionMap[inj.InjectionId] = inj
		cycleQuery[inj.Cycle] = append(cycleQuery[inj.Cycle], inj.InjectionId)
	}
	return &injection.InjectionTable{InjectionMap: injectionMap, CycleQuery: cycleQuery}
}
