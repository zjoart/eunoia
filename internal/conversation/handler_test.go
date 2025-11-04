package conversation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zjoart/eunoia/internal/a2a"
	"github.com/zjoart/eunoia/internal/conversation/platforms"
)

func TestHandleA2AMessage_EmptyBody(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	req := httptest.NewRequest(http.MethodPost, "/a2a/agent/eunoia", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp a2a.A2AResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error response, got nil")
	}

	if resp.Error.Code != a2a.ParseError {
		t.Errorf("expected error code %d, got %d", a2a.ParseError, resp.Error.Code)
	}
}

func TestHandleA2AMessage_EmptyJSON(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	req := httptest.NewRequest(http.MethodPost, "/a2a/agent/eunoia", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp a2a.A2AResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error response, got nil")
	}

	if resp.Error.Code != a2a.InvalidRequest {
		t.Errorf("expected error code %d, got %d", a2a.InvalidRequest, resp.Error.Code)
	}

	if resp.Error.Message != "Invalid Request" {
		t.Errorf("expected message 'Invalid Request', got '%s'", resp.Error.Message)
	}
}

func TestHandleA2AMessage_InvalidJSONRPCVersion(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	payload := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "test-1",
		"method":  "process_message",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/agent/eunoia", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	var resp a2a.A2AResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error response")
	}

	if resp.Error.Code != a2a.InvalidRequest {
		t.Errorf("expected error code %d, got %d", a2a.InvalidRequest, resp.Error.Code)
	}
}

func TestHandleA2AMessage_MissingMessageContent(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	payload := a2a.A2ARequest{
		JSONRPC: "2.0",
		ID:      "test-1",
		Method:  "message/send",
		Params: a2a.A2AParams{
			Message: a2a.A2AMessage{
				Kind:      "message",
				Role:      "user",
				Parts:     []a2a.A2APart{},
				Metadata:  map[string]interface{}{"telex_user_id": "user-123"},
				MessageID: "msg-1",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/agent/eunoia", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	var resp a2a.A2AResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error response")
	}

	if resp.Error.Code != a2a.InvalidParams {
		t.Errorf("expected error code %d, got %d", a2a.InvalidParams, resp.Error.Code)
	}
}

func TestHandleA2AMessage_ValidRequest(t *testing.T) {
	mockService := &MockService{
		ProcessMessageFunc: func(req *ChatRequest) (*ChatResponse, error) {
			return &ChatResponse{Response: "Test response"}, nil
		},
	}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	payload := a2a.A2ARequest{
		JSONRPC: "2.0",
		ID:      "test-1",
		Method:  "message/send",
		Params: a2a.A2AParams{
			Message: a2a.A2AMessage{
				Kind: "message",
				Role: "user",
				Parts: []a2a.A2APart{
					{Kind: "text", Text: "Hello, I'm feeling great today!"},
				},
				Metadata:  map[string]interface{}{"telex_user_id": "user-123"},
				MessageID: "msg-1",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/a2a/agent/eunoia", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp a2a.A2AResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}

	if resp.Result.Task.ID == "" {
		t.Error("expected task ID, got empty string")
	}

	if resp.Result.Task.Status != "completed" {
		t.Errorf("expected task status 'completed', got '%s'", resp.Result.Task.Status)
	}

	if resp.Result.Message.Parts[0].Text != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", resp.Result.Message.Parts[0].Text)
	}
}

func TestHandleA2AMessage_WrongHTTPMethod(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	req := httptest.NewRequest(http.MethodGet, "/a2a/agent/eunoia", nil)
	w := httptest.NewRecorder()

	handler.HandleA2AMessage(w, req)

	var resp a2a.A2AResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error response")
	}

	if resp.Error.Code != a2a.InvalidRequest {
		t.Errorf("expected error code %d, got %d", a2a.InvalidRequest, resp.Error.Code)
	}
}

func TestHandleHealthCheck(t *testing.T) {
	mockService := &MockService{}
	platform := platforms.NewPlatform("telex")
	handler := NewHandler(mockService, platform)

	req := httptest.NewRequest(http.MethodGet, "/agent/health", nil)
	w := httptest.NewRecorder()

	handler.HandleHealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%v'", resp["status"])
	}

	if resp["agent"] != "eunoia" {
		t.Errorf("expected agent 'eunoia', got '%v'", resp["agent"])
	}
}

// MockService implements the service interface for testing
type MockService struct {
	ProcessMessageFunc func(*ChatRequest) (*ChatResponse, error)
}

func (m *MockService) ProcessMessage(req *ChatRequest) (*ChatResponse, error) {
	if m.ProcessMessageFunc != nil {
		return m.ProcessMessageFunc(req)
	}
	return &ChatResponse{Response: "mock response"}, nil
}

func (m *MockService) GetConversationHistory(platformUserID string, limit int) ([]*ConversationMessage, error) {
	return nil, nil
}
