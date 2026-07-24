package journal

import "errors"

var (
	ErrTitleEmpty    = errors.New("title can't be empty")
	ErrDescription   = errors.New("description can't be empty")
	ErrInvalidStatus = errors.New("invalid journal status")
	ErrForbidden     = errors.New("resource does not belong to authenticated user")
	ErrImageNotFound = errors.New("image not found in storage")
	ErrInvalidKey    = errors.New("malformed storage key")
)
