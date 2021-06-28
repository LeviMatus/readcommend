package postgres

import (
	"github.com/pkg/errors"
)

const (
	defaultLimit uint64 = 25
)

var (
	ErrInvalidDependency = errors.New("expected a non-nil sql Database connection")
)
