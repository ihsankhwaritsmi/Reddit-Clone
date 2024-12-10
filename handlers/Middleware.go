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

func MethodOverride() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			override := c.PostForm("_method")
			if override != "" {
				c.Request.Method = override
			}
		}
		c.Next()
	}
}
