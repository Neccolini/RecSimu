package network

import (
	"github.com/Neccolini/RecSimu/cmd/network"
)

const (
	RecResendMax = 3
)

type RecState struct {
	on                bool
	resend            int
	childRequestIndex int
	childList         []string
	prevParentId      string
	waiting           bool
	isParentAlive     bool
}

func (r *RecState) NextChild() bool {
	if r.childRequestIndex < len(r.childList) {
		r.childRequestIndex++
	}
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
	r.recState.prevParentId = r.pId
	r.pId = ""
	r.joined = false
	r.recState.on = true

	p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preqR"}
	return []network.Pair{{Data: p.Serialize(), ToId: BroadCastId}}
}

func (r *RF) reconfigure() []network.Pair {
	// 子がブロードキャストしてくれるのを待つ
	if r.recState.waiting {
		return []network.Pair{}
	}
	// todo ダメそうだったら子ノードに配信を依頼する
	if r.recState.resend == RecResendMax {
		r.recState.resend = 0
		// leaf node なら
		if len(r.recState.childList) == 0 {
			// 自分の親に対してfailedを送る
			p := Packet{r.id, r.recState.prevParentId, r.id, r.recState.prevParentId, "fail"}
			r.recState.waiting = true
			return []network.Pair{{Data: p.Serialize(), ToId: r.recState.prevParentId}}
		}
		childId := r.recState.childList[r.recState.childRequestIndex]

		p := Packet{r.id, childId, r.id, childId, "rreq"}
		r.recState.waiting = true
		return []network.Pair{{Data: p.Serialize(), ToId: childId}}
	}
	r.recState.resend++

	p := Packet{r.id, BroadCastId, r.id, BroadCastId, "preqR"}
	return []network.Pair{{Data: p.Serialize(), ToId: BroadCastId}}
}
