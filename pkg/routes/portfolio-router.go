package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var PortfolioRoutes = func(router *gin.Engine) {
	router.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live and accepting connections",
		})
	})
}
