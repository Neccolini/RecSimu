package routing

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/debug"
)

const (
	BytePerFlit   = 16
	Router        = "Router"
	Coordinator   = "Coordinator"
	BroadCastId   = "BroadCast"
	CoordinatorId = "0"
)

type RF struct {
	id       string
	nodeType string
	pId      string
	joined   bool
	table    map[string]string // adjacentId -> DistId
}

type Packet struct {
	FromId string
	DistId string
	PrevId string
	NextId string
	Data   string
}

func NewRoutingFunction(id string, nodeType string) *RF {
	r := RF{}
	r.id = id
	r.nodeType = nodeType
	r.pId = ""
	r.table = map[string]string{}

	if r.nodeType == Coordinator {
		r.joined = true
	}
	return &r
}
func (r *RF) Init() []Pair {
	if r.nodeType == Router {
		if r.pId == "" {
			// parent request 送信
			p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preq"}

			return []Pair{
				{p.Serialize(), BroadCastId},
			}
		} else if !r.joined {
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			return []Pair{{jreq.Serialize(), r.pId}}
		}
	}
	return []Pair{}
}

func (r *RF) IsJoined() bool {
	return r.joined
}

func (r *RF) ParentId() string {
	return r.pId
}

func (r *RF) Reset() {

}

func (r *RF) GenMessageFromM(received []byte) []Pair {
	packet := DeserializeFrom(received)

	pair := []Pair{}

	if packet.NextId != r.id && packet.NextId != BroadCastId {
		return pair
	}
	debug.Debug.Printf("id:%s pid:%s packet: %v\n", r.id, r.pId, packet)

	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			r.table[packet.FromId] = packet.FromId
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "pack"}

			pair = []Pair{{reply.Serialize(), packet.FromId}}
		} else if packet.Data == "jreq" {
			// jackを来た方向に返す
			jack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "jack"}
			pair = []Pair{{jack.Serialize(), packet.PrevId}}
		}
	} else {
		if r.IsJoined() {
			if packet.Data == "preq" {
				r.table[packet.FromId] = packet.FromId
				pack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "pack"}
				pair = []Pair{{pack.Serialize(), packet.PrevId}}
			} else {
				r.table[packet.FromId] = packet.PrevId
				debug.Debug.Printf("%s %v\n", r.id, r.table)
				sendPacket := r.routingPacket(packet)
				pair = []Pair{{sendPacket.Serialize(), sendPacket.NextId}}
			}
		} else if packet.Data == "pack" && r.pId == "" {
			r.table[packet.FromId] = packet.FromId
			r.table[CoordinatorId] = packet.FromId
			r.pId = packet.FromId
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			pair = []Pair{{jreq.Serialize(), r.pId}}
		} else if packet.Data == "jack" {
			r.joined = true
			fmt.Printf("%s joined Network\n", r.id)
			return pair
		}
	}
	return pair
}

func (r *RF) GenMessageFromI(inst []byte) []Pair {
	return []Pair{}
}

func (r *RF) routingPacket(p Packet) *Packet {
	var neighborDistId string
	// テーブルに存在したら
	if val, ok := r.table[p.DistId]; ok {
		neighborDistId = val
	} else { // テーブルに存在しない場合
		neighborDistId = r.pId
		r.table[p.DistId] = p.PrevId
	}
	routingPacket := Packet{p.FromId, p.DistId, r.id, neighborDistId, p.Data}
	return &routingPacket
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
		log.Fatalf("error during packet deserialization: %v", err)
	}
	return packet
}
