package loginback

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// sendResetPasswordEmail sends a reset password email to the user.
func sendResetPasswordEmail(email, token string) error {
	from := "mori.team.contact@gmail.com"
	password := "qeey kngz gmyn bzwi"

	to := []string{email}
	smtpHost := "smtp.gmail.com" // SMTP host
	smtpPort := "587"            // SMTP port

	// Email body with HTML and inline CSS
	subject := "Reset Your Password"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            padding: 20px;
            background-color: #2c2c2c;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            color: #e4e4e4;
        }
        .header {
            text-align: center;
            padding: 10px 0;
        }
        .header h1 {
            font-size: 24px;
            margin: 0;
            color: #9146bc;
        }
		.header img {
            width: 100px;
            height: auto;
            margin-bottom: 10px;
            border-radius: 180px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .content {
            margin: 20px 0;
            font-size: 16px;
            line-height: 1.6;
        }
        .content p {
            color: #e4e4e4 !important;
            font-size: 16px !important;
        }
        .button {
            display: inline-block;
            padding: 12px 20px;
            margin: 20px 0;
            font-size: 16px;
            color: #fff !important;
            background-color: #9146bc;
            border-radius: 5px;
            text-decoration: none;
        }
        .footer {
            margin-top: 80px;
            text-align: center;
            font-size: 12px;
            color: #aaa;
        }
        .footer a {
            color: #9146bc;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://via.placeholder.com/100x100.png?text=Mori" alt="Mori Logo">
            <h1>Mori Team</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>We received a request to reset your password. If you made this request, click the button below to reset your password:</p>
            <p>
                <a href="http://localhost:8080/reset-password-form?token=%s" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, copy and paste the link below into your browser:</p>
            <p>http://localhost:8080/reset-password-form?token=%s</p>
            <p>If you did not request a password reset, you can safely ignore this email.</p>
            <p>Thank you,<br>Mori Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Mori Team. All Rights Reserved.</p>
            <p>
                <a href="#">Privacy Policy</a> | <a href="#">Contact Us</a>
            </p>
        </div>
    </div>
</body>
</html>
`, token, token)

	// Create the email message with MIME headers for HTML content
	message := fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", subject, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
}

// validatePassword checks if the password meets the requirements.
func validatePassword(password string) error {
	if len(password) < 7 {
		return errors.New("password must be at least 7 characters long")
	}

	var hasUppercase, hasNumber bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	if !hasUppercase {
		return errors.New("password must contain at least one uppercase")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// ResetPasswordHandler handles the reset password request and returns JSON responses.
func ResetPasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email is required"})
			return
		}

		// Check if the email exists in the database
		var userExists bool
		query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
		err := db.QueryRow(query, email).Scan(&userExists)
		if err != nil {
			log.Printf("Error checking email existence: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		if !userExists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email not found"})
			return
		}

		// Generate reset token
		token, err := GenerateToken()
		if err != nil {
			log.Printf("Error generating token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Store the token in the database
		storeTokenQuery := `UPDATE users SET reset_token = $1 WHERE email = $2`
		_, err = db.Exec(storeTokenQuery, token, email)
		if err != nil {
			log.Printf("Error storing reset token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Send reset password email
		err = sendResetPasswordEmail(email, token) // This is where the function is used
		if err != nil {
			log.Printf("Error sending reset password email: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to send email"})
			return
		}

		// Success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Password reset email sent successfully"})
	}
}

// ServeResetPasswordForm serves the password reset form when accessed with a valid token.
func ServeResetPasswordForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the token from the query parameters
		token := r.URL.Query().Get("token")
		if token == "" {
			http.ServeFile(w, r, "../frontend/reset_password/invalid_token.html") // Serve an invalid token page
			return
		}

		// Validate the token
		var email string
		query := `SELECT email FROM users WHERE reset_token = $1`
		err := db.QueryRow(query, token).Scan(&email)
		if err == sql.ErrNoRows {
			http.ServeFile(w, r, "../frontend/reset_password/invalid_token.html") // Serve an invalid token page
			return
		} else if err != nil {
			log.Printf("Error verifying reset token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// If the token is valid, serve the password reset form
		http.ServeFile(w, r, "../frontend/reset_password/reset_password_form.html")
	}
}

// VerifyResetTokenHandler verifies the reset token and updates the password.
func VerifyResetTokenHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
			return
		}

		// Check if token and password are provided
		if data.Token == "" || data.NewPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token and new password are required"})
			return
		}

		// Validate password strength
		if err := validatePassword(data.NewPassword); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Validate token
		var email string
		query := `SELECT email FROM users WHERE reset_token = $1`
		err := db.QueryRow(query, data.Token).Scan(&email)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired token"})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Update password and clear token in the database
		updateQuery := `UPDATE users SET password = $1, reset_token = NULL WHERE email = $2`
		_, err = db.Exec(updateQuery, hashedPassword, email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
	}
}
