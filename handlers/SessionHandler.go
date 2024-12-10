package handlers

import (
	"context"
	"net/http"
	sqlc "reddit-clone/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func RegisterHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Userusername: username,
		Useremail:    email,
		Userpassword: password,
	})

	// Handle database errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// If Success redirect to login
	c.Redirect(http.StatusFound, "/login")
}

func LoginHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check username and password
	user, err := queries.AuthenticateUser(ctx, sqlc.AuthenticateUserParams{
		Userusername: username,
		Userpassword: password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Get session and set token
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get session",
		})
		return
	}

	session.Values["UserID"] = user.Userid
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save session",
		})
		return
	}

	// If Success Render the login template
	c.Redirect(http.StatusFound, "/home")
}
