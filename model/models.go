// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package model

import (
	"database/sql"
)

type Post struct {
	Postid      int64
	Posttitle   string
	Postbody    string
	Likecount   sql.NullInt64
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	UsersUserid int64
}

type User struct {
	Userid       int64
	Userusername string
	Useremail    string
	Userpassword string
}