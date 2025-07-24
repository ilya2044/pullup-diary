package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ilya2044/pullup-diary/db"
)

type CreateWorkoutDayRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Date       string `json:"date"`
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

	userID, err := db.GetUserIDByTelegramID(req.TelegramID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	day, err := db.CreateWorkoutDay(userID, date)
	if err != nil {
		http.Error(w, "Ошибка при создании тренировочного дня", http.StatusInternalServerError)
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

func GetWorkoutDaysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET", http.StatusMethodNotAllowed)
		return
	}

	telegramIDStr := r.URL.Query().Get("telegram_id")
	if telegramIDStr == "" {
		http.Error(w, "telegram_id обязателен", http.StatusBadRequest)
		return
	}

	telegramID, err := ParseInt64(telegramIDStr)
	if err != nil {
		http.Error(w, "Неверный telegram_id", http.StatusBadRequest)
		return
	}

	days, err := db.GetWorkoutDaysByTelegramID(telegramID)
	if err != nil {
		http.Error(w, "Ошибка получения тренировочных дней", http.StatusInternalServerError)
		return
	}

	respDays := make([]map[string]interface{}, 0, len(days))
	for _, d := range days {
		respDays = append(respDays, map[string]interface{}{
			"id":      d.ID,
			"user_id": d.UserID,
			"date":    d.Date.Format("2006-01-02"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": respDays,
	})
}

func ParseInt64(s string) (int64, error) {
	var v int64
	_, err := fmt.Sscan(s, &v)
	return v, err
}
