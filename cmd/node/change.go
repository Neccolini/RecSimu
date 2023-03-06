package node

import (
	"log"
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
