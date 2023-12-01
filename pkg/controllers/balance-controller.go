package controllers

import (
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/models"
	"gorm.io/gorm"
	"net/http"
)

func TestController(c *gin.Context, db *gorm.DB) {
	users := models.GetAllUsers(db)
	// res, _ := json.Marshal(users)

	c.JSON(http.StatusOK, gin.H{
		"message": users,
	})
}
