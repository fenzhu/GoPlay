package main

import (
	"fmt"
	"math"
)

const ok = 9

const Big = 1 << 100
const Small = Big >> 99

func needInt(x int) int           { return x*10 + 1 }
func needFloat(x float64) float64 { return x * 0.1 }

func basic() {
	fmt.Println("basic")
	fmt.Println(math.Pi)

	//short variable assignment, can be used inside function
	k := 3
	fmt.Println(k)

	i := 42
	//need explicit conversion
	// var f float64 = i
	f := float64(i)
	fmt.Println(i)
	fmt.Println(f)

	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))

	//wrong, greater than int64
	// fmt.Println(needInt(Big))
	fmt.Println(needFloat(Big))

	fmt.Printf("\n")
}
