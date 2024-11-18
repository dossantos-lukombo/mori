package login

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"mori/database"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		usernameOrEmail := r.FormValue("username_email")
		password := r.FormValue("password")

		var user database.User
		query := `SELECT id, username, email, password, session FROM users WHERE username = $1 OR email = $2;`
		err := db.QueryRow(query, usernameOrEmail, usernameOrEmail).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Session)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Invalid username/email or password"}`, http.StatusUnauthorized)
				return
			}
			log.Printf("Error querying user: %v", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
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

		// Insert the new user into the database
		query := `INSERT INTO users (username, email, password, session) VALUES ($1, $2, $3, $4) RETURNING id;`
		sessionToken := fmt.Sprintf("%x", time.Now().UnixNano())
		var userID uint
		err = db.QueryRow(query, username, email, string(hashedPassword), sessionToken).Scan(&userID)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{"message": "Registration successful"}`)
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
