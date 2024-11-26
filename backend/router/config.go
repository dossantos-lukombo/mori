package router

import (
	"database/sql"
	"mori/app"
	"mori/captcha"
	"mori/home"
	loginback "mori/loginBack"
	"net/http"

	"github.com/gorilla/mux"
)

// Router will be used across the project
var Router *mux.Router

// InitializeRouter sets up the router
func InitializeRouter(db *sql.DB) {
	// Create a new router
	Router = mux.NewRouter()

	// Define API routes
	Router.HandleFunc("/login", loginback.LoginHandler(db))
	Router.HandleFunc("/register", loginback.RegisterHandler(db))
	Router.HandleFunc("/captcha", captcha.ServeCaptcha)
	Router.HandleFunc("/verify-email", loginback.VerifyEmailHandler(db))
	Router.HandleFunc("/reset-password", loginback.ResetPasswordHandler(db))
	Router.HandleFunc("/reset-password-form", loginback.ServeResetPasswordForm(db))
	Router.HandleFunc("/verify-reset-token", loginback.VerifyResetTokenHandler(db))
	Router.HandleFunc("/home/{sessionToken}/chat/{chatSessionToken}", home.HomeHandler).Methods("GET", "POST")
	// Protected routes
	protectedRoutes := Router.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(loginback.AuthMiddleware(db))

	// refreshRoutes := Router.PathPrefix("/auth/refresh").Subrouter()
	// refreshRoutes.Use(middleware.RefreshMiddleware)

	// Static file serving
	Router.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Root route
	Router.HandleFunc("/", app.LoginPageHandler)
}
