package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func flow() {
	fmt.Println("flow")
	sum := 0

	for i := 0; i < 10; i++ {
		sum += i
	}
	for sum < 50 {
		sum += sum
	}
	fmt.Println(sum)

	fmt.Println(sqrt(2), sqrt(-4))
	fmt.Println(pow(3, 2, 10), pow(3, 3, 20))

	fmt.Println(mySqrt(5))
	fmt.Println(math.Sqrt(5))

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.\n", os)
	}

	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	defer_func()
	defer_stack()
}

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		//v live only in the scope of if block
		return v
	} else {
		return lim
	}
}

func mySqrt(x float64) float64 {
	z := 1.0
	delta := 0.001
	//newton's method
	for math.Abs(z*z-x) > delta {
		z -= (z*z - x) / (2 * z)
		fmt.Printf("z: %f ", z)
	}
	fmt.Printf("\n")
	return z
}

func defer_func() {
	defer fmt.Println("world")
	fmt.Println("hello")
}

func defer_stack() {
	//defered calls are executed in last-in-first-out order
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("done")
}
