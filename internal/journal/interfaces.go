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
	) ([]*Journal, error)
	GetJournalByID(ctx context.Context, userID, journalID string) (*Journal, error)
	UpdateJournalByID(ctx context.Context, inp Journal) (*Journal, error)
	SoftDeleteJournal(ctx context.Context, journalID, userID string) error
}
