package conversation

import (
	"database/sql"
	"time"

	"github.com/zjoart/eunoia/pkg/logger"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveMessage(message *ConversationMessage) error {
	logger.Info("saving conversation message", logger.Fields{
		"user_id":      message.UserID,
		"message_role": message.MessageRole,
	})

	query := `INSERT INTO conversation_history (id, user_id, message_role, message_content, context_data, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, message.ID, message.UserID, message.MessageRole,
		message.MessageContent, message.ContextData, message.CreatedAt)

	if err != nil {
		logger.Error("failed to save conversation message", logger.WithError(err))
		return err
	}

	logger.Info("conversation message saved successfully", logger.Fields{"message_id": message.ID})
	return nil
}

func (r *Repository) GetConversationHistory(userID string, limit int) ([]*ConversationMessage, error) {
	logger.Info("fetching conversation history", logger.Fields{
		"user_id": userID,
		"limit":   limit,
	})

	query := `SELECT id, user_id, message_role, message_content, context_data, created_at
			  FROM conversation_history
			  WHERE user_id = ?
			  ORDER BY created_at DESC
			  LIMIT ?`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		logger.Error("failed to fetch conversation history", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	var messages []*ConversationMessage
	for rows.Next() {
		message := &ConversationMessage{}
		var contextData sql.NullString
		err := rows.Scan(&message.ID, &message.UserID, &message.MessageRole,
			&message.MessageContent, &contextData, &message.CreatedAt)
		if err != nil {
			logger.Error("failed to scan conversation message", logger.WithError(err))
			return nil, err
		}
		if contextData.Valid {
			message.ContextData = contextData.String
		}
		messages = append(messages, message)
	}

	logger.Info("conversation history fetched successfully", logger.Fields{"count": len(messages)})
	return messages, nil
}

func (r *Repository) GetRecentMessages(userID string, minutes int) ([]*ConversationMessage, error) {
	startTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	query := `SELECT id, user_id, message_role, message_content, context_data, created_at
			  FROM conversation_history
			  WHERE user_id = ? AND created_at >= ?
			  ORDER BY created_at ASC`

	rows, err := r.db.Query(query, userID, startTime)
	if err != nil {
		logger.Error("failed to fetch recent messages", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	var messages []*ConversationMessage
	for rows.Next() {
		message := &ConversationMessage{}
		var contextData sql.NullString
		err := rows.Scan(&message.ID, &message.UserID, &message.MessageRole,
			&message.MessageContent, &contextData, &message.CreatedAt)
		if err != nil {
			logger.Error("failed to scan message", logger.WithError(err))
			return nil, err
		}
		if contextData.Valid {
			message.ContextData = contextData.String
		}
		messages = append(messages, message)
	}

	return messages, nil
}
