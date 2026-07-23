package journal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/WazedKhan/Solace/internal/pagination"
	storage "github.com/WazedKhan/Solace/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	repo    JournalRepository
	storage storage.Storage
}

func NewService(repo JournalRepository, store storage.Storage) *Service {
	return &Service{
		repo:    repo,
		storage: store,
	}
}

func (s *Service) CreateJournal(ctx context.Context, req CreateJournalRequest, userId string) (*Journal, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrTitleEmpty
	}
	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, ErrDescription
	}

	status := JournalStatus(StatusDraft)
	if req.Status != nil {
		status = JournalStatus(strings.TrimSpace(*req.Status))
	}
	if !status.IsValidStatus() {
		return nil, ErrInvalidStatus
	}

	journal := Journal{
		ID:          uuid.NewString(),
		UserID:      userId,
		MoodID:      req.MoodID,
		Title:       title,
		Description: description,
		ImageURL:    req.ImageURL,
		Status:      string(status),
		CreatedAt:   time.Now(),
	}

	res, err := s.repo.CreateJournal(ctx, journal)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetJournalsByUser can used for draft list as well by status key
func (s *Service) GetJournalsByUser(
	ctx context.Context,
	userID, status string,
	peg pagination.QueryParams,
) ([]*JournalWithMood, error) {
	trimmedStatus := JournalStatus(strings.TrimSpace(status))
	if !trimmedStatus.IsValidStatus() {
		return nil, ErrInvalidStatus
	}
	res, err := s.repo.GetJournalsByUser(ctx, userID, status, peg)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *Service) GetJournalByID(ctx context.Context, userID, journalID string) (*JournalWithMood, error) {
	res, err := s.repo.GetJournalByID(ctx, userID, journalID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) UpdateJournalByID(
	ctx context.Context,
	userID, journalID string,
	req UpdateJournalRequest,
) (*Journal, error) {
	existing, err := s.repo.GetJournalByID(ctx, userID, journalID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}

	if req.Description != nil {
		existing.Description = *req.Description
	}

	if req.MoodID != nil {
		existing.MoodID = req.MoodID
	}

	if req.Status != nil {
		if !JournalStatus(*req.Status).IsValidStatus() {
			return nil, ErrInvalidStatus
		}
		existing.Status = *req.Status
	}

	if req.ImageURL != nil {
		existing.ImageURL = req.ImageURL
	}
	now := time.Now()

	journal := Journal{
		ID:          journalID,
		UserID:      userID,
		Title:       existing.Title,
		MoodID:      existing.MoodID,
		Description: existing.Description,
		ImageURL:    existing.ImageURL,
		Status:      existing.Status,
		UpdatedAt:   &now,
	}

	res, err := s.repo.UpdateJournalByID(ctx, journal)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) SoftDeleteJournal(ctx context.Context, userID, journalID string) error {
	err := s.repo.SoftDeleteJournal(ctx, journalID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ConfirmUpload(ctx context.Context, userID, journalID, key string) error {
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return ErrInvalidKey
	}
	keyUserID := parts[1]
	if keyUserID != userID {
		return ErrForbidden
	}

	ok, err := s.storage.Exists(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to verify upload: %w", err)
	}
	if !ok {
		return ErrImageNotFound
	}

	err = s.repo.UpdateImageURL(ctx, journalID, userID, key)
	if err != nil {
		return err
	}

	return nil
}
