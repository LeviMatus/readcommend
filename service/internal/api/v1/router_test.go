package v1

import (
	"strings"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	tests := map[string]struct {
		driver         driver.Driver
		expectedRoutes []string
		errAssertion   assert.ErrorAssertionFunc
		valAssertion   assert.ValueAssertionFunc
	}{
		"error - nil driver provided": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"create v1 router": {
			driver:         &drivertest.DriverMock{},
			expectedRoutes: []string{"books", "authors", "eras", "sizes", "genres"},
			errAssertion:   assert.NoError,
			valAssertion:   assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewRouter(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, r)
			if r != nil {
				var foundPatterns = map[string]struct{}{}
				for _, pattern := range tt.expectedRoutes {
					for _, subroute := range r.Routes() {
						if strings.Contains(subroute.Pattern, pattern) {
							foundPatterns[pattern] = struct{}{}
						}
					}
					_, found := foundPatterns[pattern]
					assert.True(t, found)
				}
			}
		})
	}
}
