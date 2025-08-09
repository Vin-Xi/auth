package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vin-Xi/auth/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid Email or Password")
)

type Service interface {
	Register(ctx context.Context, email string, password string, fName string, lName string) (*user.User, error)
	Login(ctx context.Context, email string, password string) (*user.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type service struct {
	repo user.UserRepository
}

func NewService(repo user.UserRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(ctx context.Context, email string, password string, fName string, lName string) (*user.User, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("could not hash password")
	}

	user := &user.User{
		Email:        email,
		PasswordHash: string(hashedPass),
		FirstName:    pgtype.Text{String: fName, Valid: true},
		LastName:     pgtype.Text{String: lName, Valid: true},
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email string, password string) (*user.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		fmt.Println(err)
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
