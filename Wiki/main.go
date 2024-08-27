package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"example.com/wiki/database"
	"example.com/wiki/page"
	"example.com/wiki/service"
)

var validPath = regexp.MustCompile("^/(save|view)/([a-zA-Z0-9]+)$")

func main() {
	dbOption := &database.DataOption{
		User:   "root",
		Passwd: "123456",
		Addr:   "127.0.0.1:3306",
		DBName: "wiki",
	}
	database.Center.CreateDatabase(dbOption)

	p1 := &page.Page{Title: "TestPage2", Body: "This is a sample Page"}
	err := p1.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	p2, err := page.LoadPage(p1.Title)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", makeHandler(service.ViewHandler))
	http.HandleFunc("/save/", makeHandler(service.SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
