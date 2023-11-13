package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/configs"
	"github.com/oxbase/portfolio_svc/pkg/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {

	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	r := gin.Default()
	routes.PortfolioRoutes(r)

	// Database connection string
	// Note: Replace with your database details
	dsn := "host=localhost user=postgres password=postgres dbname=discuss_dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	// Open the connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Error getting generic database object: %v", err)
		}
		sqlDB.Close()
	}()

	// Your database operations go here
	fmt.Println("Successfully connected to the database")

	if err := r.Run(configs.EnvConfigs.Port); err != nil {
		log.Println("Server failed to start ", err)
	}
}
