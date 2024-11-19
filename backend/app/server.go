package app

import (
	"fmt"
	"net/http"
)

// handler serves the login.html file as the first page
func handler(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	http.ServeFile(w, r, "frontend/login/login.html")
=======
	http.ServeFile(w, r, "../frontend/login/login.html")
>>>>>>> daryl
}

// StartServer starts the HTTP server and routes
func StartServer() {
	// Serve static files (CSS, JS, images, etc.) from the public directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Root route serves login.html
	http.HandleFunc("/", handler)

	// Start the server
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
