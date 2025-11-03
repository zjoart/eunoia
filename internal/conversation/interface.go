package conversation

// ServiceInterface defines the methods needed by the handler
type ServiceInterface interface {
	ProcessMessage(req *ChatRequest) (*ChatResponse, error)
	GetConversationHistory(platformUserID string, limit int) ([]*ConversationMessage, error)
}
