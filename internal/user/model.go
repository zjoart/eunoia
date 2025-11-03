package user

import "time"

type User struct {
	ID          string    `json:"id"`
	TelexUserID string    `json:"telex_user_id"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserPreferences struct {
	UserID            string    `json:"user_id"`
	ReminderTime      string    `json:"reminder_time,omitempty"`
	ReminderFrequency string    `json:"reminder_frequency"`
	PreferredTone     string    `json:"preferred_tone"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
