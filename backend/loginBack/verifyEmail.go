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

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendVerificationEmail(email, name, token string) error {
	from := "mori.team.contact@gmail.com"
	password := "qeey kngz gmyn bzwi"

	to := []string{email}
	smtpHost := "smtp.gmail.com" // SMTP host
	smtpPort := "587"            // SMTP port

	// Email body with HTML and inline CSS
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
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
        .header img {
            width: 100px;
            height: auto;
            margin-bottom: 10px;
            border-radius: 180px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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
            <img src="https://via.placeholder.com/100x100.png?text=Mori" alt="Mori Logo">
            <h1>Mori Team</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>Thank you for signing up with Mori! Please verify your email address by clicking the button below.</p>
            <p>
                <a href="http://localhost:8080/verify-email?token=%s" class="button">Verify Email</a>
            </p>
            <p>If the button doesn't work, please copy and paste the following link into your browser.</p>
            <p>http://localhost:8080/verify-email?token=%s</p>
            <p>Thank you, and welcome to Mori.</p>
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
`, name, token, token)

	// Create the email message with MIME headers for HTML content
	message := fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", subject, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
}

func VerifyEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.ServeFile(w, r, "../frontend/email_validation/invalid_token.html") // Serve invalid token page
			return
		}

		// Update the user's `verified` status
		query := `UPDATE users SET verified = TRUE, verification_token = NULL WHERE verification_token = $1;`
		result, err := db.Exec(query, token)
		if err != nil {
			log.Printf("Error verifying email: %v", err)
			http.ServeFile(w, r, "../frontend/email_validation/error.html") // Serve error page
			return
		}

		// Check if any rows were updated
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.ServeFile(w, r, "../frontend/email_validation/invalid_token.html") // Serve invalid token page
			return
		}

		// Serve success page
		http.ServeFile(w, r, "../frontend/email_validation/email_validation.html")
	}
}
