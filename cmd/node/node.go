package node

import (
	"fmt"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
	"github.com/Neccolini/RecSimu/cmd/routing"
)

type NodeInit struct {
	Id           string
	Type         string
	Instructions []instruction.Instruction
}

type Node struct {
	nodeId    string
	nodeType  string
	nodeAlive bool
	joined    bool
	nodeState NodeState

	receiveMessages message.MessageQueue
	sendMessages    message.MessageQueue
	instructions    instruction.InstructionQueue

	SendingMessage      message.Message
	ReceivingMessage    message.Message
	CommunicatingNodeId string
	RoutingFunction     routing.RoutingFunction

	waitRetries int
}

func NewNode(id string, nodeType string, instructions []instruction.Instruction) (*Node, error) {
	n := &Node{}
	n.nodeId = id
	n.nodeType = nodeType
	n.nodeAlive = true
	n.nodeState = NodeState{}
	n.receiveMessages = *message.NewMessageQueue(100)      // todo: the number of initial capacity
	n.sendMessages = *message.NewMessageQueue(100)         // todo: the number of initial capacity
	n.instructions = *instruction.NewInstructionQueue(100) // todo: the number of initial capacity
	for _, inst := range instructions {
		n.instructions.Push(inst)
	}
	n.RoutingFunction = &routing.RF{}
	n.Init()

	return n, nil
}

func (n *Node) Init() error {
	// 開始メッセージ生成
	packets, distId := n.RoutingFunction.Init(n.nodeId, n.nodeType)

	for _, packet := range packets {
		n.sendMessages.Push(*message.NewMessage(n.nodeId, distId, packet))
	}
	return nil
}

func (n *Node) Id() string {
	return n.nodeId
}

func (n *Node) Type() string {
	return n.nodeType
}

func (n *Node) State() *NodeState {
	return &n.nodeState
}

func (n *Node) Alive() bool {
	return n.nodeAlive
}

func (n *Node) IsJoined() bool {
	return n.RoutingFunction.IsJoined()
}

func (n *Node) Reset() error {
	n.nodeAlive = false
	n.joined = false
	n.receiveMessages.Clear()
	n.sendMessages.Clear()
	n.SendingMessage.Clear()
	n.ReceivingMessage.Clear()
	n.RoutingFunction.Reset()
	n.nodeState.ResetAll()
	return nil
}

func (n *Node) processInstruction() error {
	i, err := n.instructions.Front()
	if err != nil {
		return err
	}
	n.instructions.Pop()

	packets, distId := n.RoutingFunction.GenMessageFromI(i.Data)

	for _, packet := range packets {
		m := *message.NewMessage(n.nodeId, distId, packet)
		n.sendMessages.Push(m)
	}
	return nil
}

func (n *Node) sendProcess() bool {
	if n.sendMessages.IsEmpty() {
		return false
	}
	message, err := n.sendMessages.Front()
	if err != nil {
		return false
	}
	// n.sendMessages.Pop()

	n.SendingMessage = message
	n.nodeState.SendStart(message.Cycles())
	return true
}

func (n *Node) receiveProcess() bool {
	if n.receiveMessages.IsEmpty() {
		return false
	}

	message, err := n.receiveMessages.Front()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	n.receiveMessages.Pop()
	n.ReceivingMessage = message
	n.CommunicatingNodeId = n.ReceivingMessage.Id()
	n.nodeState.RecieveStart(message.Cycles())
	return true
}

func (n *Node) receiveComplete() bool {
	// 受信中のメッセージが存在するが，ノードの状態は受信中じゃない -> 受信完了
	if !n.ReceivingMessage.IsEmpty() && !n.nodeState.IsReceiving() {
		packets, distId := n.RoutingFunction.GenMessageFromM(n.ReceivingMessage.Data)

		for _, packet := range packets {
			m := *message.NewMessage(n.nodeId, distId, packet)
			n.sendMessages.Push(m)
		}
		// 現在のメッセージを破棄
		n.ReceivingMessage.Clear()
		return true
	}
	return false
}

func (n *Node) CycleSend() bool {
	// 送信中でない場合
	if n.nodeState.IsIdle() {
		return n.sendProcess()
	}
	return false
}

func (n *Node) CycleReceive() bool {
	// 受信中でない場合
	if n.nodeState.IsIdle() {
		return n.receiveProcess()
	}
	return false
}

func (n *Node) SimulateCycle() error {
	// 送信処理がうまくいったかどうか
	// 送信するメッセージが存在し，待機中でない
	if !n.SendingMessage.IsEmpty() && n.nodeState.IsSending() {
		// この場合送信がうまくいっているので，queueから削除
		if err := n.sendMessages.Pop(); err != nil {
			return err
		}
		// communicating nodeを設定
		n.CommunicatingNodeId = n.SendingMessage.ToId()
		n.waitRetries = 0
	}
	n.SendingMessage.Clear() // 送信中メッセージは削除

	// 状態を進める
	if err := n.nodeState.Next(); err != nil {
		return err
	}

	// 受信したメッセージを処理する
	n.receiveComplete()
	// 命令を実行
	n.processInstruction()

	n.endCommunication()
	return nil
}

func (n *Node) Receive(m message.Message) {
	n.receiveMessages.Push(m)
}

func (n *Node) Wait() error {
	if !n.nodeState.IsSending() {
		return fmt.Errorf("Node %s is not sending a message", n.nodeId)
	}
	if n.waitRetries >= 10 {
		// 送信失敗...
		n.nodeState.ResetAll()
		return nil
	}
	n.waitRetries++
	n.nodeState.Wait(n.waitRetries)
	n.nodeState.sending = State{false, 0} // send中止

	return nil
}

func (n *Node) endCommunication() {
	if n.CommunicatingNodeId == "" {
		n.nodeState.ResetCommunication()
	} else if n.nodeState.receiving.remaining == 0 &&
		n.nodeState.sending.remaining == 0 {
		n.CommunicatingNodeId = ""
	}
}

func (n *Node) String() string {

	condition := "inactive"
	if n.nodeAlive {
		condition = "active"
	}

	state := "idle"
	packetInfo := ""
	if n.nodeState.IsReceiving() {
		state = "receiving"
		packetInfo = ", packet from: " + fmt.Sprint(n.ReceivingMessage.Id())

		totalCycles := n.ReceivingMessage.Cycles()
		curCycles := totalCycles - n.nodeState.receiving.remaining
		packetInfo += fmt.Sprintf(", flit: %d/%d", curCycles, totalCycles)
	} else if n.nodeState.IsSending() {
		state = "sending"
		packetInfo = ", packet to: " + fmt.Sprint(n.SendingMessage.ToId())

		totalCycles := n.SendingMessage.Cycles()
		curCycles := totalCycles - n.nodeState.sending.remaining
		packetInfo += fmt.Sprintf(", flit: %d/%d", curCycles, totalCycles)
	} else if n.nodeState.IsWaiting() {
		state = "waiting"
		curCycles := n.nodeState.waiting.remaining
		packetInfo += fmt.Sprintf(",remaining wait cycle: %d", curCycles)
	}

	return fmt.Sprintf("id: %s, type: %s, condition: %s, state: %s%s", n.nodeId, n.nodeType, condition, state, packetInfo)
}
