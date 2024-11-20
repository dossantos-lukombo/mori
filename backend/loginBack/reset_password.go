package loginback

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
            <h1>Password Reset</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>We received a request to reset your password. If you made this request, click the button below to reset your password:</p>
            <p>
                <a href="http://localhost:8080/reset-password?token=%s" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, copy and paste the link below into your browser:</p>
            <p>http://localhost:8080/reset-password?token=%s</p>
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
		err = sendResetPasswordEmail(email, token)
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
