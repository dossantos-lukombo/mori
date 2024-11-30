package app

import (
    "fmt"
    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

func StartServer() {
    // Clé secrète pour CSRF
    csrfKey := []byte("32-byte-long-auth-key")

    csrfMiddleware := csrf.Protect(csrfKey,
        csrf.Secure(false), // Passez à true en production (HTTPS requis)
        csrf.HttpOnly(true),
        csrf.SameSite(csrf.SameSiteStrictMode),
    )

    // Créer le routeur
    r := mux.NewRouter()

    // Utiliser le middleware CSRF pour protéger les routes
    r.Use(csrfMiddleware)

    // Route pour obtenir le token CSRF
    r.HandleFunc("/api/csrf-token", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"csrfToken": "` + csrf.Token(r) + `"}`))
    }).Methods("GET")

    fmt.Println("Server started on port :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
