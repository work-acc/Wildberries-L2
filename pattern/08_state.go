package pattern

import "fmt"

// Паттерн «состояние».

type State interface {
	Handle()
}

type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle() {
	fmt.Println("Processing in the state A")
}

type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle() {
	fmt.Println("Processing in the state B")
}

type Context struct {
	state State
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Request() {
	c.state.Handle()
}

func main() {
	context := &Context{}

	stateA := &ConcreteStateA{}
	context.SetState(stateA)
	context.Request()

	stateB := &ConcreteStateB{}
	context.SetState(stateB)
	context.Request()
}
