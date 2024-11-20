package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"mori/captcha"
	"mori/database"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var registerMutex sync.Mutex

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registerMutex.Lock()
		defer registerMutex.Unlock()

		// Collect form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		captchaInput := strings.TrimSpace(r.FormValue("captcha_input"))

		// Validate password confirmation
		if password != confirmPassword {
			http.Error(w, `{"error": "Passwords do not match"}`, http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//log.Printf("Error hashing password: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Generate a verification token
		token, err := generateToken()
		if err != nil {
			//log.Printf("Error generating token: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Insert the new user into the database
		query := `INSERT INTO users (username, email, password, session, verification_token) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
		sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
		var userID uint
		err = db.QueryRow(query, username, email, string(hashedPassword), sessionToken, token).Scan(&userID)

		// Handle database constraint violations (e.g., duplicate username/email)
		if err != nil {
			//log.Printf("Error inserting new user: %v", err)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique constraint violation
				//log.Printf("Unique constraint violation: %v", pqErr.Constraint)
				if pqErr.Constraint == "users_username_key" {
					http.Error(w, `{"error": "Username is already taken"}`, http.StatusConflict)
					return
				}
				if pqErr.Constraint == "users_email_key" {
					http.Error(w, `{"error": "Email is already registered"}`, http.StatusConflict)
					return
				}
			}
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Validate captcha AFTER database insertion passes
		captchaID, err := r.Cookie("captcha_id")
		if err != nil || captchaID.Value == "" {
			http.Error(w, `{"error": "Captcha is required"}`, http.StatusBadRequest)
			return
		}

		if !captcha.VerifyCaptcha(captchaID.Value, captchaInput) {
			// If captcha fails, delete the newly inserted user to ensure consistency
			deleteQuery := `DELETE FROM users WHERE id = $1;`
			_, delErr := db.Exec(deleteQuery, userID)
			if delErr != nil {
				log.Printf("Error deleting user after captcha failure: %v", delErr)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error": "Invalid captcha", "reloadCaptcha": true}`)
			return
		}

		// Send verification email
		err = sendVerificationEmail(email, username, token)
		if err != nil {
			//log.Printf("Error sending email: %v", err)

			// If email sending fails, delete the newly inserted user to ensure consistency
			deleteQuery := `DELETE FROM users WHERE id = $1;`
			_, delErr := db.Exec(deleteQuery, userID)
			if delErr != nil {
				log.Printf("Error deleting user after email failure: %v", delErr)
			}

			http.Error(w, `{"error": "Failed to send verification email"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{"message": "Registration successful, please check your email to verify your account."}`)
	}
}

// Middleware to protect routes
func AuthMiddleware(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user database.User
		query := `SELECT id, username, email FROM users WHERE session = $1;`
		err = db.QueryRow(query, cookie.Value).Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			//log.Printf("Error querying user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
