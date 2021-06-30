package v1

import (
	"context"
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

func TestNewSizeHandler(t *testing.T) {

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
			h, err := NewSizeHandler(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestSizeHandler_List(t *testing.T) {
	driverMock := drivertest.DriverMock{}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Size
		sendRequest     func(string) (*http.Response, error)
		expectedHandler string
		expectedBody    string
	}{
		"search all sizes": {
			expectedHandler: "ListSizes",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			driverReturn: []entity.Size{
				{ID: 0, Title: "Any"},
				{ID: 1, Title: "Short", MaxPages: util.Int16Ptr(34)},
				{ID: 6, Title: "Monument", MinPages: util.Int16Ptr(800)},
			},
			expectedBody: `[{"id":0,"title":"Any"},{"id":1,"title":"Short","maxPages":34},{"id":6,"title":"Monument","minPages":800}]`,
		},
		"invalid http method": {
			expectedHandler: "ListSizes",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
			expectedBody: "HTTP method POST is not allowed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler := sizeHandler{driver: &driverMock}

			r := sizeRoutes(&handler)

			server := httptest.NewServer(r)
			defer server.Close()

			driverMock.
				On(tt.expectedHandler, mock.MatchedBy(func(_ context.Context) bool { return true })).
				Return(tt.driverReturn, nil)

			resp, err := tt.sendRequest(fmt.Sprintf("%s%s", server.URL, tt.target))
			assert.NoError(t, err)

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody+"\n", string(body))
		})
	}
}