package users

import (
	"context"

	"github.com/umeshj346/helloWorldServer/domain"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *domain.UserData) error
	GetUserByName(ctx context.Context, firstName, lastName string) (*domain.User, error)
	GetCountOfUsers(ctx context.Context) (int, error)
	Shutdown()
}