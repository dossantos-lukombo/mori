package main

import (
	"fmt"
	"log"
	"mori/app"
	"mori/captcha"
	"mori/database"
	login "mori/loginBack"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// Initialize the router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/login", login.LoginHandler(db))
	router.HandleFunc("/register", login.RegisterHandler(db))
	router.HandleFunc("/captcha", captcha.ServeCaptcha)
	router.HandleFunc("/verify-email", login.VerifyEmailHandler(db))
	router.HandleFunc("/reset-password", login.ResetPasswordHandler(db))
	router.HandleFunc("/reset-password-form", login.ServeResetPasswordForm(db))
	router.HandleFunc("/verify-reset-token", login.VerifyResetTokenHandler(db))

	// Protected routes
	protectedRoutes := router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(login.AuthMiddleware(db))
	protectedRoutes.HandleFunc("", protectedHandler)

	// Static file serving
	router.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Root route
	router.HandleFunc("/", app.LoginPageHandler)

	// Pass the router to server.go to start the server
	app.StartServer(router)
}

// Example protected route
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have accessed a protected route!")
}
