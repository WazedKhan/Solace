package habit

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateHabit(ctx context.Context, req HabitRequest, userId string) (*Habit, error) {
	if userId == "" {
		return nil, errors.New("user not found")
	}
	habit := Habit{
		ID:        uuid.NewString(),
		UserID:    userId,
		Title:     req.Title,
		ImageUrl:  req.ImageUrl,
		CreatedAt: time.Now(),
	}

	res, err := s.repo.CreateHabit(ctx, habit)
	if err != nil {
		return nil, err
	}

	return res, nil
}
