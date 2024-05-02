package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		UserId           int       `gorm:"primaryKey" json:"user_id"`
		Username         string    `gorm:"type:varchar(255);not null" json:"username"`
		Email            string    `gorm:"type:varchar(255);unique;not null" json:"email"`
		PublicKey        string    `gorm:"type:varchar(255)" json:"public_key"`
		HashedPublicKey  string    `gorm:"type:text" json:"hashed_public_key"`
		OtherUserDetails string    `gorm:"type:text" json:"other_user_details"`
		SignupDate       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"signup_date"`
	}
)

func (User) TableName() string {
	return "users"
}

// create user
func CreateUser(tx *gorm.DB, user *User) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// get user by GetUserById
func GetUserById(tx *gorm.DB, userId int) (*User, error) {
	user := &User{}

	err := tx.Where("user_id = ?", userId).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// get user by PublicKey
func GetUserByPublicKey(tx *gorm.DB, publicKey string) (*User, error) {
	user := &User{}

	err := tx.Where("public_key = ?", publicKey).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func UpdateUser(tx *gorm.DB, user *User) error {
	if result := tx.Save(user); result.Error != nil {
		return result.Error
	}

	return nil
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
