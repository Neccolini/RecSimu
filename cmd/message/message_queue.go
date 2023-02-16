package message

import "fmt"

type MessageQueue struct {
	data []Message
	size int
}

func NewMessageQueue(cap int) *MessageQueue {
	if cap <= 0 {
		cap = 1
	}
	return &MessageQueue{data: make([]Message, 0, cap), size: 0}
}

func (mq *MessageQueue) Push(message Message) error {
	mq.data = append(mq.data, message)
	mq.size++
	return nil
}

func (mq *MessageQueue) Pop() error {
	if mq.size == 0 {
		return fmt.Errorf("Message Queue is empty")
	}
	mq.size--
	mq.data = mq.data[1:]
	return nil
}

func (mq *MessageQueue) Front() (Message, error) {
	if mq.size == 0 {
		return Message{}, fmt.Errorf("Message Queue is empty")
	}
	return mq.data[0], nil
}

func (mq *MessageQueue) IsEmpty() bool {
	return mq.size == 0
}

func (mq *MessageQueue) Size() int {
	return mq.size
}

func (mq *MessageQueue) Clear() error {
	mq.data = nil
	mq.size = 0
	return nil
}
