package node

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
	state "github.com/Neccolini/RecSimu/cmd/node/state"
	"github.com/Neccolini/RecSimu/cmd/random"
	"github.com/Neccolini/RecSimu/cmd/routing"
)

const (
	waitRetriesMax = 8
	sendTryingMax  = 2
	idleMax        = 20
)

type NodeInit struct {
	Id           string
	Type         string
	Instructions []instruction.Instruction
}

type Node struct {
	nodeId       string
	nodeType     string
	nodeAlive    bool
	context      state.Context
	curCount     int
	curMax       int
	sendMessages message.MessageQueue
	instructions instruction.InstructionQueue

	SendingMessage   message.Message
	ReceivingMessage message.Message

	RoutingFunction routing.RoutingFunction

	waitRetries int
}

func NewNode(id string, nodeType string, instructions []instruction.Instruction) (*Node, error) {
	n := &Node{}
	n.nodeId = id
	n.nodeType = nodeType
	n.nodeAlive = true
	n.context = *state.NewContext()

	n.sendMessages = *message.NewMessageQueue(50)         // todo: the number of initial capacity
	n.instructions = *instruction.NewInstructionQueue(50) // todo: the number of initial capacity
	for _, inst := range instructions {
		n.instructions.Push(inst)
	}
	n.RoutingFunction = routing.NewRoutingFunction(n.nodeId, n.nodeType)
	n.Init()

	return n, nil
}

func (n *Node) Init() error {
	// 開始メッセージ生成
	packets := n.RoutingFunction.Init()

	for _, packet := range packets {
		n.sendMessages.Push(*message.NewMessage(n.nodeId, packet.ToId, packet.Data))
	}
	return nil
}

func (n *Node) Id() string {
	return n.nodeId
}

func (n *Node) Type() string {
	return n.nodeType
}

func (n *Node) State() string {
	return n.context.GetState()
}

func (n *Node) Alive() bool {
	return n.nodeAlive
}

func (n *Node) IsJoined() bool {
	return n.RoutingFunction.IsJoined()
}

func (n *Node) Reset() error {
	n.nodeAlive = false

	n.sendMessages.Clear()
	n.SendingMessage.Clear()
	n.ReceivingMessage.Clear()
	n.RoutingFunction.Reset()
	n.context.Reset()

	return nil
}

func (n *Node) SimulateCycle() {
	// 最後に状態を更新
	defer n.context.Handle()
	// 現在の状態に従い，context.nextを更新
	switch n.context.GetState() {
	case state.Idle:
		{
			n.curCount++
			// 送信途中のメッセージがある場合
			if !n.SendingMessage.IsEmpty() {
				n.curMax = rand.Intn(sendTryingMax) + 1
				n.curCount = 0

				n.context.SetNext(state.Sendtrying)
			} else if !n.sendMessages.IsEmpty() {
				n.SendingMessage, _ = n.sendMessages.Front()
				n.sendMessages.Pop()

				n.curMax = rand.Intn(sendTryingMax) + 1
				n.curCount = 0

				n.context.SetNext(state.Sendtrying)
			} else if n.curCount >= n.curMax {
				n.Init()
			}
		}
	case state.Sendtrying:
		{
			n.curCount++
			if n.curCount > n.curMax {
				log.Fatalf("sendtrying: n.curCount %d > n.curMax %d", n.curCount, n.curMax)
			}
			if n.curCount == n.curMax {
				// wait状態に移行
				if n.waitRetries < waitRetriesMax {
					n.waitRetries++
				}
				n.curMax = n.calcWaitCycle()
				n.curCount = 0
				n.context.SetNext(state.Waiting)
			}
		}
	case state.Sending:
		{
			// 送信に成功しているので待ち時間はリセット
			n.waitRetries = 0
			n.curCount++
			if n.curCount > n.curMax {
				log.Fatalf("sending: n.curCount %d > n.curMax %d", n.curCount, n.curMax)
			}
			if n.curCount == n.curMax {
				n.curMax = rand.Intn(idleMax) + idleMax
				n.curCount = 0

				n.SendingMessage.Clear()
				n.context.SetNext(state.Idle)
			}
		}
	case state.Receiving:
		{
			n.curCount++
			if n.curCount > n.curMax {
				log.Fatalf("receiving: n.curCount %d > n.curMax %d", n.curCount, n.curMax)
			}
			if n.curCount == n.curMax {
				n.curMax = rand.Intn(idleMax) + idleMax
				n.curCount = 0
				// process Message
				n.processReceivedMessage()

				n.ReceivingMessage.Clear()
				n.context.SetNext(state.Idle)
			}
		}
	case state.Waiting:
		{
			n.curCount++
			if n.curCount > n.curMax {
				log.Fatalf("waiting: n.curCount %d > n.curMax %d", n.curCount, n.curMax)
			}
			if n.curCount == n.curMax {
				n.curMax = sendTryingMax
				n.curCount = 0
				n.context.SetNext(state.Sendtrying)
			}
		}
	}
}

