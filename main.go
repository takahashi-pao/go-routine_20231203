package main

import (
	"fmt"
	"time"
)

// メイン関数
func main() {
	start := time.Now()
	done := make(chan interface{})
	defer close(done)

	intCh := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intCh, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
	fmt.Println("took: ", time.Since(start))
}

// intチャネル生成
func generator(done <-chan interface{}, integers ...int) <-chan int {
	intCh := make(chan int, len(integers))
	go func() {
		defer close(intCh)
		for _, i := range integers {
			select {
			case <-done: // doneチャネルが閉じた場合
				return
			case intCh <- i:
			}
		}
	}()
	return intCh
}

// 掛け算
func multiply(
	done <-chan interface{},
	intCh <-chan int,
	multiplier int,
) <-chan int {
	multipliedCh := make(chan int)
	go func() {
		defer close(multipliedCh)
		for i := range intCh {
			select {
			case <-done: // doneチャネルが閉じた場合
				return
			case multipliedCh <- i * multiplier:
			}
		}
	}()
	return multipliedCh
}

// 足し算
func add(
	done <-chan interface{},
	intCh <-chan int,
	addtive int,
) <-chan int {
	addCh := make(chan int)
	go func() {
		defer close(addCh)
		for i := range intCh {
			select {
			case <-done: // doneチャネルが閉じた場合
				return
			case addCh <- i + addtive:
			}
		}
	}()
	return addCh
}
