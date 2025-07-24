package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	telegram_id BIGINT UNIQUE NOT NULL,
	reminder_period INTEGER DEFAULT 180,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_days (
	id SERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id),
	date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS sets (
	id SERIAL PRIMARY KEY,
	day_id BIGINT NOT NULL REFERENCES workout_days(id),
	reps INTEGER DEFAULT 0,
	note TEXT,
	created_at TIMESTAMP DEFAULT NOW()
);
`

func Init() {
	var err error

	DB, err = sql.Open("postgres", "user=postgres password=1101 dbname=pullups sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка при подключении к базе:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("База не отвечает:", err)
	}

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatal("Ошибка создания таблицы users:", err)
	}

	log.Println("Успешное подключение к базе данных")

}
