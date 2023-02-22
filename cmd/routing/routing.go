package routing

import (
	"encoding/json"
	"log"
)

const (
	Router        = "Router"
	Coordinator   = "Coordinator"
	BroadCastId   = "BroadCast"
	CoordinatorId = "0"
)

type RF struct {
	id  string
	nodeType string
	joined   bool
	table    map[int][]int
}

type Packet struct {
	Id     string
	DistId string
	Data   string
}

func (r *RF) Init(id string, nodeType string) ([][]byte, string) {
	r.id = id
	r.nodeType = nodeType
	if nodeType == Router {
		// parent request 送信

		p := Packet{r.id, BroadCastId, "preq"}
		return [][]byte{p.Serialize()}, "BroadCast"
	}

	r.joined = true
	return nil, ""
}

func (r *RF) IsJoined() bool {
	return r.joined
}

func (r *RF) GenMessageFromM(received []byte) ([][]byte, string) {
	packet := DeserializeFrom(received)
	if packet.DistId != r.id && packet.DistId != BroadCastId {
		return nil, ""
	}
	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			reply := Packet{r.id, packet.Id, "pack"}
			return [][]byte{reply.Serialize()}, packet.Id
		}
	}
	return nil, ""
}

func (r *RF) GenMessageFromI(inst []byte) ([][]byte, string) {
	return nil, ""
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
