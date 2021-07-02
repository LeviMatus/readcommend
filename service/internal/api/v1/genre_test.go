package v1

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre/genretest"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewGenreHandler(t *testing.T) {

	tests := map[string]struct {
		driver       genre.Driver
		errAssertion assert.ErrorAssertionFunc
		valAssertion assert.ValueAssertionFunc
	}{
		"nil driver provided": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"handler created": {
			driver:       &genretest.DriverMock{},
			errAssertion: assert.NoError,
			valAssertion: assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h, err := NewGenreHandler(tt.driver, zap.NewNop())
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestGenreHandler_List(t *testing.T) {
	expectedJson := `[{"id":1,"title":"Fantasy/SciFy"}]`

	mockGenre := entity.Genre{
		ID:    1,
		Title: "Fantasy/SciFy",
	}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Genre
		expectedHandler string
		expectedBody    string
		expectedCode    int
		expectedErr     error
		sendRequest     func(string) (*http.Response, error)
	}{
		"search all genres": {
			expectedHandler: "ListGenres",
			target:          "/",
			driverReturn:    []entity.Genre{mockGenre},
			expectedBody:    expectedJson,
			expectedCode:    200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid http method": {
			expectedHandler: "ListGenres",
			target:          "/",
			driverReturn:    []entity.Genre{mockGenre},
			expectedBody:    `{"message":"HTTP method POST is not allowed"}`,
			expectedCode:    400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
		},
		"driver returns error": {
			expectedHandler: "ListGenres",
			target:          "/",
			driverReturn:    []entity.Genre{},
			expectedBody:    `{"message":"Internal Server Error"}`,
			expectedCode:    400,
			expectedErr:     errors.New("mock error returned from driver"),
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			driverMock := genretest.DriverMock{}
			handler := genreHandler{driver: &driverMock, logger: zap.NewNop()}

			r := genreRoutes(&handler)

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
