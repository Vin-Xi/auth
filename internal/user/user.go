package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`

	Username    pgtype.Text        `json:"username"`
	FirstName   pgtype.Text        `json:"firstName"`
	LastName    pgtype.Text        `json:"lastName"`
	LastLoginAt pgtype.Timestamptz `json:"lastLoginAt"`

	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
