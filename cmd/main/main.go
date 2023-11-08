package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/routes"
	"log"
)

func main() {
	fmt.Println("Hello World!")

	r := gin.Default()
	routes.PortfolioRoutes(r)

	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	err := r.Run()
	if err != nil {
		log.Println("Server failed to start ", err)
	}
}
