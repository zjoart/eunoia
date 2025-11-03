package conversation

import (
	"encoding/json"
	"net/http"

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
	Message   string                 `json:"message"`
	UserID    string                 `json:"userId"`
	ChannelID string                 `json:"channelId"`
	MessageID string                 `json:"messageId"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type A2AResponse struct {
	Response  string                 `json:"response"`
	MessageID string                 `json:"messageId,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (h *Handler) HandleA2AMessage(w http.ResponseWriter, r *http.Request) {
	logger.Info("received A2A message request", logger.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req A2ARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode A2A request", logger.WithError(err))
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("processing A2A message", logger.Fields{
		"user_id":    req.UserID,
		"channel_id": req.ChannelID,
		"message_id": req.MessageID,
	})

	if req.Message == "" {
		h.sendError(w, "message field is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		h.sendError(w, "userId field is required", http.StatusBadRequest)
		return
	}

	chatReq := &ChatRequest{
		TelexUserID: req.UserID,
		Message:     req.Message,
	}

	chatResp, err := h.service.ProcessMessage(chatReq)
	if err != nil {
		logger.Error("failed to process message", logger.WithError(err))
		h.sendError(w, "failed to process message", http.StatusInternalServerError)
		return
	}

	response := A2AResponse{
		Response:  chatResp.Response,
		MessageID: chatResp.MessageID,
		Metadata: map[string]interface{}{
			"agent":     "eunoia",
			"timestamp": req.Timestamp,
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

func (h *Handler) sendError(w http.ResponseWriter, message string, status int) {
	logger.Error("sending error response", logger.Fields{
		"message": message,
		"status":  status,
	})

	errResp := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Status:  status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}
