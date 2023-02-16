package routing

type RF struct {
	id int
	nodeType string
	table map[int][]int
}

const (
	Router = "router"
	Coordinator = "Coordinator"
)

func (r *RF) Init(id int, nodeType string) ([][]byte, error) {
	r.id = id
	r.nodeType = nodeType
	return nil, nil
}
func (r *RF) GenMessageFromM(received []byte) ([][]byte, error) {
	return nil, nil
}

func (r *RF) GenMessageFromI(inst []byte) ([][]byte, error) {
	return nil, nil
}
