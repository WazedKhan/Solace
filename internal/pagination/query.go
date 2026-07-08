package pagination

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

type QueryParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
