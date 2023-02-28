package routing

type Pair struct {
	Data []byte
	ToId string
}

type RoutingFunction interface {
	Init() []Pair
	GenMessageFromM(packet []byte) []Pair
	GenMessageFromI(i []byte) []Pair
	IsJoined() bool
	ParentId() string
	Reset()
}
