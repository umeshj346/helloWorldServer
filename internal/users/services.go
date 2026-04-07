package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"

	"github.com/umeshj346/helloWorldServer/domain"
)

type Service struct {
	repo UserRepository
}

func NewService(r UserRepository) *Service {
	return &Service{r}
}

func (s *Service) GetUser(ctx context.Context, firstName, lastName string) (*domain.User, error) {
	return s.repo.GetUserByName(ctx, firstName, lastName)
}

func (s *Service) AddUser(ctx context.Context, user *domain.UserData) error {
	if user.FirstName == "" {
		return fmt.Errorf("empty first name")
	}
	if user.LastName == "" {
		return fmt.Errorf("empty last name")
	}

	_, err := s.repo.GetUserByName(ctx, user.FirstName, user.LastName)
	if err != nil && !errors.Is(err, domain.ErrNoResultFound) {
		return fmt.Errorf("error checking if user is already present, err: %v", err)
	}

	if err == nil {
		return domain.ErrUserAlreadyExists
	}
	
	_, err = mail.ParseAddress(user.Email)
	if err != nil { 
		return domain.ErrInvalidEmail
	}
	
	err = s.repo.InsertUser(ctx, user)
	if err != nil {
		slog.Error("error inserting into database", "err", err)
	}

	return nil
}

func (s *Service) Shutdown() {
	s.repo.Shutdown()
}