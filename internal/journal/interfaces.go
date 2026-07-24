package journal

import (
	"context"

	"github.com/WazedKhan/Solace/internal/pagination"
)

type JournalRepository interface {
	CreateJournal(ctx context.Context, inp Journal) (*Journal, error)
	GetJournalsByUser(
		ctx context.Context,
		userID, status string,
		pag pagination.QueryParams,
	) ([]*JournalWithMood, error)
	GetJournalByID(ctx context.Context, userID, journalID string) (*JournalWithMood, error)
	UpdateJournalByID(ctx context.Context, inp Journal) (*Journal, error)
	SoftDeleteJournal(ctx context.Context, journalID, userID string) error
	GetMoodById(ctx context.Context, moodID string) (*Mood, error)
	UpdateImageURL(ctx context.Context, journalID, userID, key string) error
}
