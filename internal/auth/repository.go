package auth

import (
	"context"
	"log"

	"github.com/WazedKhan/Solace/internal/utils"
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
		INSERT INTO users(id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, password, created_at
	`
	var created User
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.PasswordHash,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	return &created, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return nil, utils.MapPostgresError(err)
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE id=$1`

	var user User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return nil, utils.MapPostgresError(err)
	}
	return &user, nil
}
