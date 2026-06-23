package habit

import "errors"

var (
	ErrInvalidHabitID = errors.New("no habit found with given id")
	ErrAlreadyChecked = errors.New("already checked today")
	ErrNotFound       = errors.New("resource not found")
)
