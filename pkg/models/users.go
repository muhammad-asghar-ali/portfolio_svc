package models

import (
	"gorm.io/gorm"
	"log"
)

type User struct {
	UserId   string `gorm:"type:varchar(15);unique_index" json:"user_id"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

func GetAllUsers(db *gorm.DB) []User {

	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("Error retrieving users: %v", result.Error)
	}

	// Print out the retrieved users
	// for _, user := range users {
	// 	fmt.Printf("User ID: %v,  User Name: %v, Email: %v\n", user.UserId, user.UserName, user.Email)
	// }
	return users
}
