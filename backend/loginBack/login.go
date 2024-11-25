package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"mori/database"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var rateLimit = make(map[string]int)
var rateLimitMutex sync.Mutex

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionToken string
		var chatSessionToken string
		if r.Method == http.MethodPost {

			// Rate limiting check
			ip := r.RemoteAddr
			rateLimitMutex.Lock()

			if rateLimit[ip] >= 5 {
				http.Error(w, `{"error": "Too many attempts. Please try again later."}`, http.StatusTooManyRequests)
				rateLimitMutex.Unlock()
				return
			}

			// Verify the CSRF token (we'll handle this in a separate function)
			if err := VerifyCSRFToken(r); err != nil {
				http.Error(w, "Forbidden: "+err.Error(), http.StatusForbidden)
				return
			}

			usernameOrEmail := r.FormValue("username_email")
			password := r.FormValue("password")

			var user database.User
			query := `SELECT id, username, email, password, session, verified FROM users WHERE username = $1 OR email = $2;`
			err := db.QueryRow(query, usernameOrEmail, usernameOrEmail).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Session, &user.Verified)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error": "Invalid username/email or password"}`, http.StatusUnauthorized)
					return
				}
				log.Printf("Error querying user: %v", err)
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}

			if !user.Verified {
				http.Error(w, `{"error": "Please verify your email before logging in"}`, http.StatusForbidden)
				return
			}

			// Compare hashed password
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				http.Error(w, `{"error": "Invalid username/email or password"}`, http.StatusUnauthorized)
				return
			}

			// Generate a new session token
			// sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
			sessionToken = uuid.New().String()
			chatSessionToken = uuid.New().String()
			userID := user.ID
			updateUserStmt, err := db.Prepare(`UPDATE users SET session = $1 WHERE id = $2;`)
			if err != nil {
				fmt.Println("Error preparing update user statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			defer updateUserStmt.Close()

			// Exécuter la requête préparée pour l'utilisateur
			_, err = updateUserStmt.Exec(sessionToken, userID)
			if err != nil {
				fmt.Println("Error executing update user statement:", err)
				http.Error(w, `{"error": "can not execute request"}`, http.StatusInternalServerError)
				return
			}

			// Préparer la requête pour mettre à jour la conversation de l'utilisateur
			updateConversationStmt, err := db.Prepare(`UPDATE conversations SET conversation_uuid = $1 WHERE user_id = $2;`)
			if err != nil {
				fmt.Println("Error preparing update conversation statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			defer updateConversationStmt.Close()

			// Exécuter la requête préparée pour la conversation
			_, err = updateConversationStmt.Exec(chatSessionToken, userID)
			if err != nil {
				fmt.Println("Error executing update conversation statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// fmt.Fprintln(w, `{"message": "Login successful"}`)

		// Set the session token as a cookie

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/home/" + sessionToken, // Cookie is valid for all paths
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		})

		// Redirect to the home page
		// w.Header().Set("Content-Type", "text/html")
		http.Redirect(w, r, "/home/"+sessionToken+"/chat/"+chatSessionToken, http.StatusFound)

	}
}
