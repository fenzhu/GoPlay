package main

import "fmt"

func channel() {
	fmt.Printf("\n-- channel --\n")

	c := make(chan int, 5)
	go generate(c)
	for i := range c {
		fmt.Println(i)
	}

	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}

func generate(c chan int) {
	for i := 0; i < cap(c); i++ {
		c <- i
	}
	close(c)
}
