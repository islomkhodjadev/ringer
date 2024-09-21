package database

import (
	"database/sql"
	"fmt"
	"time"

	"ringer/models"
)

// CreateConversationIfNotExists checks if a conversation exists for the user,
// creates one if it doesn't, and returns the conversation ID.
func CreateConversationIfNotExists(userID int) (int, error) {
	var conversationID int

	// Check if a conversation already exists for the user
	query := `
		SELECT id FROM conversations
		WHERE user_id = $1
	`
	err := DB.QueryRow(query, userID).Scan(&conversationID)

	// If a conversation exists, return it
	if err == nil {
		return conversationID, nil
	}

	// If no conversation exists, create a new one
	if err == sql.ErrNoRows {
		insertQuery := `
			INSERT INTO conversations (user_id, created_at)
			VALUES ($1, NOW())
			RETURNING id
		`
		err = DB.QueryRow(insertQuery, userID).Scan(&conversationID)
		if err != nil {
			return 0, fmt.Errorf("error creating conversation: %v", err)
		}
		return conversationID, nil
	}

	return 0, fmt.Errorf("error checking conversation existence: %v", err)
}

// GetConversationByUser fetches the conversation for a specific user.
func GetConversationByUser(userID int) (models.Conversation, error) {
	var conversation models.Conversation

	query := `
		SELECT id, user_id, created_at FROM conversations
		WHERE user_id = $1
	`
	err := DB.QueryRow(query, userID).Scan(&conversation.ID, &conversation.UserID, &conversation.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return conversation, fmt.Errorf("no conversation found for user")
		}
		return conversation, fmt.Errorf("error fetching conversation: %v", err)
	}

	return conversation, nil
}

// GetAllMessages fetches all messages for a specific conversation.
func GetAllMessages(conversationID int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, message, is_user_message, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Message, &msg.IsUserMessage, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// AddNewMessage adds a new message to a specific conversation.
func AddNewMessage(conversationID int, message string, isUserMessage bool) error {
	query := `
		INSERT INTO messages (conversation_id, message, is_user_message, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := DB.Exec(query, conversationID, message, isUserMessage)
	if err != nil {
		return fmt.Errorf("error adding new message: %v", err)
	}
	return nil
}

// ExampleUsage demonstrates how the conversation functions can be used together.
func ExampleUsage(userID int) {
	// Step 1: Create a new conversation if it doesn't already exist
	conversationID, err := CreateConversationIfNotExists(userID)
	if err != nil {
		fmt.Println("Error creating or getting conversation:", err)
		return
	}
	fmt.Println("Conversation ID:", conversationID)

	// Step 2: Add a user message to the conversation
	err = AddNewMessage(conversationID, "Hello AI, how are you?", true) // User message
	if err != nil {
		fmt.Println("Error adding user message:", err)
		return
	}

	// Step 3: Add an AI message to the conversation
	err = AddNewMessage(conversationID, "I am doing great! How can I assist you today?", false) // AI message
	if err != nil {
		fmt.Println("Error adding AI message:", err)
		return
	}

	// Step 4: Get all messages for the conversation
	messages, err := GetAllMessages(conversationID)
	if err != nil {
		fmt.Println("Error fetching conversation messages:", err)
		return
	}

	// Step 5: Display the conversation history
	for _, msg := range messages {
		sender := "AI"
		if msg.IsUserMessage {
			sender = "User"
		}
		fmt.Printf("[%s] %s: %s\n", msg.CreatedAt.Format("2006-01-02 15:04:05"), sender, msg.Message)
	}
}

func GetConversationAsString(conversationID int) (string, error) {
	query := `
		SELECT message, is_user_message, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, conversationID)
	if err != nil {
		return "", fmt.Errorf("error fetching messages: %v", err)
	}
	defer rows.Close()

	var conversationString string

	for rows.Next() {
		var message string
		var isUserMessage bool
		var createdAt time.Time

		err := rows.Scan(&message, &isUserMessage, &createdAt)
		if err != nil {
			return "", fmt.Errorf("error scanning message row: %v", err)
		}

		if isUserMessage {
			conversationString += fmt.Sprintf("user: %s\n", message)
		} else {
			conversationString += fmt.Sprintf("ai: %s\n", message)
		}
	}

	// Return the combined conversation string
	return conversationString, nil
}
