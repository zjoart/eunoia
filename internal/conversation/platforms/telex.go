package platforms

import (
	"errors"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/a2a"
)

type TelexPlatform struct{}

func NewTelexPlatform() *TelexPlatform {
	return &TelexPlatform{}
}

func (t *TelexPlatform) Name() string {
	return "telex"
}

func (t *TelexPlatform) ExtractMessage(parts []a2a.A2APart) string {
	// For Telex, we need to extract only the last text message
	// as it sends conversation history + the actual new message
	var lastText string

	for _, part := range parts {
		if part.Kind == "text" && part.Text != "" {
			lastText = part.Text
		} else if part.Kind == "data" && len(part.Data) > 0 {
			// Extract text from data parts (nested structure)
			for _, dataPart := range part.Data {
				if dataPart.Kind == "text" && dataPart.Text != "" {
					// Clean HTML tags if present
					text := strings.ReplaceAll(dataPart.Text, "<p>", "")
					text = strings.ReplaceAll(dataPart.Text, "</p>", "")
					lastText = text
				}
			}
		}
	}

	return strings.TrimSpace(lastText)
}

func (t *TelexPlatform) ExtractUserID(metadata map[string]interface{}) (string, error) {
	if userID, ok := metadata["telex_user_id"].(string); ok && userID != "" {
		return userID, nil
	}
	return "", errors.New("telex_user_id is required")
}

func (t *TelexPlatform) ExtractChannelID(metadata map[string]interface{}) (string, error) {
	if channelID, ok := metadata["telex_channel_id"].(string); ok && channelID != "" {
		return channelID, nil
	}
	return "", nil
}

func (t *TelexPlatform) ValidateRequest(req *a2a.A2ARequest) error {
	if req.Method != "message/send" {
		return errors.New("method not supported")
	}
	return nil
}

func (t *TelexPlatform) BuildResponse(id string, response *a2a.ChatResponse) *a2a.A2AResponse {
	return &a2a.A2AResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: a2a.A2AResult{
			Message: a2a.A2AMessageResult{
				Kind: "message",
				Role: "assistant",
				Parts: []a2a.A2APart{
					{
						Kind: "text",
						Text: response.Response,
					},
				},
				Metadata: map[string]interface{}{
					"agent":     "eunoia",
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				},
				MessageID: response.MessageID,
			},
		},
	}
}
