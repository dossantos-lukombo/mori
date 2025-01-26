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
		INSERT INTO conversations (user_id, conversation_id, user_request, llm_response, new_conversation) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := repo.DB.Exec(query, convo.UserID, convo.ConversationID, convo.UserRequest, convo.LLMResponse, convo.NewConversation)
	return err
}

// Get all conversations for a specific chat
func (repo *LLMConvoRepository) GetAllConvo(convo models.Conversation) ([]models.Conversation, error) {
	query := `
		SELECT user_id, conversation_id, user_request, llm_response, new_conversation
		FROM conversations
		WHERE user_id = $1
	`
	rows, err := repo.DB.Query(query, convo.UserID)
	if err != nil {
		return nil, err
	}

	var convos []models.Conversation
	for rows.Next() {
		var convo models.Conversation
		if err := rows.Scan(&convo.UserID, &convo.ConversationID, &convo.UserRequest, &convo.LLMResponse, &convo.NewConversation); err != nil {
			return nil, err
		}
		convos = append(convos, convo)
	}
	defer rows.Close()
	return convos, rows.Err()
}

// get the last conversation_id from the conversations table
func (repo *LLMConvoRepository) GetLastConvoID() (string, error) {
	query := `
		SELECT conversation_id
		FROM conversations
		ORDER BY conversation_id DESC
		LIMIT 1
	`
	var convoID string
	err := repo.DB.QueryRow(query).Scan(&convoID)
	return convoID, err
}
