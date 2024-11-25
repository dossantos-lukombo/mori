package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"mori/database"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var rateLimit = make(map[string]int)
var rateLimitMutex sync.Mutex

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Rate limiting check
		ip := r.RemoteAddr
		rateLimitMutex.Lock()
		if rateLimit[ip] >= 5 {
			http.Error(w, `{"error": "Too many attempts. Please try again later."}`, http.StatusTooManyRequests)
			rateLimitMutex.Unlock()
			return
		}
		rateLimitMutex.Unlock()

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
				rateLimitMutex.Lock()
				rateLimit[ip]++
				rateLimitMutex.Unlock()
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
			rateLimitMutex.Lock()
			rateLimit[ip]++
			rateLimitMutex.Unlock()
			return
		}

		// Generate a new session token
		sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
		updateQuery := `UPDATE users SET session = $1 WHERE id = $2;`
		_, err = db.Exec(updateQuery, sessionToken, user.ID)
		if err != nil {
			log.Printf("Error updating session: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Set the session cookie with the session token (valid for 3 days)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Now().Add(72 * time.Hour), // Cookie expires in 3 days
			HttpOnly: true,                           // Prevent access to the cookie via JavaScript
			Secure:   false,                          // Set to true in production (requires HTTPS)
			SameSite: http.SameSiteStrictMode,        // Strict cookie policy
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Login successful"}`)
	}
}
