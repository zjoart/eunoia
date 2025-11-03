package user

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	user := &User{
		ID:             "user-123",
		PlatformUserID: "platform-456",
		Username:       "testuser",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.PlatformUserID, user.Username, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUserByPlatformID_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	platformID := "platform-456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "platform_user_id", "username", "created_at", "updated_at"}).
		AddRow("user-123", platformID, "testuser", now, now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE platform_user_id").
		WithArgs(platformID).
		WillReturnRows(rows)

	user, err := repo.GetUserByPlatformID(platformID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != "user-123" {
		t.Errorf("expected ID 'user-123', got '%s'", user.ID)
	}

	if user.PlatformUserID != platformID {
		t.Errorf("expected platform ID '%s', got '%s'", platformID, user.PlatformUserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUserByPlatformID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	platformID := "non-existent"

	mock.ExpectQuery("SELECT (.+) FROM users WHERE platform_user_id").
		WithArgs(platformID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByPlatformID(platformID)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetOrCreateUser_ExistingUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	platformID := "platform-456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "platform_user_id", "username", "created_at", "updated_at"}).
		AddRow("user-123", platformID, "testuser", now, now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE platform_user_id").
		WithArgs(platformID).
		WillReturnRows(rows)

	user, err := repo.GetOrCreateUser(platformID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.PlatformUserID != platformID {
		t.Errorf("expected platform ID '%s', got '%s'", platformID, user.PlatformUserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetOrCreateUser_NewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	platformID := "new-platform-789"

	// First query returns not found
	mock.ExpectQuery("SELECT (.+) FROM users WHERE platform_user_id").
		WithArgs(platformID).
		WillReturnError(sql.ErrNoRows)

	// Then insert new user
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), platformID, "", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := repo.GetOrCreateUser(platformID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.PlatformUserID != platformID {
		t.Errorf("expected platform ID '%s', got '%s'", platformID, user.PlatformUserID)
	}

	if user.ID == "" {
		t.Error("expected generated ID, got empty string")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userID := "user-123"
	username := "newusername"

	mock.ExpectExec("UPDATE users SET username").
		WithArgs(username, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUsername(userID, username)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
