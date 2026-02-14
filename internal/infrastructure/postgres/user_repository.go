package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) port.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Name, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at FROM users WHERE id = $1`,
		id,
	)
	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("user")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return domain.ReconstructUser(u.ID, u.Name, u.CreatedAt, u.UpdatedAt), nil
}
