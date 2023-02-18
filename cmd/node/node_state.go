package node

import "fmt"

type NodeState struct {
	sending   State
	receiving State
	// 状態が増えたときはここに追加する
}

type State struct {
	state     bool
	remaining int // number of cycles
}

func (n *NodeState) Next() error {
	// receiving
	if n.receiving.state {
		n.receiving.remaining--
		if n.receiving.remaining < 0 {
			n.receiving.state = false
		}
	} else if n.sending.state {
		n.sending.remaining--
		if n.sending.remaining < 0 {
			n.sending.state = false
		}
	}
	return nil
}

func (n *NodeState) IsSending() bool {
	return n.sending.state
}

func (n *NodeState) IsReceiving() bool {
	return n.receiving.state
}

func (n *NodeState) IsIdle() bool {
	return !n.IsSending() && !n.IsReceiving()
}

func (n *NodeState) SendStart(cycles int) error {
	if n.sending.state {
		return fmt.Errorf("sending another message")
	}
	n.sending.state = true
	n.sending.remaining = cycles
	return nil
}

func (n *NodeState) RecieveStart(cycles int) error {
	if n.receiving.state {
		return fmt.Errorf("receiving another message")
	}
	n.receiving.state = true
	n.receiving.remaining = cycles
	return nil
}
