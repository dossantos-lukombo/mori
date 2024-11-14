package app

import (
    "fmt"
    "net/http"
)


func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Bienvenue sur mon serveur en Go!")
}

// Fonction pour démarrer le serveur
func StartServer() {
    http.HandleFunc("/", handler)    

    port := ":8080"
    fmt.Printf("Serveur démarré sur le port %s\n", port)

    if err := http.ListenAndServe(port, nil); err != nil {
        fmt.Printf("Erreur lors du démarrage du serveur: %s\n", err)
    }
}
