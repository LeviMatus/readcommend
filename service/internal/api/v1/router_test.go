package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetBooks(t *testing.T) {
	driverMock := drivertest.DriverMock{}

	tests := map[string]struct {
		target                string
		method                string
		expectedCode          int
		sendRequest           func(url string) (*http.Response, error)
		setDriverExpectations func(d *drivertest.DriverMock)
	}{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setDriverExpectations(&driverMock)
			router, err := NewRouter(&driverMock)
			assert.NoError(t, err)

			testServer := httptest.NewServer(router)
			defer testServer.Close()

			resp, err := tt.sendRequest(fmt.Sprintf("%s%s", testServer.URL, tt.target))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
