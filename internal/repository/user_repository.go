package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/manish-npx/go-auth-api/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func User(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) New(ctx context.Context, u *model.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`,
		u.Name, u.Email, u.PasswordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *UserRepository) Get(ctx context.Context, email string) (*model.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	)
	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
