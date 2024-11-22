package app

import (
	"fmt"
	"log"
	"net/http"
)

// LoginPageHandler serves the login.html file as the first page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/login/login.html")
}

// StartServer starts the HTTP server with the given router
func StartServer(router http.Handler) {
	// Server port
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)

	// Start the server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
