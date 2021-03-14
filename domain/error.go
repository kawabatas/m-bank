package domain

import "errors"

// predefined errors.
var (
	ErrDuplicateEntity = errors.New("duplicate entity")
	ErrNoSuchEntity    = errors.New("no such entity")
	ErrInvalidUUID     = errors.New("invalid uuid")
	ErrShortBalance    = errors.New("short balance")
	ErrInvalidParam    = errors.New("invalid param")
)
