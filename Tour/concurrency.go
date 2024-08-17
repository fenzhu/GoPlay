package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/tour/tree"
)

func concurrency() {
	go say("world")
	say("hello")

	s := []int{1, 2, 3, 4, 5}
	c := make(chan int)
	half := len(s) / 2
	go sum(s[:half], c)
	go sum(s[half:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	//overfill the channel, wait for the receiver, main thread sleep
	//deadlock
	// ch <- 3
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	fc := make(chan int, 10)
	go co_fibonacci(cap(fc), fc)
	for i := range fc {
		fmt.Printf("%d\n", i)
	}

	cs := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-cs)
		}
		quit <- 0
	}()
	select_fibonacci(cs, quit)

	select_default()

	ct := make(chan int)
	go Walk(tree.New(1), ct)
	for i := 0; i < 10; i++ {
		v := <-ct
		fmt.Printf("%d ", v)
	}
	fmt.Printf("\n")

	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))

	counter := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go counter.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(counter.Value("somekey"))
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func co_fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func select_fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		//select wait on communication operations
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func select_default() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	if t.Left != nil {
		Walk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)

	go Walk(t1, c1)
	go Walk(t2, c2)

	for i := 0; i < 10; i++ {
		v1 := <-c1
		v2 := <-c2
		if v1 != v2 {
			return false
		}
	}
	return true
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()

	defer c.mu.Unlock()
	return c.v[key]
}
