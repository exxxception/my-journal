package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Post struct {
	ID        int64
	UserID    int64
	Subject   string
	Event     string
	CreatedAt time.Time
}

func GetAllPosts(userID int64) ([]Post, error) {
	var posts []Post

	var query = `SELECT id, user_id, subject, event, created_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Subject, &post.Event, &post.CreatedAt)
		if err != nil {
			log.Printf("failed to scan thread row: %v", err) // LOG
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(post *Post) (int64, error) {
	var query = `INSERT INTO posts (user_id, subject, event, created_at) VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(query, post.UserID, post.Subject, post.Event, post.CreatedAt)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query: %w", err) // ERROR
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve last insert id: %w", err) // ERROR
	}

	return id, nil
}

func JournalUpdateHandler(w http.ResponseWriter, r *http.Request) {
	subject := r.FormValue("subject")
	event := r.FormValue("event")

	session, err := GetSessionFromRequest(r)
	if err != nil {
		// error
		return
	}

	post := &Post{
		UserID:    session.User.ID,
		Subject:   subject,
		Event:     event,
		CreatedAt: time.Now(),
	}

	_, err = CreatePost(post)
	if err != nil {
		http.Error(w, "failed to create post row", http.StatusInternalServerError) // HTTP ERROR
		log.Printf("failed to create post row: %v", err)                           // LOG
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	Header(w, r, true, session)
}

func SuccessPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	Header(w, r, true, session)
}

func JournalPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	posts, err := GetAllPosts(session.User.ID)
	if err != nil {
		log.Println(err)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/journal.html"))

	data := map[string]interface{}{
		"Number":      len(posts),
		"User":        &User{Username: session.User.Username},
		"Posts":       posts,
		"CurrentYear": time.Now().Year(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
