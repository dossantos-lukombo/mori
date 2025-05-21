package db

import (
	"database/sql"

	"mori/pkg/models"
)

type GroupRepository struct {
	DB *sql.DB
}

// GetAllAndRelations retrieves all groups and their relations to the user.
func (repo *GroupRepository) GetAllAndRelations(userID string) ([]models.Group, error) {
	query := `
		SELECT 
			group_id, 
			name, 
			(SELECT COUNT(*) 
			 FROM group_users 
			 WHERE group_users.group_id = groups.group_id 
			   AND group_users.user_id = $1) AS member, 
			(administrator = $1)::int AS admin
		FROM groups;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		var member, admin int
		if err := rows.Scan(&group.ID, &group.Name, &member, &admin); err != nil {
			return nil, err
		}
		group.Member = member > 0
		group.Administrator = admin > 0
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

// GetUserGroups retrieves the groups a user belongs to or administers.
func (repo *GroupRepository) GetUserGroups(userID string) ([]models.Group, error) {
	query := `
		SELECT 
			group_id, 
			name, 
			(administrator = $1)::int AS admin
		FROM groups 
		WHERE 
			(SELECT COUNT(*) 
			 FROM group_users 
			 WHERE group_users.group_id = groups.group_id 
			   AND group_users.user_id = $1) = 1 
			OR administrator = $1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		var admin int
		if err := rows.Scan(&group.ID, &group.Name, &admin); err != nil {
			return nil, err
		}
		group.Administrator = admin > 0
		group.Member = !group.Administrator
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

// New creates a new group in the database.
func (repo *GroupRepository) New(group models.Group) error {
	query := `
		INSERT INTO groups (group_id, name, description, administrator) 
		VALUES ($1, $2, $3, $4);
	`
	_, err := repo.DB.Exec(query, group.ID, group.Name, group.Description, group.AdminID)
	return err
}

// GetData retrieves details of a specific group.
func (repo *GroupRepository) GetData(groupId string) (models.Group, error) {
	query := `
		SELECT name, description, administrator 
		FROM groups 
		WHERE group_id = $1;
	`
	var group models.Group
	err := repo.DB.QueryRow(query, groupId).Scan(&group.Name, &group.Description, &group.AdminID)
	group.ID = groupId
	return group, err
}

// GetMembers retrieves the members of a specific group.
func (repo *GroupRepository) GetMembers(groupId string) ([]models.User, error) {
	query := `
		SELECT user_id,
			COALESCE(nickname, first_name || ' ' || last_name) AS display_name, 
			image 
		FROM users
		WHERE 
			user_id = (SELECT administrator FROM groups WHERE group_id = $1)
			OR (
				(SELECT COUNT(*) 
				 FROM group_users 
				 WHERE group_id = $1 
				   AND user_id = users.user_id) = 1
			);
	`
	rows, err := repo.DB.Query(query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var member models.User
		if err := rows.Scan(&member.ID, &member.Nickname, &member.ImagePath); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

// IsMember checks if a user is a member of a specific group.
func (repo *GroupRepository) IsMember(groupId, userId string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM group_users 
		WHERE group_id = $1 
		  AND user_id = $2;
	`
	var member int
	err := repo.DB.QueryRow(query, groupId, userId).Scan(&member)
	return member > 0, err
}

// IsAdmin checks if a user is an administrator of a specific group.
func (repo *GroupRepository) IsAdmin(groupId, userId string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM groups 
		WHERE group_id = $1 
		  AND administrator = $2;
	`
	var admin int
	err := repo.DB.QueryRow(query, groupId, userId).Scan(&admin)
	return admin > 0, err
}

// GetAdmin retrieves the administrator of a specific group.
func (repo *GroupRepository) GetAdmin(groupId string) (string, error) {
	query := `
		SELECT administrator 
		FROM groups 
		WHERE group_id = $1; 
	`
	var admin string
	err := repo.DB.QueryRow(query, groupId).Scan(&admin)
	return admin, err
}

// SaveMember adds a user as a member of a specific group.
func (repo *GroupRepository) SaveMember(userId, groupId string) error {
	query := `
		INSERT INTO group_users (group_id, user_id) 
		VALUES ($1, $2);
	`
	_, err := repo.DB.Exec(query, groupId, userId)
	return err
}

// RemoveMember removes a user from a specific group.
func (repo *GroupRepository) RemoveMember(userId, groupId string) error {
	// Begin a transaction to ensure all operations complete successfully
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	
	// We can't anonymize messages as there's no sender_name column
	// Instead, we'll just remove the user from the group
	// Remove the user from the group
	_, err = tx.Exec(`
		DELETE FROM group_users 
		WHERE group_id = $1 AND user_id = $2;
	`, groupId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// Commit the transaction
	return tx.Commit()
}

// DeleteGroup deletes a group and all associated group_users entries.
func (repo *GroupRepository) DeleteGroup(groupId string) error {
	// Try a simpler approach - delete one by one without a transaction first
	// This will help identify which delete is causing problems
	
	// 1. Delete group_messages entries
	deleteGroupMsgsQuery := `DELETE FROM group_messages WHERE message_id IN 
		(SELECT message_id FROM messages WHERE receiver_id = $1 AND type = 'GROUP')`
	if _, err := repo.DB.Exec(deleteGroupMsgsQuery, groupId); err != nil {
		return err
	}
	
	// 2. Delete messages
	deleteMsgsQuery := `DELETE FROM messages WHERE receiver_id = $1 AND type = 'GROUP'`
	if _, err := repo.DB.Exec(deleteMsgsQuery, groupId); err != nil {
		return err
	}
	
	// 3. Delete notifications (fix column names: user_id not target_id)
	deleteNotifsQuery := `DELETE FROM notifications WHERE user_id = $1 OR content = $1`
	if _, err := repo.DB.Exec(deleteNotifsQuery, groupId); err != nil {
		return err
	}
	
	// 4. Delete group memberships
	deleteUsersQuery := `DELETE FROM group_users WHERE group_id = $1`
	if _, err := repo.DB.Exec(deleteUsersQuery, groupId); err != nil {
		return err
	}
	
	// 5. Delete the group itself
	deleteGroupQuery := `DELETE FROM groups WHERE group_id = $1`
	if _, err := repo.DB.Exec(deleteGroupQuery, groupId); err != nil {
		return err
	}
	
	return nil
}
