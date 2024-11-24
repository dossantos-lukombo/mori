package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"mori/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
		updateQuery := `UPDATE users SET session = $1 WHERE id = $2;`
		_, err = db.Exec(updateQuery, sessionToken, user.ID)
		if err != nil {
			log.Printf("Error updating session: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Login successful"}`)
	}
}
