package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilya2044/pullup-diary/db"
)

type CreateWorkoutDayRequest struct {
	UserID int64  `json:"user_id"`
	Date   string `json:"date"`
}

func WorkoutDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}

	var req CreateWorkoutDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Неверный формат даты. Используйте YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	day, err := db.CreateWorkoutDay(req.UserID, date)
	if err != nil {
		http.Error(w, "Пользователь не зарегестрирован", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Тренировочный день добавлен",
		"data": map[string]interface{}{
			"id":      day.ID,
			"user_id": day.UserID,
			"date":    day.Date.Format("2006-01-02"),
		},
	})
}
