package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	// Set up routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/home", authMiddleware(homeHandler))
	http.HandleFunc("/me", authMiddleware(meHandler))
	http.HandleFunc("/logout", authMiddleware(logoutHandler))

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database/app.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if not exist
	createTables := `
	CREATE TABLE IF NOT EXISTS Users (
		UserID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserUsername TEXT NOT NULL UNIQUE,
		UserEmail TEXT NOT NULL UNIQUE,
		UserPassword TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS Posts (
		PostID INTEGER PRIMARY KEY AUTOINCREMENT,
		PostTitle TEXT NOT NULL,
		PostBody TEXT NOT NULL,
		LikeCount INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		Users_UserID INTEGER NOT NULL,
		FOREIGN KEY (Users_UserID) REFERENCES Users(UserID)
	);
	`
	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := db.Exec("INSERT INTO Users (UserUsername, UserEmail, UserPassword) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		err := db.QueryRow("SELECT UserID FROM Users WHERE UserUsername = ? AND UserPassword = ?", username, password).Scan(&userID)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set session token (in a real app, use a secure random token)
		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: fmt.Sprintf("%d", userID),
			Path:  "/",
		})

		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get logged-in user's ID from the session token
	userID := cookie.Value

	// Fetch user's posts
	myPosts := []struct {
		PostID    int
		PostTitle string
		PostBody  string
		LikeCount int
	}{}
	rows, err := db.Query("SELECT PostID, PostTitle, PostBody, LikeCount FROM Posts WHERE Users_UserID = ?", userID)
	if err != nil {
		http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post struct {
			PostID    int
			PostTitle string
			PostBody  string
			LikeCount int
		}
		if err := rows.Scan(&post.PostID, &post.PostTitle, &post.PostBody, &post.LikeCount); err != nil {
			http.Error(w, "Failed to parse your posts", http.StatusInternalServerError)
			return
		}
		myPosts = append(myPosts, post)
	}

	// Fetch all posts
	allPosts := []struct {
		PostID    int
		PostTitle string
		PostBody  string
		LikeCount int
	}{}
	rows, err = db.Query("SELECT PostID, PostTitle, PostBody, LikeCount FROM Posts")
	if err != nil {
		http.Error(w, "Failed to fetch all posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post struct {
			PostID    int
			PostTitle string
			PostBody  string
			LikeCount int
		}
		if err := rows.Scan(&post.PostID, &post.PostTitle, &post.PostBody, &post.LikeCount); err != nil {
			http.Error(w, "Failed to parse all posts", http.StatusInternalServerError)
			return
		}
		allPosts = append(allPosts, post)
	}

	// Prepare data for template
	data := struct {
		MyPosts []struct {
			PostID    int
			PostTitle string
			PostBody  string
			LikeCount int
		}
		AllPosts []struct {
			PostID    int
			PostTitle string
			PostBody  string
			LikeCount int
		}
	}{
		MyPosts:  myPosts,
		AllPosts: allPosts,
	}

	// Render template
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render home page", http.StatusInternalServerError)
	}
}

type Post struct {
	Title string
	Body  string
}

func fetchPosts(rows *sql.Rows) []Post {
	var posts []Post
	for rows.Next() {
		var post Post
		rows.Scan(&post.Title, &post.Body)
		posts = append(posts, post)
	}
	return posts
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Convert the session token to user ID
	userID := cookie.Value

	// Fetch user details from the database
	var username, email string
	err = db.QueryRow("SELECT UserUsername, UserEmail FROM Users WHERE UserID = ?", userID).Scan(&username, &email)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		Username string
		Email    string
	}{
		Username: username,
		Email:    email,
	}

	// Render the profile page
	tmpl := template.Must(template.ParseFiles("templates/me.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Delete the cookie
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}
