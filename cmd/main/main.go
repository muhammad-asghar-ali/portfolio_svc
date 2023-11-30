package main

import (
	"log"

	"github.com/0xbase-Corp/portfolio_svc/pkg/configs"
	"github.com/0xbase-Corp/portfolio_svc/pkg/routes"
	"github.com/gin-gonic/gin"
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
