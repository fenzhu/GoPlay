package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

func main() {
	PrintWorker()
	SliceMemory()
	MemoryLeak()
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

func SliceMemory() {
	var a = make([]int, 8)
	fmt.Printf("len %v cap %v\n", len(a), cap(a))
	a = append(a, 1)
	fmt.Printf("len %v cap %v\n", len(a), cap(a))
	a = make([]int, 1024)
	fmt.Printf("len %v cap %v\n", len(a), cap(a))
	a = append(a, 1)
	fmt.Printf("len %v cap %v\n", len(a), cap(a))

}

func MemoryLeak() {
	fmt.Printf("Goroutines at start: %d\n", runtime.NumGoroutine())
	ch := make(chan int, 1)

	// 1. go routine
	// 该go routine阻塞在一个不会有数据写入的channel上
	// 泄漏：go routine的栈内存、其引用的堆内存不会被回收
	// 思路：确保每个go routine都有明确的退出路径
	// 最佳实践： 使用context管理go routine的生命周期
	go func(ctx context.Context) {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-ctx.Done():
			fmt.Println("context done")
		}
	}(context.Background())
	fmt.Printf("Goroutines at end: %d\n", runtime.NumGoroutine())

	// 2. 切片
	big := make([]byte, 1024*1024)
	small := big[:10]
	// 只引用了big的10个字节，但是它指向原始的底层数组，所以big的1MB内存不会被回收
	fmt.Println(small)

	// 3. time.Ticker

	// 4. 全局变量不断添加数据，不清理

	// 5. 闭包
	large := make([]byte, 1024*1024)
	closure := func() {
		_ = large
	}
	// closure存活，那么闭包引用的1MB内存就不会被释放
	// 闭包不要无意中捕获大的、不必要的变量，如果只是需要变量中的一部分数据，可以将数据作为参数传递
	_ = closure

	// 6. 排查内存泄漏
	// 使用/pprof/heap
	// 使用/pprof/goroutine
}
