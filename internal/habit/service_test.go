package habit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mock lives in the test file
type mockHabitRepo struct {
	// store what each method should return
	habit *Habit
	err   error
}

func (m *mockHabitRepo) GetHabitByID(ctx context.Context, habitID string) (*Habit, error) {
	return m.habit, m.err
}

func (m *mockHabitRepo) CreateHabit(ctx context.Context, habit Habit) (*Habit, error) {
	return m.habit, m.err
}

func (m *mockHabitRepo) GetHabitsByUserID(ctx context.Context, qParams HabitQueryRequest) ([]*Habit, error) {
	return nil, m.err
}

func (m *mockHabitRepo) CheckInTx(ctx context.Context, habitID string, newStreak int, now time.Time) (*int, error) {
	return &newStreak, m.err
}

func TestHabit_CheckIn(t *testing.T) {
	t.Run("first ever check-in", func(t *testing.T) {
		repo := &mockHabitRepo{
			habit: &Habit{
				ID:            "habit-1",
				UserID:        "user-1",
				CurrentStreak: 0,
				LastCheckedAt: nil,
			},
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), repo.habit.UserID, repo.habit.ID)
		assert.Nil(t, err)
		assert.NotNil(t, streak)
		assert.Equal(t, 1, *streak)
	})

	t.Run("already checked in today", func(t *testing.T) {
		now := time.Now()
		repo := &mockHabitRepo{
			habit: &Habit{
				ID:            "habit-1",
				UserID:        "user-1",
				CurrentStreak: 1,
				LastCheckedAt: &now,
			},
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), repo.habit.UserID, repo.habit.ID)
		assert.Nil(t, streak)
		assert.Equal(t, ErrAlreadyChecked, err)
	})

	t.Run("checked in yesterday increments streak", func(t *testing.T) {
		yesterday := time.Now().AddDate(0, 0, -1)
		repo := &mockHabitRepo{
			habit: &Habit{
				ID:            "habit-1",
				UserID:        "user-1",
				CurrentStreak: 3,
				LastCheckedAt: &yesterday,
			},
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), repo.habit.UserID, repo.habit.ID)
		assert.Nil(t, err)
		assert.NotNil(t, streak)
		assert.Equal(t, 4, *streak)
	})

	t.Run("streak broken resets to 1", func(t *testing.T) {
		twoDaysAgo := time.Now().AddDate(0, 0, -2)
		repo := &mockHabitRepo{
			habit: &Habit{
				ID:            "habit-1",
				UserID:        "user-1",
				CurrentStreak: 5,
				LastCheckedAt: &twoDaysAgo,
			},
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), repo.habit.UserID, repo.habit.ID)
		assert.Nil(t, err)
		assert.NotNil(t, streak)
		assert.Equal(t, 1, *streak)
	})

	t.Run("habit not found", func(t *testing.T) {
		repo := &mockHabitRepo{
			err: ErrNotFound,
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), "user-1", "missing-habit")
		assert.Nil(t, streak)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("habit belongs to different user", func(t *testing.T) {
		repo := &mockHabitRepo{
			habit: &Habit{
				ID:     "habit-1",
				UserID: "other-user",
			},
		}
		service := NewService(repo)

		streak, err := service.CheckIn(context.Background(), "user-1", repo.habit.ID)
		assert.Nil(t, streak)
		assert.Equal(t, ErrInvalidHabitID, err)
	})
}
