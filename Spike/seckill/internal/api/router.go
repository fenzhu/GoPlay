package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// Use gin.New() instead of gin.Default() to avoid default gin.Logger() middleware
	router := gin.New()

	// Add Recovery middleware to recover from panics, if desired
	router.Use(gin.Recovery())

	router.POST("/seckill", SeckillHandler)
	router.GET("/products", GetProductsHandler)

	// Admin routes
	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/reset", AdminResetHandler)
	}

	return router
}