package routing

type RoutingFunction interface {
	Init(id string, nodeType string) ([][]byte, string)
	GenMessageFromM(packet []byte) ([][]byte, string)
	GenMessageFromI(i []byte) ([][]byte, string)
	IsJoined() bool
}
