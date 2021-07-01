package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {

	tests := map[string]struct {
		input          driver.Driver
		expectedRoutes []string
		errAssertion   assert.ErrorAssertionFunc
		valAssertion   assert.ValueAssertionFunc
	}{
		"code contract not met - nil driver": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"create server": {
			input:          &drivertest.DriverMock{},
			expectedRoutes: []string{"api"},
			errAssertion:   assert.NoError,
			valAssertion:   assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server, err := New(tt.input)
			tt.errAssertion(t, err)
			tt.valAssertion(t, server)
			if server != nil {
				var foundPatterns = map[string]struct{}{}
				for _, pattern := range tt.expectedRoutes {
					for _, subroute := range server.mux.Routes() {
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

func TestServer_Serve(t *testing.T) {

	tests := map[string]struct {
		input           *drivertest.DriverMock
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
				Limit:            util.Uint64Ptr(50),
			},
			target:      "/api/v1/books?title=The+Silmarillion&max-year=1980&min-year=1970&max-pages=400&min-pages=300&genres=2&genres=6&authors=42&authors=43&limit=50",
			expectedErr: nil,
			expect:      []entity.Book{{}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			apiServer, err := New(tt.input)
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
		apiServer, err := New(&dbDriver)
		assert.NoError(t, err)
		assert.NotNil(t, apiServer)
		testServer := httptest.NewServer(apiServer.mux)
		response, err := testServer.Client().Head(fmt.Sprintf("%s/api/v1/books", testServer.URL))
		assert.NoError(t, err)
		assert.Equal(t, 400, response.StatusCode)
	})
}
