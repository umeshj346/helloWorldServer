package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/mail"

	"github.com/umeshj346/helloWorldServer/domain"
)

type UserPgRepository struct {
	db *sql.DB
}

func NewUserPgRepository(db *sql.DB) *UserPgRepository {
	return &UserPgRepository{db}
}

func (r *UserPgRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err:= r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("error executing the query", "err", err)
		return nil, err
	}
	
	defer rows.Close()
	
	result := make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		var email string
		err = rows.Scan(
			&t.ID,
			&t.FirstName,
			&t.LastName,
			&email,
		)
		
		if err != nil {
			slog.Error("error scanning a row", "err", err)
			return nil, err
		}
		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			slog.Error("error parsing the email", "err", err)
			return nil, err
		}
		t.Email = *parsedEmail
		result = append(result, t)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	
	return result, nil
}

func (r *UserPgRepository) InsertUser(ctx context.Context, user *domain.UserData) error {
	query := `
		INSERT INTO users (first_name, last_name, email)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email)
	return err
}

func (r *UserPgRepository) GetUserByName(ctx context.Context, firstName, lastName string) (res *domain.User, err error) {
	query := `
		SELECT id, first_name, last_name, email 
		from users
		where first_name = $1 and last_name = $2
	`

	resList, err := r.fetch(ctx, query, firstName, lastName)
	if err != nil {
		return
	}

	if len(resList) < 1 {
		return res, domain.ErrNoResultFound
	}
	return &resList[0], nil
}

func (r *UserPgRepository) GetCountOfUsers(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM users
	`
	var usersCnt int
	row := r.db.QueryRowContext(ctx, query)
	err := row.Scan(&usersCnt)
	if err != nil {
		return usersCnt, err
	}
	return usersCnt, nil
	
}

func (r *UserPgRepository) Shutdown() {
	slog.Info("user UserPgRepository shutting down")
	err := r.db.Close()
	if err != nil {
		slog.Error("got error when closing dbconnection", "err", err)
	}
	slog.Info("user UserPgRepository shutting down complete")
}