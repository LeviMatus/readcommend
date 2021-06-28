package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt16Ptr(t *testing.T) {
	var expected int16 = 15
	actual := Int16Ptr(expected)
	if actual == nil {
		t.FailNow()
	}
	assert.Equal(t, expected, *actual)
}
