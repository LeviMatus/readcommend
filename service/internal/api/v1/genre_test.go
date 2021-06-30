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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGenreHandler(t *testing.T) {

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
			h, err := NewGenreHandler(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestGenreHandler_List(t *testing.T) {
	driverMock := drivertest.DriverMock{}

	expectedJson := `[{"id":1,"title":"Fantasy/SciFy"}]`

	mockGenre := entity.Genre{
		ID:    1,
		Title: "Fantasy/SciFy",
	}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Genre
		sendRequest     func(string) (*http.Response, error)
		expectedHandler string
		expectedBody    string
	}{
		"search all genres": {
			expectedHandler: "ListGenres",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			driverReturn: []entity.Genre{mockGenre},
			expectedBody: expectedJson,
		},
		"invalid http method": {
			expectedHandler: "ListGenres",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
			driverReturn: []entity.Genre{mockGenre},
			expectedBody: "HTTP method POST is not allowed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler := genreHandler{driver: &driverMock}

			r := genreRoutes(&handler)

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
