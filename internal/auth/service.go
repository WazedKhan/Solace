package auth

import (
	"context"
	"time"

	jwt_token "github.com/WazedKhan/Solace/internal/auth/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	generator *jwt_token.Generator
}

func NewService(repo *Repository, generator *jwt_token.Generator) *Service {
	return &Service{
		repo:      repo,
		generator: generator,
	}
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generator.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
