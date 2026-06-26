package habit

import (
	"context"
	"time"
)

type HabitRepository interface {
	CreateHabit(ctx context.Context, habit Habit) (*Habit, error)
	GetHabitByID(ctx context.Context, habitID string) (*Habit, error)
	GetHabitsByUserID(ctx context.Context, qParams HabitQueryRequest) ([]*Habit, error)
	CheckInTx(ctx context.Context, habitID string, newStreak int, now time.Time) (*int, error)
}
