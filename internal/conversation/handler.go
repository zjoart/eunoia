package conversation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zjoart/eunoia/internal/a2a"
	"github.com/zjoart/eunoia/internal/conversation/platforms"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Handler struct {
	service          *Service
	platformRegistry *platforms.PlatformRegistry
}

func NewHandler(service *Service, platformRegistry *platforms.PlatformRegistry) *Handler {
	return &Handler{
		service:          service,
		platformRegistry: platformRegistry,
	}
}

func (h *Handler) HandleA2AMessage(w http.ResponseWriter, r *http.Request) {
	logger.Info("received A2A message request", logger.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		h.sendA2AError(w, -32600, "Invalid Request", "method not allowed")
		return
	}

	var req a2a.A2ARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode A2A request", logger.WithError(err))
		h.sendA2AError(w, -32700, "Parse error", "invalid JSON")
		return
	}

	// Validate JSON-RPC 2.0 format
	if req.JSONRPC != "2.0" {
		h.sendA2AError(w, -32600, "Invalid Request", "jsonrpc version must be 2.0")
		return
	}

	// Determine platform from metadata or use default (telex)
	platformName := "telex" // default platform
	if platform, ok := req.Params.Message.Metadata["platform"].(string); ok {
		platformName = platform
	}

	platform, exists := h.platformRegistry.GetPlatform(platformName)
	if !exists {
		h.sendA2AError(w, -32601, "Method not found", "unsupported platform: "+platformName)
		return
	}

	// Validate request with platform-specific logic
	if err := platform.ValidateRequest(&req); err != nil {
		h.sendA2AError(w, -32601, "Method not found", err.Error())
		return
	}

	// Extract user ID using platform-specific logic
	userID, err := platform.ExtractUserID(req.Params.Message.Metadata)
	if err != nil {
		h.sendA2AError(w, -32602, "Invalid params", err.Error())
		return
	}

	// Extract channel ID (optional)
	channelID, _ := platform.ExtractChannelID(req.Params.Message.Metadata)

	// Extract message content from parts
	messageText := h.extractMessageText(req.Params.Message.Parts)

	logger.Info("processing A2A message", logger.Fields{
		"platform":   platformName,
		"user_id":    userID,
		"channel_id": channelID,
		"message_id": req.Params.Message.MessageID,
		"message":    messageText,
	})

	if messageText == "" {
		h.sendA2AError(w, -32602, "Invalid params", "message content is required")
		return
	}

	chatReq := &ChatRequest{
		TelexUserID: userID,
		Message:     messageText,
	}

	chatResp, err := h.service.ProcessMessage(chatReq)
	if err != nil {
		logger.Error("failed to process message", logger.WithError(err))
		h.sendA2AError(w, -32603, "Internal error", "failed to process message")
		return
	}

	// Build platform-specific response
	response := platform.BuildResponse(req.ID, &a2a.ChatResponse{
		Response:  chatResp.Response,
		MessageID: chatResp.MessageID,
	})

	logger.Info("A2A message processed successfully", logger.Fields{
		"platform":   platformName,
		"message_id": chatResp.MessageID,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"agent":   "eunoia",
		"service": "mental wellbeing assistant",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) extractMessageText(parts []a2a.A2APart) string {
	var messageText string

	for _, part := range parts {
		if part.Kind == "text" && part.Text != "" {
			messageText += part.Text + " "
		} else if part.Kind == "data" && len(part.Data) > 0 {
			// Extract text from data parts (nested structure)
			for _, dataPart := range part.Data {
				if dataPart.Kind == "text" && dataPart.Text != "" {
					// Clean HTML tags if present
					text := strings.ReplaceAll(dataPart.Text, "<p>", "")
					text = strings.ReplaceAll(text, "</p>", "")
					messageText += text + " "
				}
			}
		}
	}

	return strings.TrimSpace(messageText)
}

func (h *Handler) sendA2AError(w http.ResponseWriter, code int, message string, data interface{}) {
	logger.Error("sending A2A error response", logger.Fields{
		"code":    code,
		"message": message,
		"data":    data,
	})

	errResp := a2a.A2AResponse{
		JSONRPC: "2.0",
		Error: &a2a.A2AError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC 2.0 always returns 200 with error in body
	json.NewEncoder(w).Encode(errResp)
}
