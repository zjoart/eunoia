package platforms

import (
	"errors"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/a2a"
)

type PlatformImpl struct {
	name string
}

func NewPlatform(name string) *PlatformImpl {
	return &PlatformImpl{name: name}
}

func (p *PlatformImpl) Name() string {
	return p.name
}

func (p *PlatformImpl) ExtractMessage(parts []a2a.A2APart) string {
	var lastText string

	for _, part := range parts {
		if part.Kind == "text" && part.Text != "" {
			lastText = part.Text
		} else if part.Kind == "data" && len(part.Data) > 0 {
			for _, dataPart := range part.Data {
				if dataPart.Kind == "text" && dataPart.Text != "" {
					text := strings.ReplaceAll(dataPart.Text, "<p>", "")
					text = strings.ReplaceAll(text, "</p>", "")
					lastText = text
				}
			}
		}
	}

	return strings.TrimSpace(lastText)
}

func (p *PlatformImpl) ExtractUserID(metadata map[string]interface{}) (string, error) {
	userIDKeys := []string{"platform_user_id", "telex_user_id", "user_id"}

	for _, key := range userIDKeys {
		if userID, ok := metadata[key].(string); ok && userID != "" {
			return userID, nil
		}
	}

	return "", errors.New("user_id is required in metadata")
}

func (p *PlatformImpl) ExtractChannelID(metadata map[string]interface{}) (string, error) {
	channelIDKeys := []string{"platform_channel_id", "telex_channel_id", "channel_id"}

	for _, key := range channelIDKeys {
		if channelID, ok := metadata[key].(string); ok && channelID != "" {
			return channelID, nil
		}
	}

	return "", nil
}

func (p *PlatformImpl) ValidateRequest(req *a2a.A2ARequest) error {
	if req.Method != "message/send" {
		return errors.New("method not supported")
	}
	return nil
}

func (p *PlatformImpl) BuildResponse(id string, response *a2a.ChatResponse) *a2a.A2AResponse {
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
