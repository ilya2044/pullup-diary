package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ilya2044/pullup-diary/db"
)

type SetRequest struct {
	WorkoutDayID int64  `json:"day_id"`
	Reps         int    `json:"reps"`
	Note         string `json:"note"`
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

	set, err := db.AddSet(req.WorkoutDayID, req.Reps, req.Note)
	if err != nil {
		http.Error(w, "Пользователь не зарегистрирован", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Тренировочный подход добавлен",
		"data":    set,
	})

}
