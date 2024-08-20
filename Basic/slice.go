package main

import (
	"fmt"
	"os"
	"regexp"
)

func slice() {
	//slice is an abstraction on top of Go's array

	//array type specifies a length and an element type
	//an array's size if fixed
	var a [4]int
	a[0] = 1
	i := a[0]
	fmt.Println(i)
	//zero value of an array is an array whose elements are themselved zeroed
	fmt.Println(a[1])

	//array are values, not a pointer to the first array element
	//pass an array value will make a copy of its contents

	//literal
	b := [2]string{"Penn", "Teller"}
	c := [...]string{"Peter", "Paul"}
	fmt.Println(b)
	fmt.Println(c)

	//slice, leave out the element count
	letters := []string{"a", "b", "c"}
	fmt.Println(letters)
	var s []byte
	//func make([]T, len, cap) []T
	s = make([]byte, 5, 5)
	//equal to s = make([]byte, 5)
	//len(s) == 5
	//cap(s) == 5

	//create by "slicing" an existing slice or array
	x := s[1:]

	fmt.Println(s, x)

	//slcie is a descriptor of an array segment
	//1 pointer to array
	//2 length of the segment
	//3 capacity (number of elements in the underlying array beginning at the slice pointer)

	//slice cannot grow beyond its capacity, attempting to do so cause a runtime panic
	//increase capacity: create a new slice and copy the contents of the original slice into it
	for i := range s {
		s[i] = 1
	}
	t := make([]byte, len(s), (cap(s)+1)*2)
	// for i := range s {
	// 	t[i] = s[i]
	// }
	copy(t, s)
	s = t
	fmt.Println(s, cap(s))

	s = append(s, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	s = append(s, x...)
	fmt.Println(s)

	var p []byte
	// for _, v := range s {
	// 	p = append(p, v)
	// }
	//zero value of slice acts like a zero-length slice
	p = append(p, s...)
	//copy do not allocate new memory
	// copy(p, s)
	fmt.Println(p)
}

var digitRegexp = regexp.MustCompile("[0-0]+")

// return a slice reference the entire file, can not release the file array
func FindDigits(filename string) []byte {
	b, _ := os.ReadFile(filename)
	return digitRegexp.Find(b)
}

// copy data to a new array
func CopyDigits(filename string) []byte {
	b, _ := os.ReadFile(filename)
	b = digitRegexp.Find(b)

	c := make([]byte, len(b))
	copy(c, b)
	return c
}
