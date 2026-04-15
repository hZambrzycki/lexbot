package shared

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrEmptyField         = errors.New("required field is empty")
	ErrInvalidState       = errors.New("invalid state")
	ErrNotFound           = errors.New("entity not found")
	ErrInvalidAssociation = errors.New("invalid association")
)
