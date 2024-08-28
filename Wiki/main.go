package main

import (
	"log"
	"net/http"

	"example.com/wiki/database"
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

	http.HandleFunc("/view/", service.ViewHandler)
	http.HandleFunc("/save/", service.SaveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
