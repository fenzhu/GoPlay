package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page")}
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

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
