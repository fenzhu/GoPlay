package main

import (
	"fmt"
	"strings"

	"golang.org/x/tour/wc"
)

type Vertex struct {
	X int
	Y int
}

func types() {
	fmt.Println("types")

	i := 42
	p := &i
	fmt.Println(p)
	fmt.Println(*p)
	//go has no pointer arithmetic

	fmt.Println(Vertex{1, 2})
	v := Vertex{2, 4}
	v.X = 4
	fmt.Println(v)

	pv := &v
	pv.X = 6
	fmt.Println(v)

	v2 := Vertex{X: 1}
	fmt.Println(v2)

	var a [2]string
	a[0] = "hello"
	a[1] = "world"
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)
	//slice, get [1, 4) of primes
	//slice just describes a section of an underlying array
	var s []int = primes[1:4]
	fmt.Println(s)

	//slice literals
	q := []int{1, 2, 3, 4}
	fmt.Println(q)
	//q[:] == q[:4] == q[0:] == q[0:4]
	fmt.Println(q[:])

	printSlice(q)
	printSlice(q[:0])
	printSlice(q[:4])
	printSlice(q[2:])
	//nil slice
	var si []int
	printSlice(si)

	//make allocate a zeroed array and return a slice that refers to that array
	ms := make([]int, 5)
	msb := make([]int, 0, 5)
	printSlice(ms)
	printSlice(msb)
	ms = append(ms, 1)
	printSlice(ms)

	for i, v := range ms {
		fmt.Printf("%d : %d\n", i, v)
	}

	sliceExercise()

	m := make(map[string]Vertex)
	m["Bell"] = Vertex{40, -74}
	fmt.Println(m)

	m1 := map[string]Vertex{
		"b": Vertex{1, 2},
		"c": {3, 4},
	}
	fmt.Println(m1)

	delete(m1, "b")

	v, mok := m["b"]
	fmt.Println(v, mok)

	wc.Test(mapExercise)

	closure()
	closureExercise()

	fmt.Printf("\n")
}

func printSlice(s []int) {
	//len: number of elements the slice contains
	//cap: number of elements in the underlying array, counting from the first element in slice
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func sliceExercise() {
	// pic.Show(myPic)
}

func myPic(dx, dy int) [][]uint8 {
	image := make([][]uint8, dy)
	for i := range image {
		image[i] = make([]uint8, dx)
		for j := range image[i] {
			image[i][j] = uint8(i * j)
		}
	}
	return image
}

func mapExercise(s string) map[string]int {
	m := make(map[string]int)
	words := strings.Fields(s)
	for _, w := range words {
		m[w] += 1
	}

	return m
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func closure() {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}
}

func fibonacci() func() int {
	past := 0
	cur := 1
	return func() int {
		res := past
		past = cur
		cur += res
		return res
	}
}

func closureExercise() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
