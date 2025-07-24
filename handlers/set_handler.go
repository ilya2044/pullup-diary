package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilya2044/pullup-diary/db"
)

type SetRequest struct {
	TelegramID string `json:"telegram_id"`
	Date       string `json:"date"`
	Reps       int    `json:"reps"`
	Note       string `json:"note"`
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}

	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	telegramID, err := ParseInt64(req.TelegramID)
	if err != nil {
		http.Error(w, "Неверный telegram_id", http.StatusBadRequest)
		return
	}

	var date time.Time
	if req.Date == "" {
		// Если дата не указана — ставим сегодняшнюю
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			http.Error(w, "Неверный формат даты. Используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	userID, err := db.GetUserIDByTelegramID(telegramID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	dayID, err := db.GetWorkoutDayIDByUserIDAndDate(userID, date)
	if err != nil {
		http.Error(w, "Тренировочный день не найден", http.StatusBadRequest)
		return
	}

	set, err := db.AddSet(dayID, req.Reps, req.Note)
	if err != nil {
		http.Error(w, "Ошибка при добавлении подхода", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Тренировочный подход добавлен",
		"data":    set,
	})
}

func GetSetHandler(w http.ResponseWriter, r *http.Request) {
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

	sets, err := db.GetSetByTelegramID(telegramID)
	if err != nil {
		http.Error(w, "Ошибка получения подходов", http.StatusInternalServerError)
		return
	}

	respSets := make([]map[string]interface{}, 0, len(sets))
	for _, s := range sets {
		respSets = append(respSets, map[string]interface{}{
			"id":     s.ID,
			"day_id": s.DayID,
			"reps":   s.Reps,
			"note":   s.Note,
			"date":   s.CreatedAt.Format("2006-01-02"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": respSets,
	})
}
