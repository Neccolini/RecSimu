package routing

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math"
	"math/big"
)

const (
	Router        = "Router"
	Coordinator   = "Coordinator"
	BroadCastId   = -1
	CoordinatorId = 0
)

type RF struct {
	id       int
	nodeType string
	joined   bool
	table    map[int][]int
}

type Packet struct {
	Id     int
	DistId int
	Data   string
}

func (r *RF) Init(nodeType string) ([][]byte, int) {
	r.nodeType = nodeType
	if nodeType == Router {
		// parent request 送信
		randomId, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
		if err != nil {
			log.Fatal(err)
		}
		r.id = int(randomId.Int64())

		p := Packet{r.id, BroadCastId, "preq"}
		return [][]byte{p.Serialize()}, BroadCastId
	}

	r.id = CoordinatorId
	r.joined = true
	return nil, 0
}

func (r *RF) IsJoined() bool {
	return r.joined
}

func (r *RF) GenMessageFromM(received []byte) ([][]byte, int) {
	packet := DeserializeFrom(received)
	if packet.DistId != r.id && packet.DistId != BroadCastId {
		return nil, 0
	}
	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			reply := Packet{r.id, packet.Id, "pack"}
			return [][]byte{reply.Serialize()}, packet.Id
		}
	}
	return nil, 0
}

func (r *RF) GenMessageFromI(inst []byte) ([][]byte, int) {
	return nil, 0
}

func (p *Packet) Serialize() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("error during packet serialization: %v", err)
	}
	return jsonData
}

func DeserializeFrom(data []byte) Packet {
	var packet Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		log.Fatalf("error during packet deserialization %v", err)
	}
	return packet
}
