package storage

import "errors"

var (
	ErrUserExists   = errors.New("already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrItemConflict = errors.New("item conflict")
)
