package network

import (
	"github.com/Neccolini/RecSimu/cmd/debug"

	"github.com/Neccolini/RecSimu/cmd/network"
	"github.com/Neccolini/RecSimu/cmd/set"
)

const (
	RecResendMax = 3
)

type RecState struct {
	on                 bool
	resend             int
	childRequestIndex  int
	childList          []string
	prevParentId       string
	waiting            bool
	isParentAlive      bool
	broadcastedRecFlag bool
	isUpNode           set.Set[string]
}

func NewRecState() *RecState {
	return &RecState{false, 0, 0, []string{}, "", false, false, false, *set.NewSet("--1")}
}

func (r *RecState) NextChild() bool {
	if r.childRequestIndex < len(r.childList) {
		r.childRequestIndex++
	}
	// すべての子を見終わっていたらtrueを返す
	return r.childRequestIndex >= len(r.childList)
}

func (r *RecState) ChildListContains(targetId string) bool {
	for _, id := range r.childList {
		if id == targetId {
			return true
		}
	}
	return false
}

func (r *RecState) Reset() {
	r.on = false
	r.resend = 0
	r.childRequestIndex = 0
	r.prevParentId = ""
	r.waiting = false
	r.isParentAlive = false
}

func (r *RF) InitReconfiguration() []network.Pair {
	if r.recState.prevParentId == "" {
		r.recState.prevParentId = r.pId
	}
	r.pId = ""
	r.recState.on = true
	r.recState.waiting = false

	p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preqR"}
	return []network.Pair{{Data: p.Serialize(), ToId: BroadCastId}}
}

func (r *RF) reconfigure() []network.Pair {
	r.recState.on = true
	// 子がブロードキャストしてくれるのを待つ
	// fmt.Println(r.id, r.recState.waiting, r.recState.broadcastedRecFlag, r.recState.resend)
	if r.recState.waiting {
		return []network.Pair{}
	}
	if !r.recState.broadcastedRecFlag {
		packetList := r.multiCastChildren("rec")
		r.recState.broadcastedRecFlag = true
		return packetList
	}
	// todo ダメそうだったら子ノードに配信を依頼する
	if r.recState.resend == RecResendMax {
		r.recState.resend = 0
		// fmt.Printf("%s %d %v\n", r.id, r.recState.childRequestIndex, r.recState.childList)
		// leaf node なら
		if len(r.recState.childList) == 0 || r.recState.childRequestIndex >= len(r.recState.childList) {
			// 自分の親に対してfailedを送る
			if r.recState.isParentAlive {
				p := Packet{r.id, r.recState.prevParentId, r.id, r.recState.prevParentId, "fail"}
				r.recState.waiting = true
				r.recState.childRequestIndex = 0
				return []network.Pair{{Data: p.Serialize(), ToId: r.recState.prevParentId}}
			} else if len(r.recState.childList) > 0 {
				r.recState.childRequestIndex = 0
			} else {
				return []network.Pair{}
			}
		}

		childId := r.recState.childList[r.recState.childRequestIndex]
		debug.Debug.Println(r.id, "rreq 送信 to", childId)
		p := Packet{r.id, childId, r.id, childId, "rreq"}
		r.recState.waiting = true

		r.recState.NextChild()

		return []network.Pair{{Data: p.Serialize(), ToId: childId}}
	}
	r.recState.resend++

	p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preqR"}
	return []network.Pair{{Data: p.Serialize(), ToId: BroadCastId}}
}

func (r *RF) failReceive() []network.Pair {
	r.recState.waiting = false
	return r.reconfigure()
}

func (r *RF) updateTableValue(from string, to string) {
	for k, v := range r.table {
		if v == from {
			r.table[k] = to
		}
	}
}
