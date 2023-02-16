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
	nodeId           int
	nodeType         string
	nodeAlive        bool
	nodeState        NodeState
	receivedMessages message.MessageQueue
	sendingMessages  message.MessageQueue
	instructions     instruction.InstructionQueue
	SendMessage      message.Message

	RoutingFunction routing.RoutingFunction
}

func NewNode(id int, nodeType string, instructions []instruction.Instruction) (*Node, error) {
	n := &Node{}
	n.nodeId = id
	n.nodeType = nodeType
	n.nodeAlive = true
	n.nodeState = NodeState{Idle}
	n.receivedMessages = *message.NewMessageQueue(100)     // todo: the number of initial capacity
	n.sendingMessages = *message.NewMessageQueue(100)      // todo: the number of initial capacity
	n.instructions = *instruction.NewInstructionQueue(100) // todo: the number of initial capacity
	for _, inst := range instructions {
		n.instructions.Push(inst)
	}

	n.RoutingFunction = &routing.RF{}
	packets, err := n.RoutingFunction.Init(nodeType)
	if err != nil {
		log.Fatalf("node initialization failed: %v", err)
	}
	for _, packet := range packets {
		n.sendingMessages.Push(*message.NewMessage(id, true, packet))
	}

	return n, nil
}
func (n *Node) Id() int {
	return n.nodeId
}

func (n *Node) Type() string {
	return n.nodeType
}

func (n *Node) State() NodeState {
	return n.nodeState
}

func (n *Node) Alive() bool {
	return n.nodeAlive
}

func (n *Node) processReceivedMessage() error {
	m, err := n.receivedMessages.Front()
	if err != nil {
		return err
	}
	n.receivedMessages.Pop()

	if m.IsValid() {
		// todo ここでRFが受信メッセージを読んで送信メッセージを生成する
		packets, err := n.RoutingFunction.GenMessageFromM(m.Data)
		if err != nil {
			return err
		}
		if len(packets) > 0 {
			for _, packet := range packets {
				m := *message.NewMessage(n.nodeId, true, packet)
				n.sendingMessages.Push(m)
			}
		}
	}
	return nil
}
func (n *Node) processInstruction() error {
	// todo processReceivedMessgeのような構成にする
	// 122-139行目を参考に
	return nil
}

func (n *Node) send() error {
	// sendingMessagesからとってきて，送る
	if n.sendingMessages.IsEmpty() {
		return fmt.Errorf("Send message queue is empty")
	}
	message, err := n.sendingMessages.Front()
	if err != nil {
		return err
	}

	n.sendingMessages.Pop()
	n.SendMessage = message
	return nil
}

func (n *Node) SimulateCycle() error {
	switch n.nodeState.state {
	case Sending:
		{
			n.nodeState.transit()
		}
	case Idle:
		{
			if !n.instructions.IsEmpty() {
				// todo タスクの処理
				instruction, err := n.instructions.Front()
				if err != nil {
					return err
				}
				packets, err := n.RoutingFunction.GenMessageFromI(instruction.Data)
				if err != nil {
					return err
				}
				if len(packets) > 0 {
					for _, packet := range packets {
						m := *message.NewMessage(n.nodeId, true, packet)
						n.sendingMessages.Push(m)
					}
					n.nodeState.transit()
				}
			} else if !n.sendingMessages.IsEmpty() {
				// 送信
				if err := n.send(); err != nil {
					return err
				}
				n.nodeState.transit()
			} else if !n.receivedMessages.IsEmpty() {
				// 受信メッセージの処理
				if err := n.processReceivedMessage(); err != nil {
					return err
				}
				n.nodeState.transit()
			} else {
				// 何もしない
			}
		}
	default:
	}
	return nil
}

func (n *Node) Receive(message message.Message) error {
	n.receivedMessages.Push(message)
	return nil
}
