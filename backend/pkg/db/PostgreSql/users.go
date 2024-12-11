package db

import (
	"database/sql"

	"mori/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Add inserts a new user into the database.
func (repo *UserRepository) Add(user models.User) error {
	query := `
		INSERT INTO users(user_id, email, first_name, last_name, nickname, about, password, birthday, image) 
		VALUES($1, $2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9);
	`
	_, err := repo.DB.Exec(query, user.ID, user.Email, user.FirstName, user.LastName, user.Nickname, user.About, user.Password, user.DateOfBirth, user.ImagePath)
	return err
}

// EmailNotTaken checks if an email is not already registered.
func (repo *UserRepository) EmailNotTaken(email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE email = $1;
	`
	var count int
	err := repo.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// FindUserByEmail retrieves a user ID and password by email.
func (repo *UserRepository) FindUserByEmail(email string) (models.User, error) {
	query := `
		SELECT user_id, password 
		FROM users 
		WHERE email = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Password)
	return user, err
}

// GetAllAndFollowing retrieves a list of users with follower and following status.
func (repo *UserRepository) GetAllAndFollowing(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name), 
		       (SELECT COUNT(*) FROM followers WHERE followers.user_id = $1 AND follower_id = users.user_id) AS follower, 
		       (SELECT COUNT(*) FROM followers WHERE followers.user_id = users.user_id AND follower_id = $1) AS following, 
		       image 
		FROM users 
		WHERE user_id != $1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var follower, following int
		if err := rows.Scan(&user.ID, &user.Nickname, &follower, &following, &user.ImagePath); err != nil {
			return nil, err
		}
		user.Follower = follower > 0
		user.Following = following > 0
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetDataMin retrieves basic user information (ID, nickname, image).
func (repo *UserRepository) GetDataMin(userID string) (models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name), 
		       image 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.ID, &user.Nickname, &user.ImagePath)
	return user, err
}

// IsFollowing checks if the current user is following another user.
func (repo *UserRepository) IsFollowing(userID, currentUserID string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM followers 
		WHERE user_id = $1 AND follower_id = $2;
	`
	var count int
	err := repo.DB.QueryRow(query, userID, currentUserID).Scan(&count)
	return count > 0, err
}

// ProfileStatus retrieves the status of a user profile.
func (repo *UserRepository) ProfileStatus(userID string) (string, error) {
	query := `
		SELECT status 
		FROM users 
		WHERE user_id = $1;
	`
	var status string
	err := repo.DB.QueryRow(query, userID).Scan(&status)
	return status, err
}

// GetProfileMax retrieves full user information.
func (repo *UserRepository) GetProfileMax(userID string) (models.User, error) {
	query := `
		SELECT COALESCE(nickname, first_name || ' ' || last_name), 
		       first_name, 
		       last_name, 
		       image, 
		       email, 
		       TO_CHAR(birthday, 'DD.MM.YYYY'), 
		       about 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.Nickname, &user.FirstName, &user.LastName, &user.ImagePath, &user.Email, &user.DateOfBirth, &user.About)
	user.ID = userID
	return user, err
}

// GetProfileMin retrieves minimal user profile information.
func (repo *UserRepository) GetProfileMin(userID string) (models.User, error) {
	query := `
		SELECT COALESCE(nickname, first_name || ' ' || last_name), 
		       image 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.Nickname, &user.ImagePath)
	user.ID = userID
	return user, err
}

// GetFollowers retrieves the followers of a user.
func (repo *UserRepository) GetFollowers(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name) 
		FROM users 
		WHERE (SELECT COUNT(*) FROM followers WHERE followers.user_id = $1 AND follower_id = users.user_id) = 1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetFollowing retrieves the users the current user is following.
func (repo *UserRepository) GetFollowing(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name) 
		FROM users 
		WHERE (SELECT COUNT(*) FROM followers WHERE followers.follower_id = $1 AND user_id = users.user_id) = 1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetStatus retrieves the current status of a user.
func (repo *UserRepository) GetStatus(userID string) (string, error) {
	query := `
		SELECT status 
		FROM users 
		WHERE user_id = $1;
	`
	var status string
	err := repo.DB.QueryRow(query, userID).Scan(&status)
	return status, err
}

// SetStatus updates the status of a user.
func (repo *UserRepository) SetStatus(user models.User) error {
	query := `
		UPDATE users 
		SET status = $1 
		WHERE user_id = $2;
	`
	_, err := repo.DB.Exec(query, user.Status, user.ID)
	return err
}

// SaveFollower adds a follower to a user.
func (repo *UserRepository) SaveFollower(userID, followerID string) error {
	query := `
		INSERT INTO followers(user_id, follower_id) 
		VALUES ($1, $2);
	`
	_, err := repo.DB.Exec(query, userID, followerID)
	return err
}

// DeleteFollower removes a follower from a user.
func (repo *UserRepository) DeleteFollower(userID, followerID string) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2;
	`
	_, err := repo.DB.Exec(query, userID, followerID)
	return err
}
