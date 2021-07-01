package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt16InRange(t *testing.T) {
	tests := map[string]struct {
		target int16
		min    int16
		max    int16
		expect assert.BoolAssertionFunc
	}{
		"target in range": {
			target: 4,
			min:    3,
			max:    5,
			expect: assert.True,
		},
		"target below range": {
			target: 2,
			min:    3,
			max:    5,
			expect: assert.False,
		},
		"target above range": {
			target: 6,
			min:    3,
			max:    5,
			expect: assert.False,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.expect(t, Int16InRange(tt.target, tt.min, tt.max))
		})
	}
}
