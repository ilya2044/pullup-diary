package db

import "github.com/ilya2044/pullup-diary/models"

func AddSet(workoutDayID int64, reps int, note string) (Set, error) {
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
