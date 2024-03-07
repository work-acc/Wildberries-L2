package pattern

import "fmt"

// Паттерн «комманда».

type Command interface {
	Execute()
}

type ConcreteCommandA struct {
	receiver *Receiver
}

func (c *ConcreteCommandA) Execute() {
	c.receiver.ActionA()
}

type ConcreteCommandB struct {
	receiver *Receiver
}

func (c *ConcreteCommandB) Execute() {
	c.receiver.ActionB()
}

type Receiver struct{}

func (r *Receiver) ActionA() {
	fmt.Println("Receiver: Performing an action A")
}

func (r *Receiver) ActionB() {
	fmt.Println("Receiver: Performing an action B")
}

type Invoker struct {
	command Command
}

func (i *Invoker) SetCommand(command Command) {
	i.command = command
}

func (i *Invoker) ExecuteCommand() {
	i.command.Execute()
}

func main() {
	receiver := &Receiver{}

	commandA := &ConcreteCommandA{receiver}
	commandB := &ConcreteCommandB{receiver}

	invoker := &Invoker{}

	invoker.SetCommand(commandA)
	invoker.ExecuteCommand()

	invoker.SetCommand(commandB)
	invoker.ExecuteCommand()
}
