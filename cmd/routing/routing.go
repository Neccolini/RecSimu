package routing

import (
	"encoding/json"
	"fmt"
	"log"
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

func (r *RF) Init(id string, nodeType string) ([][]byte, string) {
	r.id = id
	r.nodeType = nodeType
	r.pId = ""
	r.table = map[string]string{}
	if nodeType == Router {
		// parent request 送信
		p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preq"}
		return [][]byte{p.Serialize()}, "BroadCast"
	}

	r.joined = true
	return nil, ""
}

func (r *RF) IsJoined() bool {
	return r.joined
}

func (r *RF) Reset() {

}

func (r *RF) GenMessageFromM(received []byte) ([][]byte, string) {
	packet := DeserializeFrom(received)
	if packet.NextId != r.id && packet.NextId != BroadCastId {
		return nil, ""
	}
	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			r.table[packet.FromId] = packet.FromId
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "pack"}
			return [][]byte{reply.Serialize()}, packet.FromId
		} else if packet.Data == "jreq" {
			// jackを来た方向に返す
			jack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "jack"}
			return [][]byte{jack.Serialize()}, packet.PrevId
		}
	} else {
		if r.IsJoined() {
			if packet.Data == "preq" {
				r.table[packet.FromId] = packet.FromId
				pack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "pack"}
				return [][]byte{pack.Serialize()}, packet.PrevId
			} else {
				r.table[packet.FromId] = packet.PrevId
				sendPacket := r.routingPacket(packet)
				return [][]byte{sendPacket.Serialize()}, sendPacket.NextId
			}
		}
		if packet.Data == "pack" && r.pId == "" {
			r.table[packet.FromId] = packet.FromId
			r.table[CoordinatorId] = packet.FromId
			r.pId = packet.FromId
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			return [][]byte{jreq.Serialize()}, r.pId
		} else if packet.Data == "jack" {
			r.joined = true
			fmt.Printf("%s joined Network\n", r.id)
			return nil, ""
		}
	}
	return nil, ""
}

func (r *RF) GenMessageFromI(inst []byte) ([][]byte, string) {
	return nil, ""
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
		log.Fatalf("error during packet deserialization %v", err)
	}
	return packet
}
