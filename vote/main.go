package main

import (
	"net/http"
	"os"

	"example.com/vote/database"
	"example.com/vote/service"
	"github.com/gin-gonic/gin"
)

func main() {
	database.CreateRedis(&database.Option{
		Name: "vote",
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

	r.POST("/article/post", func(c *gin.Context) {
		user := c.Query("user")
		title := c.Query("title")
		link := c.Query("link")
		message := service.ArticlePost(user, title, link)
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})
	r.POST("/article/vote", func(c *gin.Context) {
		user := c.Query("user")
		article := c.Query("article")
		err := service.ArtivleVote(user, article)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		}
	})
	r.GET("/article/:page", func(c *gin.Context) {
		page := c.GetInt64("page")
		articles := service.GetArticles(page)

		c.JSON(http.StatusOK, gin.H{
			"message": articles,
		})
	})

	r.Run()
}
