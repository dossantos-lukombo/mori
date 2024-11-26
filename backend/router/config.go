package router

import (
	"database/sql"
	"mori/app"
	"mori/captcha"
	"mori/home"
	loginback "mori/loginBack"
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
	Router.HandleFunc("/login", loginback.LoginHandler(db)).Methods("GET", "POST")
	Router.HandleFunc("/home/{sessionToken}/chat/{chatSessionToken}", home.HomeHandler).Methods("GET", "POST")
	Router.HandleFunc("/register", loginback.RegisterHandler(db)).Methods("POST", "GET")
	Router.HandleFunc("/captcha", captcha.ServeCaptcha).Methods("GET", "POST")
	Router.HandleFunc("/verify-email", loginback.VerifyEmailHandler(db)).Methods("POST", "GET")
	Router.HandleFunc("/reset-password", loginback.ResetPasswordHandler(db)).Methods("POST", "GET")
	Router.HandleFunc("/reset-password-form", loginback.ServeResetPasswordForm(db)).Methods("POST", "GET")
	Router.HandleFunc("/verify-reset-token", loginback.VerifyResetTokenHandler(db)).Methods("POST", "GET")
	// Protected routes
	protectedRoutes := Router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(loginback.AuthMiddleware(db))

	// Static file serving
	Router.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Root route
	Router.HandleFunc("/", app.LoginPageHandler)
}
