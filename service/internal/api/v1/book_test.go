package v1

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/drivertest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBookHandler(t *testing.T) {

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
			h, err := NewBookHandler(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestBookHandler_List(t *testing.T) {
	expectedJson := `[{"id":1,"title":"The Silmarillion","yearPublished":1977,"rating":3.9,"pages":365,"genre":{"id":2,"title":"Fantasy/SciFi"},"author":{"id":42,"firstName":"John","lastName":"Tolkien"}}]`

	mockBook := entity.Book{
		ID:            1,
		Title:         "The Silmarillion",
		YearPublished: 1977,
		Rating:        3.9,
		Pages:         365,
		Genre: entity.Genre{
			ID:    2,
			Title: "Fantasy/SciFi",
		},
		Author: entity.Author{
			ID:        42,
			FirstName: "John",
			LastName:  "Tolkien",
		},
	}

	tests := map[string]struct {
		target          string
		expectedParams  book.SearchInput
		driverReturn    []entity.Book
		expectedHandler string
		expectedBody    string
		expectedCode    int
		expectedErr     error
		sendRequest     func(string) (*http.Response, error)
	}{
		"search all books": {
			expectedHandler: "SearchBooks",
			target:          "/",
			driverReturn:    []entity.Book{mockBook},
			expectedBody:    expectedJson,
			expectedCode:    200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"driver returns error": {
			expectedHandler: "SearchBooks",
			target:          "/",
			driverReturn:    []entity.Book{mockBook},
			expectedBody:    "internal server error",
			expectedCode:    400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			expectedErr: errors.New("mock internal error from driver"),
		},
		"search for specific books": {
			expectedHandler: "SearchBooks",
			target:          "/?title=The+Silmarillion&max-year=1980&min-year=1970&max-pages=400&min-pages=300&genres=2&genres=6&authors=42&authors=43&limit=50",
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
			driverReturn: []entity.Book{mockBook},
			expectedBody: expectedJson,
			expectedCode: 200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid http method": {
			expectedHandler: "ListAuthors",
			target:          "/",
			expectedBody:    "HTTP method POST is not allowed",
			expectedCode:    400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
		},
		"invalid books param - min-pages": {
			target:       "/?min-pages=0",
			expectedBody: "invalid URL query parameter provided: min-pages is 0 but should be in range [1,10000]",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - max-pages": {
			target:       "/?max-pages=10001",
			expectedBody: "invalid URL query parameter provided: max-pages is 10001 but should be in range [1,10000]",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - min-year": {
			target:       "/?min-year=1799",
			expectedBody: "invalid URL query parameter provided: min-year is 1799 but should be in range [1800,2100]",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - max-year": {
			target:       "/?max-year=2101",
			expectedBody: "invalid URL query parameter provided: max-year is 2101 but should be in range [1800,2100]",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - limit": {
			target:       "/?limit=0",
			expectedBody: "invalid URL query parameter provided: limit is 0 but should be greater than 0",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - authors": {
			target:       "/books?authors=1,beta,3",
			expectedBody: "invalid URL query parameter provided: recieved wrong type for parameter authors",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid books param - genres": {
			target:       "/books?genres=1,beta,3",
			expectedBody: "invalid URL query parameter provided: recieved wrong type for parameter genres",
			expectedCode: 400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			driverMock := drivertest.DriverMock{}
			handler := bookHandler{driver: &driverMock}

			r := bookRoutes(&handler)

			server := httptest.NewServer(r)
			defer server.Close()

			driverMock.
				On(tt.expectedHandler, mock.MatchedBy(func(_ context.Context) bool { return true }), tt.expectedParams).
				Return(tt.driverReturn, tt.expectedErr)

			resp, err := tt.sendRequest(fmt.Sprintf("%s%s", server.URL, tt.target))
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody+"\n", string(body))
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}

	t.Run("error when nil not provided", func(t *testing.T) {
		driverMock := drivertest.DriverMock{}
		handler := bookHandler{driver: &driverMock}
		driverMock.On("SearchBooks",
			mock.MatchedBy(func(_ context.Context) bool { return true }),
			book.SearchInput{}).Return([]entity.Book{}, nil)

		req := httptest.
			NewRequest("GET", "/", nil).
			WithContext(context.WithValue(context.Background(), bookSearchParamKey, (*BookQueryParams)(nil)))

		w := httptest.NewRecorder()
		handler.List(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "internal server error\n", string(body))
	})
}
