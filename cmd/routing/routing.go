package routing

import (
	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
)

type RF struct {
	id int
	nodeType string
	table map[int][]int
}

func (r *RF) Init(id int, nodeType string) ([]message.Message, error) {
	r.id = id
	r.nodeType = nodeType
	
	return nil, nil
}
func (r *RF) GenMessageFromM(received message.Message) (message.Message, error) {
	m := *message.NewMessage(r.id, false, []byte{})
	return m, nil
}

func (r *RF) GenMessageFromI(i instruction.Instruction) (message.Message, error) {
	return message.Message{}, nil
}
