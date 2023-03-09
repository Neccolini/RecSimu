package node

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/message"
)

func (n *Node) Remove() error {
	if !n.Alive() {
		log.Fatalf("node %s is already removed", n.Id())
	}

	n.Reset()

	return nil
}

func (n *Node) Rejoin() error {
	if n.Alive() {
		log.Fatalf("node %s is not removed", n.Id())
	}

	n.nodeAlive = true

	return nil
}

func (n *Node) InitReconfiguration(id string) {
	// 自身の親が削除ノードなら，その情報を削除しブロードキャスト
	pairs := n.RoutingFunction.InitReconfiguration()
	for _, pair := range pairs {
		n.sendMessages.Push(*message.NewMessage(fmt.Sprintf("Rec_%s_%s_%d", n.nodeId, pair.ToId, n.curCycle), n.nodeId, pair.ToId, pair.Data))
	}
}
