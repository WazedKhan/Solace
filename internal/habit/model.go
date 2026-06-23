package habit

import (
	"time"

	"github.com/WazedKhan/Solace/internal/pagination"
)

type Habit struct {
	ID            string
	Title         string
	UserID        string
	ImageUrl      *string
	CurrentStreak int
	LastCheckedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type HabitResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	ImageUrl      *string    `json:"image_url,omitempty"`
	CurrentStreak int        `json:"current_streak"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type HabitRequest struct {
	Title    string  `json:"title"`
	ImageUrl *string `json:"image_url"`
}

type CheckInResponse struct {
	CurrentStreak int `json:"current_streak"`
}

type HabitQueryRequest struct {
	QueryParams pagination.QueryParams
	UserID      string
}

type CheckInRequest struct {
	HabitID string `json:"habit_id"`
}
