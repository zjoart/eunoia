package platforms

import (
	"github.com/zjoart/eunoia/internal/a2a"
)

type Platform interface {
	Name() string
	ExtractUserID(metadata map[string]interface{}) (string, error)
	ExtractChannelID(metadata map[string]interface{}) (string, error)
	ExtractMessage(parts []a2a.A2APart) string
	ExtractHistory(parts []a2a.A2APart, currentMessageID string) []a2a.A2AMessageResult
	ValidateRequest(req *a2a.A2ARequest) error
	BuildResponse(requestId, messageId string, history []a2a.A2AMessageResult, response *a2a.ChatResponse) *a2a.A2AResponse
}
