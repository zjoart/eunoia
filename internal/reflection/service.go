package reflection

import (
	"fmt"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/agent"
	"github.com/zjoart/eunoia/internal/pkg/id"
	"github.com/zjoart/eunoia/internal/user"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Service struct {
	repo          *Repository
	userRepo      *user.Repository
	geminiService *agent.GeminiService
}

func NewService(repo *Repository, userRepo *user.Repository, geminiService *agent.GeminiService) *Service {
	return &Service{
		repo:          repo,
		userRepo:      userRepo,
		geminiService: geminiService,
	}
}

func (s *Service) CreateReflection(req *CreateReflectionRequest) (*Reflection, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("reflection content cannot be empty")
	}

	userRecord, err := s.userRepo.GetOrCreateUser(req.PlatformUserID)
	if err != nil {
		logger.Error("failed to get or create user", logger.WithError(err))
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	sentiment, err := s.geminiService.AnalyzeSentiment(req.Content)
	if err != nil {
		logger.Warn("failed to analyze sentiment", logger.WithError(err))
		sentiment = "unknown"
	}

	keyThemes, err := s.geminiService.ExtractKeyThemes(req.Content)
	if err != nil {
		logger.Warn("failed to extract key themes", logger.WithError(err))
		keyThemes = ""
	}

	aiAnalysis, err := s.generateReflectionAnalysis(req.Content, sentiment, keyThemes)
	if err != nil {
		logger.Warn("failed to generate AI analysis", logger.WithError(err))
		aiAnalysis = "Analysis unavailable at this time."
	}

	reflection := &Reflection{
		ID:         id.Generate(),
		UserID:     userRecord.ID,
		Content:    req.Content,
		Sentiment:  strings.TrimSpace(sentiment),
		KeyThemes:  strings.TrimSpace(keyThemes),
		AIAnalysis: aiAnalysis,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateReflection(reflection); err != nil {
		return nil, fmt.Errorf("failed to create reflection: %w", err)
	}

	return reflection, nil
}

func (s *Service) GetReflectionHistory(platformUserID string, limit int) ([]*Reflection, error) {
	userRecord, err := s.userRepo.GetUserByPlatformID(platformUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.repo.GetReflectionsByUserID(userRecord.ID, limit)
}

func (s *Service) generateReflectionAnalysis(content, sentiment, themes string) (string, error) {
	systemPrompt := `You are a thoughtful companion helping someone process their inner experience.

Respond with warmth and insight:
- Acknowledge what stands out in their reflection
- Notice patterns or connections they might not see
- Validate the complexity of their feelings
- Offer a gentle perspective or question for further reflection
- Keep it brief (under 80 words) and genuine`

	userPrompt := fmt.Sprintf(`They reflected: "%s"

The emotional tone seems %s, touching on: %s

Offer a brief, supportive response that honors their experience:`, content, sentiment, themes)

	analysis, err := s.geminiService.GenerateContent(systemPrompt, userPrompt, nil)
	if err != nil {
		return "", err
	}

	return analysis, nil
}
