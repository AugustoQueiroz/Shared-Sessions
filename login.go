package main

// Login Service
//
// The purpose of this service is simply making the
// user log into spotify and, after success, redirecting
// them to the next service

import (
    // Standard Packages
    "fmt"
    "log"
    "net/http"
    "html/template"

    // External Packages
    "golang.org/x/oauth2"
    "github.com/zmb3/spotify"
    "github.com/gorilla/websocket"
)

type User struct {
    Token       oauth2.Token
    Client      spotify.Client
    LoggedIn    bool
}

var (
    clients = make(map[*websocket.Conn]bool)
    authenticator       spotify.Authenticator
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
    url := auth.AuthURL(state)

    t, _ := template.ParseFiles("index.html")
    t.Execute(w, Index_Context{LoginUrl: url})
}

func login(w http.ResponseWriter, r *http.Request) {
    _, err := authenticator.Token(state, r)
    if err != nil {
        log.Fatal("Error trying to authenticate user:", err)
    }

    if st := r.FormValue("state"); st != state {
        log.Fatal("State mismatch: %s != %s\n", st, state)
    }

    log.Println("Logged In!")
    http.Redirect(w, r, "http://localhost:8888/sessions", 303)
}

func webSocketConnect(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }

    // register client
    clients[ws] = true
    go webSocketHandler(ws)
}

func webSocketHandler(connection *websocket.Conn) {
    for {
        _, messageBody, err := connection.ReadMessage()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("%s\n", messageBody)
    }
}

func main() {
    auth = spotify.NewAuthenticator(redirectURI,
                                    spotify.ScopeUserReadCurrentlyPlaying,
                                    spotify.ScopeUserReadPlaybackState,
                                    spotify.ScopeUserModifyPlaybackState)

    http.HandleFunc("/callback", login)
    http.HandleFunc("/ws", webSocketConnect)
    http.HandleFunc("/", serveLogin)

    http.ListenAndServe(":8888", nil)
}
