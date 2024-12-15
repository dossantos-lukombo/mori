package db

import (
	"database/sql"

	"mori/pkg/models"
)

type MsgRepository struct {
	DB *sql.DB
}

// Save inserts a new message into the messages table.
func (repo *MsgRepository) Save(msg models.ChatMessage) error {
	query := `
		INSERT INTO messages (message_id, sender_id, receiver_id, type, content) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := repo.DB.Exec(query, msg.ID, msg.SenderId, msg.ReceiverId, msg.Type, msg.Content)
	return err
}

// SaveGroupMsg inserts a new message into the group_messages table.
func (repo *MsgRepository) SaveGroupMsg(msg models.ChatMessage) error {
	query := `
		INSERT INTO group_messages (message_id, receiver_id) 
		VALUES ($1, $2);
	`
	_, err := repo.DB.Exec(query, msg.ID, msg.ReceiverId)
	return err
}

// GetAll retrieves all messages between a sender and a receiver.
func (repo *MsgRepository) GetAll(msgIn models.ChatMessage) ([]models.ChatMessage, error) {
	query := `
		SELECT message_id, sender_id, receiver_id, type, content, created_at, is_read 
		FROM messages 
		WHERE 
			(receiver_id = $1 AND sender_id = $2) 
			OR (receiver_id = $2 AND sender_id = $1) 
		ORDER BY created_at ASC;
	`
	rows, err := repo.DB.Query(query, msgIn.ReceiverId, msgIn.SenderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.SenderId, &msg.ReceiverId, &msg.Type, &msg.Content, &msg.CreatedAt, &msg.IsRead); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// GetAllGroup retrieves all messages in a group for a user.
func (repo *MsgRepository) GetAllGroup(userId, groupId string) ([]models.ChatMessage, error) {
	query := `
		SELECT message_id, sender_id, receiver_id, type, content, created_at 
		FROM messages 
		WHERE 
			(sender_id = $1 AND receiver_id = $2) 
			OR (
				receiver_id = $2 
				AND (
					(SELECT COUNT(*) FROM groups WHERE group_id = $2 AND administrator = $1) = 1 
					OR (SELECT COUNT(*) FROM group_users WHERE group_id = $2 AND user_id = $1) = 1
				)
			) 
		ORDER BY created_at ASC;
	`
	rows, err := repo.DB.Query(query, userId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.SenderId, &msg.ReceiverId, &msg.Type, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// MarkAsRead updates a message as read.
func (repo *MsgRepository) MarkAsRead(msg models.ChatMessage) error {
	query := `
		UPDATE messages 
		SET is_read = $1 
		WHERE message_id = $2 AND receiver_id = $3;
	`
	_, err := repo.DB.Exec(query, 1, msg.ID, msg.ReceiverId)
	return err
}

// MarkAsReadGroup updates a group message as read.
func (repo *MsgRepository) MarkAsReadGroup(msg models.ChatMessage) error {
	query := `
		UPDATE group_messages 
		SET is_read = $1 
		WHERE message_id = $2 AND receiver_id = $3;
	`
	_, err := repo.DB.Exec(query, 1, msg.ID, msg.ReceiverId)
	return err
}

// GetUnread retrieves unread messages for a user.
func (repo *MsgRepository) GetUnread(userId string) ([]models.ChatStats, error) {
	query := `
		SELECT sender_id, type, COUNT(*) 
		FROM messages 
		WHERE receiver_id = $1 AND is_read = 0 
		GROUP BY sender_id, type;
	`
	rows, err := repo.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatStats
	for rows.Next() {
		var msg models.ChatStats
		if err := rows.Scan(&msg.ID, &msg.Type, &msg.UnreadMsgCount); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// GetUnreadGroup retrieves unread group messages for a user.
func (repo *MsgRepository) GetUnreadGroup(userId string) ([]models.ChatStats, error) {
	query := `
		SELECT receiver_id, type, COUNT(*) 
		FROM messages 
		WHERE type = 'GROUP' 
			AND (
				(SELECT administrator FROM groups WHERE group_id = messages.receiver_id) = $1 
				OR (SELECT COUNT(*) FROM group_users WHERE group_id = messages.receiver_id AND user_id = $1) = 1
			) 
			AND (SELECT is_read FROM group_messages WHERE message_id = messages.message_id AND receiver_id = $1) = 0 
		GROUP BY receiver_id, type;
	`
	rows, err := repo.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatStats
	for rows.Next() {
		var msg models.ChatStats
		if err := rows.Scan(&msg.ID, &msg.Type, &msg.UnreadMsgCount); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// GetChatHistoryIds retrieves all chat history IDs for a user.
func (repo *MsgRepository) GetChatHistoryIds(userId string) (map[string]bool, error) {
	idmap := make(map[string]bool)

	queryReceiver := `
		SELECT sender_id 
		FROM messages 
		WHERE receiver_id = $1 AND type = 'PERSON';
	`
	rowsReceiver, err := repo.DB.Query(queryReceiver, userId)
	if err != nil {
		return idmap, err
	}
	defer rowsReceiver.Close()

	for rowsReceiver.Next() {
		var id string
		if err := rowsReceiver.Scan(&id); err != nil {
			return idmap, err
		}
		idmap[id] = true
	}

	querySender := `
		SELECT receiver_id 
		FROM messages 
		WHERE sender_id = $1 AND type = 'PERSON';
	`
	rowsSender, err := repo.DB.Query(querySender, userId)
	if err != nil {
		return idmap, err
	}
	defer rowsSender.Close()

	for rowsSender.Next() {
		var id string
		if err := rowsSender.Scan(&id); err != nil {
			return idmap, err
		}
		idmap[id] = true
	}

	return idmap, nil
}

// HasHistory checks if a chat history exists between two users.
func (repo *MsgRepository) HasHistory(senderId, receiverId string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM messages 
		WHERE 
			(sender_id = $1 AND receiver_id = $2) 
			OR (sender_id = $2 AND receiver_id = $1);
	`
	var result int
	err := repo.DB.QueryRow(query, senderId, receiverId).Scan(&result)
	return result > 0, err
}
