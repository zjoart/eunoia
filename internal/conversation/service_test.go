package conversation

import (
	"strings"
	"testing"
)

func TestDetectMoodIntent_HappyMoods(t *testing.T) {
	service := &Service{}

	tests := []struct {
		message       string
		expectedScore int
		expectedLabel string
	}{
		{"I'm feeling amazing today", 9, "joyful"},
		{"feeling great about my progress", 8, "happy"},
		{"I feel happy", 8, "happy"},
		{"I'm good", 7, "content"},
		{"feeling excited about tomorrow", 8, "happy"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			score, label := service.detectMoodIntent(strings.ToLower(tt.message))
			if score != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, score)
			}
			if label != tt.expectedLabel {
				t.Errorf("expected label %s, got %s", tt.expectedLabel, label)
			}
		})
	}
}

func TestDetectMoodIntent_LowMoods(t *testing.T) {
	service := &Service{}

	tests := []struct {
		message       string
		expectedScore int
		expectedLabel string
	}{
		{"I'm feeling terrible", 2, "very low"},
		{"feeling sad today", 3, "sad"},
		{"I feel anxious about work", 3, "anxious"},
		{"feeling stressed", 3, "anxious"},
		{"I'm depressed", 2, "very low"},
		{"feeling down", 3, "sad"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			score, label := service.detectMoodIntent(strings.ToLower(tt.message))
			if score != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, score)
			}
			if label != tt.expectedLabel {
				t.Errorf("expected label %s, got %s", tt.expectedLabel, label)
			}
		})
	}
}

func TestDetectMoodIntent_NoMood(t *testing.T) {
	service := &Service{}

	tests := []string{
		"Hello, how are you?",
		"I like pizza",
		"What's the weather today?",
		"Can you help me?",
	}

	for _, message := range tests {
		t.Run(message, func(t *testing.T) {
			score, label := service.detectMoodIntent(strings.ToLower(message))
			if score != 0 {
				t.Errorf("expected no mood detected (score 0), got score %d", score)
			}
			if label != "" {
				t.Errorf("expected no label, got %s", label)
			}
		})
	}
}

func TestIsReflectionIntent(t *testing.T) {
	service := &Service{}

	tests := []struct {
		message  string
		expected bool
	}{
		{"today i realized that i need to take better care of myself", true},
		{"i've been thinking about my career goals lately", true},
		{"looking back, i can see how much i've grown", true},
		{"lately i've been feeling more confident", true},
		{"i noticed that i'm more patient now", true},
		{"reflecting on my journey so far", true},
		{"i wonder what the future holds for me", true},
		{"Hello", false},
		{"How are you?", false},
		{"I'm happy", false},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := service.isReflectionIntent(strings.ToLower(tt.message))
			if result != tt.expected {
				t.Errorf("expected %v, got %v for message: %s", tt.expected, result, tt.message)
			}
		})
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name          string
		context       string
		shouldHave    []string
		shouldNotHave []string
	}{
		{
			name:    "with user context",
			context: "Recent check-ins: 5 entries\nLatest mood: 7/10 (content)",
			shouldHave: []string{
				"Eunoia",
				"empathetic",
				"Context about this person:",
				"Recent check-ins",
				"Latest mood",
			},
			shouldNotHave: []string{},
		},
		{
			name:    "without context",
			context: "",
			shouldHave: []string{
				"Eunoia",
				"empathetic",
				"Listen with genuine curiosity",
			},
			shouldNotHave: []string{"Context about this person:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := service.buildSystemPrompt(tt.context)

			for _, phrase := range tt.shouldHave {
				if !strings.Contains(prompt, phrase) {
					t.Errorf("expected prompt to contain '%s'", phrase)
				}
			}

			for _, phrase := range tt.shouldNotHave {
				if strings.Contains(prompt, phrase) {
					t.Errorf("expected prompt NOT to contain '%s'", phrase)
				}
			}
		})
	}
}

func TestConvertToGeminiHistory(t *testing.T) {
	service := &Service{}

	messages := []*ConversationMessage{
		{MessageContent: "Hello"},
		{MessageContent: "Hi there!"},
		{MessageContent: "How are you?"},
		{MessageContent: "I'm doing well"},
	}

	history := service.convertToGeminiHistory(messages)

	if len(history) != 4 {
		t.Errorf("expected 4 messages in history, got %d", len(history))
	}

	if history[0] != "Hello" {
		t.Errorf("expected first message 'Hello', got '%s'", history[0])
	}
}

func TestConvertToGeminiHistory_LimitTo10(t *testing.T) {
	service := &Service{}

	messages := make([]*ConversationMessage, 15)
	for i := 0; i < 15; i++ {
		messages[i] = &ConversationMessage{
			MessageContent: "Message " + string(rune('A'+i)),
		}
	}

	history := service.convertToGeminiHistory(messages)

	if len(history) != 10 {
		t.Errorf("expected history limited to 10, got %d", len(history))
	}
}

func TestBuildUserContext(t *testing.T) {

	t.Skip("TODO: implement mocking for repos to test various context scenarios")
}
