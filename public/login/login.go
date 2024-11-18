package login

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"mori/database"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendVerificationEmail(email, token string) error {
	from := "botformori@gmail.com"
	password := "brkn xnnh wokk hrzo"

	to := []string{email}
	smtpHost := "smtp.gmail.com" // SMTP host
	smtpPort := "587"            // SMTP port

	// Email body
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`Hello, Please verify your email address by clicking the link below: http://localhost:8080/verify-email?token=%s Thank you!`, token)

	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
}

func VerifyEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, `{"error": "Invalid verification token"}`, http.StatusBadRequest)
			return
		}

		// Update the user's `verified` status
		query := `UPDATE users SET verified = TRUE, verification_token = NULL WHERE verification_token = $1;`
		result, err := db.Exec(query, token)
		if err != nil {
			log.Printf("Error verifying email: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Check if any rows were updated
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Email verified successfully!"}`)
	}
}
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			http.Error(w, `{"error": "Passwords do not match"}`, http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Generate a verification token
		token, err := generateToken()
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Insert the new user into the database
		query := `INSERT INTO users (username, email, password, session, verification_token) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
		sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
		var userID uint
		err = db.QueryRow(query, username, email, string(hashedPassword), sessionToken, token).Scan(&userID)
		if err != nil {
			// Check for unique constraint violation
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation
				if pqErr.Constraint == "users_username_key" {
					http.Error(w, `{"error": "Username is already taken"}`, http.StatusConflict)
					return
				}
				if pqErr.Constraint == "users_email_key" {
					http.Error(w, `{"error": "Email is already registered"}`, http.StatusConflict)
					return
				}
			}
			log.Printf("Error inserting new user: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Send verification email
		err = sendVerificationEmail(email, token)
		if err != nil {
			log.Printf("Error sending email: %v", err)
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
			log.Printf("Error querying user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
