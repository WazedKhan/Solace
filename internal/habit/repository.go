package habit

import (
	"context"
	"errors"
	"time"

	"github.com/WazedKhan/Solace/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrFailedWriting = errors.New("failed to write into db")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateHabit(ctx context.Context, habit Habit) (*Habit, error) {
	query := `
		INSERT INTO habits (id, user_id, title, image_url)
		VALUES($1, $2, $3, $4)
		RETURNING id, title, user_id, image_url, current_streak, last_checked_at, created_at, updated_at
	`
	var habitRes Habit
	err := r.db.QueryRow(
		ctx,
		query,
		habit.ID,
		habit.UserID,
		habit.Title,
		habit.ImageUrl,
	).Scan(
		&habitRes.ID,
		&habitRes.Title,
		&habitRes.UserID,
		&habitRes.ImageUrl,
		&habitRes.CurrentStreak,
		&habitRes.LastCheckedAt,
		&habitRes.CreatedAt,
		&habitRes.UpdatedAt,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}

	return &habitRes, nil
}

func (r *Repository) GetHabitsByUserID(ctx context.Context, qParams HabitQueryRequest) ([]*Habit, error) {
	query := `
		SELECT id, title, user_id, image_url, current_streak, last_checked_at, created_at, updated_at
		FROM habits
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(
		ctx,
		query,
		qParams.UserID,
		qParams.QueryParams.Limit,
		qParams.QueryParams.Offset,
	)
	if err != nil {
		return []*Habit{}, utils.MapPostgresError(err)
	}
	defer rows.Close()

	var habits []*Habit
	for rows.Next() {
		var habit Habit
		err := rows.Scan(
			&habit.ID,
			&habit.Title,
			&habit.UserID,
			&habit.ImageUrl,
			&habit.CurrentStreak,
			&habit.LastCheckedAt,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return []*Habit{}, utils.MapPostgresError(err)
		}
		habits = append(habits, &habit)
	}
	if err := rows.Err(); err != nil {
		return []*Habit{}, utils.MapPostgresError(err)
	}

	return habits, nil
}

func (r *Repository) GetHabitByID(ctx context.Context, habitId string) (*Habit, error) {
	query := `
		SELECT id, title, user_id, image_url, current_streak, last_checked_at, created_at, updated_at
		FROM habits
		WHERE id=$1
	`
	var habit Habit
	err := r.db.QueryRow(ctx, query, habitId).Scan(
		&habit.ID,
		&habit.Title,
		&habit.UserID,
		&habit.ImageUrl,
		&habit.CurrentStreak,
		&habit.LastCheckedAt,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	return &habit, err
}

func (r *Repository) CheckHabitByID(ctx context.Context, habitID string, cStreak int) (*int, error) {
	query := `
		UPDATE habits
		SET last_checked_at=$1,
			current_streak=$2,
			updated_at=$3
		WHERE id=$4
		RETURNING current_streak
	`
	var currentStreak int
	err := r.db.QueryRow(
		ctx, query,
		time.Now().UTC(),
		cStreak,
		time.Now(),
		habitID,
	).Scan(
		&currentStreak,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	return &currentStreak, nil
}

func (r *Repository) CreateHabitCheckingLog(ctx context.Context, habitID string) error {
	query := `
		INSERT INTO habit_checking(habit_id, checked_date)
		VALUES($1, $2)
	`
	_, err := r.db.Exec(ctx, query, habitID, time.Now())
	if err != nil {
		return utils.MapPostgresError(err)
	}
	return nil
}
