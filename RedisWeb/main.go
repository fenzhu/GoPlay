package main

import "example.com/redisweb/database"

func main() {
	database.CreateRedis(&database.Option{
		Name: "redisweb",
		Addr: "127.0.0.1:6379",
	})
}
