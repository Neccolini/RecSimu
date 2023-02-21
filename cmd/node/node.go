package node

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
	"github.com/Neccolini/RecSimu/cmd/routing"
)

type NodeInit struct {
	Id           int
	Type         string
	Instructions []instruction.Instruction
}

type Node struct {
	nodeId    int
	nodeType  string
	nodeAlive bool
	nodeState NodeState

	receiveMessages message.MessageQueue
	sendMessages    message.MessageQueue
	instructions    instruction.InstructionQueue

	SendingMessage   message.Message
	ReceivingMessage message.Message

	RoutingFunction routing.RoutingFunction
}

func NewNode(id int, nodeType string, instructions []instruction.Instruction) (*Node, error) {
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
	// 開始メッセージ生成
	packets, err := n.RoutingFunction.Init(nodeType)
	if err != nil {
		log.Fatalf("node initialization failed: %v", err)
	}
	for _, packet := range packets {
		n.sendMessages.Push(*message.NewMessage(id, true, packet))
	}

	return n, nil
}

func (n *Node) Id() int {
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

func (n *Node) processInstruction() error {
	i, err := n.instructions.Front()
	if err != nil {
		return err
	}
	n.instructions.Pop()

	packets, err := n.RoutingFunction.GenMessageFromI(i.Data)
	if err != nil {
		return err
	}

	for _, packet := range packets {
		m := *message.NewMessage(n.nodeId, true, packet)
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
	n.nodeState.RecieveStart(message.Cycles())
	return true
}

func (n *Node) receiveComplete() bool {
	// 受信中のメッセージが存在するが，ノードの状態は受信中じゃない -> 受信完了
	if !n.ReceivingMessage.IsEmpty() && !n.nodeState.IsReceiving() {
		packets, err := n.RoutingFunction.GenMessageFromM(n.ReceivingMessage.Data)
		if err != nil {
			return false // todo エラー処理必要？
		}
		for _, packet := range packets {
			m := *message.NewMessage(n.nodeId, true, packet)
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

	return nil
}

func (n *Node) Receive(m message.Message) {
	n.receiveMessages.Push(m)
}

func (n *Node) Wait() error {
	if !n.nodeState.IsSending() {
		return fmt.Errorf("Node %d is not sending a message", n.nodeId)
	}

	n.nodeState.Wait()
	n.nodeState.sending = State{false, 0} // send中止

	return nil
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
		packetInfo = ", packet from: " + fmt.Sprint(n.SendingMessage.Id())

		totalCycles := n.ReceivingMessage.Cycles()
		curCycles := totalCycles - n.nodeState.receiving.remaining
		packetInfo = fmt.Sprintf(", flit: %d/%d", curCycles, totalCycles)
	} else if n.nodeState.IsSending() {
		state = "sending"
		totalCycles := n.SendingMessage.Cycles()
		curCycles := totalCycles - n.nodeState.sending.remaining
		packetInfo = fmt.Sprintf(", flit: %d/%d", curCycles, totalCycles)
	}

	return fmt.Sprintf("id: %d, type: %s, condition: %s, state: %s%s", n.nodeId, n.nodeType, condition, state, packetInfo)
}
