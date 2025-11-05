package checkin

import (
	"fmt"
	"time"

	"github.com/zjoart/eunoia/internal/user"
	"github.com/zjoart/eunoia/pkg/id"
	"github.com/zjoart/eunoia/pkg/logger"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
}

func NewService(repo *Repository, userRepo *user.Repository) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateCheckIn(req *CreateCheckInRequest) (*EmotionalCheckIn, error) {
	if req.MoodScore < 1 || req.MoodScore > 10 {
		return nil, fmt.Errorf("mood score must be between 1 and 10")
	}

	userRecord, err := s.userRepo.GetOrCreateUser(req.PlatformUserID)
	if err != nil {
		logger.Error("failed to get or create user", logger.WithError(err))
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	checkIn := &EmotionalCheckIn{
		ID:          id.Generate(),
		UserID:      userRecord.ID,
		MoodScore:   req.MoodScore,
		MoodLabel:   req.MoodLabel,
		Description: req.Description,
		CheckInDate: time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateCheckIn(checkIn); err != nil {
		return nil, fmt.Errorf("failed to create check-in: %w", err)
	}

	return checkIn, nil
}

func (s *Service) GetCheckInHistory(platformUserID string, limit int) ([]*EmotionalCheckIn, error) {
	userRecord, err := s.userRepo.GetUserByPlatformID(platformUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.repo.GetCheckInsByUserID(userRecord.ID, limit)
}

func (s *Service) GetCheckInStats(platformUserID string, days int) (*CheckInStats, error) {
	userRecord, err := s.userRepo.GetUserByPlatformID(platformUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.repo.GetCheckInStats(userRecord.ID, days)
}

func (s *Service) GetTodayCheckIn(platformUserID string) (*EmotionalCheckIn, error) {
	userRecord, err := s.userRepo.GetUserByPlatformID(platformUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.repo.GetTodayCheckIn(userRecord.ID)
}

func (s *Service) GenerateMoodInsight(stats *CheckInStats) string {
	if stats.TotalCheckIns == 0 {
		return "Welcome! Start tracking your emotional wellbeing by sharing how you're feeling today."
	}

	insight := fmt.Sprintf("Over the past period, your average mood has been %.1f/10. ", stats.AverageMoodScore)

	switch stats.MoodTrend {
	case "improving":
		insight += "Your mood is trending upward, which is wonderful to see!"
	case "declining":
		insight += "I notice your mood has been declining. Remember, it's okay to have difficult days."
	case "stable":
		insight += "Your mood has been stable, which shows consistency."
	case "new":
		insight += "Keep tracking to see patterns over time."
	}

	return insight
}
