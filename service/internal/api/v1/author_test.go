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

func TestNewAuthorHandler(t *testing.T) {

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
			h, err := NewAuthorHandler(tt.driver)
			tt.errAssertion(t, err)
			tt.valAssertion(t, h)
		})
	}
}

func TestAuthorHandler_List(t *testing.T) {
	driverMock := drivertest.DriverMock{}

	expectedJson := `[{"id":1,"firstName":"John","lastName":"Tolkien"}]`

	mockAuthor := entity.Author{
		ID:        1,
		FirstName: "John",
		LastName:  "Tolkien",
	}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Author
		sendRequest     func(string) (*http.Response, error)
		expectedHandler string
		expectedBody    string
	}{
		"search all authors": {
			expectedHandler: "ListAuthors",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
			driverReturn: []entity.Author{mockAuthor},
			expectedBody: expectedJson,
		},
		"invalid http method": {
			expectedHandler: "ListAuthors",
			target:          "/",
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
			driverReturn: []entity.Author{mockAuthor},
			expectedBody: "HTTP method POST is not allowed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler := authorHandler{driver: &driverMock}

			r := authorRoutes(&handler)

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
