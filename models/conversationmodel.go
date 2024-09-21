package models

import "time"

// Conversation represents a conversation between a user and the AI
type Conversation struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	Message        string    `json:"message"`
	IsUserMessage  bool      `json:"is_user_message"`
	CreatedAt      time.Time `json:"created_at"`
}
