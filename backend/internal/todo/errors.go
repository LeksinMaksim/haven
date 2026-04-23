package todo

import "errors"

var (
	ErrNotFound   = errors.New("todo not found")
	ErrEmptyTitle = errors.New("title cannot be empty")
)
