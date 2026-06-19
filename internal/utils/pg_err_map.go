package utils

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrForeignKey    = errors.New("related resource not found")
	ErrInvalidInput  = errors.New("invalid input")
)

func MapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrAlreadyExists

		case "23503":
			return ErrForeignKey

		case "22P02":
			return ErrInvalidInput
		}
	}

	return err
}
