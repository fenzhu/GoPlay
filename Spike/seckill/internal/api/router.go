
package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/seckill", SeckillHandler)
	router.GET("/products", GetProductsHandler)

	// Admin routes
	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/reset", AdminResetHandler)
	}

	return router
}
