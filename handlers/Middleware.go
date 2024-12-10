package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session-name")
		if err != nil || session.Values["UserID"] == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Proceed with request
		c.Next()
	}
}
