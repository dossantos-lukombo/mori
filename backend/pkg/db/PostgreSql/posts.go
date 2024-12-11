package db

import (
	"database/sql"

	"mori/pkg/models"
)

type PostRepository struct {
	DB *sql.DB
}

// GetAll retrieves all posts visible to a user.
func (repo *PostRepository) GetAll(userID string) ([]models.Post, error) {
	query := `
		SELECT post_id, created_by, content, image 
		FROM posts 
		WHERE 
			(group_id IS NULL AND visibility = 'PUBLIC') 
			OR (group_id IS NULL AND visibility = 'PRIVATE' AND 
				(SELECT COUNT(*) FROM followers WHERE posts.created_by = followers.user_id AND followers.follower_id = $1) = 1) 
			OR (group_id IS NULL AND visibility = 'ALMOST_PRIVATE' AND 
				(SELECT COUNT(*) FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = $1) = 1) 
			OR (group_id IS NULL AND created_by = $1)
		ORDER BY created_at DESC;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// GetUserPosts retrieves all posts of a user visible to the current user.
func (repo *PostRepository) GetUserPosts(userID, currentUserID string) ([]models.Post, error) {
	query := `
		SELECT post_id, created_by, content, image 
		FROM posts 
		WHERE 
			(group_id IS NULL AND visibility = 'PUBLIC' AND created_by = $1) 
			OR (group_id IS NULL AND visibility = 'PRIVATE' AND created_by = $1 AND 
				(SELECT COUNT(*) FROM followers WHERE posts.created_by = followers.user_id AND followers.follower_id = $2) = 1) 
			OR (group_id IS NULL AND visibility = 'ALMOST_PRIVATE' AND created_by = $1 AND 
				(SELECT COUNT(*) FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = $2) = 1) 
			OR (group_id IS NULL AND created_by = $1 AND $1 = $2) 
		ORDER BY created_at DESC;
	`
	rows, err := repo.DB.Query(query, userID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// GetGroupPosts retrieves all posts of a specific group.
func (repo *PostRepository) GetGroupPosts(groupID string) ([]models.Post, error) {
	query := `
		SELECT post_id, created_by, content, image 
		FROM posts 
		WHERE group_id = $1 
		ORDER BY created_at DESC;
	`
	rows, err := repo.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// New creates a new post in the database.
func (repo *PostRepository) New(post models.Post) error {
	query := `
		INSERT INTO posts (post_id, group_id, created_by, content, image, visibility) 
		VALUES ($1, (NULLIF($2, '')), $3, $4, $5, $6);
	`
	_, err := repo.DB.Exec(query, post.ID, post.GroupID, post.AuthorID, post.Content, post.ImagePath, post.Visibility)
	return err
}

// SaveAccess grants access to a post for a specific user.
func (repo *PostRepository) SaveAccess(postId, userId string) error {
	query := `
		INSERT INTO almost_private (post_id, user_id) 
		VALUES ($1, $2);
	`
	_, err := repo.DB.Exec(query, postId, userId)
	return err
}
