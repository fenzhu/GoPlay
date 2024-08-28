package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/wiki/database"
	"example.com/wiki/page"
	"example.com/wiki/service"
)

func main() {
	dbOption := &database.DataOption{
		User:   "root",
		Passwd: "123456",
		Addr:   "127.0.0.1:3306",
		DBName: "wiki",
	}
	database.Center.CreateDatabase(dbOption)

	cacheOption := &database.CacheOption{
		Name: "wiki",
	}
	database.Center.CreateCache(cacheOption)

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

	http.HandleFunc("/view/", service.ViewHandler)
	http.HandleFunc("/save/", service.SaveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
