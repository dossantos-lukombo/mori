package router

import (
    "database/sql"
    "mori/app"
    "mori/captcha"
    "mori/home"
    loginback "mori/loginBack"
    "net/http"
    "sync"

    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
)

// Router will be used across the project
var Router *mux.Router
var Mu sync.Mutex

// corsMiddleware handles CORS settings for incoming requests
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight (OPTIONS) request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Pass the request to the next handler
        next.ServeHTTP(w, r)
    })
}

// InitializeRouter sets up the router
func InitializeRouter(db *sql.DB) {
    // Clé secrète pour CSRF
    csrfKey := []byte("32-byte-long-auth-key")

    // Créer une nouvelle instance du middleware CSRF
    csrfMiddleware := csrf.Protect(csrfKey,
        csrf.Secure(false), // Passez à true en production pour HTTPS
        csrf.HttpOnly(true),
        csrf.SameSite(csrf.SameSiteStrictMode),
    )

    // Create a new router
    Router = mux.NewRouter()

    // Apply CORS middleware to all routes
    Router.Use(corsMiddleware)

    // Define API routes with "/api" prefix
    api := Router.PathPrefix("/api").Subrouter()

    // Ajoutez le middleware CSRF à toutes les routes protégées
    api.Use(csrfMiddleware)

    // Route pour obtenir le token CSRF
    api.HandleFunc("/csrf-token", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"csrfToken": "` + csrf.Token(r) + `"}`))
    }).Methods("GET")

    // Définir les routes API
    api.HandleFunc("/login", loginback.LoginHandler(db)).Methods("POST")
    api.HandleFunc("/home/{sessionToken}/chat/{chatSessionToken}", home.HomeHandler).Methods("GET", "POST")
    api.HandleFunc("/register", loginback.RegisterHandler(db)).Methods("POST")
    api.HandleFunc("/captcha", captcha.ServeCaptcha).Methods("GET")
    api.HandleFunc("/verify-email", loginback.VerifyEmailHandler(db)).Methods("POST")
    api.HandleFunc("/reset-password", loginback.ResetPasswordHandler(db)).Methods("POST")
    api.HandleFunc("/reset-password-form", loginback.ServeResetPasswordForm(db)).Methods("GET")
    api.HandleFunc("/verify-reset-token", loginback.VerifyResetTokenHandler(db)).Methods("POST")
    api.HandleFunc("/profile/avatar", loginback.UpdateAvatarHandler(db)).Methods("POST")
    api.HandleFunc("/profile/update", loginback.UpdateUserInfoHandler(db)).Methods("PUT")
    api.HandleFunc("/profile/change-password", loginback.ChangePasswordHandler(db)).Methods("PUT")
    api.HandleFunc("/profile/delete-account", loginback.DeleteAccountHandler(db)).Methods("DELETE")
    api.HandleFunc("/profile/delete-favorites", loginback.DeleteFavoritesHandler(db)).Methods("DELETE")
    api.HandleFunc("/profile/delete-conversations", loginback.DeleteConversationsHandler(db)).Methods("DELETE")

    // Protected routes
    protectedRoutes := api.PathPrefix("/protected").Subrouter()
    protectedRoutes.Use(loginback.AuthMiddleware(db))

    // Root route
    Router.HandleFunc("/", app.LoginPHandler)
}
