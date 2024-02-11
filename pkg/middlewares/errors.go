package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errors "github.com/0xbase-Corp/portfolio_svc/pkg/shared"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case *errors.APIError:
				c.AbortWithStatusJSON(e.Code, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
