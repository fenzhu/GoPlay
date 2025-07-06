package main

import (
	"fmt"
	"sync"
)

func main() {
	PrintWorker()
}

func PrintWorker() {
	const limit = 100

	// WaitGroup 用来等待两个协程执行结束
	var wg sync.WaitGroup
	wg.Add(2)

	// 创建两个 channel 用于两个协程间通信
	chOdd := make(chan struct{}, 1) // 奇数协程
	chEven := make(chan struct{})   // 偶数协程

	// 打印奇数的协程
	go func() {
		defer wg.Done()
		for num := 1; num <= limit; num += 2 {
			<-chOdd // 等待轮到自己
			fmt.Printf("%d\n", num)
			chEven <- struct{}{} // 通知偶数协程
		}
	}()

	// 打印偶数的协程
	go func() {
		defer wg.Done()
		for num := 2; num <= limit; num += 2 {
			<-chEven // 等待轮到自己
			fmt.Printf("%d\n", num)
			chOdd <- struct{}{} // 通知奇数协程
		}
	}()

	fmt.Println("start print")
	chOdd <- struct{}{} // 启动第一个协程

	wg.Wait() // 等待所有协程执行完毕
	fmt.Println("done")
}
