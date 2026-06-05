package main

import (
	"log"
	"ludo-server/handler"
	"net/http"
)

func main() {
	// API routes
	http.HandleFunc("/ws", handler.ServeWS)

	// Serve static files for frontend later
	// fs := http.FileServer(http.Dir("../web/dist"))
	// http.Handle("/", fs)

	log.Println("Server starting on :18189")
	if err := http.ListenAndServe(":18189", nil); err != nil {
		log.Fatal(err)
	}
}
