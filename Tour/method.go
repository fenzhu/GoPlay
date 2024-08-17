package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/tour/reader"
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

	if err := RunError(); err != nil {
		//priority, error() > string()
		fmt.Println(err)
	}

	fmt.Println(sqrtChecked(2))
	fmt.Println(sqrtChecked(-2))

	read()
	reader.Validate(&MyReader{})

	s13 := strings.NewReader("Lbh penpxrq gur pbqr!")
	r13 := rot13Reader{s13}
	io.Copy(os.Stdout, &r13)

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

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func (e *MyError) String() string {
	return fmt.Sprintf("string at %v, %s", e.When, e.What)
}

func RunError() *MyError {
	return &MyError{time.Now(), "it didn't work"}
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	// convert to float64, otherwise call Error() infinitely
	// fmt.Println(float64(e))
	return fmt.Sprintf("cannot Sqrt negative number: %f", e)
}

func sqrtChecked(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	} else {
		z := 1.0
		delta := 0.001
		//newton's method
		for math.Abs(z*z-x) > delta {
			z -= (z*z - x) / (2 * z)
			fmt.Printf("z: %f ", z)
		}
		fmt.Printf("\n")
		return z, nil
	}
}

func read() {
	r := strings.NewReader("Hello, Reader!")
	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v, err = %v, b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}

type MyReader struct{}

func (r *MyReader) Read(b []byte) (int, error) {
	for i, _ := range b {
		b[i] = 'A'
	}
	return len(b), nil
}

type rot13Reader struct {
	r io.Reader
}

func (r *rot13Reader) Read(b []byte) (int, error) {
	len, ok := r.r.Read(b)
	if ok != io.EOF {
		for i := range b {
			b[i] = rot13(b[i])
		}
	}
	return len, ok
}

func rot13(x byte) byte {
	input := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	output := []byte("NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm")
	match := bytes.IndexByte(input, x)
	if match == -1 {
		return x
	}
	return output[match]
}
