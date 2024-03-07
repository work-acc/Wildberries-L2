package pattern

import "fmt"

// Паттерн «строитель».

type Product struct {
	Part1 string
	Part2 string
}

type Builder interface {
	BuildPart1()
	BuildPart2()
	GetProduct() *Product
}

type ConcreteBuilder struct {
	product *Product
}

func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{
		product: &Product{},
	}
}

func (b *ConcreteBuilder) BuildPart1() {
	b.product.Part1 = "The first part is built"
}

func (b *ConcreteBuilder) BuildPart2() {
	b.product.Part2 = "The second part is built"
}

func (b *ConcreteBuilder) GetProduct() *Product {
	return b.product
}

type Director struct {
	builder Builder
}

func NewDirector(builder Builder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) Construct() *Product {
	d.builder.BuildPart1()
	d.builder.BuildPart2()
	return d.builder.GetProduct()
}

func main() {
	concreteBuilder := NewConcreteBuilder()
	director := NewDirector(concreteBuilder)
	product := director.Construct()

	fmt.Printf("The first part: %s\n", product.Part1)
	fmt.Printf("The second part: %s\n", product.Part2)
}
