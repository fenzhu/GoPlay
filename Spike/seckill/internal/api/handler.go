package api

import (
	"net/http"
	"strconv"

	"seckill/internal/service"

	"github.com/gin-gonic/gin"
)

func SeckillHandler(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err := service.Seckill(productID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "seckill successful"})
}

func GetProductsHandler(c *gin.Context) {
	products, err := service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductHandler(c *gin.Context) {
	product, err := service.GetProduct(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func AdminResetHandler(c *gin.Context) {
	if err := service.ResetSystem(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System reset successful"})
}
