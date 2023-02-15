package run

type Message struct {
	Id      int
	Data    string
}

type MessageBuffer struct {
	Messages []Message
}

func (m *Message) IsEmpty() bool {
	return len(m.Data) == 0
}
