Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error

В функции test создается блок кода, внутри которого "что-то делается".
После этого блока возвращается nil, что является значением типа указателя на customError.

В функции main объявляется переменная err типа error. Затем вызывается функция test(),
и результат ее выполнения присваивается переменной err.
Но err не будет равна nil, т.к. err будет хранить указатель на тип customError,
даже если внутри функции test возвращается nil. Поэтому выведется error.

```
