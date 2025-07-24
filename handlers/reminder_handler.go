package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ilya2044/pullup-diary/db"
)

type ReminderRequest struct {
	TelegramID string `json:"telegram_id"`
	Period     int    `json:"period"`
}

func ReminderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}

	var req ReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	telegramID, err := strconv.ParseInt(req.TelegramID, 10, 64)
	if err != nil {
		http.Error(w, "Неверный telegram_id", http.StatusBadRequest)
		return
	}

	if req.Period < 1 {
		http.Error(w, "Период должен быть больше 0", http.StatusBadRequest)
		return
	}

	err = db.UpdateReminderPeriod(telegramID, req.Period)
	if err != nil {
		http.Error(w, "Ошибка обновления периода напоминаний", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Период напоминаний обновлен",
		"period":  req.Period,
	})
}

func GetReminderPeriodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	telegramIDStr := r.URL.Query().Get("telegram_id")
	if telegramIDStr == "" {
		http.Error(w, "telegram_id is required", http.StatusBadRequest)
		return
	}

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid telegram_id", http.StatusBadRequest)
		return
	}

	period, err := db.GetReminderPeriod(telegramID)
	if err != nil {
		http.Error(w, "Failed to get reminder period: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"period": period,
	})
}
