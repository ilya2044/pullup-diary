package main

import (
	"log"
	"net/http"

	"github.com/ilya2044/pullup-diary/db"
	"github.com/ilya2044/pullup-diary/handlers"
)

func main() {
	db.Init()

	http.HandleFunc("/users", handlers.UsersHandler)
	http.HandleFunc("/workout_day", handlers.WorkoutDayHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
