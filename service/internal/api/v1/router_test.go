package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	}{
		"invalid books method": {
			target:       "/books",
			method:       "POST",
			expectedCode: 405,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Head(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {},
		},
		"get books method": {
			target:       "/books",
			method:       "GET",
			expectedCode: 200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - min-pages": {
			target:       "/books?min-pages=0",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - max-pages": {
			target:       "/books?max-pages=10001",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - min-year": {
			target:       "/books?min-year=1799",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - max-year": {
			target:       "/books?min-year=2101",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - limit": {
			target:       "/books?limit=0",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - authors": {
			target:       "/books?authors=1,beta,3",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
		"invalid books param - genres": {
			target:       "/books?genres=1,beta,3",
			method:       "GET",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			setDriverExpectations: func(d *drivertest.DriverMock) {
				d.On("SearchBooks",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					book.SearchInput{}).Return(([]entity.Book)(nil), nil)
			},
		},
	}

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
