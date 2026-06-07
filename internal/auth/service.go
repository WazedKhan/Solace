package auth

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	if req.Email == "" {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, err
	}
	user := User{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}
