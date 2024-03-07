package main

import (
	"fmt"
	"time"
)

// Реализовать функцию, которая будет объединять один или более done-каналов в single-канал, если один из его составляющих каналов закроется.
// Очевидным вариантом решения могло бы стать выражение при использованием select, которое бы реализовывало эту связь,
// однако иногда неизвестно общее число done-каналов, с которыми вы работаете в рантайме.
// В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or-каналов, реализовывала бы весь функционал.

// Определение функции:
// var or func(channels ...<- chan interface{}) <- chan interface{}
// Пример использования функции:
// sig := func(after time.Duration) <- chan interface{} {
// 	c := make(chan interface{})
// 	go func() {
// 		defer close(c)
// 		time.Sleep(after)
// }()
// return c
// }

// start := time.Now()
// <-or (
// 	sig(2*time.Hour),
// 	sig(5*time.Minute),
// 	sig(1*time.Second),
// 	sig(1*time.Hour),
// 	sig(1*time.Minute),
// )

// fmt.Printf(“fone after %v”, time.Since(start))

func main() {
	single := func(after time.Duration) <-chan interface{} {
		channel := make(chan interface{})
		go func() {
			defer close(channel)
			time.Sleep(after)
		}()
		return channel
	}

	start := time.Now()
	<-channelJoin(
		single(2*time.Hour),
		single(5*time.Minute),
		single(1*time.Second),
		single(1*time.Hour),
		single(1*time.Minute),
	)

	fmt.Printf("Closed after: %v\n", time.Since(start))
}

func channelJoin(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		channel := make(chan interface{})
		close(channel)
		return channel
	case 1:
		return channels[0]
	}

	combChannel := make(chan interface{})
	go func() {
		defer close(combChannel)

		select {
		case <-channels[0]:
			return
		case <-channels[1]:
			return
		case <-channelJoin(channels[2:]...):
			return
		}
	}()

	return combChannel
}
