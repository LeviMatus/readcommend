package v1

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewEraHandler(t *testing.T) {

	tests := map[string]struct {
		driver       driver.Driver
		errAssertion assert.ErrorAssertionFunc
		valAssertion assert.ValueAssertionFunc
	}{
		"nil driver provided": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"handler created": {
			driver:       &drivertest.DriverMock{},
			errAssertion: assert.NoError,
			valAssertion: assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h, err := NewEraHandler(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestEraHandler_List(t *testing.T) {
	mockEra := entity.Era{
		ID:    1,
		Title: "Fantasy/SciFy",
	}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Era
		expectedHandler string
		expectedBody    string
		expectedCode    int
		expectedErr     error
		sendRequest     func(string) (*http.Response, error)
	}{
		"search all eras": {
			expectedHandler: "ListEras",
			target:          "/",
			driverReturn: []entity.Era{
				{ID: 1, Title: "Any"},
				{ID: 2, Title: "Classical", MaxYear: util.Int16Ptr(1969)},
				{ID: 3, Title: "Modern", MinYear: util.Int16Ptr(1970)},
			},
			expectedBody: `[{"id":1,"title":"Any"},{"id":2,"title":"Classical","maxYear":1969},{"id":3,"title":"Modern","minYear":1970}]`,
			expectedCode: 200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid http method": {
			expectedHandler: "ListEras",
			target:          "/",
			driverReturn:    []entity.Era{mockEra},
			expectedBody:    "HTTP method POST is not allowed",
			expectedCode:    400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
		},
		"driver returns error": {
			expectedHandler: "ListEras",
			target:          "/",
			driverReturn:    []entity.Era{},
			expectedBody:    "internal server error",
			expectedCode:    400,
			expectedErr:     errors.New("mock error returned from driver"),
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			driverMock := drivertest.DriverMock{}
			handler := eraHandler{driver: &driverMock}

			r := eraRoutes(&handler)

			server := httptest.NewServer(r)
			defer server.Close()

			driverMock.
				On(tt.expectedHandler, mock.MatchedBy(func(_ context.Context) bool { return true })).
				Return(tt.driverReturn, tt.expectedErr)

			resp, err := tt.sendRequest(fmt.Sprintf("%s%s", server.URL, tt.target))
			assert.NoError(t, err)

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody+"\n", string(body))
		})
	}
}
