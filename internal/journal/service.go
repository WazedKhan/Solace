package journal

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo JournalRepository
}

func NewService(repo JournalRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateJournal(ctx context.Context, req CreateJournalRequest, userId string) (*Journal, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrTitleEmpty
	}
	if strings.TrimSpace(req.Description) == "" {
		return nil, ErrDescription
	}

	status := "draft"
	if req.Status != nil {
		status = *req.Status
	}
	switch status {
	case "draft", "published":
	default:
		return nil, ErrInvalidStatus
	}

	moodID := req.MoodID
	if moodID != nil {
		mood, err := s.repo.GetMoodById(ctx, *req.MoodID)
		if err != nil {
			return nil, err
		}
		moodID = &mood.ID
	}

	journal := Journal{
		ID:          uuid.NewString(),
		UserID:      userId,
		MoodID:      moodID,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	res, err := s.repo.CreateJournal(ctx, journal)
	if err != nil {
		return nil, err
	}

	return res, nil
}
