package main

import (
	"log"
	"net/http"

	"github.com/ilya2044/pullup-diary/db"
	"github.com/ilya2044/pullup-diary/handlers"
	"github.com/ilya2044/pullup-diary/telegram"
)

func main() {
	db.Init()

	http.HandleFunc("/users", handlers.UsersHandler)
	http.HandleFunc("/workout_day", handlers.WorkoutDayHandler)
	http.HandleFunc("/workout_days", handlers.GetWorkoutDaysHandler)
	http.HandleFunc("/set", handlers.SetHandler)
	http.HandleFunc("/sets", handlers.GetSetHandler)
	http.HandleFunc("/reminder", handlers.ReminderHandler)

	go telegram.RunBot()

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
