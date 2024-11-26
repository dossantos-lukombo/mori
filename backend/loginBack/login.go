package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"mori/database"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Handle login logic
		if r.Method == http.MethodPost {
			// Verify the CSRF token
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

			// Generate session tokens
			sessionToken := uuid.New().String()
			chatSessionToken := uuid.New().String()

			userID := user.ID
			updateUserStmt, err := db.Prepare(`UPDATE users SET session = $1 WHERE id = $2;`)
			if err != nil {
				fmt.Println("Error preparing update user statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			defer updateUserStmt.Close()

			_, err = updateUserStmt.Exec(sessionToken, userID)
			if err != nil {
				fmt.Println("Error executing update user statement:", err)
				http.Error(w, `{"error": "Cannot execute request"}`, http.StatusInternalServerError)
				return
			}

			updateConversationStmt, err := db.Prepare(`UPDATE conversations SET conversation_uuid = $1 WHERE user_id = $2;`)
			if err != nil {
				fmt.Println("Error preparing update conversation statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			defer updateConversationStmt.Close()

			_, err = updateConversationStmt.Exec(chatSessionToken, userID)
			if err != nil {
				fmt.Println("Error executing update conversation statement:", err)
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Path:     "/home/" + sessionToken,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
				Secure:   true,
			})

			http.Redirect(w, r, "/home/"+sessionToken+"/chat/"+chatSessionToken, http.StatusFound)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
