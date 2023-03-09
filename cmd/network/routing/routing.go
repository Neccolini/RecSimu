package network

import (
	"github.com/Neccolini/RecSimu/cmd/debug"
	"github.com/Neccolini/RecSimu/cmd/network"
)

const (
	BytePerFlit   = 16
	Router        = "Router"
	Coordinator   = "Coordinator"
	BroadCastId   = "BroadCast"
	Joined        = "Joined"
	CoordinatorId = "0"
)

type RF struct {
	id       string
	nodeType string
	pId      string
	joined   bool
	table    map[string]string // adjacentId -> DistId
	recState RecState
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

func (r *RF) Init() []network.Pair {
	if r.recState.on {
		return r.reconfigure()
	}
	if r.nodeType == Router {
		if r.pId == "" {
			// parent request 送信
			p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preq"}

			return []network.Pair{
				{Data: p.Serialize(), ToId: BroadCastId},
			}
		} else if !r.joined {
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			return []network.Pair{{Data: jreq.Serialize(), ToId: r.pId}}
		}
	}
	return []network.Pair{}
}

func (r *RF) IsJoined() bool {
	return r.joined
}

func (r *RF) ParentId() string {
	return r.pId
}

func (r *RF) Reset() {

}

func (r *RF) ProcessMessage(m string) []network.Pair {
	return []network.Pair{}
}

func (r *RF) GenMessageFromM(received []byte) []network.Pair {
	packet := DeserializeFrom(received)
	debug.Debug.Printf("b id:%s pid:%s packet: %v\n", r.id, r.pId, packet)
	pair := []network.Pair{}

	if packet.NextId != r.id && packet.NextId != BroadCastId {
		return pair
	}
	if r.drainPacket(packet) {
		return pair
	}
	if packet.DistId == r.id && packet.Data == "" {
		return r.ProcessMessage(packet.Data)
	}
	debug.Debug.Printf("a id:%s pid:%s packet: %v\n", r.id, r.pId, packet)

	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "pack"}
			pair = []network.Pair{{Data: reply.Serialize(), ToId: packet.FromId}}
		} else if packet.Data == "preqR" {
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "packR"}
			pair = []network.Pair{{Data: reply.Serialize(), ToId: packet.FromId}}
		} else if packet.Data[:4] == "jreq" {
			r.table[packet.FromId] = packet.PrevId
			// jackを来た方向に返す
			jack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "jack"}
			if len(packet.Data) == 5 {
				jack = Packet{r.id, packet.FromId, r.id, packet.PrevId, "jackR"}
			}
			pair = []network.Pair{{Data: jack.Serialize(), ToId: packet.PrevId}}
		} else {
			sendPacket := r.routingPacket(packet)
			if sendPacket != nil {
				pair = []network.Pair{{Data: sendPacket.Serialize(), ToId: sendPacket.NextId}}
			}
		}
	} else {
		if r.IsJoined() {
			if packet.Data == "preq" {
				pack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "pack"}
				pair = []network.Pair{{Data: pack.Serialize(), ToId: packet.PrevId}}
			} else if packet.Data == "preqR" && packet.FromId != r.pId && !r.recState.on {
				pack := Packet{r.id, packet.FromId, r.id, packet.PrevId, "packR"}
				pair = []network.Pair{{Data: pack.Serialize(), ToId: packet.PrevId}}
			} else if packet.Data == "rreq" {
				r.recState.isParentAlive = true
				return r.InitReconfiguration()
			} else if packet.Data == "fail" {
				// 別の子に対してブロードキャストを依頼する
				if !r.recState.NextChild() {
					r.recState.childRequestIndex = 0
					return r.InitReconfiguration()
				}
				r.recState.waiting = true
				return r.reconfigure()
			} else {
				// childList に追加
				if packet.Data == "jreqR" && packet.PrevId == packet.FromId && !r.recState.ChildListContains(packet.FromId) {
					r.recState.childList = append(r.recState.childList, packet.FromId)
				}

				r.table[packet.FromId] = packet.PrevId
				sendPacket := r.routingPacket(packet)
				if sendPacket != nil {
					pair = []network.Pair{{Data: sendPacket.Serialize(), ToId: sendPacket.NextId}}
				}
			}
		} else if packet.Data[:4] == "pack" && r.pId == "" {
			r.recState.on = false
			r.table[packet.FromId] = packet.FromId
			r.table[CoordinatorId] = packet.FromId
			r.pId = packet.FromId
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			if len(packet.Data) == 5 {
				jreq = Packet{r.id, CoordinatorId, r.id, r.pId, "jreqR"}
			}
			pair = []network.Pair{{Data: jreq.Serialize(), ToId: r.pId}}
		} else if packet.Data[:4] == "jack" {
			r.joined = true
			if r.recState.isParentAlive {
				p := Packet{r.id, r.recState.prevParentId, r.id, r.recState.prevParentId, "packR"}
				pair = append(pair, network.Pair{Data: p.Serialize(), ToId: r.recState.prevParentId})
			}

			debug.Debug.Printf("%s joined Network\n", r.id)
			pair = append(pair, network.Pair{Data: nil, ToId: Joined})
		}
	}
	return pair
}

func (r *RF) GenMessageFromI(distId string, data string) []network.Pair {
	packet := Packet{r.id, distId, r.id, r.table[distId], data}
	return []network.Pair{{Data: packet.Serialize(), ToId: r.table[distId]}}
}

func (r *RF) routingPacket(p Packet) *Packet {
	if p.DistId == BroadCastId {
		return nil
	}
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

func (r *RF) drainPacket(p Packet) bool {
	if _, ok := r.table[p.FromId]; ok &&
		(p.Data == "jreq" || p.Data == "preq") {
		return true
	}
	return false
}
