package loginback

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
			http.ServeFile(w, r, "public/error/invalid_token.html") // Serve a custom error page for invalid tokens
			return
		}

		// Update the user's `verified` status
		query := `UPDATE users SET verified = TRUE, verification_token = NULL WHERE verification_token = $1;`
		result, err := db.Exec(query, token)
		if err != nil {
			log.Printf("Error verifying email: %v", err)
			http.ServeFile(w, r, "public/error/internal_error.html") // Serve an error page for internal server errors
			return
		}

		// Check if any rows were updated
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.ServeFile(w, r, "public/error/expired_token.html") // Serve an error page for expired tokens
			return
		}

		// Serve the success page
		http.ServeFile(w, r, "public/login/verifyEmailFront/email_verified.css")
	}
}
