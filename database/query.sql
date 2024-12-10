-- name: GetAllPosts :many
SELECT PostID, PostTitle, PostBody
FROM Posts
ORDER BY created_at DESC;

-- name: CreateUser :execresult
INSERT INTO Users (UserUsername, UserEmail, UserPassword)
VALUES (?, ?, ?);

-- name: AuthenticateUser :one
SELECT UserID, UserUsername, UserEmail, UserPassword
FROM Users
WHERE UserUsername = ?
AND UserPassword = ?;