package user

import "time"

type User struct {
	ID             string    `json:"id"`
	PlatformUserID string    `json:"platform_user_id"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
