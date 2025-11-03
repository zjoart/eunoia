package checkin

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

func (r *Repository) CreateCheckIn(checkIn *EmotionalCheckIn) error {
	logger.Info("creating emotional check-in", logger.Fields{
		"user_id":    checkIn.UserID,
		"mood_score": checkIn.MoodScore,
	})

	query := `INSERT INTO emotional_checkins (id, user_id, mood_score, mood_label, description, check_in_date, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, checkIn.ID, checkIn.UserID, checkIn.MoodScore, checkIn.MoodLabel,
		checkIn.Description, checkIn.CheckInDate, checkIn.CreatedAt)

	if err != nil {
		logger.Error("failed to create emotional check-in", logger.WithError(err))
		return err
	}

	logger.Info("emotional check-in created successfully", logger.Fields{"check_in_id": checkIn.ID})
	return nil
}

func (r *Repository) GetCheckInsByUserID(userID string, limit int) ([]*EmotionalCheckIn, error) {
	logger.Info("fetching check-ins for user", logger.Fields{
		"user_id": userID,
		"limit":   limit,
	})

	query := `SELECT id, user_id, mood_score, mood_label, description, check_in_date, created_at
			  FROM emotional_checkins
			  WHERE user_id = ?
			  ORDER BY check_in_date DESC, created_at DESC
			  LIMIT ?`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		logger.Error("failed to fetch check-ins", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	var checkIns []*EmotionalCheckIn
	for rows.Next() {
		checkIn := &EmotionalCheckIn{}
		err := rows.Scan(&checkIn.ID, &checkIn.UserID, &checkIn.MoodScore, &checkIn.MoodLabel,
			&checkIn.Description, &checkIn.CheckInDate, &checkIn.CreatedAt)
		if err != nil {
			logger.Error("failed to scan check-in row", logger.WithError(err))
			return nil, err
		}
		checkIns = append(checkIns, checkIn)
	}

	logger.Info("check-ins fetched successfully", logger.Fields{"count": len(checkIns)})
	return checkIns, nil
}

func (r *Repository) GetCheckInStats(userID string, days int) (*CheckInStats, error) {
	logger.Info("calculating check-in stats", logger.Fields{
		"user_id": userID,
		"days":    days,
	})

	startDate := time.Now().AddDate(0, 0, -days)

	query := `SELECT AVG(mood_score) as avg_score, COUNT(*) as total_count
			  FROM emotional_checkins
			  WHERE user_id = ? AND check_in_date >= ?`

	var avgScore sql.NullFloat64
	var totalCount int

	err := r.db.QueryRow(query, userID, startDate).Scan(&avgScore, &totalCount)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("failed to calculate stats", logger.WithError(err))
		return nil, err
	}

	stats := &CheckInStats{
		TotalCheckIns: totalCount,
	}

	if avgScore.Valid {
		stats.AverageMoodScore = avgScore.Float64
	}

	checkIns, err := r.GetCheckInsByUserID(userID, 2)
	if err != nil {
		return stats, nil
	}

	if len(checkIns) > 0 {
		stats.LastCheckIn = checkIns[0]

		if len(checkIns) > 1 {
			if checkIns[0].MoodScore > checkIns[1].MoodScore {
				stats.MoodTrend = "improving"
			} else if checkIns[0].MoodScore < checkIns[1].MoodScore {
				stats.MoodTrend = "declining"
			} else {
				stats.MoodTrend = "stable"
			}
		} else {
			stats.MoodTrend = "new"
		}
	}

	logger.Info("check-in stats calculated", logger.Fields{
		"avg_score": stats.AverageMoodScore,
		"total":     stats.TotalCheckIns,
	})

	return stats, nil
}

func (r *Repository) GetTodayCheckIn(userID string) (*EmotionalCheckIn, error) {
	today := time.Now().Format("2006-01-02")

	query := `SELECT id, user_id, mood_score, mood_label, description, check_in_date, created_at
			  FROM emotional_checkins
			  WHERE user_id = ? AND DATE(check_in_date) = ?
			  ORDER BY created_at DESC
			  LIMIT 1`

	checkIn := &EmotionalCheckIn{}
	err := r.db.QueryRow(query, userID, today).Scan(
		&checkIn.ID, &checkIn.UserID, &checkIn.MoodScore, &checkIn.MoodLabel,
		&checkIn.Description, &checkIn.CheckInDate, &checkIn.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		logger.Error("failed to get today's check-in", logger.WithError(err))
		return nil, err
	}

	return checkIn, nil
}
