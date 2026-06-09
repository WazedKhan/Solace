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
		return nil, mapPostgresError(err)
	}
	return &created, nil
}

func (r *Repository) GetUsers(ctx context.Context, q GetUserQuery) ([]User, error) {
	query := `
	SELECT id, name, email, created_at
	FROM users
	WHERE (
		$1 = ''
		OR name ILIKE '%' || $1 || '%'
		OR email ILIKE '%' || $1 || '%'
	)
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
`
	rows, err := r.db.Query(ctx, query, q.Search, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
