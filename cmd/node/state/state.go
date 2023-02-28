package node

const (
	Sending    string = "Sending"
	Sendtrying string = "Sendtrying"
	Receiving  string = "Receiving"
	Idle       string = "Idle"
	Waiting    string = "Waiting"
)

type State interface {
	handle(context *Context)
	getConcreteState() string
}

type concreteState struct {
	state string
}

func (c *concreteState) getConcreteState() string {
	return c.state
}

// 送信中
type ConcreteStateSending struct {
	*concreteState
}

func NewConcreteStateSending() *ConcreteStateSending {
	return &ConcreteStateSending{
		concreteState: &concreteState{
			state: Sending,
		},
	}
}

func (c *ConcreteStateSending) handle(context *Context) {
	switch context.next {
	case Idle:
		{
			context.SetState(NewConcreteStateIdle())
		}
	}
}

// 受信中
type ConcreteStateReceiving struct {
	*concreteState
}

func NewConcreteStateReceiving() *ConcreteStateReceiving {
	return &ConcreteStateReceiving{
		concreteState: &concreteState{
			state: Receiving,
		},
	}
}

func (c *ConcreteStateReceiving) handle(context *Context) {
	switch context.next {
	case Idle:
		{
			context.SetState(NewConcreteStateIdle())
		}
	}
}

// 送信前待機（ack待ち）
type ConcreteStateSendtrying struct {
	*concreteState
}

func NewConcreteStateSendtrying() *ConcreteStateSendtrying {
	return &ConcreteStateSendtrying{
		concreteState: &concreteState{
			state: Sendtrying,
		},
	}
}

func (c *ConcreteStateSendtrying) handle(context *Context) {
	switch context.next {
	case Sending:
		{
			context.SetState(NewConcreteStateSending())
		}
	case Waiting:
		{
			context.SetState(NewConcreteStateWaiting())
		}
	}
}

// 何もしていない
type ConcreteStateIdle struct {
	*concreteState
}

func NewConcreteStateIdle() *ConcreteStateIdle {
	return &ConcreteStateIdle{
		concreteState: &concreteState{
			state: Idle,
		},
	}
}

func (c *ConcreteStateIdle) handle(context *Context) {
	switch context.next {
	case Sendtrying:
		{
			context.SetState(NewConcreteStateSendtrying())
		}
	case Receiving:
		{
			context.SetState(NewConcreteStateReceiving())
		}
	}
}

// 送信に失敗した場合の待機状態
type ConcreteStateWaiting struct {
	*concreteState
}

func NewConcreteStateWaiting() *ConcreteStateWaiting {
	return &ConcreteStateWaiting{
		concreteState: &concreteState{
			state: Waiting,
		},
	}
}

func (c *ConcreteStateWaiting) handle(context *Context) {
	switch context.next {
	case Sendtrying:
		{
			context.SetState(NewConcreteStateSendtrying())
		}
	case Receiving:
		{
			context.SetState(NewConcreteStateReceiving())
		}
	}
}
