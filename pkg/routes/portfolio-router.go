package routes

import (
	"github.com/gin-gonic/gin"
	Controllers "github.com/oxbase/portfolio_svc/pkg/controllers"
)

var PortfolioRoutes = func(router *gin.Engine) {
	router.GET("/healthy", Controllers.HealthCheck)

	v1 := router.Group("/api/v1")

	v1.GET("/test", Controllers.TestController)
}
