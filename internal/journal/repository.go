package journal

import (
	"context"

	"github.com/WazedKhan/Solace/internal/pagination"
	"github.com/WazedKhan/Solace/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateJournal(ctx context.Context, inp Journal) (*Journal, error) {
	query := `
		INSERT INTO journals(id, user_id, mood_id, title, description, status, image_url, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, mood_id, title, description, status, image_url, created_at
	`
	var journal Journal
	err := r.db.QueryRow(
		ctx,
		query,
		inp.ID,
		inp.UserID,
		inp.MoodID,
		inp.Title,
		inp.Description,
		inp.Status,
		inp.ImageURL,
		inp.CreatedAt,
	).Scan(
		&journal.ID,
		&journal.UserID,
		&journal.MoodID,
		&journal.Title,
		&journal.Description,
		&journal.Status,
		&journal.ImageURL,
		&journal.CreatedAt,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}

	return &journal, nil
}

func (r *Repository) GetJournalsByUser(
	ctx context.Context,
	userID, status string,
	peg pagination.QueryParams,
) ([]*JournalWithMood, error) {
	query := `
	SELECT
		j.id,
		j.user_id,
		j.title,
		j.description,
		j.image_url,
		j.status,
		j.created_at,
		j.updated_at,

		m.id,
		m.name

	FROM journals j
	LEFT JOIN moods m
		ON j.mood_id = m.id

	WHERE j.user_id = $1
	AND j.deleted_at IS NULL
	AND j.status = $2

	ORDER BY j.created_at DESC
	LIMIT $3 OFFSET $4;
	`

	rows, err := r.db.Query(ctx, query, userID, status, peg.Limit, peg.Offset)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	var journals []*JournalWithMood

	for rows.Next() {

		var (
			j        JournalWithMood
			moodID   *string
			moodName *string
		)
		err := rows.Scan(
			&j.ID,
			&j.UserID,
			&j.Title,
			&j.Description,
			&j.ImageURL,
			&j.Status,
			&j.CreatedAt,
			&j.UpdatedAt,
			&moodID,
			&moodName,
		)
		if err != nil {
			return nil, utils.MapPostgresError(err)
		}

		if moodID != nil {
			j.Mood = &Mood{
				ID:   *moodID,
				Name: *moodName,
			}
		}
		journals = append(journals, &j)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.MapPostgresError(err)
	}

	return journals, nil
}

func (r *Repository) GetJournalByID(
	ctx context.Context,
	userID, journalID string,
) (*JournalWithMood, error) {
	query := `
	SELECT
		j.id,
		j.user_id,
		j.title,
		j.description,
		j.image_url,
		j.status,
		j.created_at,
		j.updated_at,

		m.id,
		m.name

	FROM journals j
	LEFT JOIN moods m
		ON j.mood_id = m.id
	WHERE
		j.user_id = $1
		AND j.id = $2
		AND j.deleted_at IS NULL;
	`
	var (
		j        JournalWithMood
		moodID   *string
		moodName *string
	)
	err := r.db.QueryRow(ctx, query, userID, journalID).Scan(
		&j.ID,
		&j.UserID,
		&j.Title,
		&j.Description,
		&j.ImageURL,
		&j.Status,
		&j.CreatedAt,
		&j.UpdatedAt,
		&moodID,
		&moodName,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	if moodID != nil {
		j.Mood = &Mood{
			ID:   *moodID,
			Name: *moodName,
		}
	}

	return &j, nil
}

func (r *Repository) UpdateJournalByID(ctx context.Context, inp Journal) (*Journal, error) {
	query := `
	UPDATE journals
	SET
		title = $1,
		description = $2,
		status = $3,
		mood_id = $4,
		image_url = $5,
		updated_at = $6
	WHERE user_id = $7 AND id = $8 AND deleted_at IS NULL
	RETURNING id, user_id, mood_id, title, description, image_url, status, created_at, updated_at;
	`
	var j Journal
	err := r.db.QueryRow(
		ctx,
		query,
		inp.Title,
		inp.Description,
		inp.Status,
		inp.MoodID,
		inp.ImageURL,
		inp.UpdatedAt,
		inp.UserID,
		inp.ID,
	).Scan(
		&j.ID,
		&j.UserID,
		&j.MoodID,
		&j.Title,
		&j.Description,
		&j.ImageURL,
		&j.Status,
		&j.CreatedAt,
		&j.UpdatedAt,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}

	return &j, nil
}

func (r *Repository) SoftDeleteJournal(ctx context.Context, journalID, userID string) error {
	query := `
	UPDATE journals
	SET
		deleted_at = NOW()
	WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL;
	`
	cmdTag, err := r.db.Exec(ctx, query, userID, journalID)
	if err != nil {
		return utils.MapPostgresError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return utils.ErrNotFound
	}
	return nil
}

func (r *Repository) GetMoodById(ctx context.Context, moodID string) (*Mood, error) {
	query := `
	SELECT id, name FROM moods WHERE id=$1;
	`
	var mood Mood
	err := r.db.QueryRow(ctx, query, moodID).Scan(
		&mood.ID,
		&mood.Name,
	)
	if err != nil {
		return nil, utils.MapPostgresError(err)
	}
	return &mood, nil
}
