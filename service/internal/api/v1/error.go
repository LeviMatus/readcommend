package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

const (
	methodNotAllowed    = "HTTP method %s is not allowed"
	internalServerError = "Internal Server Error"
)

// ErrorResponse is used to wrap errors and status codes for faulty requests and responses.
type ErrorResponse struct {
	// Err is the internal error.
	Err error `json:"-"`

	// StatusCode is the status code attached to this particular error.
	StatusCode int `json:"-"`

	// ErrorString is the message to be displayed to the client.
	ErrorString string `json:"message"`
}

// Render sets the status code for the request.
func (e *ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

// ErrBadRequest converts an error to a ErrorResponse for easy rendering on the client-side.
func ErrBadRequest(err error) render.Renderer {
	return &ErrorResponse{
		Err:         err,
		StatusCode:  http.StatusBadRequest,
		ErrorString: err.Error(),
	}
}

// ErrInternalServer returns a 400 status code with a string "Internal Server Error."
// For the purposes of the challenge, everything is a 400. In a real-world application this would use
// http.StatusInternalServerError.
func ErrInternalServer(err error) render.Renderer {
	return &ErrorResponse{
		Err:         err,
		StatusCode:  http.StatusBadRequest,
		ErrorString: internalServerError,
	}
}

// ErrMethodNotAllowed returns a 400 status code with a string specifying what method was rejected.
// For the purposes of the challenge, everything is a 400. In a real-world application this would use
// http.StatusMethodNotAllowed.
func ErrMethodNotAllowed(method string) render.Renderer {
	return ErrBadRequest(fmt.Errorf(methodNotAllowed, method))
}
