package habit

import (
	"context"
	"errors"
	"time"

	"github.com/WazedKhan/Solace/internal/pagination"
	"github.com/WazedKhan/Solace/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo HabitRepository
}

func NewService(repo HabitRepository) *Service {
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

func (s *Service) GetHabitsByUserID(
	ctx context.Context,
	userId string,
	params pagination.QueryParams,
) ([]*Habit, error) {
	qParams := HabitQueryRequest{
		QueryParams: params,
		UserID:      userId,
	}
	res, err := s.repo.GetHabitsByUserID(ctx, qParams)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *Service) CheckIn(ctx context.Context, userId, habitId string) (*int, error) {
	habit, err := s.repo.GetHabitByID(ctx, habitId)
	if err != nil {
		return nil, err
	}

	// check if the habit belong to requested user
	if habit.UserID != userId {
		return nil, ErrInvalidHabitID
	}

	// fetch the last_checked_at
	today := time.Now()
	switch {
	case habit.LastCheckedAt == nil:
		habit.CurrentStreak = 1

	case utils.IsSameDay(*habit.LastCheckedAt, today):
		return nil, ErrAlreadyChecked

	case utils.IsSameDay(*habit.LastCheckedAt, today.AddDate(0, 0, -1)):
		habit.CurrentStreak++

	default:
		habit.CurrentStreak = 1
	}

	currentStreak, err := s.repo.CheckInTx(ctx, habitId, habit.CurrentStreak, today)
	if err != nil {
		return nil, err
	}

	return currentStreak, nil
}
