package main

import (
	"fmt"
	"sync"
)

func mmap() {
	//hash table, offer fast lookup, add, delete
	//map[KeyType]ValueType

	//nil map behaves like empty map when reading
	//attempt to write to nil map cause a runtime panic

	fmt.Printf("\n-- map --\n")
	var m map[string]int
	m = make(map[string]int)
	m["route"] = 66

	n := len(m)
	fmt.Println(n)
	i, ok := m["route"]
	fmt.Println(i, ok)

	delete(m, "route")
	i, ok = m["route"]
	fmt.Println(i, ok)

	m["1"] = 1
	m["2"] = 2
	for k, v := range m {
		fmt.Println(k, v)
	}
	commits := map[string]int{
		"rsc": 3711,
		"r":   2138,
	}
	fmt.Println(commits)

	//a map of boolean values can be used as a set-like data structure
	visited := make(map[string]bool)
	fmt.Println(visited)

	//slice, map, function is not comparable

	//map is not safe for concurrent use
	var counter = struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	go func() {
		counter.Lock()
		counter.m["key"] = 2
		counter.Unlock()
	}()

	counter.RLock()
	fmt.Println(counter.m["key"])
	counter.RUnlock()

	//map iterating order is not stable
}
