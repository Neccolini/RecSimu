package node

import "fmt"

func (n *Node) Remove() error {
	if !n.Alive() {
		return fmt.Errorf("node %s is already removed", n.Id())
	}

	n.Reset()

	return nil
}

func (n *Node) Rejoin() error {
	if n.Alive() {
		return fmt.Errorf("node %s is not removed", n.Id())
	}

	n.nodeAlive = true

	return nil
}
