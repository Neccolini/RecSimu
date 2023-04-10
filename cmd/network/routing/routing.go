package network

import (
	"strings"

	"github.com/Neccolini/RecSimu/cmd/debug"
	"github.com/Neccolini/RecSimu/cmd/network"
)

const (
	BytePerFlit   = 64
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
	r.recState = *NewRecState()
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
	pair := []network.Pair{}

	if packet.NextId != r.id && packet.NextId != BroadCastId {
		return pair
	}
	if r.drainPacket(packet) {
		return pair
	}
	debug.Debug.Printf("%s: %v\n", r.id, packet)
	if packet.DistId == r.id && packet.Data == "" {
		return r.ProcessMessage(packet.Data)
	}

	if r.nodeType == Coordinator {
		if packet.Data == "preq" {
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "pack"}
			pair = []network.Pair{{Data: reply.Serialize(), ToId: packet.FromId}}
		} else if packet.Data == "preqR" && !r.recState.isUpNode.Contains(packet.FromId) {
			reply := Packet{r.id, packet.FromId, r.id, packet.FromId, "packR"}
			pair = []network.Pair{{Data: reply.Serialize(), ToId: packet.FromId}}
		} else if len(packet.Data) >= 4 && packet.Data[:4] == "jreq" {
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

				r.InitReconfiguration()
				return r.reconfigure()
			} else if packet.Data == "fail" {
				return r.failReceive()
			} else if packet.Data == "rec" {
				r.recState.on = true
				r.recState.waiting = true
				r.recState.broadcastedRecFlag = true
				packetList := []network.Pair{}
				for _, distId := range r.recState.childList {
					p := Packet{r.id, distId, r.id, distId, "rec"}
					packetList = append(packetList, network.Pair{Data: p.Serialize(), ToId: distId})
				}
				r.recState.broadcastedRecFlag = true
				return packetList
			} else if packet.Data == "packR" && packet.DistId == r.id {
				// r.recState.on = false
				r.pId = packet.FromId
				jreq := Packet{r.id, CoordinatorId, r.id, packet.PrevId, "jreqR"}
				pair = []network.Pair{{Data: jreq.Serialize(), ToId: packet.PrevId}}
			} else if len(packet.Data) >= 5 && packet.Data[:5] == "jackR" && packet.DistId == r.id {
				if packet.FromId != CoordinatorId && !r.recState.on {
					return []network.Pair{}
				}
				debug.Debug.Println("OK", r.recState.prevParentId, r.recState.isParentAlive)
				// 自分宛にjackRが届いた
				if r.recState.isParentAlive {
					p := Packet{r.id, r.recState.prevParentId, r.id, r.recState.prevParentId, packet.Data}
					pair = append(pair, network.Pair{Data: p.Serialize(), ToId: r.recState.prevParentId})
				}

				if len(packet.Data) >= 7 {
					r.recState.isUpNode.Reset()
					arr := strings.Split(packet.Data[5:], "/")
					for _, upId := range arr {
						r.recState.isUpNode.Add(upId)
					}
				}

				// 親を設定
				r.pId = packet.PrevId

				r.updateTableValue(r.recState.prevParentId, r.pId)

				// 子に対して再構成が完了したことを伝える．
				pair = append(pair, r.multiCastChildren(packet.Data)...)

				// 再構成終了
				r.recState.Reset()

				debug.Debug.Printf("%s rejoined Network\n", r.id)
				pair = append(pair, network.Pair{Data: nil, ToId: Joined})
				return pair
			} else {
				// childList に追加
				if packet.Data == "jreq" && packet.PrevId == packet.FromId && !r.recState.ChildListContains(packet.FromId) {
					r.recState.childList = append(r.recState.childList, packet.FromId)
				}

				// r.table[packet.FromId] = packet.PrevId
				sendPacket := r.routingPacket(packet)
				if sendPacket != nil {
					pair = []network.Pair{{Data: sendPacket.Serialize(), ToId: sendPacket.NextId}}
				}
			}
		} else if len(packet.Data) >= 4 && packet.Data[:4] == "pack" && r.pId == "" {
			// r.recState.on = false
			// r.table[CoordinatorId] = packet.PrevId
			r.pId = packet.FromId
			jreq := Packet{r.id, CoordinatorId, r.id, r.pId, "jreq"}
			if len(packet.Data) == 5 {
				jreq = Packet{r.id, CoordinatorId, r.id, r.pId, "jreqR"}
			}

			pair = []network.Pair{{Data: jreq.Serialize(), ToId: r.pId}}
		} else if len(packet.Data) >= 4 && packet.Data[:4] == "jack" {
			r.joined = true
			if r.recState.isParentAlive {
				p := Packet{r.id, r.recState.prevParentId, r.id, r.recState.prevParentId, "packR"}
				pair = append(pair, network.Pair{Data: p.Serialize(), ToId: r.recState.prevParentId})
			}

			if len(packet.Data) >= 6 {
				arr := strings.Split(packet.Data[4:], "/")
				for _, upId := range arr {
					r.recState.isUpNode.Add(upId)
				}
			}

			debug.Debug.Printf("%s joined Network\n", r.id)
			pair = append(pair, network.Pair{Data: nil, ToId: Joined})
		}
	}
	return pair
}

func (r *RF) GenMessageFromI(distId string, data string) []network.Pair {
	nextId := r.pId
	if val, ok := r.table[distId]; ok {
		nextId = val
	}
	packet := Packet{r.id, distId, r.id, nextId, data}
	return []network.Pair{{Data: packet.Serialize(), ToId: nextId}}
}

func (r *RF) routingPacket(p Packet) *Packet {
	if p.DistId == BroadCastId {
		return nil
	}
	if p.FromId == r.id {
		return nil
	}
	if p.DistId == r.id {
		return nil
	}
	if len(p.Data) >= 4 && p.Data[:4] == "jreq" {
		r.table[p.FromId] = p.PrevId
	}

	var neighborDistId string
	// テーブルに存在したら
	if val, ok := r.table[p.DistId]; ok {
		neighborDistId = val
	} else { // テーブルに存在しない場合
		neighborDistId = r.pId
	}
	// 新規ノードにとって，Up方向のノード番号がわかるようにデータ部分に自身のIDを追加する
	if len(p.Data) >= 4 && p.Data[:4] == "jack" {
		p.Data += "/" + r.id
	}
	routingPacket := Packet{p.FromId, p.DistId, r.id, neighborDistId, p.Data}
	debug.Debug.Printf("%s routing %v -> %v\n", r.id, p, routingPacket)
	return &routingPacket
}

func (r *RF) drainPacket(p Packet) bool {
	if _, ok := r.table[p.FromId]; ok &&
		(p.Data == "jreq" || p.Data == "preq") {
		return true
	}
	return false
}

func (r *RF) multiCastChildren(msg string) []network.Pair {
	packetList := []network.Pair{}
	for _, distId := range r.recState.childList {
		p := Packet{r.id, distId, r.id, distId, msg}
		packetList = append(packetList, network.Pair{Data: p.Serialize(), ToId: distId})
	}
	return packetList
}
