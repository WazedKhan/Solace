package journal

import "errors"

var (
	ErrTitleEmpty    = errors.New("title can't be empty")
	ErrDescription   = errors.New("description can't be empty")
	ErrInvalidStatus = errors.New("invalid journal status")
)
