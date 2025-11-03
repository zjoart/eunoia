package conversation

import (
	"fmt"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/agent"
	"github.com/zjoart/eunoia/internal/checkin"
	"github.com/zjoart/eunoia/internal/reflection"
	"github.com/zjoart/eunoia/internal/user"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Service struct {
	repo           *Repository
	userRepo       *user.Repository
	checkInRepo    *checkin.Repository
	reflectionRepo *reflection.Repository
	geminiService  *agent.GeminiService
}

func NewService(
	repo *Repository,
	userRepo *user.Repository,
	checkInRepo *checkin.Repository,
	reflectionRepo *reflection.Repository,
	geminiService *agent.GeminiService,
) *Service {
	return &Service{
		repo:           repo,
		userRepo:       userRepo,
		checkInRepo:    checkInRepo,
		reflectionRepo: reflectionRepo,
		geminiService:  geminiService,
	}
}

func (s *Service) ProcessMessage(req *ChatRequest) (*ChatResponse, error) {
	logger.Info("processing chat message", logger.Fields{"telex_user_id": req.TelexUserID})

	if strings.TrimSpace(req.Message) == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	userRecord, err := s.userRepo.GetOrCreateUser(req.TelexUserID)
	if err != nil {
		logger.Error("failed to get or create user", logger.WithError(err))
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	userMessageID := fmt.Sprintf("%d", time.Now().UnixNano())
	userMessage := &ConversationMessage{
		ID:             userMessageID,
		UserID:         userRecord.ID,
		MessageRole:    "user",
		MessageContent: req.Message,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveMessage(userMessage); err != nil {
		logger.Warn("failed to save user message", logger.WithError(err))
	}

	context, err := s.buildUserContext(userRecord.ID)
	if err != nil {
		logger.Warn("failed to build user context", logger.WithError(err))
		context = ""
	}

	conversationHistory, err := s.repo.GetRecentMessages(userRecord.ID, 30)
	if err != nil {
		logger.Warn("failed to get conversation history", logger.WithError(err))
		conversationHistory = []*ConversationMessage{}
	}

	geminiHistory := s.convertToGeminiHistory(conversationHistory)

	systemPrompt := s.buildSystemPrompt(context)

	response, err := s.geminiService.GenerateContent(systemPrompt, req.Message, geminiHistory)
	if err != nil {
		logger.Error("failed to generate response", logger.WithError(err))
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	assistantMessageID := fmt.Sprintf("%d", time.Now().UnixNano()+1)
	assistantMessage := &ConversationMessage{
		ID:             assistantMessageID,
		UserID:         userRecord.ID,
		MessageRole:    "assistant",
		MessageContent: response,
		ContextData:    context,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveMessage(assistantMessage); err != nil {
		logger.Warn("failed to save assistant message", logger.WithError(err))
	}

	logger.Info("chat message processed successfully", logger.Fields{"message_id": assistantMessageID})

	return &ChatResponse{
		Response:  response,
		MessageID: assistantMessageID,
	}, nil
}

func (s *Service) buildSystemPrompt(userContext string) string {
	prompt := `You are Eunoia, a compassionate AI assistant focused on mental wellbeing and emotional support. 

Your role is to:
- Perform daily emotional check-ins with users
- Listen actively and provide empathetic responses
- Help users reflect on their emotions and experiences
- Offer supportive, non-judgmental guidance
- Recognize patterns in mood and emotional states
- Provide gentle suggestions for self-care when appropriate

Guidelines:
- Always be warm, empathetic, and supportive
- Never provide medical advice or diagnosis
- If someone is in crisis, encourage them to seek professional help
- Keep responses conversational and human-like
- Reference past conversations and patterns when relevant
- Validate feelings without minimizing concerns
- Use simple, clear language
- Less Than 150 words

`
	if userContext != "" {
		prompt += "\nUser Context:\n" + userContext
	}

	return prompt
}

func (s *Service) buildUserContext(userID string) (string, error) {
	var contextParts []string

	checkIns, err := s.checkInRepo.GetCheckInsByUserID(userID, 5)
	if err == nil && len(checkIns) > 0 {
		contextParts = append(contextParts, fmt.Sprintf("Recent check-ins: %d entries", len(checkIns)))
		if checkIns[0] != nil {
			contextParts = append(contextParts, fmt.Sprintf("Latest mood: %d/10 (%s)",
				checkIns[0].MoodScore, checkIns[0].MoodLabel))
		}
	}

	reflections, err := s.reflectionRepo.GetReflectionsByUserID(userID, 3)
	if err == nil && len(reflections) > 0 {
		contextParts = append(contextParts, fmt.Sprintf("Recent reflections: %d entries", len(reflections)))
		if reflections[0] != nil && reflections[0].Sentiment != "" {
			contextParts = append(contextParts, fmt.Sprintf("Latest sentiment: %s", reflections[0].Sentiment))
		}
	}

	stats, err := s.checkInRepo.GetCheckInStats(userID, 7)
	if err == nil && stats.TotalCheckIns > 0 {
		contextParts = append(contextParts, fmt.Sprintf("7-day mood average: %.1f/10", stats.AverageMoodScore))
		if stats.MoodTrend != "" && stats.MoodTrend != "new" {
			contextParts = append(contextParts, fmt.Sprintf("Mood trend: %s", stats.MoodTrend))
		}
	}

	if len(contextParts) == 0 {
		return "New user - no previous history", nil
	}

	return strings.Join(contextParts, "\n"), nil
}

func (s *Service) convertToGeminiHistory(messages []*ConversationMessage) []string {
	var history []string

	for _, msg := range messages {
		if len(history) >= 10 {
			break
		}

		history = append(history, msg.MessageContent)
	}

	return history
}

func (s *Service) GetConversationHistory(telexUserID string, limit int) ([]*ConversationMessage, error) {
	userRecord, err := s.userRepo.GetUserByTelexID(telexUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	messages, err := s.repo.GetConversationHistory(userRecord.ID, limit)
	if err != nil {
		return nil, err
	}

	reversedMessages := make([]*ConversationMessage, len(messages))
	for i, msg := range messages {
		reversedMessages[len(messages)-1-i] = msg
	}

	return reversedMessages, nil
}
