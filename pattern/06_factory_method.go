package pattern

import "fmt"

// Паттерн «фабричный метод».

type Product interface {
	Use() string
}

type Factory interface {
	CreateProduct() Product
}

type ConcreteProductA struct{}

func (p *ConcreteProductA) Use() string {
	return "Product A"
}

type ConcreteFactoryA struct{}

func (f *ConcreteFactoryA) CreateProduct() Product {
	return &ConcreteProductA{}
}

type ConcreteProductB struct{}

func (p *ConcreteProductB) Use() string {
	return "Product B"
}

type ConcreteFactoryB struct{}

func (f *ConcreteFactoryB) CreateProduct() Product {
	return &ConcreteProductB{}
}

func main() {
	factoryA := &ConcreteFactoryA{}
	productA := factoryA.CreateProduct()
	fmt.Println(productA.Use())

	factoryB := &ConcreteFactoryB{}
	productB := factoryB.CreateProduct()
	fmt.Println(productB.Use())
}
