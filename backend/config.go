package main

import (
	"database/sql"
	"mori/app"
	"mori/captcha"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Router will be used across the project
var Router *mux.Router
var Mu sync.Mutex

// InitializeRouter sets up the router
func InitializeRouter(db *sql.DB) {
	// Create a new router
	Router = mux.NewRouter()

	// Define API routes
	Router.HandleFunc("/login", LoginHandler(db)).Methods("GET", "POST")
	Router.HandleFunc("/home/{sessionToken}/chat/{chatSessionToken}", HomeHandler).Methods("GET", "POST")
	Router.HandleFunc("/register", RegisterHandler(db)).Methods("GET", "POST")
	Router.HandleFunc("/captcha", captcha.ServeCaptcha).Methods("GET", "POST")
	Router.HandleFunc("/verify-email", VerifyEmailHandler(db)).Methods("GET", "POST")
	Router.HandleFunc("/reset-password", ResetPasswordHandler(db)).Methods("GET", "POST")
	Router.HandleFunc("/reset-password-form", ServeResetPasswordForm(db)).Methods("GET", "POST")
	Router.HandleFunc("/verify-reset-token", VerifyResetTokenHandler(db)).Methods("GET", "POST")

	// Protected routes
	protectedRoutes := Router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(AuthMiddleware(db))
	// llmProtectedRoutes := Router.PathPrefix("/llm-protected").Subrouter()
	// llmProtectedRoutes.Use(SendRequestWithToken())

	// Static file serving
	Router.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Root route
	Router.HandleFunc("/", app.LoginPageHandler)
}
