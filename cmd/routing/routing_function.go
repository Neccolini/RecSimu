package routing

type RoutingFunction interface {
	Init(nodeType string) ([][]byte, error)
	GenMessageFromM(packet []byte) ([][]byte, error)
	GenMessageFromI(i []byte) ([][]byte, error)
}
