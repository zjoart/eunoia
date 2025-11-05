package conversation

import (
	"encoding/json"
	"net/http"

	"github.com/zjoart/eunoia/internal/a2a"
	"github.com/zjoart/eunoia/internal/conversation/platforms"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Handler struct {
	service  ServiceInterface
	platform platforms.Platform
}

func NewHandler(service ServiceInterface, platform platforms.Platform) *Handler {
	return &Handler{
		service:  service,
		platform: platform,
	}
}

func (h *Handler) HandleA2AMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendA2AError(w, a2a.InvalidRequest, "Invalid Request", "method not allowed")
		return
	}

	var req a2a.A2ARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode A2A request", logger.WithError(err))
		h.sendA2AError(w, a2a.ParseError, "Parse error", "invalid JSON")
		return
	}

	// validate JSON-RPC 2.0 format
	if req.JSONRPC != "2.0" {
		h.sendA2AError(w, a2a.InvalidRequest, "Invalid Request", "jsonrpc version must be 2.0")
		return
	}

	// use the configured platform
	platform := h.platform
	platformName := platform.Name()

	// validate request with platform-specific logic
	if err := platform.ValidateRequest(&req); err != nil {
		h.sendA2AError(w, a2a.MethodNotFound, "Method not found", err.Error())
		return
	}

	// extract user ID using platform-specific logic
	userID, err := platform.ExtractUserID(req.Params.Message.Metadata)
	if err != nil {
		h.sendA2AError(w, a2a.InvalidParams, "Invalid params", err.Error())
		return
	}

	// extract channel ID
	channelID, _ := platform.ExtractChannelID(req.Params.Message.Metadata)

	// extract message content from parts
	messageText := platform.ExtractMessage(req.Params.Message.Parts)

	messageId := req.Params.Message.MessageID

	// extract conversation history from parts
	history := platform.ExtractHistory(req.Params.Message.Parts, messageId)

	logger.Info("processing A2A message", logger.Fields{
		"platform":   platformName,
		"user_id":    userID,
		"channel_id": channelID,
		"message_id": messageId,
	})

	if messageText == "" {
		h.sendA2AError(w, a2a.InvalidParams, "Invalid params", "message content is required")
		return
	}

	chatReq := &ChatRequest{
		PlatformUserID: userID,
		Message:        messageText,
		MessageID:      messageId,
	}

	chatResp, err := h.service.ProcessMessage(chatReq)
	if err != nil {
		logger.Error("failed to process message", logger.WithError(err))
		h.sendA2AError(w, a2a.InternalError, "Internal error", "failed to process message")
		return
	}

	// build platform-specific response with history
	response := platform.BuildResponse(req.ID, messageId, history, &a2a.ChatResponse{
		Response: chatResp.Response,
	})

	logger.Info("A2A message processed successfully", logger.Fields{
		"platform":   platformName,
		"message_id": messageId,
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errResp)
}
