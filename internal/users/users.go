package users

import (
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
	users []User
	rw sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{}
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

	newUser := User {
		FirstName: firstName,
		LastName: lastName,
		Email: *parsedEmail,
	}
	m.users = append(m.users, newUser)
	return nil
}

func (m *Manager) getUserByNameUnsafe(firstName, lastName string) (*User, error) {
	for i, user:= range m.users {
		if user.FirstName == firstName && user.LastName == lastName {
			return &m.users[i], nil
		}
	}
	return nil, ErrNoResultFound
}

func (m *Manager) GetUserByName(firstName, lastName string) (*User, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for i, user:= range m.users {
		if user.FirstName == firstName && user.LastName == lastName {
			return &m.users[i], nil
		}
	}
	return nil, ErrNoResultFound
}

func (m *Manager) Shutdown() {
	slog.Info("user manager shutting down")
	time.Sleep(2*time.Second)
	slog.Info("user manager shutting down complete")
}