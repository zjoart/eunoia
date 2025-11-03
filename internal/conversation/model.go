package conversation

import "time"

type ConversationMessage struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	MessageRole    string    `json:"message_role"`
	MessageContent string    `json:"message_content"`
	ContextData    string    `json:"context_data,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ChatRequest struct {
	TelexUserID string `json:"telex_user_id"`
	Message     string `json:"message"`
}

type ChatResponse struct {
	Response  string `json:"response"`
	MessageID string `json:"message_id"`
}
