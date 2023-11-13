package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetDB() *gorm.DB {

	// Database connection string
	// Note: Replace with your database details
	dsn := "host=localhost user=postgres password=postgres dbname=sales_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	// Open the connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting generic database object: %v", err)
	}
	// defer sqlDB.Close()
	// Ping the database to check for connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Your database operations go here
	fmt.Println("Successfully connected to the database")

	return db
}
