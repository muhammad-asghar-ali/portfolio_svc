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

	err := r.Run()
	if err != nil {
		log.Println("Server failed to start ", err)
	}
}
