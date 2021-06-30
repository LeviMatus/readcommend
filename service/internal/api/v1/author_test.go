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
	"github.com/pkg/errors"
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
	expectedJson := `[{"id":1,"firstName":"John","lastName":"Tolkien"}]`

	mockAuthor := entity.Author{
		ID:        1,
		FirstName: "John",
		LastName:  "Tolkien",
	}

	tests := map[string]struct {
		target          string
		driverReturn    []entity.Author
		expectedHandler string
		expectedBody    string
		expectedCode    int
		expectedErr     error
		sendRequest     func(string) (*http.Response, error)
	}{
		"search all authors": {
			expectedHandler: "ListAuthors",
			target:          "/",
			driverReturn:    []entity.Author{mockAuthor},
			expectedBody:    expectedJson,
			expectedCode:    200,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Get(url)
			},
		},
		"invalid http method": {
			expectedHandler: "ListAuthors",
			target:          "/",
			driverReturn:    []entity.Author{mockAuthor},
			expectedBody:    `{"message":"HTTP method POST is not allowed"}`,
			expectedCode:    400,
			sendRequest: func(url string) (*http.Response, error) {
				return http.Post(url, "application/json", nil)
			},
		},
		"driver returns error": {
			expectedHandler: "ListAuthors",
			target:          "/",
			driverReturn:    []entity.Author{},
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
			driverMock := drivertest.DriverMock{}
			handler := authorHandler{driver: &driverMock}

			r := authorRoutes(&handler)

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
