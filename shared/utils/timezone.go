package utils

import (
	"time"
)

// GetDBTime returns the current time in 'America/Toronto' timezone.
func GetDBTime() (time.Time, error) {
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(location), nil
}
