package main

import (
	"fmt"
	"log"

	"example/greetings"
)

func main() {
	log.SetPrefix("greetings: ")

	names := []string{"baker", "go", "sb"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
