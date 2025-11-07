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
			name:    "with_user_context",
			context: "Recent check-ins: 3 entries\nLatest mood: 7/10 (content)",
			shouldHave: []string{
				"Eunoia",
				"empathetic",
				"Background context:",
				"Recent check-ins: 3 entries",
			},
			shouldNotHave: []string{},
		},
		{
			name:    "without_context",
			context: "",
			shouldHave: []string{
				"Eunoia",
				"empathetic",
				"Respond DIRECTLY",
				"genuine and conversational",
			},
			shouldNotHave: []string{"Background context:"},
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
		{MessageRole: "user", MessageContent: "I'm feeling stressed about my internship"},
		{MessageRole: "assistant", MessageContent: "It sounds like the internship is weighing on you. What specifically feels most stressful right now?"},
		{MessageRole: "user", MessageContent: "The pace is fast, tasks drop in a minute and the next we have to submit"},
		{MessageRole: "assistant", MessageContent: "That sounds incredibly overwhelming. The constant pressure with such quick turnarounds must feel exhausting. How does that tend to show up for you?"},
	}

	history := service.convertToGeminiHistory(messages)

	if len(history) != 4 {
		t.Errorf("expected 4 messages in history, got %d", len(history))
	}

	// check format includes role prefix
	if !strings.Contains(history[0], "User: I'm feeling stressed") {
		t.Errorf("expected first message to contain stress mention, got '%s'", history[0])
	}

	if !strings.Contains(history[1], "Eunoia:") {
		t.Errorf("expected second message to have 'Eunoia:' prefix, got '%s'", history[1])
	}

	// verify conversation flow is preserved
	if !strings.Contains(history[2], "User:") && !strings.Contains(history[2], "pace is fast") {
		t.Errorf("expected third message to contain follow-up about pace, got '%s'", history[2])
	}
}

func TestConvertToGeminiHistory_LimitTo10(t *testing.T) {
	service := &Service{}

	// simulate a long conversation (15 messages total)
	messages := make([]*ConversationMessage, 15)
	conversationTopics := []string{
		"I'm excited about my presentation tomorrow",
		"That's wonderful! What are you most looking forward to sharing?",
		"The new features I built for the project",
		"That sounds great! How are you feeling about it overall?",
		"A bit nervous but excited",
		"It's natural to feel both. Let's do a quick check-in.",
		"Okay, sounds good",
		"On a scale of 1-10, how would you rate your current mood?",
		"I'd say about a 7",
		"That's a solid place to be! What's contributing to that 7?",
		"The excitement is helping, but also feeling the pressure",
		"The pressure is real. How does it usually show up for you?",
		"I tend to overthink and get anxious about small details",
		"That's a really common experience. What helps you when that happens?",
		"Usually talking it through helps me gain perspective",
	}

	for i := 0; i < 15; i++ {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		messages[i] = &ConversationMessage{
			MessageRole:    role,
			MessageContent: conversationTopics[i],
		}
	}

	history := service.convertToGeminiHistory(messages)

	// Should only get last 10 messages (excludes first 5)
	if len(history) != 10 {
		t.Errorf("expected history limited to 10, got %d", len(history))
	}

	// First message in history should be the 6th message (index 5)
	// "It's natural to feel both. Let's do a quick check-in."
	if !strings.Contains(history[0], "natural to feel both") {
		t.Errorf("expected first message in history to be from position 5, got '%s'", history[0])
	}

	// last message should be the most recent
	if !strings.Contains(history[9], "talking it through") {
		t.Errorf("expected last message to be the final message, got '%s'", history[9])
	}
}

func TestBuildUserContext(t *testing.T) {

	t.Skip("TODO: implement mocking for repos to test various context scenarios")
}
