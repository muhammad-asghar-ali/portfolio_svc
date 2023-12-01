package main

import (
	"github.com/gin-gonic/gin"
	"github.com/0xbase-Corp/portfolio_svc/pkg/configs"
	"github.com/0xbase-Corp/portfolio_svc/pkg/routes"
	"log"
	"github.com/joho/godotenv"
)

func main() {
		// Load environment variables file from app.env
		err := godotenv.Load("app.env")
		if err != nil {
			log.Fatal("Error loading app.env file")
		}

		//print out message to confirm app.env is loaded
		log.Println("app.env file loaded successfully")

	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	db := configs.GetDB()

	r := gin.Default()
	routes.PortfolioRoutes(r, db)

	if err := r.Run(configs.EnvConfigs.Port); err != nil {
		log.Println("Server failed to start ", err)
	}
}
