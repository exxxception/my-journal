package main

import (
	"log"
	"net/http"
	"strings"
)

func HandleAPIRequest(w http.ResponseWriter, r *http.Request, path string) {
	switch {
	case StartsWith(path, "/user"):
		switch path[len("/user"):] {
		case "/signin":
			UserSigninHandler(w, r)
		case "/signup":
			UserSignupHandler(w, r)
		case "/logout":
			UserLogoutHandler(w, r)
		}
	case StartsWith(path, "/journal"):
		switch path[len("/journal"):] {
		case "/update":
			JournalUpdateHandler(w, r)
		}
	}
}

func HandlePageRequest(w http.ResponseWriter, r *http.Request, path string) {
	switch path {
	default:
		IndexPageHandler(w, r)
	case "/signup":
		SignupPageHandler(w, r)
	case "signin":
		SigninPageHandler(w, r)
	case "/update":
		UpdatePageHandler(w, r)
	case "/success":
		SuccessPageHandler(w, r)
	case "/journal":
		JournalPageHandler(w, r)
	case "/faq":
		FAQPageHandler(w, r)
	case "contact":
		ContactPageHandler(w, r)
	}
}

type Router struct{}

func (rt *Router) RouterFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	default:
		HandlePageRequest(w, r, path)
	case StartsWith(path, "/api"):
		HandleAPIRequest(w, r, path[len("/api"):])
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/styles/"):
		http.StripPrefix("/styles/", http.FileServer(http.Dir("./html/styles"))).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/img/"):
		http.StripPrefix("/img/", http.FileServer(http.Dir("./html/img"))).ServeHTTP(w, r)
	default:
		rt.RouterFunc(w, r)
	}
}

func main() {
	if err := OpenDB("db"); err != nil {
		log.Fatalf("Failed to open DB: %v", err) // LOG
	}
	defer CloseDB()

	router := &Router{}
	log.Println("Server start...")                  // LOG
	log.Fatal(http.ListenAndServe(":8080", router)) // LOG
}
