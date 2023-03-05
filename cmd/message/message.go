package message

import "github.com/Neccolini/RecSimu/cmd/routing"

const bytePerFlit = routing.BytePerFlit

type Packet []byte
type Message struct {
	MessageId string
	fromId    string
	toId      string
	ready     bool
	cycles    int
	Data      Packet
}

func NewMessage(mId string, fromId string, toId string, data []byte) *Message {
	cycles := (len(data) + bytePerFlit) / bytePerFlit
	return &Message{mId, fromId, toId, true, cycles, data}
}

func (m *Message) IsEmpty() bool {
	return len(m.Data) == 0
}

func (m *Message) Id() string {
	return m.fromId
}

func (m *Message) ToId() string {
	return m.toId
}
func (m *Message) Cycles() int {
	return m.cycles
}

func (m *Message) IsReady() bool {
	return m.ready
}

func (m *Message) Ready() {
	m.ready = true
}

func (m *Message) Clear() error {
	m.Data = []byte{}
	return nil
}
