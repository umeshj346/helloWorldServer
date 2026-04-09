package mocks

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/umeshj346/helloWorldServer/domain"
)

type MockRepo struct {
	users []domain.User
}

func (mr *MockRepo) InsertUser(ctx context.Context, user *domain.UserData) error {
	email, err := mail.ParseAddress(user.Email)
	if err != nil {
		return domain.ErrInvalidEmail
	}
	mr.users = append(mr.users, 
		domain.User{
			FirstName: user.FirstName,
			LastName: user.LastName,
			Email: *email,
		},
	)
	return nil
}

func (mr *MockRepo) GetUserByName(ctx context.Context, firstName, lastName string) (*domain.User, error) {
	for i, user := range mr.users {
		if user.FirstName == firstName && user.LastName == lastName {
			return &mr.users[i], nil
		}
	}
	return nil, domain.ErrNoResultFound
}

func (mr *MockRepo) Shutdown() {
	fmt.Println("Mock Repo shut down")
}

func (mr *MockRepo) GetCountOfUsers(ctx context.Context) (int , error) {
	return len(mr.users), nil
}