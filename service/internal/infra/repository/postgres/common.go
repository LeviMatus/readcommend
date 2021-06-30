package postgres

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidDependency = errors.New("expected a non-nil sql Database connection")
)
