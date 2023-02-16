package routing

type RoutingFunction interface {
	Init(id int, nodeType string) ([][]byte, error)
	GenMessageFromM(flit []byte) ([][]byte, error)
	GenMessageFromI(i []byte) ([][]byte, error)
}
