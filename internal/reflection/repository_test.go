package reflection

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateReflection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	reflection := &Reflection{
		ID:         "reflection-123",
		UserID:     "user-456",
		Content:    "Today I realized I need to focus more on self-care",
		Sentiment:  "positive",
		KeyThemes:  "self-care, health",
		AIAnalysis: "User is showing awareness of their needs",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO reflections").
		WithArgs(reflection.ID, reflection.UserID, reflection.Content, reflection.Sentiment,
			reflection.KeyThemes, reflection.AIAnalysis, reflection.CreatedAt, reflection.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateReflection(reflection)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetReflectionsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "content", "sentiment", "key_themes", "ai_analysis", "created_at", "updated_at"}).
		AddRow("ref-1", userID, "Reflection 1", "positive", "growth", "Analysis 1", now, now).
		AddRow("ref-2", userID, "Reflection 2", "neutral", "work", "Analysis 2", now, now)

	mock.ExpectQuery("SELECT (.+) FROM reflections").
		WithArgs(userID, 5).
		WillReturnRows(rows)

	reflections, err := repo.GetReflectionsByUserID(userID, 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(reflections) != 2 {
		t.Errorf("expected 2 reflections, got %d", len(reflections))
	}

	if reflections[0].Sentiment != "positive" {
		t.Errorf("expected sentiment 'positive', got '%s'", reflections[0].Sentiment)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
