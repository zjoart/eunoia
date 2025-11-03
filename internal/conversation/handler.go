package conversation

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/zjoart/eunoia/pkg/logger"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type A2ARequest struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"`
	Method  string    `json:"method"`
	Params  A2AParams `json:"params"`
}

type A2AParams struct {
	Message       A2AMessage `json:"message"`
	Configuration A2AConfig  `json:"configuration"`
}

type A2AMessage struct {
	Kind      string                 `json:"kind"`
	Role      string                 `json:"role"`
	Parts     []A2APart              `json:"parts"`
	Metadata  map[string]interface{} `json:"metadata"`
	MessageID string                 `json:"messageId"`
}

type A2APart struct {
	Kind string    `json:"kind"`
	Text string    `json:"text,omitempty"`
	Data []A2APart `json:"data,omitempty"`
}

type A2AConfig struct {
	AcceptedOutputModes    []string            `json:"acceptedOutputModes"`
	HistoryLength          int                 `json:"historyLength"`
	PushNotificationConfig A2APushNotification `json:"pushNotificationConfig"`
	Blocking               bool                `json:"blocking"`
}

type A2APushNotification struct {
	URL            string                 `json:"url"`
	Token          string                 `json:"token"`
	Authentication map[string]interface{} `json:"authentication"`
}

type A2AResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"`
	Result  A2AResult `json:"result,omitempty"`
	Error   *A2AError `json:"error,omitempty"`
}

type A2AResult struct {
	Message A2AMessageResult `json:"message"`
}

type A2AMessageResult struct {
	Kind      string                 `json:"kind"`
	Role      string                 `json:"role"`
	Parts     []A2APart              `json:"parts"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	MessageID string                 `json:"messageId"`
}

type A2AError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

	var req A2ARequest
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

	if req.Method != "message/send" {
		h.sendA2AError(w, -32601, "Method not found", "method not supported")
		return
	}

	// Extract message content from parts
	messageText := h.extractMessageText(req.Params.Message.Parts)
	userID := req.Params.Message.Metadata["telex_user_id"].(string)

	logger.Info("processing A2A message", logger.Fields{
		"user_id":    userID,
		"channel_id": req.Params.Message.Metadata["telex_channel_id"],
		"message_id": req.Params.Message.MessageID,
		"message":    messageText,
	})

	if messageText == "" {
		h.sendA2AError(w, -32602, "Invalid params", "message content is required")
		return
	}

	if userID == "" {
		h.sendA2AError(w, -32602, "Invalid params", "telex_user_id is required")
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

	// Build A2A response
	response := A2AResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: A2AResult{
			Message: A2AMessageResult{
				Kind: "message",
				Role: "assistant",
				Parts: []A2APart{
					{
						Kind: "text",
						Text: chatResp.Response,
					},
				},
				Metadata: map[string]interface{}{
					"agent":     "eunoia",
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				},
				MessageID: chatResp.MessageID,
			},
		},
	}

	logger.Info("A2A message processed successfully", logger.Fields{
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

func (h *Handler) extractMessageText(parts []A2APart) string {
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

	errResp := A2AResponse{
		JSONRPC: "2.0",
		Error: &A2AError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC 2.0 always returns 200 with error in body
	json.NewEncoder(w).Encode(errResp)
}
