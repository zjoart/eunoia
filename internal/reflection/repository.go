package reflection

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

func (r *Repository) CreateReflection(reflection *Reflection) error {
	logger.Info("creating reflection", logger.Fields{"user_id": reflection.UserID})

	query := `INSERT INTO reflections (id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, reflection.ID, reflection.UserID, reflection.Content, reflection.Sentiment,
		reflection.KeyThemes, reflection.AIAnalysis, reflection.CreatedAt, reflection.UpdatedAt)

	if err != nil {
		logger.Error("failed to create reflection", logger.WithError(err))
		return err
	}

	logger.Info("reflection created successfully", logger.Fields{"reflection_id": reflection.ID})
	return nil
}

func (r *Repository) GetReflectionsByUserID(userID string, limit int) ([]*Reflection, error) {
	logger.Info("fetching reflections for user", logger.Fields{
		"user_id": userID,
		"limit":   limit,
	})

	query := `SELECT id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at
			  FROM reflections
			  WHERE user_id = ?
			  ORDER BY created_at DESC
			  LIMIT ?`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		logger.Error("failed to fetch reflections", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	var reflections []*Reflection
	for rows.Next() {
		reflection := &Reflection{}
		err := rows.Scan(&reflection.ID, &reflection.UserID, &reflection.Content, &reflection.Sentiment,
			&reflection.KeyThemes, &reflection.AIAnalysis, &reflection.CreatedAt, &reflection.UpdatedAt)
		if err != nil {
			logger.Error("failed to scan reflection row", logger.WithError(err))
			return nil, err
		}
		reflections = append(reflections, reflection)
	}

	logger.Info("reflections fetched successfully", logger.Fields{"count": len(reflections)})
	return reflections, nil
}

func (r *Repository) GetRecentReflections(userID string, days int) ([]*Reflection, error) {
	logger.Info("fetching recent reflections", logger.Fields{
		"user_id": userID,
		"days":    days,
	})

	startDate := time.Now().AddDate(0, 0, -days)

	query := `SELECT id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at
			  FROM reflections
			  WHERE user_id = ? AND created_at >= ?
			  ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, startDate)
	if err != nil {
		logger.Error("failed to fetch recent reflections", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	var reflections []*Reflection
	for rows.Next() {
		reflection := &Reflection{}
		err := rows.Scan(&reflection.ID, &reflection.UserID, &reflection.Content, &reflection.Sentiment,
			&reflection.KeyThemes, &reflection.AIAnalysis, &reflection.CreatedAt, &reflection.UpdatedAt)
		if err != nil {
			logger.Error("failed to scan reflection row", logger.WithError(err))
			return nil, err
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}
