package v1

import (
	"context"
	"io/ioutil"
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
	driverMock := drivertest.DriverMock{}

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
		config          config.API
		expectedHandler string
		expectedBody    string
	}{
		"search all books": {
			expectedHandler: "SearchBooks",
			target:          "/api/v1/books",
			driverReturn:    []entity.Book{mockBook},
			expectedBody:    expectedJson,
		},
		"search for specific books": {
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
			target:       "/api/v1/books?title=The+Silmarillion&max_year=1980&min_year=1970&max_pages=400&min_pages=300&genres=2&genres=6&authors=42&authors=43&limit=50",
			driverReturn: []entity.Book{mockBook},
			expectedBody: expectedJson,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler := bookHandler{driver: &driverMock}

			driverMock.On(tt.expectedHandler,
				mock.MatchedBy(func(_ context.Context) bool { return true }),
				tt.expectedParams).Return(tt.driverReturn, nil)

			req := httptest.
				NewRequest("GET", tt.target, nil).
				WithContext(context.WithValue(context.Background(), bookSearchParamKey, &BookQueryParams{
					Title:            tt.expectedParams.Title,
					MaxYearPublished: tt.expectedParams.MaxYearPublished,
					MinYearPublished: tt.expectedParams.MinYearPublished,
					MaxPages:         tt.expectedParams.MaxPages,
					MinPages:         tt.expectedParams.MinPages,
					GenreIDs:         tt.expectedParams.GenreIDs,
					AuthorIDs:        tt.expectedParams.AuthorIDs,
					Limit:            tt.expectedParams.Limit,
				}))

			w := httptest.NewRecorder()

			handler.List(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody+"\n", string(body))
		})
	}
}
