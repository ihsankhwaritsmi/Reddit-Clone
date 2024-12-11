-- name: CreateUser :execresult
INSERT INTO Users (UserUsername, UserEmail, UserPassword)
VALUES (?, ?, ?);

-- name: AuthenticateUser :one
SELECT UserID, UserUsername, UserEmail, UserPassword
FROM Users
WHERE UserUsername = ?
AND UserPassword = ?;

-- name: GetUser :one
SELECT UserID, UserUsername, UserEmail
FROM Users
WHERE UserID = ?;

-- name: GetAllPosts :many
SELECT 
    PostID, 
    PostTitle, 
    PostBody,
    UserUsername
FROM Posts
JOIN Users
ON Posts.Users_UserID = Users.UserID
ORDER BY created_at DESC;

-- name: GetUserPosts :many
SELECT 
    PostID, 
    PostTitle, 
    PostBody,
    UserUsername
FROM Posts
JOIN Users
ON Posts.Users_UserID = Users.UserID
WHERE Users_UserID = ?
ORDER BY created_at DESC;

-- name: CreatePost :execresult
INSERT INTO Posts (PostTitle, PostBody, Users_UserID)
VALUES (?, ?, ?);

-- name: DeletePost :exec
DELETE FROM Posts WHERE PostID = ?
AND Users_UserID = ?;

-- name: UpdatePost :exec
UPDATE Posts
SET PostTitle = ?, PostBody = ?
WHERE PostID = ?
AND Users_UserID = ?;

-- name: GetPost :one
SELECT 
    PostID, 
    PostTitle, 
    PostBody
FROM Posts
WHERE PostID = ?;

-- name: GetComments :many
SELECT 
    com.CommentID,
    com.CommentBody,
    com.created_at,
    usr.UserID,
    usr.UserUsername,
    pst.PostID,
    pst.PostTitle
FROM 
    Comments AS com
JOIN 
    Users AS usr ON com.Users_UserID = usr.UserID
JOIN 
    Posts AS pst ON com.Posts_PostID = pst.PostID
WHERE 
    pst.PostID = ?;

-- name: CreateComment :exec
INSERT INTO Comments (CommentBody, Users_UserID, Posts_PostID)
VALUES (?, ?, ?);

