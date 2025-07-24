package db

import (
	"database/sql"
	"fmt"

	"github.com/ilya2044/pullup-diary/models"
)

type UserReminder struct {
	TelegramID     int64
	ReminderPeriod int
}

func GetUsersWithReminderPeriod() ([]UserReminder, error) {
	rows, err := DB.Query("SELECT telegram_id, reminder_period FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserReminder
	for rows.Next() {
		var u UserReminder
		if err := rows.Scan(&u.TelegramID, &u.ReminderPeriod); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

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

func GetUserIDByTelegramID(telegramID int64) (int64, error) {
	var userID int64
	err := DB.QueryRow("SELECT id FROM users WHERE telegram_id = $1", telegramID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func UpdateReminderPeriod(telegramID int64, period int) error {
	userID, err := GetUserIDByTelegramID(telegramID)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET reminder_period = $1 WHERE id = $2", period, userID)
	return err
}
