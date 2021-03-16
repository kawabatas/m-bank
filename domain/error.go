package domain

import "errors"

// predefined errors.
var (
	ErrNoSuchEntity  = errors.New("no such entity")
	ErrDuplicateUUID = errors.New("duplicate uuid")
	ErrInvalidUUID   = errors.New("invalid uuid")
	ErrShortBalance  = errors.New("short balance")
	ErrInvalidParam  = errors.New("invalid param")
)
