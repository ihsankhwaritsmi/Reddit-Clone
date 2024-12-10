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
WHERE Users_UserID = ?;