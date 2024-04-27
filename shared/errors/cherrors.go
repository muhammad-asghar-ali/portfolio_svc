package errors

import (
	"github.com/gin-gonic/gin"
)

func HandleChannelErrors(errorCh <-chan error, c *gin.Context) []string {
	var errs []error

	for err := range errorCh {
		errs = append(errs, err)
	}

	errorMessages := make([]string, 0, len(errs))
	if len(errs) > 0 {
		for _, err := range errs {
			errorMessages = append(errorMessages, err.Error())
		}

		return errorMessages
	}

	return errorMessages
}
