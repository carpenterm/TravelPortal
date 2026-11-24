package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skandpurohit/TravelPortal/config"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", config.DBFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL CHECK(LENGTH(title) <= 255),
		comment TEXT CHECK(LENGTH(comment) <= 1024),
		repeat TEXT CHECK(LENGTH(repeat) <= 255)
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date)`)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	return db
}
