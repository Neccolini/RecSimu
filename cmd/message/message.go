package message

const bytePerFlit = 8 // todo

type Packet []byte
type Message struct {
	id     int
	valid  bool
	ready bool
	cycles int
	Data   Packet
}

func NewMessage(id int, valid bool, data []byte) *Message {
	cycles := (len(data) + bytePerFlit) / bytePerFlit
	return &Message{id, valid, false, cycles, data}
}

func (m *Message) IsEmpty() bool {
	return len(m.Data) == 0
}

func (m *Message) Id() int {
	return m.id
}

func (m *Message) Cycles() int {
	return m.cycles
}

func (m *Message) IsValid() bool {
	return m.valid
}

func (m *Message) Invalidate() error {
	m.valid = false
	return nil
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
