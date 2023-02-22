package routing

type RoutingFunction interface {
	Init(nodeType string) ([][]byte, int)
	GenMessageFromM(packet []byte) ([][]byte, int)
	GenMessageFromI(i []byte) ([][]byte, int)
	IsJoined() bool
}
