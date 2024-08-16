package main

import "fmt"

func flow() {
	fmt.Println("flow")
	sum := 0

	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
}
