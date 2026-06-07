package auth

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidInput       = errors.New("invalid user input")
)

func mapPostgresError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
	}

	return err
}
