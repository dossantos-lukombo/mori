package db

import (
	"database/sql"

	"mori/pkg/models"
)

type LLMConvoRepository struct {
	DB *sql.DB
}

// Save inserts a new message into the conversations table.
func (repo *LLMConvoRepository) SaveConvo(convo models.Conversation) error {
	query := `
		INSERT INTO conversations (user_id, conversation_id, user_request, llm_response, new_conversation, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := repo.DB.Exec(query, convo.UserID, convo.ConversationID, convo.UserRequest, convo.LLMResponse, convo.NewConversation, convo.CreatedAt, convo.UpdateAt)
	return err
}

// Get all conversations for a specific chat
func (repo *LLMConvoRepository) GetAllConvo(convo models.Conversation) ([]models.Conversation, error) {
	query := `
		SELECT user_id, conversation_id, user_request, llm_response, new_conversation, created_at, updated_at
		FROM conversations
		WHERE user_id = $1
	`
	rows, err := repo.DB.Query(query, convo.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convos []models.Conversation
	for rows.Next() {
		var convo models.Conversation
		if err := rows.Scan(&convo.UserID, &convo.ConversationID, &convo.UserRequest, &convo.LLMResponse, &convo.NewConversation, &convo.CreatedAt, &convo.UpdateAt); err != nil {
			return nil, err
		}
		convos = append(convos, convo)
	}
	return convos, rows.Err()
}

// // SaveHistory inserts a new message into the conversations table.
// func (repo *LLMConvoRepository) SaveHistory(convo models.Conversation) error {
// 	query := `
// 		INSERT INTO conversations (history)
// 		VALUES ($1);
// 	`
// 	_, err := repo.DB.Exec(query, convo.History)
// 	return err
// }
