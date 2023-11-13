package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/models"
	"gorm.io/gorm"
	"net/http"
)

func TestController(c *gin.Context, db *gorm.DB) {
	models.GetAllUsers(db)
	c.JSON(http.StatusOK, gin.H{
		"message": "Test API Request Successfull",
	})
}
