package db

import (
	"database/sql"
	"fmt"

	"github.com/ilya2044/pullup-diary/models"
)

func CreateUser(telegramID int64) error {
	_, err := DB.Exec("INSERT INTO users (telegram_id) VALUES ($1)", telegramID)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}
	return nil
}

func DeleteUser(telegramID int64) error {
	_, err := DB.Exec("DELETE FROM users WHERE telegram_id = ($1)", telegramID)
	if err != nil {
		return err
	}
	return nil
}

func FindUserByTelegramID(telegramID int64) (*models.User, error) {
	user := models.User{}
	row := DB.QueryRow("SELECT id, telegram_id, reminder_period, created_at FROM users WHERE telegram_id = $1", telegramID)
	err := row.Scan(&user.ID, &user.TelegramID, &user.ReminderPeriod, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