func (n *Node) processReceivedMessage() {
	// 受信したメッセージを読みこみ，
	// RFに投げて，
	// 受信メッセージが帰ってくればそれをキューにプッシュ
	packets := n.RoutingFunction.GenMessageFromM(n.ReceivingMessage.Data)
	for _, packet := range packets {
		n.sendMessages.Push(*message.NewMessage(n.nodeId, packet.ToId, packet.Data))
	}
}

func (n *Node) SetSending() {
	defer n.context.Handle()

	if n.context.GetState() != state.Sendtrying {
		log.Fatalf("the state of node %s is %s: cannot transit to Sending State", n.nodeId, n.context.GetState())
	}
	if n.SendingMessage.IsEmpty() {
		log.Fatal("the sending message is empty")
	}
	if n.SendingMessage.Id() != n.nodeId {
		log.Fatalf("the id in the message is incorrect %s %s", n.SendingMessage.Id(), n.nodeId)
	}
	n.curMax = n.SendingMessage.Cycles()
	n.curCount = 0

	// sendtrying -> sending
	n.context.SetNext(state.Sending)
}

func (n *Node) SetReceiving(msg message.Message) {
	defer n.context.Handle()

	if n.context.GetState() != state.Idle && n.context.GetState() != state.Waiting {
		log.Fatalf("the state of node %s is %s: cannot transit to Receiving State", n.nodeId, n.context.GetState())
	}
	if msg.IsEmpty() {
		log.Fatal("the receiving message is empty")
	}

	n.ReceivingMessage = msg

	n.curMax = msg.Cycles()
	n.curCount = 0

	// idle | waiting -> receiving
	n.context.SetNext(state.Receiving)
}

func (n *Node) calcWaitCycle() int {
	return random.RandomInt(1, int(math.Pow(2, float64(n.waitRetries))))
}

// sending
func (n *Node) IsSending() bool {
	return n.context.GetState() == state.Sending
}

// sendtrying
func (n *Node) IsSendTrying() bool {
	return n.context.GetState() == state.Sendtrying
}

// receiving
func (n *Node) IsReceiving() bool {
	return n.context.GetState() == state.Receiving
}

// idle
func (n *Node) IsIdle() bool {
	return n.context.GetState() == state.Idle
}

// waiting
func (n *Node) IsWaiting() bool {
	return n.context.GetState() == state.Waiting
}

func (n *Node) String() string {
	res := ""
	switch n.context.GetState() {
	case state.Idle:
		{
			res = fmt.Sprintf("%s %s %d/%d", n.nodeId, n.context.GetState(), n.curCount, n.curMax)
		}
	case state.Sendtrying, state.Sending, state.Waiting:
		{
			res = fmt.Sprintf("%s %s to %s %d/%d", n.nodeId, n.context.GetState(), n.SendingMessage.ToId(), n.curCount, n.curMax)
		}
	case state.Receiving:
		{
			res = fmt.Sprintf("%s %s from %s %d/%d", n.nodeId, n.context.GetState(), n.ReceivingMessage.Id(), n.curCount, n.curMax)
		}
	default:
		{
			res = ""
		}
	}
	if n.IsJoined() {
		res += " joined"
	}
	if n.RoutingFunction.ParentId() != "" {
		res += " " + n.RoutingFunction.ParentId()
	}
	return res
}
