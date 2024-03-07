package pattern

import "fmt"

// Паттерн «цепочка вызовов».

type Handler interface {
	SetNext(handler Handler)
	HandleRequest(request int)
}

type ConcreteHandlerA struct {
	next Handler
}

func (h *ConcreteHandlerA) SetNext(handler Handler) {
	h.next = handler
}

func (h *ConcreteHandlerA) HandleRequest(request int) {
	if request < 10 {
		fmt.Println("ConcreteHandlerA processes the request")
	} else if h.next != nil {
		h.next.HandleRequest(request)
	}
}

type ConcreteHandlerB struct {
	next Handler
}

func (h *ConcreteHandlerB) SetNext(handler Handler) {
	h.next = handler
}

func (h *ConcreteHandlerB) HandleRequest(request int) {
	if request >= 10 && request < 20 {
		fmt.Println("ConcreteHandlerB processes the request")
	} else if h.next != nil {
		h.next.HandleRequest(request)
	}
}

func main() {
	handlerA := &ConcreteHandlerA{}
	handlerB := &ConcreteHandlerB{}

	handlerA.SetNext(handlerB)

	requests := []int{5, 12, 15, 25}

	for _, req := range requests {
		handlerA.HandleRequest(req)
	}
}
