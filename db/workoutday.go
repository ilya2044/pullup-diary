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

func GetWorkoutDaysByUserID(userID int64) ([]models.Workout_Day, error) {
	rows, err := DB.Query("SELECT id, user_id, date FROM workout_days WHERE user_id = $1 ORDER BY date", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []models.Workout_Day
	for rows.Next() {
		var day models.Workout_Day
		if err := rows.Scan(&day.ID, &day.UserID, &day.Date); err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, nil
}

func GetWorkoutDayIDByUserIDAndDate(userID int64, date time.Time) (int64, error) {
	var dayID int64
	err := DB.QueryRow(
		"SELECT id FROM workout_days WHERE user_id = $1 AND date = $2",
		userID, date).Scan(&dayID)
	if err != nil {
		return 0, err
	}
	return dayID, nil
}

func GetWorkoutDaysByTelegramID(telegramID int64) ([]models.Workout_Day, error) {
	userID, err := GetUserIDByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT id, user_id, date FROM workout_days WHERE user_id = $1 ORDER BY date", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []models.Workout_Day
	for rows.Next() {
		var day models.Workout_Day
		if err := rows.Scan(&day.ID, &day.UserID, &day.Date); err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, nil
}
