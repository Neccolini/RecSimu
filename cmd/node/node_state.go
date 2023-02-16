package node

type NodeState struct {
	state State
}

type State string

const (
	Sending State = "Sending"
	Idle    State = "Idle"
)

func (ns *NodeState) transit() {
	switch ns.state {
	case Sending:
		{
			ns.state = Idle
		}
	case Idle:
		{
			ns.state = Sending
		}
	default:
	}
}
