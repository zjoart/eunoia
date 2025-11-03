package platforms

import (
	"github.com/zjoart/eunoia/internal/a2a"
)

type Platform interface {
	Name() string
	ExtractUserID(metadata map[string]interface{}) (string, error)
	ExtractChannelID(metadata map[string]interface{}) (string, error)
	ValidateRequest(req *a2a.A2ARequest) error
	BuildResponse(id string, response *a2a.ChatResponse) *a2a.A2AResponse
}

type PlatformRegistry struct {
	platforms map[string]Platform
}

func NewPlatformRegistry() *PlatformRegistry {
	return &PlatformRegistry{
		platforms: make(map[string]Platform),
	}
}

func (pr *PlatformRegistry) Register(platform Platform) {
	pr.platforms[platform.Name()] = platform
}

func (pr *PlatformRegistry) GetPlatform(name string) (Platform, bool) {
	platform, exists := pr.platforms[name]
	return platform, exists
}

func (pr *PlatformRegistry) GetAllPlatforms() map[string]Platform {
	return pr.platforms
}
