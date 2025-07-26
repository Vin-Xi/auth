package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vin-Xi/auth/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid Email or Password")
)

type Service interface {
	Register(ctx context.Context, email string, password string) (*user.User, error)
	Login(ctx context.Context, email string, password string) (*user.User, error)
}

type service struct {
	repo user.UserRepository
}

func NewService(repo user.UserRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(ctx context.Context, email string, password string) (*user.User, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("could not hash password")
	}

	user := &user.User{
		Email:        email,
		PasswordHash: string(hashedPass),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email string, password string) (*user.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.PasswordHash))

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
