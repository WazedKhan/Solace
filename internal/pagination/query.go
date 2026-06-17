package pagination

type QueryParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
