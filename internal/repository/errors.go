package repository

import "errors"

var (
	// ErrNotFound is returned when an entity with the given ID does not exist.
	ErrNotFound = errors.New("repository: not found")

	// ErrAlreadyExists is returned when creating an entity whose ID already exists.
	ErrAlreadyExists = errors.New("repository: already exists")
)
