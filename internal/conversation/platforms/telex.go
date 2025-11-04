package platforms

import (
	"errors"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/a2a"
	"github.com/zjoart/eunoia/internal/pkg/id"
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

	// Iterate through all parts to find the last text message
	for _, part := range parts {
		if part.Kind == "text" && part.Text != "" {
			lastText = part.Text
		} else if part.Kind == "data" && len(part.Data) > 0 {
			// Check data parts for the last message
			for _, dataPart := range part.Data {
				if dataPart.Kind == "text" && dataPart.Text != "" {
					text := strings.ReplaceAll(dataPart.Text, "<p>", "")
					text = strings.ReplaceAll(text, "</p>", "")
					text = strings.TrimSpace(text)
					if text != "" {
						lastText = text
					}
				}
			}
		}
	}

	return strings.TrimSpace(lastText)
} // ExtractHistory extracts conversation history from the data parts
func (p *PlatformImpl) ExtractHistory(parts []a2a.A2APart, currentMessageID string) []a2a.A2AMessageResult {
	var history []a2a.A2AMessageResult

	for _, part := range parts {
		if part.Kind == "data" && len(part.Data) > 0 {
			// Each pair of texts in data represents a conversation turn
			// Odd indices are typically user messages, even are agent responses
			for i, dataPart := range part.Data {
				if dataPart.Kind == "text" && dataPart.Text != "" {
					text := strings.ReplaceAll(dataPart.Text, "<p>", "")
					text = strings.ReplaceAll(text, "</p>", "")
					text = strings.TrimSpace(text)

					if text == "" {
						continue
					}

					// Determine role based on pattern (this is a simple heuristic)
					role := "user"
					if i%2 == 1 {
						role = "agent"
					}

					history = append(history, a2a.A2AMessageResult{
						MessageID: currentMessageID,
						Role:      role,
						TaskID:    id.Generate(),
						Parts: []a2a.A2APart{
							{
								Kind: "text",
								Text: text,
							},
						},
						Kind: "message",
					})
				}
			}
		}
	}

	return history
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

func (p *PlatformImpl) BuildResponse(requestId, messageID string, history []a2a.A2AMessageResult, response *a2a.ChatResponse) *a2a.A2AResponse {
	taskID := id.Generate()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Create the new agent response message (using messageID from request)
	newMessage := a2a.A2AMessageResult{
		MessageID: messageID,
		Role:      "agent",
		Parts: []a2a.A2APart{
			{
				Kind: "text",
				Text: response.Response,
			},
		},
		Kind:   "message",
		TaskID: taskID,
		Metadata: map[string]interface{}{
			"agent": "eunoia",
		},
	}

	// Append the new agent response to history
	updatedHistory := append(history, newMessage)

	return &a2a.A2AResponse{
		JSONRPC: "2.0",
		ID:      requestId,
		Result: a2a.A2AResult{
			ID:        taskID,
			ContextID: id.Generate(),
			Status: a2a.A2ATaskStatus{
				State:     "completed",
				Timestamp: timestamp,
				Message:   newMessage,
			},
			Artifacts: []a2a.A2AArtifact{},
			History:   updatedHistory,
			Kind:      "task",
		},
	}
}
