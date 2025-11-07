package conversation

import (
	"fmt"
	"strings"
	"time"

	"github.com/zjoart/eunoia/internal/agent"
	"github.com/zjoart/eunoia/internal/checkin"
	"github.com/zjoart/eunoia/internal/reflection"
	"github.com/zjoart/eunoia/internal/user"
	"github.com/zjoart/eunoia/pkg/id"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Service struct {
	repo              *Repository
	userRepo          *user.Repository
	checkInRepo       *checkin.Repository
	reflectionRepo    *reflection.Repository
	checkInService    *checkin.Service
	reflectionService *reflection.Service
	geminiService     *agent.GeminiService
}

func NewService(
	repo *Repository,
	userRepo *user.Repository,
	checkInRepo *checkin.Repository,
	reflectionRepo *reflection.Repository,
	geminiService *agent.GeminiService,
) *Service {
	checkInService := checkin.NewService(checkInRepo, userRepo)
	reflectionService := reflection.NewService(reflectionRepo, userRepo, geminiService)

	return &Service{
		repo:              repo,
		userRepo:          userRepo,
		checkInRepo:       checkInRepo,
		reflectionRepo:    reflectionRepo,
		checkInService:    checkInService,
		reflectionService: reflectionService,
		geminiService:     geminiService,
	}
}

func (s *Service) ProcessMessage(req *ChatRequest) (*ChatResponse, error) {
	if strings.TrimSpace(req.Message) == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	userRecord, err := s.userRepo.GetOrCreateUser(req.PlatformUserID)
	if err != nil {
		logger.Error("failed to get or create user", logger.WithError(err))
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	userMessage := &ConversationMessage{
		ID:             id.Generate(),
		UserID:         userRecord.ID,
		MessageRole:    "user",
		MessageContent: req.Message,
		MessageID:      req.MessageID,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveMessage(userMessage); err != nil {
		logger.Warn("failed to save user message", logger.WithError(err))
	}

	s.detectAndHandleIntents(req.PlatformUserID, req.Message)

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

	assistantMessage := &ConversationMessage{
		ID:             id.Generate(),
		UserID:         userRecord.ID,
		MessageRole:    "assistant",
		MessageContent: response,
		MessageID:      req.MessageID,
		ContextData:    context,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveMessage(assistantMessage); err != nil {
		logger.Warn("failed to save assistant message", logger.WithError(err))
	}

	return &ChatResponse{
		Response: response,
	}, nil
}

func (s *Service) buildSystemPrompt(userContext string) string {
	prompt := `You are Eunoia, a warm and empathetic mental wellbeing companion.

Core principles:
- Respond DIRECTLY to what the user just shared - don't deflect or redirect
- Build on the conversation naturally, don't restart each time
- Show you're listening by referencing what they said
- Be genuine and conversational, not formulaic
- Match their energy - if they're sharing, engage; if they're asking, answer

Your style:
- Natural and warm, like a trusted friend
- Acknowledge their emotions authentically
- Ask ONE good follow-up question when relevant (not every time)
- Offer gentle reflections or perspective when helpful
- Keep it brief (2-3 sentences usually enough)
- NO generic phrases like "Let's check in" or "I'm here to help" - just engage naturally

Important:
- Read the conversation history to maintain continuity
- Don't repeat yourself or use the same opening patterns
- If they express emotion, acknowledge SPECIFICALLY what they said
- If they ask for help, provide actual guidance
- Crisis indicators should prompt gentle encouragement for professional support

You're a companion on their journey, not a script following a checklist.
`
	if userContext != "" {
		prompt += "\n\nBackground context:\n" + userContext + "\n\nUse this to inform your responses, but stay focused on the current conversation."
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

	// get last 5 message pairs (10 messages total) for better context
	startIndex := 0
	if len(messages) > 10 {
		startIndex = len(messages) - 10
	}

	for i := startIndex; i < len(messages); i++ {
		msg := messages[i]
		// format with role prefix for clearer context
		rolePrefix := "User"
		if msg.MessageRole == "assistant" {
			rolePrefix = "Eunoia"
		}
		history = append(history, fmt.Sprintf("%s: %s", rolePrefix, msg.MessageContent))
	}

	return history
}

func (s *Service) GetConversationHistory(platformUserID string, limit int) ([]*ConversationMessage, error) {
	userRecord, err := s.userRepo.GetUserByPlatformID(platformUserID)
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

func (s *Service) detectAndHandleIntents(platformUserID, message string) {
	messageLower := strings.ToLower(message)
	messageLen := len(strings.Fields(message))

	moodScore, moodLabel := s.detectMoodIntent(messageLower)
	if moodScore > 0 {
		checkInReq := &checkin.CreateCheckInRequest{
			PlatformUserID: platformUserID,
			MoodScore:      moodScore,
			MoodLabel:      moodLabel,
			Description:    message,
		}
		if _, err := s.checkInService.CreateCheckIn(checkInReq); err != nil {
			logger.Warn("failed to auto-create check-in", logger.WithError(err))
		} else {
			logger.Info("auto-created check-in from conversation", logger.Fields{"mood": moodLabel})
		}
	}

	if messageLen >= 15 && s.isReflectionIntent(messageLower) {
		reflectionReq := &reflection.CreateReflectionRequest{
			PlatformUserID: platformUserID,
			Content:        message,
		}
		if _, err := s.reflectionService.CreateReflection(reflectionReq); err != nil {
			logger.Warn("failed to auto-create reflection", logger.WithError(err))
		} else {
			logger.Info("auto-created reflection from conversation")
		}
	}
}

func (s *Service) detectMoodIntent(messageLower string) (int, string) {
	moodPatterns := map[string]struct {
		score int
		label string
	}{
		"amazing":    {9, "joyful"},
		"fantastic":  {9, "joyful"},
		"wonderful":  {9, "joyful"},
		"great":      {8, "happy"},
		"good":       {7, "content"},
		"happy":      {8, "happy"},
		"joyful":     {9, "joyful"},
		"excited":    {8, "happy"},
		"okay":       {5, "neutral"},
		"fine":       {6, "content"},
		"alright":    {5, "neutral"},
		"meh":        {4, "low"},
		"tired":      {4, "low"},
		"stressed":   {3, "anxious"},
		"anxious":    {3, "anxious"},
		"worried":    {3, "anxious"},
		"sad":        {3, "sad"},
		"down":       {3, "sad"},
		"depressed":  {2, "very low"},
		"terrible":   {2, "very low"},
		"awful":      {2, "very low"},
		"horrible":   {2, "very low"},
		"struggling": {3, "struggling"},
	}

	for keyword, mood := range moodPatterns {
		if strings.Contains(messageLower, "feel "+keyword) ||
			strings.Contains(messageLower, "feeling "+keyword) ||
			strings.Contains(messageLower, "i'm "+keyword) ||
			strings.Contains(messageLower, "i am "+keyword) {
			return mood.score, mood.label
		}
	}

	return 0, ""
}

func (s *Service) isReflectionIntent(messageLower string) bool {
	reflectionIndicators := []string{
		"today i", "i've been thinking", "i realized", "i noticed",
		"looking back", "i feel like", "lately i've", "i've noticed",
		"been feeling", "it's been", "struggling with", "grateful for",
		"thinking about", "i wonder", "reflecting on",
	}

	for _, indicator := range reflectionIndicators {
		if strings.Contains(messageLower, indicator) {
			return true
		}
	}

	return false
}
