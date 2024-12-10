package handlers

import (
	"context"
	"net/http"
	sqlc "reddit-clone/model"

	"github.com/gin-gonic/gin"
)

func AllPostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {

	posts, err := queries.GetAllPosts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get posts",
		})
		return
	}

	data := struct {
		AllPosts []sqlc.GetAllPostsRow
	}{
		AllPosts: posts,
	}

	c.HTML(http.StatusOK, "home.html", data)

}

func PersonalPostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {

	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	userID := session.Values["UserID"].(int64)

	posts, err := queries.GetUserPosts(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get posts",
		})
		return
	}

	data := struct {
		MyPosts []sqlc.GetUserPostsRow
	}{
		MyPosts: posts,
	}

	c.HTML(http.StatusOK, "myposts.html", data)
}
