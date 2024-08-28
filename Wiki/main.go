package main

import (
	"log"
	"net/http"
	"os"

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

	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/view/", service.ViewHandler)
	http.HandleFunc("/save/", service.SaveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
