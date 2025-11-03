package checkin

import "time"

type EmotionalCheckIn struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	MoodScore   int       `json:"mood_score"`
	MoodLabel   string    `json:"mood_label"`
	Description string    `json:"description"`
	CheckInDate time.Time `json:"check_in_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateCheckInRequest struct {
	PlatformUserID string `json:"platform_user_id"`
	MoodScore      int    `json:"mood_score"`
	MoodLabel      string `json:"mood_label"`
	Description    string `json:"description"`
}

type CheckInResponse struct {
	CheckIn *EmotionalCheckIn `json:"check_in"`
	Message string            `json:"message"`
}

type CheckInStats struct {
	AverageMoodScore float64           `json:"average_mood_score"`
	TotalCheckIns    int               `json:"total_check_ins"`
	LastCheckIn      *EmotionalCheckIn `json:"last_check_in,omitempty"`
	MoodTrend        string            `json:"mood_trend"`
}
