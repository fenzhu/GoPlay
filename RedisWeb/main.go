package main

import (
	"net/http"
	"os"

	"example.com/redisweb/database"
	"example.com/redisweb/service"
	"github.com/gin-gonic/gin"
)

func main() {
	database.CreateRedis(&database.Option{
		Name: "redisweb",
		Addr: "127.0.0.1:6379",
	})

	logFile, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = logFile

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		user := c.Query("user")
		token := c.Query("token")
		item := c.Query("item")

		service.UpdateToken(token, user, item)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	r.Run(":8080")
}
