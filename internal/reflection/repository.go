package reflection

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReflection(reflection *Reflection) error {
	query := `INSERT INTO reflections (id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, reflection.ID, reflection.UserID, reflection.Content, reflection.Sentiment,
		reflection.KeyThemes, reflection.AIAnalysis, reflection.CreatedAt, reflection.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetReflectionsByUserID(userID string, limit int) ([]*Reflection, error) {
	query := `SELECT id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at
			  FROM reflections
			  WHERE user_id = ?
			  ORDER BY created_at DESC
			  LIMIT ?`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reflections []*Reflection
	for rows.Next() {
		reflection := &Reflection{}
		err := rows.Scan(&reflection.ID, &reflection.UserID, &reflection.Content, &reflection.Sentiment,
			&reflection.KeyThemes, &reflection.AIAnalysis, &reflection.CreatedAt, &reflection.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}

func (r *Repository) GetRecentReflections(userID string, days int) ([]*Reflection, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	query := `SELECT id, user_id, content, sentiment, key_themes, ai_analysis, created_at, updated_at
			  FROM reflections
			  WHERE user_id = ? AND created_at >= ?
			  ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reflections []*Reflection
	for rows.Next() {
		reflection := &Reflection{}
		err := rows.Scan(&reflection.ID, &reflection.UserID, &reflection.Content, &reflection.Sentiment,
			&reflection.KeyThemes, &reflection.AIAnalysis, &reflection.CreatedAt, &reflection.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}
