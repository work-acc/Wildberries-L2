package pattern

import "fmt"

// Паттерн «стратегия».

type Strategy interface {
	ExecuteStrategy()
}

type ConcreteStrategyA struct{}

func (s *ConcreteStrategyA) ExecuteStrategy() {
	fmt.Println("Strategy execution A")
}

type ConcreteStrategyB struct{}

func (s *ConcreteStrategyB) ExecuteStrategy() {
	fmt.Println("Strategy execution B")
}

type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *Context) Execute() {
	c.strategy.ExecuteStrategy()
}

func main() {
	context := &Context{}

	strategyA := &ConcreteStrategyA{}
	context.SetStrategy(strategyA)
	context.Execute()

	strategyB := &ConcreteStrategyB{}
	context.SetStrategy(strategyB)
	context.Execute()
}
