package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"sync"
	"time"
)

var ErrNoResultFound = errors.New("no result found")

type User struct {
	FirstName string
	LastName  string
	Email     mail.Address
}

type Manager struct {
	db *sql.DB
	rw sync.RWMutex
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) AddUser(firstName, lastName, email string) error {
	if firstName == "" {
		return fmt.Errorf("empty first name")
	}
	if lastName == "" {
		return fmt.Errorf("empty last name")
	}

	m.rw.Lock()
	defer m.rw.Unlock()
	foundUser, err := m.getUserByNameUnsafe(firstName, lastName)
	if err != nil && !errors.Is(err, ErrNoResultFound) {
		return fmt.Errorf("error checking if user is already present")
	}

	if foundUser != nil {
		return fmt.Errorf("user with this name already exists")
	}
	
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil { 
		return fmt.Errorf("invalid email: %s", email)
	}
	
	query := `
		INSERT INTO users (first_name, last_name, email)
		VALUES ($1, $2, $3)
	`

	_, err = m.db.Exec(query, firstName, lastName, parsedEmail.Address)
	if err != nil {
		return fmt.Errorf("error inserting into database: %s", err)
	}

	return nil
}

func (m *Manager) getUserByNameUnsafe(firstName, lastName string) (*User, error) {
	query := `
		SELECT first_name, last_name, email 
		from users
		where first_name = $1 and last_name = $2
	`

	row := m.db.QueryRow(query, firstName, lastName)

	foundUser := &User{}
	var email string

	err := row.Scan(&foundUser.FirstName, &foundUser.LastName, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoResultFound
		}
		return nil, err
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing the mail Address, err: %v", err)
	}
	foundUser.Email = *parsedEmail
	return foundUser, nil
}

func (m *Manager) GetUserByName(firstName, lastName string) (*User, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	query := `
		SELECT first_name, last_name, email 
		from users
		where first_name = $1 and last_name = $2
	`

	row := m.db.QueryRow(query, firstName, lastName)

	foundUser := &User{}
	var email string

	err := row.Scan(&foundUser.FirstName, &foundUser.LastName, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoResultFound
		}
		return nil, err
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing the mail Address, err: %v", err)
	}
	foundUser.Email = *parsedEmail
	return foundUser, nil
}

func (m *Manager) Shutdown() {
	slog.Info("user manager shutting down")
	time.Sleep(2*time.Second)
	slog.Info("user manager shutting down complete")
}