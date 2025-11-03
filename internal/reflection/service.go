package reflection

import (
	"fmt"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/agent"
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
	logger.Info("processing reflection request", logger.Fields{"telex_user_id": req.TelexUserID})

	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("reflection content cannot be empty")
	}

	userRecord, err := s.userRepo.GetOrCreateUser(req.TelexUserID)
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
		ID:         fmt.Sprintf("%d", time.Now().UnixNano()),
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

	logger.Info("reflection created successfully", logger.Fields{"reflection_id": reflection.ID})
	return reflection, nil
}

func (s *Service) GetReflectionHistory(telexUserID string, limit int) ([]*Reflection, error) {
	userRecord, err := s.userRepo.GetUserByTelexID(telexUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.repo.GetReflectionsByUserID(userRecord.ID, limit)
}

func (s *Service) generateReflectionAnalysis(content, sentiment, themes string) (string, error) {
	systemPrompt := `You are a compassionate mental wellbeing assistant. Provide a brief, supportive analysis of the user's reflection. 
Keep your response under 100 words. Be empathetic, non-judgmental, and offer gentle insights or validation.`

	userPrompt := fmt.Sprintf(`User's reflection: %s
Detected sentiment: %s
Key themes: %s

Provide a supportive response:`, content, sentiment, themes)

	analysis, err := s.geminiService.GenerateContent(systemPrompt, userPrompt, nil)
	if err != nil {
		return "", err
	}

	return analysis, nil
}
