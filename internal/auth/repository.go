package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user User) (*User, error) {
	query := `
		INSERT INTO users(id, name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, created_at
	`
	var created User
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, mapPostgresError(err)
	}
	return &created, nil
}
