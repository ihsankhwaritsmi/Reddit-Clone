package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reddit-clone/handlers"
	sqlc "reddit-clone/model"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/home")
	})

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

	authGroup.POST("/createpost", func(c *gin.Context) {
		handlers.CreatePostHandler(c, ctx, queries)
	})

	authGroup.GET("/post/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_post.html", nil)
	})

	authGroup.GET("/post/edit/:id", func(c *gin.Context) {
		handlers.EditPostPlaceholderHandler(c, ctx, queries)
	})

	authGroup.POST("/post/delete/:id", func(c *gin.Context) {
		handlers.DeletePostHandler(c, ctx, queries)
	})

	authGroup.POST("/post/edit/:id", func(c *gin.Context) {
		handlers.EditPostHandler(c, ctx, queries)
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
