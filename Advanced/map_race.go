package main

import (
	"fmt"
	"sync"
	"time"
)

// 这个程序演示了在没有同步机制的情况下，
// 多个 goroutine 并发地对同一个 map 进行写操作，
// 这会大概率导致一个 "fatal error: concurrent map writes" 的 panic。

func MapWriteRace() {
	// 创建一个普通的 map
	unprotectedMap := make(map[string]int)

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup

	fmt.Println("启动10个 goroutine 并发写入 map...")

	// 启动10个 goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 每个 goroutine 尝试写入100次
			for j := 0; j < 100; j++ {
				// 生成一个唯一的 key
				key := fmt.Sprintf("goroutine_%d_key_%d", id, j)
				// 对 map 进行写操作
				unprotectedMap[key] = j

				// 短暂休眠，增加不同 goroutine 之间发生冲突的可能性
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	// 等待所有 goroutine 结束
	wg.Wait()

	// 如果程序没有 panic，我们会打印 map 的最终大小。
	// 但实际上，这个程序几乎不可能正常执行到这里。
	fmt.Println("所有 goroutine 执行完毕。")
	fmt.Printf("Map 的最终大小: %d\n", len(unprotectedMap))
}
