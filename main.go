package main

// Log Service

import (
	// Standard Packages
	"encoding/json"
	"log"
	"net/http"
	"os"

	// External Packages
	"github.com/gorilla/mux"
)

// LogBody The body of a log request
type LogBody struct {
	Device  string `json:"device"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func appendLog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var logBody LogBody
	err := json.NewDecoder(r.Body).Decode(&logBody)
	if err != nil {
		// If failed to log send back a InternalServerError (500) response and quit
		log.Println("Error decoding log request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Log the log in the format "Device: (Type) Message"
	log.Printf("%s: (%s) %s\n", logBody.Device, logBody.Type, logBody.Message)

	// If got here, send back a Created (201) response
	w.WriteHeader(http.StatusCreated)
}

func main() {
	// Set the log file as the output
	f, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Setup the router
	router := mux.NewRouter()

	router.HandleFunc("/log", appendLog)

	go http.ListenAndServe(":8888", router)

	select {}
}
