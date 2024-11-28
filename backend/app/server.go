package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// loginPageHandler serves the login page with a CSRF token
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a CSRF token
	csrfToken, err := GenerateToken()
	if err != nil {
		http.Error(w, "Error generating CSRF token", http.StatusInternalServerError)
		return
	}

	// Set CSRF token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token", // Name of the cookie
		Value:    csrfToken,    // The CSRF token value
		Path:     "/",          // Available throughout the application
		HttpOnly: false,        // Allow JS to access the cookie
		Secure:   false,        // Set to true if you're using HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Use Go templates to pass CSRF token to the HTML form
	tmpl, err := template.ParseFiles("../frontend/login/login.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Pass the CSRF token to the HTML form (used as a hidden input)
	err = tmpl.Execute(w, map[string]interface{}{"CSRFToken": csrfToken})
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
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

func StartServerLLMProtected(router http.Handler) {

	// Server port
	port := ":8000"
	fmt.Printf("Server LLM started on port %s\n", port)

	// http.HandleFunc("/auth/refresh", middleware.VerifyAndRefreshTokenHandler)

	// Start the server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
