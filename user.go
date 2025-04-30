package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64
	Username string
	Password string
}

func GenerateSessionToken() (string, error) {
	const length = 10
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err // ERROR
	}
	return hex.EncodeToString(b), nil
}

func CreateUser(user *User) (int64, error) {
	var err error
	user.Password = HashPassword(user.Password)

	var query = `INSERT INTO users (username, password) values ($1, $2)`

	result, err := db.Exec(query, user.Username, user.Password)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query: %w", err) // ERROR
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve last insert id: %w", err) // ERROR
	}

	return id, nil
}

func GetUserByUsername(username string, user *User) error {
	var query = "SELECT id, username, password FROM users WHERE username = $1"
	return db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
}

func ReportLoginError(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("html/index.html"))

	data := map[string]interface{}{
		"Error":       "Invalid username or password",
		"User":        nil,
		"CurrentYear": time.Now().Year(),
	}

	tmpl.Execute(w, data)
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("Token")
	if err != nil {
		log.Println(err)
		return
	}

	delete(Sessions, token.Value)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserSigninHandler(w http.ResponseWriter, r *http.Request) {
	// NOTE(vlad0924): check http method ?

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		log.Printf("empty username or password") // LOG [X]

		ReportLoginError(w)
		return
	}

	var user User
	if err := GetUserByUsername(username, &user); err != nil {
		log.Printf("failed to get user: %v", err) // LOG [X]

		ReportLoginError(w)
		return
	}

	if HashPassword(password) != user.Password {
		log.Println("failed qi passwords") // LOG [X]

		ReportLoginError(w)
		return
	}

	token, err := GenerateSessionToken()
	if err != nil {
		log.Fatal("failed get generate token: %w", err) // LOG [X]
		return
	}
	expiry := time.Now().Add(OneWeek)

	session := &Session{
		ID:     user.ID,
		Expiry: expiry,
		User:   user,
	}

	Sessions[token] = session

	SetCookieToken(w, r, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	// NOTE(vlad0924): check http method ?

	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if password1 != password2 {
		// TODO(vlad0924): report error ?
		return
	}

	username := r.FormValue("username")

	user := &User{
		Username: username,
		Password: password1,
	}

	_, err := CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // HTTP ERROR
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	Header(w, r, false, nil)
}

func SigninPageHandler(w http.ResponseWriter, r *http.Request) {
	Header(w, r, false, nil)
}
