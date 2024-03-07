package pattern

import "fmt"

// Паттерн «посетитель».

type Element interface {
	Accept(visitor Visitor)
}

type ConcreteElementA struct{}

func (c *ConcreteElementA) Accept(visitor Visitor) {
	visitor.VisitConcreteElementA(c)
}

type ConcreteElementB struct{}

func (c *ConcreteElementB) Accept(visitor Visitor) {
	visitor.VisitConcreteElementB(c)
}

type Visitor interface {
	VisitConcreteElementA(element *ConcreteElementA)
	VisitConcreteElementB(element *ConcreteElementB)
}

type ConcreteVisitor1 struct{}

func (c *ConcreteVisitor1) VisitConcreteElementA(element *ConcreteElementA) {
	fmt.Println("The first visitor visited ConcreteElementA")
}

func (c *ConcreteVisitor1) VisitConcreteElementB(element *ConcreteElementB) {
	fmt.Println("The first visitor visited ConcreteElementB")
}

type ConcreteVisitor2 struct{}

func (c *ConcreteVisitor2) VisitConcreteElementA(element *ConcreteElementA) {
	fmt.Println("The second visitor visitedConcreteElementA")
}

func (c *ConcreteVisitor2) VisitConcreteElementB(element *ConcreteElementB) {
	fmt.Println("The second visitor visited ConcreteElementB")
}

type ObjectStructure struct {
	elements []Element
}

func (o *ObjectStructure) Attach(element Element) {
	o.elements = append(o.elements, element)
}

func (o *ObjectStructure) Accept(visitor Visitor) {
	for _, e := range o.elements {
		e.Accept(visitor)
	}
}

func main() {
	objectStructure := ObjectStructure{}

	objectStructure.Attach(&ConcreteElementA{})
	objectStructure.Attach(&ConcreteElementB{})

	visitor1 := &ConcreteVisitor1{}
	objectStructure.Accept(visitor1)

	visitor2 := &ConcreteVisitor2{}
	objectStructure.Accept(visitor2)
}
