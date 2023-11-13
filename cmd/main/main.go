package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/configs"
	"github.com/oxbase/portfolio_svc/pkg/routes"
	"log"
)

func main() {

	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	db := configs.GetDB()

	r := gin.Default()
	routes.PortfolioRoutes(r, db)

	if err := r.Run(configs.EnvConfigs.Port); err != nil {
		log.Println("Server failed to start ", err)
	}
}
