package main

import "fmt"

var a int = 3
var b int = f()
var c int

//  = 7

func f() int {
	fmt.Printf("func f called\n")
	return 5
}

func init() {
	fmt.Printf("func init called\n")
}

func main() {
	fmt.Printf("func main called\n")
	c = 7
	fmt.Printf("%d %d %d\n", a, b, c)
}
