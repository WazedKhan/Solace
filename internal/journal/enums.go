package journal

type JournalStatus string

const (
	StatusDraft     JournalStatus = "draft"
	StatusPublished JournalStatus = "published"
)

// this is better to move into different file which is dedicated for validation in future
func (s JournalStatus) IsValidStatus() bool {
	switch s {
	case StatusDraft, StatusPublished:
		return true
	default:
		return false
	}
}
