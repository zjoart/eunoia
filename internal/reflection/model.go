package reflection

import "time"

type Reflection struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Content    string    `json:"content"`
	Sentiment  string    `json:"sentiment"`
	KeyThemes  string    `json:"key_themes"`
	AIAnalysis string    `json:"ai_analysis"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateReflectionRequest struct {
	TelexUserID string `json:"telex_user_id"`
	Content     string `json:"content"`
}

type ReflectionResponse struct {
	Reflection *Reflection `json:"reflection"`
	Insights   string      `json:"insights"`
}
