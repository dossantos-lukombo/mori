package main

import (
	"fmt"
	"log"
	"mori/app"
	"mori/database"
	"mori/public/login"
	"net/http"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// Define routes
	http.HandleFunc("/login", login.LoginHandler(db))
	http.HandleFunc("/register", login.RegisterHandler(db))
	http.HandleFunc("/protected", login.AuthMiddleware(db, protectedHandler))
	http.HandleFunc("/verify-email", login.VerifyEmailHandler(db))

	// Start the server
	app.StartServer()
}

// Example protected route
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have accessed a protected route!")
}
