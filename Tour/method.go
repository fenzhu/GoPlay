package main

import (
	"fmt"
	"math"
	"strconv"
)

type Point struct {
	X, Y float64
}

type MyFloat float64

// method is a function with a special receiver argument
func method() {
	fmt.Println("method")

	p := Point{3, 4}
	fmt.Println(p.Abs())

	f := MyFloat(-math.SqrtE)
	fmt.Println(f.Abs())

	//go interprets this statement as (&p).Scale(10)
	p.Scale(10)
	fmt.Println(p)
	//reversely, interpret  (&p).func() as p.func()

	var i any = "hello"
	s, ok := i.(string)
	fmt.Println(s, ok)
	s1, ok1 := i.(int)
	fmt.Println(s1, ok1)

	CheckType(21)
	CheckType("Hello")
	CheckType(true)

	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, &ip)
	}

	fmt.Printf("\n")
}

func (p Point) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// declare method on non-struct types
// only on types defined in the same package
func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	} else {
		return float64(f)
	}
}

func (p *Point) Scale(s float64) {
	p.X *= s
	p.Y *= s
}

// interface is a set of method signatures
// interface{} empty interface may hold values of any type
// type any = interface{}
func CheckType(v any) {
	switch val := v.(type) {
	case int:
		fmt.Printf("twice %v is %v\n", val, val*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", val, len(val))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

type IPAddr [4]byte

func (ip *IPAddr) String() string {
	str := ""
	arr := make([]string, 4)
	for i := 0; i < 4; i++ {
		str += fmt.Sprintf("%d.", ip[i])
		//do not use string(ip[i])
		arr[i] = strconv.Itoa(int(ip[i]))
	}
	fmt.Println(arr)
	return str[:len(str)-1]
}
