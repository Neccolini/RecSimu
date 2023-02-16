package instruction

import "fmt"

type InstructionQueue struct {
	data []Instruction
	size int
}

func NewInstructionQueue(cap int) *InstructionQueue {
	if cap <= 0 {
		cap = 1
	}
	return &InstructionQueue{data: make([]Instruction, 0, cap), size: 0}
}

func (iq *InstructionQueue) Push(i Instruction) error {
	iq.data = append(iq.data, i)
	iq.size++
	return nil
}

func (iq *InstructionQueue) Pop() error {
	if iq.size == 0 {
		return fmt.Errorf("Instruction Queue is empty")
	}
	iq.size--
	iq.data = iq.data[1:]
	return nil
}

func (iq *InstructionQueue) Front() (Instruction, error) {
	if iq.size == 0 {
		return Instruction{}, fmt.Errorf("Instruction Queue is Empty")
	}
	return iq.data[0], nil
}

func (iq *InstructionQueue) IsEmpty() bool {
	return iq.size == 0
}

func (iq *InstructionQueue) Size() int {
	return iq.size
}

func (iq *InstructionQueue) Clear() error {
	iq.data = nil
	iq.size = 0
	return nil
}
