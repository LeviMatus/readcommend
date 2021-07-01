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

func TestStringPtr(t *testing.T) {
	var expected = "foobar"
	actual := StringPtr(expected)
	if actual == nil {
		t.FailNow()
	}
	assert.Equal(t, expected, *actual)
}

func TestUint64Ptr(t *testing.T) {
	var expected uint64 = 1
	actual := Uint64Ptr(expected)
	if actual == nil {
		t.FailNow()
	}
	assert.Equal(t, expected, *actual)
}
