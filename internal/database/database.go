package database

import (
	"context"
	"fmt"

	"github.com/Vin-Xi/auth/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func InitDB(ctx context.Context, db_url string) (*pgxpool.Pool, error) {
	if db_url == "" {
		return nil, fmt.Errorf("%v file not found or failed to read", db_url)
	}

	pool, err := pgxpool.New(ctx, db_url)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("unable to ping, connection failed with DB due to: %v", err)
	}

	fmt.Println("Connection is successful!")

	return pool, nil
}

func NewPostresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user")
	}

	return nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users where email = $1`
	u := &user.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user %w", err)
	}

	return u, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users where id = $1`
	u := &user.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user %w", err)
	}

	return u, nil
}
