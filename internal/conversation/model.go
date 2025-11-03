package conversation

import "time"

type ConversationMessage struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	MessageRole    string    `json:"message_role"`
	MessageContent string    `json:"message_content"`
	MessageID      string    `json:"message_id,omitempty"`
	ContextData    string    `json:"context_data,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ChatRequest struct {
	PlatformUserID string `json:"platform_user_id"`
	Message        string `json:"message"`
	MessageID      string `json:"message_id"`
}

type ChatResponse struct {
	Response string `json:"response"`
}
