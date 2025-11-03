package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zjoart/eunoia/internal/pkg/id"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(user *User) error {
	query := `INSERT INTO users (id, platform_user_id, username, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.ID, user.PlatformUserID, user.Username, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByPlatformID(platformUserID string) (*User, error) {
	query := `SELECT id, platform_user_id, username, created_at, updated_at
			  FROM users
			  WHERE platform_user_id = ?`

	user := &User{}
	err := r.db.QueryRow(query, platformUserID).Scan(
		&user.ID, &user.PlatformUserID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetOrCreateUser(platformUserID string) (*User, error) {
	user, err := r.GetUserByPlatformID(platformUserID)
	if err == nil {
		return user, nil
	}

	newUser := &User{
		ID:             generateID(),
		PlatformUserID: platformUserID,
		Username:       "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *Repository) UpdateUsername(userID, username string) error {
	query := `UPDATE users SET username = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, username, time.Now(), userID)
	if err != nil {
		return err
	}

	return nil
}

func generateID() string {
	return id.Generate()
}
