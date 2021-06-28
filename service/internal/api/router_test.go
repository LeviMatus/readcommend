package api

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
	"github.com/LeviMatus/readcommend/service/pkg/config"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {

	tests := map[string]struct {
		input        driver.Driver
		config       config.API
		errAssertion assert.ErrorAssertionFunc
		valAssertion assert.ValueAssertionFunc
	}{
		"code contract not met - nil driver": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"create server": {
			input:        &drivertest.DriverMock{},
			errAssertion: assert.NoError,
			valAssertion: assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server, err := New(tt.input, tt.config)
			tt.errAssertion(t, err)
			tt.valAssertion(t, server)
		})
	}
}

func TestServer_Serve(t *testing.T) {

	tests := map[string]struct {
		input           *drivertest.DriverMock
		config          config.API
		expectedHandler string
		expectedParams  book.SearchInput
		target          string
		expectedErr     error
		expect          interface{}
	}{
		"search all books": {
			input:           &drivertest.DriverMock{},
			expectedHandler: "SearchBooks",
			target:          "/api/v1/books",
			expectedErr:     nil,
			expect:          []entity.Book{{}},
		},
		"search for specific books": {
			input:           &drivertest.DriverMock{},
			expectedHandler: "SearchBooks",
			expectedParams: book.SearchInput{
				Title:            util.StringPtr("The Silmarillion"),
				MaxYearPublished: util.Int16Ptr(1980),
				MinYearPublished: util.Int16Ptr(1970),
				MaxPages:         util.Int16Ptr(400),
				MinPages:         util.Int16Ptr(300),
				GenreIDs:         []int16{2, 6},
				AuthorIDs:        []int16{42, 43},
				Limit:            util.Uint16Ptr(50),
			},
			target:      "/api/v1/books?title=The+Silmarillion&max_year=1980&min_year=1970&max_pages=400&min_pages=300&genres=2&genres=6&authors=42&authors=43&limit=50",
			expectedErr: nil,
			expect:      []entity.Book{{}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			apiServer, err := New(tt.input, tt.config)
			assert.NoError(t, err)
			assert.NotNil(t, apiServer)

			tt.input.On(tt.expectedHandler,
				mock.MatchedBy(func(_ context.Context) bool { return true }),
				tt.expectedParams).Return(tt.expect, tt.expectedErr)

			testServer := httptest.NewServer(apiServer.mux)
			response, err := testServer.Client().Get(fmt.Sprintf("%s%s", testServer.URL, tt.target))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}

	t.Run("invalid HTTP method", func(t *testing.T) {
		dbDriver := drivertest.DriverMock{}
		apiServer, err := New(&dbDriver, config.API{})
		assert.NoError(t, err)
		assert.NotNil(t, apiServer)
		testServer := httptest.NewServer(apiServer.mux)
		response, err := testServer.Client().Head(fmt.Sprintf("%s/api/v1/books", testServer.URL))
		assert.NoError(t, err)
		assert.Equal(t, 405, response.StatusCode)
	})
}
