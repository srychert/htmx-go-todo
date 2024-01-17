package model

import (
	"database/sql"
	"log"
)

var db *sql.DB

func Setup() {
	var err error
	db, err = sql.Open("sqlite", "./todo.db")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	statement, _ := db.Prepare(
		"CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY, title TEXT, date INTEGER, done INTEGER)")

	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
