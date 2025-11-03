package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/zjoart/eunoia/pkg/logger"
	"google.golang.org/api/option"
)

type GeminiService struct {
	apiKey string
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(apiKey string) *GeminiService {
	ctx := context.Background()

	if apiKey == "" {
		logger.Error("gemini api key is empty")
		return nil
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("failed to create gemini client", logger.WithError(err))
		return nil
	}

	modelName := "gemini-2.5-flash"
	model := client.GenerativeModel(modelName)

	model.SetTemperature(0.9)

	logger.Info("gemini service initialized", logger.Fields{
		"model": modelName,
	})

	return &GeminiService{
		apiKey: apiKey,
		client: client,
		model:  model,
	}
}

func (g *GeminiService) GenerateContent(systemPrompt string, userMessage string, conversationHistory []string) (string, error) {
	ctx := context.Background()

	// Build full prompt with history embedded in text
	var promptBuilder strings.Builder
	promptBuilder.WriteString(systemPrompt)
	promptBuilder.WriteString("\n\n")

	if len(conversationHistory) > 0 {
		promptBuilder.WriteString("Previous conversation:\n")
		for i, msg := range conversationHistory {
			role := "User"
			if i%2 == 1 {
				role = "Assistant"
			}
			promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, msg))
		}
		promptBuilder.WriteString("\n")
	}

	promptBuilder.WriteString("Current message:\n")
	promptBuilder.WriteString(userMessage)

	resp, err := g.model.GenerateContent(ctx, genai.Text(promptBuilder.String()))
	if err != nil {
		logger.Error("failed to generate content", logger.WithError(err))
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		logger.Error("no candidates in gemini response", logger.Fields{
			"prompt_feedback": resp.PromptFeedback,
		})
		return "", fmt.Errorf("no candidates in response - content may have been blocked")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		logger.Error("no content parts in gemini response", logger.Fields{
			"finish_reason": resp.Candidates[0].FinishReason,
		})
		return "", fmt.Errorf("no content in response")
	}

	var responseText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		responseText.WriteString(fmt.Sprintf("%v", part))
	}

	return responseText.String(), nil
}

func (g *GeminiService) AnalyzeSentiment(text string) (string, error) {

	prompt := fmt.Sprintf(`Analyze the sentiment of the following text and respond with only one word: "positive", "negative", "neutral", or "mixed".

Text: %s

Sentiment:`, text)

	sentiment, err := g.GenerateContent("You are a sentiment analysis assistant.", prompt, []string{})
	if err != nil {
		return "", err
	}

	sentiment = strings.TrimSpace(strings.ToLower(sentiment))
	return sentiment, nil
}

func (g *GeminiService) ExtractKeyThemes(text string) (string, error) {

	prompt := fmt.Sprintf(`Extract 3-5 key themes or topics from the following text. Return them as a comma-separated list.

Text: %s

Key themes:`, text)

	themes, err := g.GenerateContent("You are a text analysis assistant.", prompt, []string{})
	if err != nil {
		return "", err
	}

	themes = strings.TrimSpace(themes)
	return themes, nil
}

func (g *GeminiService) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}
