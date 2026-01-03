package database

import (
	"log"
	"database/sql"

	_ "modernc.org/sqlite"
)

func Init() *sql.DB {
	db, err := sql.Open("sqlite", "tmp.db")
	if err != nil { log.Fatal("Failed to initialize database: ", err)}

	return db
}