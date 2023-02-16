package message

type Message struct {
	id    int
	valid bool
	data  []byte
}

func NewMessage(id int, valid bool, data []byte) *Message {
	return &Message{id, valid, data}
}

func (m *Message) IsEmpty() bool {
	return len(m.data) == 0
}

func (m *Message) Id() int {
	return m.id
}

func (m *Message) IsValid() bool {
	return m.valid
}

func (m *Message) Invalidate() error {
	m.valid = false
	return nil
}

func (m *Message) Clear() error {
	m.data = []byte{}
	return nil
}
