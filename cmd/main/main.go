package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/configs"
	"github.com/oxbase/portfolio_svc/pkg/routes"
	"log"
)

func main() {

	r := gin.Default()
	routes.PortfolioRoutes(r)

	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	err := r.Run(configs.EnvConfigs.Port)
	if err != nil {
		log.Println("Server failed to start ", err)
	}
}
