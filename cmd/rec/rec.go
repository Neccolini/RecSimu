package rec

const (
	Add    = "add"
	Rejoin = "rejoin"
	Remove = "remove"
)

type RecInfo struct {
	Id            string   `json:"id"`
	NodeType      string   `json:"nodeType"`
	Cycle         int      `json:"cycle"`
	Operation     string   `json:"operation"`
	AdjacencyList []string `json:"adjacencies"`
}
