package entity

import (
	"errors"
)

var (
	// ErrInvalidQueryParam occurs when an invalid parameter range or type was provided.
	ErrInvalidQueryParam = errors.New("invalid URL query parameter provided")
)
