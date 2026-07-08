package pagination

import (
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request) QueryParams {
	query := r.URL.Query()
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > MaxLimit {
		limit = MaxLimit
	}
	return QueryParams{Limit: limit, Offset: offset}
}
