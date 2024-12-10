package handlers

import (
	"context"
	"log"
	"net/http"
	sqlc "reddit-clone/model"
	"strconv"

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

func CreatePostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {

	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	userID := session.Values["UserID"].(int64)

	title := c.PostForm("title")
	content := c.PostForm("body")

	_, err = queries.CreatePost(ctx, sqlc.CreatePostParams{
		Posttitle:   title,
		Postbody:    content,
		UsersUserid: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

func DeletePostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	// Retrieve the session
	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is logged in
	userID, ok := session.Values["UserID"].(int64)
	if !ok {
		log.Printf("No valid user ID in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the post ID from the URL parameters
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

	// Debug logging
	log.Printf("Attempting to delete post %d for user %d", postID, userID)

	// Attempt to delete the post
	err = queries.DeletePost(ctx, sqlc.DeletePostParams{
		UsersUserid: userID,
		Postid:      postID,
	})
	if err != nil {
		log.Printf("Delete post error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete post",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Post deleted successfully")

	// Try different redirect methods
	c.Header("Location", "/myposts")
	c.Status(http.StatusSeeOther)
	c.Abort()
}

func EditPostHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {
	// Retrieve the session
	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is logged in
	userID, ok := session.Values["UserID"].(int64)
	if !ok {
		log.Printf("No valid user ID in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the post ID from the URL parameters
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

	// Debug logging
	log.Printf("Attempting to edit post %d for user %d", postID, userID)

	err = queries.UpdatePost(ctx, sqlc.UpdatePostParams{
		UsersUserid: userID,
		Postid:      postID,
		Posttitle:   c.PostForm("title"),
		Postbody:    c.PostForm("body"),
	})
	if err != nil {
		log.Printf("Update post error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update post",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Post updated successfully")
	c.Redirect(http.StatusFound, "/myposts")

}

func EditPostPlaceholderHandler(c *gin.Context, ctx context.Context, queries *sqlc.Queries) {

	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["UserID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Get the post ID from the URL parameters
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

	dataToEdit, err := queries.GetPost(ctx, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get post",
		})
		return
	}

	data := struct {
		PostToEdit sqlc.GetPostRow
	}{
		PostToEdit: dataToEdit,
	}

	c.HTML(http.StatusOK, "edit_post.html", data)
}
