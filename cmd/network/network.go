package network

type Pair struct {
	Data []byte
	ToId string
}

type RoutingFunction interface {
	Init() []Pair
	GenMessageFromM(packet []byte) []Pair
	GenMessageFromI(distId string, data string) []Pair
	IsJoined() bool
	ParentId() string
	Reset()
	Reconfigure() []Pair
	InitReConfiguration(id string) []Pair
}
