package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(15);unique_index" json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
