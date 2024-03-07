Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Выведутся все цифры с 1 до 8, а затем будет бесконечный вывод программы нулями

Все дело в функции merge, которая не завершается. Из-за этого программа не
завершится. После чтения всех значений из каналов a и b в цикле
for v := range c, программа будет ожидать завершения горутины внутри merge,
чего не произойдет. И она будет выводить дефолтные значения int - 0. Чтобы
исправить это, необходимо проверять, закрылись ли каналы a и b.

```
