package db

import (
	"database/sql"
	"log/slog"
	_ "github.com/lib/pq"
)

func NewPostgresDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("error opening db", "err", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		slog.Error("error connecting to DB", "err", err)
		return nil
	}
	return db
}