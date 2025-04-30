package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func Header(w http.ResponseWriter, r *http.Request, authorized bool, session *Session) {
	page := "html/index.html"
	if r.URL.Path != "/" {
		page = "html" + r.URL.Path + ".html"
	}

	tmpl := template.Must(template.ParseFiles(page))

	data := map[string]interface{}{
		"Error":       "",
		"User":        nil,
		"CurrentYear": time.Now().Year(),
	}

	if !authorized {
		tmpl.Execute(w, data)
		return
	}

	if session != nil {
		data["User"] = &User{Username: session.User.Username}
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		Header(w, r, false, nil)
		return
	}

	Header(w, r, true, session)
}

func FAQPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		Header(w, r, false, nil)
		return
	}

	Header(w, r, true, session)
}

func ContactPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		Header(w, r, false, nil)
		return
	}

	Header(w, r, true, session)
}
