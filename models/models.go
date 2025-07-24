package models

import "time"

type User struct {
	ID             int64
	TelegramID     int64
	ReminderPeriod int
	CreatedAt      time.Time
}

type Workout_Day struct {
	ID     int64
	UserID int64
	Date   time.Time
}

type Set struct {
	ID        int64
	DayID     int64
	Reps      int
	Note      string
	CreatedAt time.Time
}
