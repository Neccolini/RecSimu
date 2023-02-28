package node

type Context struct {
	state State
	next  string
}

func NewContext() *Context {
	return &Context{
		state: NewConcreteStateIdle(),
	}
}

func (c *Context) SetState(obj State) {
	c.state = obj
}

func (c *Context) SetNext(ns string) {
	c.next = ns
}

func (c *Context) Handle() {
	c.state.handle(c)
}

func (c *Context) GetState() string {
	return c.state.getConcreteState()
}

func (c *Context) Reset() {
	c.state = NewConcreteStateIdle()
	c.next = ""
}
