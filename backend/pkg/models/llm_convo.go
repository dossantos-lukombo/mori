package models

type Conversation struct {
	UserID          string `json:"user_id"`
	ConversationID  string `json:"conversation_id"`
	UserRequest     string `json:"user_request"`
	LLMResponse     string `json:"llm_response"`
	NewConversation bool   `json:"new_conversation"`
}

type LLMConvoRepository interface {
	SaveConvo(Conversation) error
	GetAllConvo(Conversation) ([]Conversation, error)
	GetLastConvoID() (string, error)
	GetLastConvo() (Conversation, error)
	// get all for specific chat

}
