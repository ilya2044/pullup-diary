package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ilya2044/pullup-diary/db"
)

const apiBaseURL = "http://localhost:8080"

var (
	userStates = make(map[int64]string)
	stateMu    sync.Mutex
)

func startReminderLoop(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			users, err := db.GetUsersWithReminderPeriod()
			if err != nil {
				log.Println("Ошибка получения пользователей для напоминаний:", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			for _, user := range users {
				msg := tgbotapi.NewMessage(user.TelegramID, "Напоминание")
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Ошибка отправки напоминания пользователю %d: %v", user.TelegramID, err)
				}

				time.Sleep(time.Duration(user.ReminderPeriod) * time.Minute)
			}
		}
	}()
}

func setUserState(userID int64, state string) {
	stateMu.Lock()
	defer stateMu.Unlock()
	userStates[userID] = state
}

func getUserState(userID int64) string {
	stateMu.Lock()
	defer stateMu.Unlock()
	return userStates[userID]
}

func clearUserState(userID int64) {
	stateMu.Lock()
	defer stateMu.Unlock()
	delete(userStates, userID)
}

func registerUser(telegramID int64) error {
	reqBody := map[string]int64{"telegram_id": telegramID}
	data, _ := json.Marshal(reqBody)
	resp, err := http.Post(apiBaseURL+"/users", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusInternalServerError {
		return fmt.Errorf("ошибка регистрации: статус %d", resp.StatusCode)
	}
	return nil
}

// func createWorkoutDay(telegramID int64, date string) error {
// 	reqBody := map[string]interface{}{
// 		"telegram_id": telegramID,
// 		"date":        date,
// 	}
// 	data, _ := json.Marshal(reqBody)
// 	resp, err := http.Post(apiBaseURL+"/workout_day", "application/json", bytes.NewBuffer(data))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("не удалось создать тренировочный день, статус: %d", resp.StatusCode)
// 	}
// 	return nil
// }

// func getWorkoutDays(telegramID int64) ([]map[string]interface{}, error) {
// 	resp, err := http.Get(fmt.Sprintf("%s/workout_days?telegram_id=%d", apiBaseURL, telegramID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result struct {
// 		Data []map[string]interface{} `json:"data"`
// 	}
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result.Data, nil
// }

func addSet(telegramID int64, date string, reps int, note string) error {
	reqBody := map[string]interface{}{
		"telegram_id": fmt.Sprintf("%d", telegramID),
		"date":        date,
		"reps":        reps,
		"note":        note,
	}
	data, _ := json.Marshal(reqBody)
	resp, err := http.Post(apiBaseURL+"/set", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("не удалось добавить подход, статус: %d", resp.StatusCode)
	}
	return nil
}

func getSets(telegramID int64) ([]map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/sets?telegram_id=%d", apiBaseURL, telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func RunBot() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Не задан TELEGRAM_BOT_TOKEN")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Авторизован как %s", bot.Self.UserName)

	startReminderLoop(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		err := registerUser(userID)
		if err != nil {
			log.Printf("Ошибка регистрации пользователя: %v", err)
		}

		state := getUserState(userID)

		switch {
		case strings.HasPrefix(text, "/start"):
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Добавить подход"),
					tgbotapi.NewKeyboardButton("/sets"),
					tgbotapi.NewKeyboardButton("/reminder"),
				),
			)
			msg := tgbotapi.NewMessage(chatID, "Выбери действие:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			clearUserState(userID)

		case text == "Добавить подход":
			msg := tgbotapi.NewMessage(chatID, "Кол-во_повторов заметка")
			bot.Send(msg)
			setUserState(userID, "awaiting_set")

		case state == "awaiting_set":
			parts := strings.Fields(text)
			if len(parts) < 1 {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка, введи количество повторов и (опционально) заметку"))
				continue
			}

			reps, err := strconv.Atoi(parts[0])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Первое слово - число"))
				continue
			}

			note := ""
			if len(parts) > 1 {
				note = strings.Join(parts[1:], " ")
			}

			date := time.Now().Format("2006-01-02")

			err = addSet(userID, date, reps, note)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при добавлении подхода: "+err.Error()))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Добавлен подход: %d повторов на %s", reps, date)))
			}

			clearUserState(userID)

		// case strings.HasPrefix(text, "/days"):
		// 	days, err := getWorkoutDays(userID)
		// 	if err != nil {
		// 		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении дней: "+err.Error()))
		// 		continue
		// 	}
		// 	if len(days) == 0 {
		// 		bot.Send(tgbotapi.NewMessage(chatID, "Тренировочные дни не созданы"))
		// 		continue
		// 	}
		// 	var resp strings.Builder
		// 	resp.WriteString("Тренировочные дни:\n")
		// 	for _, day := range days {
		// 		resp.WriteString(fmt.Sprintf("- %s\n", day["date"]))
		// 	}
		// 	bot.Send(tgbotapi.NewMessage(chatID, resp.String()))

		case strings.HasPrefix(text, "/sets"):
			sets, err := getSets(userID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении подходов: "+err.Error()))
				continue
			}
			if len(sets) == 0 {
				bot.Send(tgbotapi.NewMessage(chatID, "Нет подходов"))
				continue
			}
			var resp strings.Builder
			resp.WriteString("Ваши подходы:\n")
			for _, set := range sets {
				date, _ := set["date"].(string)
				repsFloat, ok := set["reps"].(float64)
				if !ok {
					resp.WriteString("- [ошибка: reps не число]\n")
					continue
				}
				reps := int(repsFloat)

				note, _ := set["note"].(string)
				if note != "" {
					resp.WriteString(fmt.Sprintf("- %s: %d повторов (%s)\n", date, reps, note))
				} else {
					resp.WriteString(fmt.Sprintf("- %s: %d повторов\n", date, reps))
				}
			}
			bot.Send(tgbotapi.NewMessage(chatID, resp.String()))

		case strings.HasPrefix(text, "/reminder"):
			args := strings.Fields(text)
			if len(args) != 2 {
				bot.Send(tgbotapi.NewMessage(chatID, "Использование: /reminder количество_минут"))
				continue
			}
			period, err := strconv.Atoi(args[1])
			if err != nil || period < 1 {
				bot.Send(tgbotapi.NewMessage(chatID, "Неверное кол-во минут"))
				continue
			}

			reqBody := map[string]interface{}{
				"telegram_id": fmt.Sprintf("%d", userID),
				"period":      period,
			}
			data, _ := json.Marshal(reqBody)
			resp, err := http.Post(apiBaseURL+"/reminder", "application/json", bytes.NewBuffer(data))
			if err != nil || resp.StatusCode != http.StatusOK {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при обновлении периода напоминаний"))
				continue
			}
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Период напоминаний установлен на %d минут", period)))

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда или сообщение. /start для меню"))
		}
	}
}
