package v1

import (
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	tests := map[string]struct {
		driver       driver.Driver
		errAssertion assert.ErrorAssertionFunc
		valAssertion assert.ValueAssertionFunc
	}{
		"error - nil driver provided": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"create v1 router": {
			driver:       &drivertest.DriverMock{},
			errAssertion: assert.NoError,
			valAssertion: assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewRouter(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, r)
		})
	}
}
