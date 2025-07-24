package db

import "github.com/ilya2044/pullup-diary/models"

func AddSet(workoutDayID int64, reps int, note string) (models.Set, error) {
	var id int64
	err := DB.QueryRow(
		"INSERT INTO sets (day_id, reps, note) VALUES ($1, $2, $3) RETURNING id",
		workoutDayID, reps, note).Scan(&id)
	return models.Set{
		ID:    id,
		DayID: workoutDayID,
		Reps:  reps,
		Note:  note,
	}, err
}

func GetSetByTelegramID(telegramID int64) ([]models.Set, error) {
	userID, err := GetUserIDByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
	SELECT sets.id, sets.day_id, sets.reps, sets.note, sets.created_at
	FROM sets
	JOIN workout_days ON sets.day_id = workout_days.id
	WHERE workout_days.user_id = $1
	ORDER BY sets.created_at DESC
`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []models.Set
	for rows.Next() {
		var set models.Set
		if err := rows.Scan(&set.ID, &set.DayID, &set.Reps, &set.Note, &set.CreatedAt); err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, nil
}
