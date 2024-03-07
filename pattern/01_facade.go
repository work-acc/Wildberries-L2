package pattern

import "fmt"

// Паттерн «фасад».

type SubsystemA struct {
}

func (s *SubsystemA) OperationA() {
	fmt.Println("Subsystem A: Operation A")
}

type SubsystemB struct {
}

func (s *SubsystemB) OperationB() {
	fmt.Println("Subsystem B: Operation B")
}

type Facade struct {
	subsystemA *SubsystemA
	subsystemB *SubsystemB
}

func NewFacade() *Facade {
	return &Facade{
		subsystemA: &SubsystemA{},
		subsystemB: &SubsystemB{},
	}
}

func (f *Facade) Operation() {
	fmt.Println("Facade: Operation")
	f.subsystemA.OperationA()
	f.subsystemB.OperationB()
}

func main() {
	facade := NewFacade()
	facade.Operation()
}
