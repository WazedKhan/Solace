package journal

import "time"

type Mood struct {
	ID   string
	Name string
}

type Journal struct {
	ID          string
	UserID      string
	MoodID      *string
	Title       string
	Description string
	ImageURL    *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type JournalResponse struct {
	ID          string     `json:"id"`
	Mood        *Mood      `json:"mood,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageURL    *string    `json:"image_url,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateJournalRequest struct {
	MoodID      *string `json:"mood_id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type UpdateJournalRequest struct {
	MoodID      *string `json:"mood_id,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	ImageURL    *string `json:"image_url"`
}
