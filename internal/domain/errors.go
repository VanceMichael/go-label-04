package domain

import "errors"

var (
	ErrNotFound  = errors.New("airbridge: not found")
	ErrConflict  = errors.New("airbridge: conflict")
	ErrForbidden = errors.New("airbridge: forbidden")
	ErrInvalid   = errors.New("airbridge: invalid input")
	ErrExpired   = errors.New("airbridge: expired")
	ErrRevoked   = errors.New("airbridge: revoked")
	ErrCapacity  = errors.New("airbridge: capacity unavailable")
	ErrState     = errors.New("airbridge: invalid state transition")
)
