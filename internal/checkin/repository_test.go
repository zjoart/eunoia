package checkin

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateCheckIn(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	checkIn := &EmotionalCheckIn{
		ID:          "checkin-123",
		UserID:      "user-456",
		MoodScore:   8,
		MoodLabel:   "happy",
		Description: "Feeling great today",
		CheckInDate: time.Now(),
		CreatedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO emotional_checkins").
		WithArgs(checkIn.ID, checkIn.UserID, checkIn.MoodScore, checkIn.MoodLabel,
			checkIn.Description, checkIn.CheckInDate, checkIn.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateCheckIn(checkIn)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetCheckInsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "mood_score", "mood_label", "description", "check_in_date", "created_at"}).
		AddRow("checkin-1", userID, 8, "happy", "Great day", now, now).
		AddRow("checkin-2", userID, 6, "content", "Okay day", now, now)

	mock.ExpectQuery("SELECT (.+) FROM emotional_checkins").
		WithArgs(userID, 5).
		WillReturnRows(rows)

	checkIns, err := repo.GetCheckInsByUserID(userID, 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(checkIns) != 2 {
		t.Errorf("expected 2 check-ins, got %d", len(checkIns))
	}

	if checkIns[0].MoodScore != 8 {
		t.Errorf("expected mood score 8, got %d", checkIns[0].MoodScore)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetCheckInStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-456"

	// Mock the stats query
	statsRows := sqlmock.NewRows([]string{"avg_score", "total_count"}).
		AddRow(7.5, 10)

	mock.ExpectQuery("SELECT AVG\\(mood_score\\) as avg_score, COUNT\\(\\*\\) as total_count").
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(statsRows)

	// Mock the recent check-ins query
	now := time.Now()
	checkInRows := sqlmock.NewRows([]string{"id", "user_id", "mood_score", "mood_label", "description", "check_in_date", "created_at"}).
		AddRow("checkin-1", userID, 8, "happy", "Great", now, now).
		AddRow("checkin-2", userID, 7, "content", "Good", now.Add(-24*time.Hour), now.Add(-24*time.Hour))

	mock.ExpectQuery("SELECT (.+) FROM emotional_checkins").
		WithArgs(userID, 2).
		WillReturnRows(checkInRows)

	stats, err := repo.GetCheckInStats(userID, 7)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if stats.TotalCheckIns != 10 {
		t.Errorf("expected 10 total check-ins, got %d", stats.TotalCheckIns)
	}

	if stats.AverageMoodScore != 7.5 {
		t.Errorf("expected average 7.5, got %f", stats.AverageMoodScore)
	}

	if stats.MoodTrend != "improving" {
		t.Errorf("expected mood trend 'improving', got '%s'", stats.MoodTrend)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetTodayCheckIn_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-456"
	now := time.Now()
	today := now.Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"id", "user_id", "mood_score", "mood_label", "description", "check_in_date", "created_at"}).
		AddRow("checkin-1", userID, 8, "happy", "Great day", now, now)

	mock.ExpectQuery("SELECT (.+) FROM emotional_checkins WHERE user_id = \\? AND DATE\\(check_in_date\\)").
		WithArgs(userID, today).
		WillReturnRows(rows)

	checkIn, err := repo.GetTodayCheckIn(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if checkIn == nil {
		t.Fatal("expected check-in, got nil")
	}

	if checkIn.MoodScore != 8 {
		t.Errorf("expected mood score 8, got %d", checkIn.MoodScore)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetTodayCheckIn_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-456"
	today := time.Now().Format("2006-01-02")

	mock.ExpectQuery("SELECT (.+) FROM emotional_checkins WHERE user_id = \\? AND DATE\\(check_in_date\\)").
		WithArgs(userID, today).
		WillReturnError(sql.ErrNoRows)

	checkIn, err := repo.GetTodayCheckIn(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if checkIn != nil {
		t.Errorf("expected nil check-in, got %v", checkIn)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
