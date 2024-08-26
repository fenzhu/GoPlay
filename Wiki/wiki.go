package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(save|view)/([a-zA-Z0-9]+)$")

func main() {
	start()

	p1 := &Page{Title: "TestPage2", Body: "This is a sample Page"}
	err := p1.save()
	if err != nil {
		fmt.Println(err)
		return
	}

	p2, err := loadPage(p1.Title)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(p2.Body))

	// http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }
