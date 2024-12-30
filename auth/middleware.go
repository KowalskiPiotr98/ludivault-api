package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// GetUserMiddleware returns gin handler function that reads user data from session store and saves it as context item.
// Note that this middleware does not prevent not logged in users from using the app.
func GetUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := RetrieveUserFromSession(c)
		if err != nil {
			log.Debugf("Failed to retrieve user from session. Possibly just not logged in: %v", err)
			return
		}

		c.Set("userId", userId)
	}
}

// IsLoggedIn returns true if the user is logged in to the application.
// Use GetUserId to retrieve the ID of the currently logged in user.
func IsLoggedIn(c *gin.Context) bool {
	_, exists := c.Get("userId")
	return exists
}

// GetUserId returns the id of the user currently logged in.
// It will return uuid.Nil if no user is logged in.
//
// Use IsLoggedIn first to check if the user is logged in.
func GetUserId(c *gin.Context) uuid.UUID {
	value, exists := c.Get("userId")
	if !exists {
		return uuid.Nil
	}
	return value.(uuid.UUID)
}
