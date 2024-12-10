package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reddit-clone/handlers"
	sqlc "reddit-clone/model"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func run() error {

	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./database/app.db")
	if err != nil {
		return err
	}

	queries := sqlc.New(db)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	authGroup := router.Group("/")
	authGroup.Use(handlers.AuthMiddleware())
	authGroup.GET("/home", func(c *gin.Context) {
		handlers.AllPostHandler(c, ctx, queries)
	})

	authGroup.GET("/profile", func(c *gin.Context) {
		handlers.ProfileHandler(c, ctx, queries)

	})

	authGroup.GET("/logout", func(c *gin.Context) {
		handlers.LogoutHandler(c, ctx, queries)
	})

	authGroup.GET("/myposts", func(c *gin.Context) {
		handlers.PersonalPostHandler(c, ctx, queries)
	})

	v := router.Group("/api")
	{
		v.POST("/register", func(c *gin.Context) {
			handlers.RegisterHandler(c, ctx, queries)
		})
		v.POST("/login", func(c *gin.Context) {
			handlers.LoginHandler(c, ctx, queries)
		})
	}

	fmt.Println("Server running on http://localhost:8080")
	router.Run(":8080")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
