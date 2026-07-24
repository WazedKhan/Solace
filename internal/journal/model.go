package journal

import "time"

type Mood struct {
	ID   string
	Name string
}

type MoodResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type JournalWithMood struct {
	Journal
	Mood *Mood
}

type JournalResponse struct {
	ID          string        `json:"id"`
	Mood        *MoodResponse `json:"mood,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ImageURL    *string       `json:"image_url,omitempty"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
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

type ConfirmUploadRequest struct {
	Key string `json:"key"`
}

type PresignUploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type PresignUploadRequest struct {
	ContentType string `json:"content_type"`
}
