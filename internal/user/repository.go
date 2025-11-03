package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zjoart/eunoia/internal/pkg/id"
	"github.com/zjoart/eunoia/pkg/logger"
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
	logger.Info("creating user", logger.Fields{"telex_user_id": user.TelexUserID})

	query := `INSERT INTO users (id, telex_user_id, username, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.ID, user.TelexUserID, user.Username, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logger.Error("failed to create user", logger.WithError(err))
		return err
	}

	logger.Info("user created successfully", logger.Fields{"user_id": user.ID})
	return nil
}

func (r *Repository) GetUserByTelexID(telexUserID string) (*User, error) {
	logger.Info("fetching user by telex id", logger.Fields{"telex_user_id": telexUserID})

	query := `SELECT id, telex_user_id, username, created_at, updated_at
			  FROM users
			  WHERE telex_user_id = ?`

	user := &User{}
	err := r.db.QueryRow(query, telexUserID).Scan(
		&user.ID, &user.TelexUserID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("user not found", logger.Fields{"telex_user_id": telexUserID})
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		logger.Error("failed to fetch user", logger.WithError(err))
		return nil, err
	}

	logger.Info("user fetched successfully", logger.Fields{"user_id": user.ID})
	return user, nil
}

func (r *Repository) GetOrCreateUser(telexUserID string) (*User, error) {
	user, err := r.GetUserByTelexID(telexUserID)
	if err == nil {
		return user, nil
	}

	logger.Info("creating new user", logger.Fields{"telex_user_id": telexUserID})

	newUser := &User{
		ID:          generateID(),
		TelexUserID: telexUserID,
		Username:    "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *Repository) UpdateUsername(userID, username string) error {
	logger.Info("updating username", logger.Fields{"user_id": userID})

	query := `UPDATE users SET username = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, username, time.Now(), userID)
	if err != nil {
		logger.Error("failed to update username", logger.WithError(err))
		return err
	}

	logger.Info("username updated successfully")
	return nil
}

func generateID() string {
	return id.Generate()
}
