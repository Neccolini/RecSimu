package routing

import (
	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/message"
)

type RoutingFunction interface {
	Init(id int, nodeType string) ([]message.Message, error)
	GenMessageFromM(m message.Message) (message.Message, error)
	GenMessageFromI(i instruction.Instruction) (message.Message, error)
}
