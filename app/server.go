package app

import (
	"fmt"
	"net/http"
)

// handler serves the login.html file as the first page
func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/login/login.html")
}

// StartServer starts the HTTP server and routes
func StartServer() {
	// Serve static files (CSS, JS, images, etc.) from the public directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Root route serves login.html
	http.HandleFunc("/", handler)

	// Start the server
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
