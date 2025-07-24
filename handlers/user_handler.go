package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ilya2044/pullup-diary/db"
)

type CreateUserRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

type DeleteUserRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := db.CreateUser(req.TelegramID)
		if err != nil {
			http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь успешно добавлен",
		})

	case r.Method == http.MethodDelete:
		var user DeleteUserRequest
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := db.DeleteUser(user.TelegramID)
		if err != nil {
			http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь успешно удален",
		})

	default:
		http.Error(w, "Only POST/DELETE allowed", http.StatusMethodNotAllowed)
	}
}
