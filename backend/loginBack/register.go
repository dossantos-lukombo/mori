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

		// Rate limiting check
		ip := r.RemoteAddr
		rateLimitMutex.Lock()
		if rateLimit[ip] >= 5 {
			http.Error(w, `{"error": "Too many attempts. Please try again later."}`, http.StatusTooManyRequests)
			rateLimitMutex.Unlock()
			return
		}
		rateLimitMutex.Unlock()

		// Collect form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		captchaInput := strings.TrimSpace(r.FormValue("captcha_input"))

		// Validate password confirmation
		if password != confirmPassword {
			http.Error(w, `{"error": "Passwords do not match"}`, http.StatusBadRequest)
			rateLimitMutex.Lock()
			rateLimit[ip]++
			rateLimitMutex.Unlock()
			return
		}

		// Validate password strength
		if err := validatePassword(password); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			rateLimitMutex.Lock()
			rateLimit[ip]++
			rateLimitMutex.Unlock()
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Generate a verification token
		token, err := GenerateToken()
		if err != nil {
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
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique constraint violation
				http.Error(w, `{"error": "Username or Email already taken"}`, http.StatusConflict)
				rateLimitMutex.Lock()
				rateLimit[ip]++
				rateLimitMutex.Unlock()
				return
			}
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Validate captcha AFTER database insertion passes
		captchaID, err := r.Cookie("captcha_id")
		if err != nil || captchaID.Value == "" {
			http.Error(w, `{"error": "Captcha is required"}`, http.StatusBadRequest)
			rateLimitMutex.Lock()
			rateLimit[ip]++
			rateLimitMutex.Unlock()
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
			rateLimitMutex.Lock()
			rateLimit[ip]++
			rateLimitMutex.Unlock()
			return
		}

		// Send verification email
		err = sendVerificationEmail(email, username, token)
		if err != nil {
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

func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
