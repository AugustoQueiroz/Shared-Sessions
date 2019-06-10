package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
)

// 128513-128591

type Session struct {
	Session_Code [3]int
}

type Index_Context struct {
	LoginUrl string
}

type Client struct {
	client    spotify.Client
	logged_in bool
}

var live_sessions []Session

const redirectURI = "http://11155126.ngrok.io/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

var (
	clients []Client
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
	clients = append(clients, Client{auth.NewClient(tok), true})

	log.Println("Logged In!")
	http.Redirect(w, r, "http://11155126.ngrok.io/joinSession", 303)
}

func joinSession(w http.ResponseWriter, r *http.Request) {
	var session_code [3]int
	s := r.FormValue(sessionCode)
	session_code = []int(r.FormValue(sessionCode))

	// Check if session exists
	valid := false
	for _, sess := range live_sessions {
		if session_code == sess.Session_Code {
			valid = true
			break
		}
	}
	if !valid {
		log.Printf("Session %v not found", session_code)
		return
	}
	log.Printf("User %s joining session %s", clients[0].client.CurrentUser().DisplayName, session_code)
	var templateData = struct {
		curSession [3]int
		users      []string
	}{}
	templateData.curSession = session_code
	for _, i := range clients {
		templateData.users = append(templateData.users, i.client.CurrentUser().DisplayName)
	}

	t, _ := template.ParseFiles("session.html")
	t.Execute(w, &templateData)

}

func session(w http.ResponseWriter, r *http.Request) {
	if !client[0].logged_in {
		http.Redirect(w, r, "http://11155126.ngrok.io", 307)
		return
	}

	session_code := r.Form.Get(sessionCode)

	var templateData = struct {
		curSession [3]int
		users      []string
	}{}
	templateData.curSession = session_code
	for _, i := range clients {
		templateData.users = append(templateData.users, i.client.CurrentUser().DisplayName)
	}

	t, _ := template.ParseFiles("session.html")
	t.Execute(w, &templateData)

	err := clients[0].client.Play()
	if err != nil {
		log.Print(err)
	}

	player, err := clients[0].client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Print(err)
		return
	}
	log.Println(player)
}

func newSession(w http.ResponseWriter, r *http.Request) {
	// Session Page:
	//	- Session Code
	//	- Change Session Button
	if !client[0].logged_in {
		http.Redirect(w, r, "http://11155126.ngrok.io", 307)
		return
	}

	log.Println(live_sessions)

	session_code := generate_session_code()

	for _, code := range session_code {
		log.Println(code)
	}
	url := fmt.Sprintf("http://11155126.ngrok.io/session?sessionCode=%v", session_code)
	http.Redirect(w, r, url, 307)

	/*err := clients[0].client.Play()
	if err != nil {
		log.Print(err)
	}

	player, err := clients[0].client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Print(err)
		return
	}
	log.Println(player)*/
}

func main() {
	viper.AutomaticEnv()

	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)

	http.HandleFunc("/callback", authenticate)
	http.HandleFunc("/newSession", newSession)
	http.HandleFunc("/joinSession", joinSession)
	http.HandleFunc("/", index)

	go http.ListenAndServe(":8888", nil)

	select {}
}
