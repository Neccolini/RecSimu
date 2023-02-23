package read

type RecInfo struct {
	Id        string   `json:"id"`
	Cycle     int      `json:"cycle"`
	Operation string   `json:"operation"`
	AdjacencyList []string `json:"adjacencies"`
}

func MapCycle2RecInfo(recInfos []RecInfo) map[int][]RecInfo{
	m := map[int][]RecInfo{}
	for _, recInfo := range recInfos {
		m[recInfo.Cycle] = append(m[recInfo.Cycle], recInfo)
	}
	return m
}
