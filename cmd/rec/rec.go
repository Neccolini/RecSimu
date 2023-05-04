package rec

const (
	Add    = "add"
	Rejoin = "rejoin"
	Remove = "remove"
)

type RecInfo struct {
	Id            string   `json:"id"`
	NodeType      string   `json:"node_type"`
	Cycle         int      `json:"cycle"`
	Operation     string   `json:"operation"`
	AdjacencyList []string `json:"adjacencies"`
}
