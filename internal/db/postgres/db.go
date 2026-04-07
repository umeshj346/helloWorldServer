package postgres

import (
	"database/sql"
	"log/slog"
	_ "github.com/lib/pq"
)

func NewPostgresDB(connStr string) *sql.DB {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("error connecting to database", "err", err)
		return nil
	}

	if err := dbConn.Ping(); err != nil {
		slog.Error("error pinging to database", "err", err)
		return nil
	}
	return dbConn
}