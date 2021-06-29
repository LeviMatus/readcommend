package entity

import (
	"errors"
)

var (
	ErrInvalidQueryParam = errors.New("invalid URL query parameter provided")
)
