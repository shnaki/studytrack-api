package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/repository/sqlcgen"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type userRepository struct {
	q *sqlcgen.Queries
}

// NewUserRepository creates a new UserRepository implementation using PostgreSQL.
func NewUserRepository(pool *pgxpool.Pool) port.UserRepository {
	return &userRepository{q: sqlcgen.New(pool)}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:        toPgUUID(user.ID),
		Name:      user.Name,
		CreatedAt: toPgTimestamptz(user.CreatedAt),
		UpdatedAt: toPgTimestamptz(user.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	row, err := r.q.GetUserByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("user")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return domain.ReconstructUser(
		fromPgUUID(row.ID),
		row.Name,
		fromPgTimestamptz(row.CreatedAt),
		fromPgTimestamptz(row.UpdatedAt),
	), nil
}
