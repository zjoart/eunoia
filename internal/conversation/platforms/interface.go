package platforms

import (
	"github.com/zjoart/eunoia/internal/a2a"
)

type Platform interface {
	Name() string
	ExtractUserID(metadata map[string]interface{}) (string, error)
	ExtractChannelID(metadata map[string]interface{}) (string, error)
	ExtractMessage(parts []a2a.A2APart) string
	ValidateRequest(req *a2a.A2ARequest) error
	BuildResponse(id string, response *a2a.ChatResponse) *a2a.A2AResponse
}
