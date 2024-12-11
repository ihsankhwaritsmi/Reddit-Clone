package handlers

import (
	"context"
	"log"
	"net/http"
	sqlc "reddit-clone/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SinglePostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	comments, err := queries.GetComments(ctx, postID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get comments",
		})
		return
	}

	singlePost, err := queries.GetPost(ctx, postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get post",
		})
		return
	}

	data := struct {
		Comments     []sqlc.GetCommentsRow
		SinglePost   sqlc.GetPostRow
		CommentCount int
		CreatedAt    string
	}{
		Comments:     comments,
		SinglePost:   singlePost,
		CommentCount: len(comments),
	}

	c.HTML(http.StatusOK, "single_post.html", data)
}

func CreateCommentsHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	userID := session.Values["UserID"].(int64)

	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid post ID: %s", postIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"details": err.Error(),
		})
		return
	}

	content := c.PostForm("comment")
	log.Printf("Received comment: %s", content) // Add this for debugging

	err = queries.CreateComment(ctx, sqlc.CreateCommentParams{
		Commentbody: content,
		UsersUserid: userID,
		PostsPostid: postID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create comment",
		})
		return
	}

	c.Redirect(http.StatusFound, "/post/"+c.Param("id"))
}
