package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// bcryptCost of 12 balances security and registration latency (~300ms on modest hardware)
const bcryptCost = 12

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	if req.Email == "" {
		return nil, ErrInvalidInput
	}

	if req.Password == "" {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	user := User{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetUsers(ctx context.Context, q GetUserQuery) ([]User, error) {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 10
	}

	if q.Offset < 0 {
		q.Offset = 0
	}

	return s.repo.GetUsers(ctx, q)
}
