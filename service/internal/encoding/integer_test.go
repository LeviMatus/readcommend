package encoding

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockScanner struct {
	value interface{}
}

func (m mockScanner) Scan(s sql.Scanner) error {
	return s.Scan(m.value)
}

func TestNullInt16_Scan(t *testing.T) {
	tests := map[string]struct {
		scanner     mockScanner
		expectedInt int16
		expectValid assert.BoolAssertionFunc
		assertErr   assert.ErrorAssertionFunc
	}{
		"int16 is nil": {
			scanner:     mockScanner{value: nil},
			expectedInt: 0,
			expectValid: assert.False,
			assertErr:   assert.NoError,
		},
		"int16 is not nil": {
			scanner:     mockScanner{value: int64(1)},
			expectedInt: 1,
			expectValid: assert.True,
			assertErr:   assert.NoError,
		},
		"scanned type is of the right type": {
			scanner:     mockScanner{value: "invalid"},
			expectedInt: 0,
			expectValid: assert.False,
			assertErr:   assert.Error,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var actual NullInt16
			err := tt.scanner.Scan(&actual)
			tt.assertErr(t, err)
			assert.Equal(t, tt.expectedInt, actual.Int16)
			tt.expectValid(t, actual.Valid)
		})
	}
}

func TestNullInt16_Value(t *testing.T) {
	tests := map[string]struct {
		valuer       driver.Valuer
		expected     driver.Value
		errAssertion assert.ErrorAssertionFunc
	}{
		"value is nil": {
			valuer:       &NullInt16{},
			expected:     driver.Value(nil),
			errAssertion: assert.NoError,
		},
		"value is not nil": {
			valuer:       &NullInt16{Int16: 1, Valid: true},
			expected:     driver.Value(int16(1)),
			errAssertion: assert.NoError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := tt.valuer.Value()
			tt.errAssertion(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
