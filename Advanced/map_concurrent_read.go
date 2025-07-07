package main

import (
	"fmt"
	"sync"
	"time"
)

// 这个程序演示了多个 goroutine 并发地“只读”同一个 map 是安全的。
// 但同时，它也展示了只要引入一个“写”操作，程序就会立刻变得不安全并崩溃。
// 最后，使用 RWMutex 来保护 map 的并发读写

type SafeMap struct {
	m map[int]int
	sync.RWMutex
}

func (s *SafeMap) Get(key int) int {
	s.RLock()
	defer s.RUnlock()

	return s.m[key]
}

func (s *SafeMap) Set(key int, value int) {
	s.Lock()
	defer s.Unlock()

	s.m[key] = value
}

func map_concurrent_read() {
	// 1. 创建并完全初始化一个 map
	// 我们假设这个 map 在程序启动后就不会再被修改
	safeMap := &SafeMap{m: make(map[int]int)}
	for i := 0; i < 100; i++ {
		safeMap.Set(i, i*i)
	}

	var wg sync.WaitGroup
	fmt.Println("启动5个 goroutine 并发读取 map...")

	// 2. 启动多个 goroutine 进行并发读
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 每个 goroutine 读取1000次
			for j := 0; j < 1000; j++ {
				// 从 map 中读取一个值
				_ = safeMap.Get(j % 100) // 使用 % 确保 key 存在
				// 短暂休眠，增加并发冲突的机会
				time.Sleep(time.Microsecond)
			}
			fmt.Printf("Gorouine %d 读取完成。", id)
		}(i)
	}

	// ------------------- 危险的写入操作 -------------------
	go func() {
		// 在其他 goroutine 正在读取时，尝试写入
		for i := 0; i < 100; i++ {
			time.Sleep(2 * time.Millisecond) // 等待一会，确保读操作正在进行
			fmt.Println("!!! 尝试进行一次并发写入... !!!")
			safeMap.Set(101, 12345)
		}
	}()
	// ----------------------------------------------------

	// 等待所有“读取”的 goroutine 完成
	wg.Wait()

	fmt.Println("所有并发读取操作成功完成，程序没有 panic。")
	fmt.Println("这证明了纯并发读取是安全的。")
	fmt.Println("请尝试取消注释中的写入代码块，再次运行以观察 panic。")
}
