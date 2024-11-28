package loginback

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// UpdateAvatar
func UpdateAvatarHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		file, _, err := r.FormFile("avatar")
		if err != nil {
			http.Error(w, "Failed to read avatar file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Assuming the avatar is stored as a file path or binary in the database (example only)
		avatarPath := "/path/to/save/avatar" // Save the file and get the path
		// Update the user's avatar path in the database
		stmt, err := db.Prepare(`UPDATE users SET avatar = $1 WHERE id = $2`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(avatarPath, userID)
		if err != nil {
			http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Avatar updated successfully"}`)
	}
}

// UpdateUserInfo
func UpdateUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		username := r.FormValue("username")
		email := r.FormValue("email")

		stmt, err := db.Prepare(`UPDATE users SET username = $1, email = $2 WHERE id = $3`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(username, email, userID)
		if err != nil {
			http.Error(w, "Failed to update user information", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "User information updated successfully"}`)
	}
}

// ChangePassword
func ChangePasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")

		// Verify current password
		var hashedPassword string
		err := db.QueryRow(`SELECT password FROM users WHERE id = $1`, userID).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
		if err != nil {
			http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
			return
		}

		// Update to the new password
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare(`UPDATE users SET password = $1 WHERE id = $2`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(newHashedPassword, userID)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Password changed successfully"}`)
	}
}

// DeleteAccount
func DeleteAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)

		stmt, err := db.Prepare(`DELETE FROM users WHERE id = $1`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(userID)
		if err != nil {
			http.Error(w, "Failed to delete account", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Account deleted successfully"}`)
	}
}

// DeleteFavorites
func DeleteFavoritesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)

		stmt, err := db.Prepare(`DELETE FROM favorites WHERE user_id = $1`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(userID)
		if err != nil {
			http.Error(w, "Failed to delete favorites", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Favorites deleted successfully"}`)
	}
}

// DeleteConversations
func DeleteConversationsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)

		stmt, err := db.Prepare(`DELETE FROM conversations WHERE user_id = $1`)
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(userID)
		if err != nil {
			http.Error(w, "Failed to delete conversations", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Conversations deleted successfully"}`)
	}
}
