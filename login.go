package main

// Login Service
//
// The purpose of this service is simply making the
// user log into spotify and, after success, redirecting
// them to the next service

import (
	// Standard Packages

	"context"
	"html/template"
	"log"
	"net/http"

	// External Packages
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// User Struct representing a user
type User struct {
	Token oauth2.Token `json:"token"`
	Room  string       `json:"room"`
}

// LoginContext Struct with the context for the login page
type LoginContext struct {
	LoginURL string
}

var redirectURI = "http://localhost:8888/callback"

var (
	authenticator spotify.Authenticator
	state         = "abc123"
)

func serveLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	url := authenticator.AuthURL(state)

	t, _ := template.ParseFiles("resources/index.html", "resources/style/stylesheet.css")
	t.Execute(w, LoginContext{LoginURL: url})
}

func login(w http.ResponseWriter, r *http.Request) {
	token, err := authenticator.Token(state, r)
	if err != nil {
		log.Fatal("Error trying to authenticate user:", err)
	}

	if st := r.FormValue("state"); st != state {
		log.Fatal("State mismatch: %s != %s\n", st, state)
	}

	user := User{*token, ""}
	userID := getUserID(user)
	log.Println(userID)
	saveNewUser(userID, user)

	log.Println("Logged In!")
	http.Redirect(w, r, "http://localhost:8888/sessions", 303)
}

func getUserID(user User) string {
	client := authenticator.NewClient(&user.Token)
	currentUser, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	return currentUser.ID
}

func saveNewUser(userID string, user User) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	dbClient, err := app.DatabaseWithURL(context.Background(), "https://sharedsessions-c125b.firebaseio.com/")
	if err != nil {
		log.Fatal(err)
	}

	ref := dbClient.NewRef("/users/" + userID)

	ref.Set(context.Background(), user)
}

func main() {
	authenticator = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState)

	router := mux.NewRouter()

	router.HandleFunc("/callback", login)
	router.HandleFunc("/", serveLogin)

	resourcesFileServer := http.FileServer(http.Dir("resources/"))
	router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources", resourcesFileServer))

	http.ListenAndServe(":8888", router)
}
