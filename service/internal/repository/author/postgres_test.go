package author

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresRepository(t *testing.T) {

	var db sql.DB

	tests := map[string]struct {
		input        *sql.DB
		expect       *authorPostgresRepo
		errAssertion assert.ErrorAssertionFunc
	}{
		"error on nil input": {
			input:        nil, // I could just omit this line, but I'll explicitly set nil for clarity.
			expect:       nil,
			errAssertion: assert.Error,
		},
		"successful create repository": {
			input:        &db,
			expect:       &authorPostgresRepo{db: &db},
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := NewPostgresRepository(tt.input)
			assert.Equal(t, tt.expect, actual)
			tt.errAssertion(t, err)
		})
	}

}
