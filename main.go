package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
)

// 128513-128591

type Session struct {
	Session_Code [3]int
}

type Index_Context struct {
	LoginUrl string
}

var live_sessions []Session

const redirectURI = "http://11155126.ngrok.io/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

var (
	client    spotify.Client
	logged_in = false
)

func generate_session_code() [3]int {
	// Generates a random session code and checks that it's not being used
	var session_code [3]int
	valid := false

	for valid == false {
		session_code = [3]int{rand.Intn(78) + 128513, rand.Intn(78) + 128513, rand.Intn(78) + 128513}
		valid = true

		// UPGRADE THIS TO BINARY SEARCH
		for _, sess := range live_sessions {
			if session_code == sess.Session_Code {
				valid = false
				break
			}
		}
		// STOP BEING LAZY
	}

	live_sessions = append(live_sessions, Session{Session_Code: session_code})
	return session_code
}

func index(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, Index_Context{LoginUrl: url})
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client = auth.NewClient(tok)
	logged_in = true

	fmt.Println("Logged In!")
	http.Redirect(w, r, "http://11155126.ngrok.io/sessions", 303)
}

func joinSession(w http.ResponseWriter, r *http.Request) {

}

func session(w http.ResponseWriter, r *http.Request) {
	// Session Page:
	//	- Session Code
	//	- Change Session Button
	if !logged_in {
		http.Redirect(w, r, "http://11155126.ngrok.io", 307)
		return
	}

	fmt.Println(live_sessions)

	session_code := generate_session_code()

	for _, code := range session_code {
		fmt.Println(code)
	}

	t, _ := template.ParseFiles("session.html")
	t.Execute(w, Session{Session_Code: session_code})

	err := client.Play()
	if err != nil {
		log.Print(err)
	}

	player, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(player)
}

func main() {
	os.Setenv("SPOTIFY_ID", "6e5bf9d9b01a4232a2f9d3b666a714b6")
	os.Setenv("SPOTIFY_SECRET", "d745ad859b08470e95e423ad94940379")

	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)

	http.HandleFunc("/callback", authenticate)
	http.HandleFunc("/sessions", session)
	http.HandleFunc("/", index)

	go http.ListenAndServe(":8888", nil)

	select {}
}
