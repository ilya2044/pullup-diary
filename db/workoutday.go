package db

import (
	"time"

	"github.com/ilya2044/pullup-diary/models"
)

func CreateWorkoutDay(userID int64, date time.Time) (*models.Workout_Day, error) {
	var id int64
	err := DB.QueryRow(
		"INSERT INTO workout_days (user_id, date) VALUES ($1, $2) RETURNING id",
		userID, date,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.Workout_Day{
		ID:     id,
		UserID: userID,
		Date:   date,
	}, nil
}
