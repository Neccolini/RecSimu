package network

import "github.com/Neccolini/RecSimu/cmd/network"

func (r *RF) Reconfigure() []network.Pair {
	return []network.Pair{}
}

func (r *RF) InitReConfiguration(id string) []network.Pair {
	// 障害を検知したら，隣接ノードにブロードキャストする．
	// 子ノードにブロードキャストし，障害が発生していることを知らせる．
	//親側ノードは親側にブロードキャストして知らせる
	if id == r.pId {
		// ブロードキャストする
		packet := Packet{r.id, BroadCastId, r.id, BroadCastId, "preq"}
	}

}