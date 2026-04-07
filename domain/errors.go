package domain

import "errors"

var (
	ErrNoResultFound = errors.New("no result found")
	ErrUserAlreadyExists = errors.New("user with this name already exists")
	ErrInvalidEmail = errors.New("invalid email")
)