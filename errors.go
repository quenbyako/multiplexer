package multiplexer

import (
	"errors"
	"fmt"
)

var (
	errPathSupportsPostOnly    = errors.New("Endpoint supports only POST requests")
	errServerJSONEncodingError = errors.New("Server can't encode response")
	errRequestBodyIsNotJSON    = errors.New("request body is not application/json")
	errTooManyURLsPerRequest   = errors.New("Too many urls per single request (max 20)")
)

type errRequestFailedDesc struct {
	err error
}

func (e *errRequestFailedDesc) Error() string { return "Request failed: " + e.err.Error() }
func (e *errRequestFailedDesc) Unwrap() error { return e.err }

type errRequestEndedWithBadStatusCode struct {
	code   int
	status string
}

func (e *errRequestEndedWithBadStatusCode) Error() string {
	return fmt.Sprintf("Request return %v code (%v)", e.code, e.status)
}

type errRequestFailedReadBody struct {
	err error
}

func (e *errRequestFailedReadBody) Error() string {
	return "Failed to read request body: " + e.err.Error()
}
func (e *errRequestFailedReadBody) Unwrap() error { return e.err }

type errFetchingEndpoint struct {
	url string
	err error
}

func (e *errFetchingEndpoint) Error() string {
	return fmt.Sprintf("Fetching %v: %v", e.url, e.err)
}
func (e *errFetchingEndpoint) Unwrap() error { return e.err }

type errParameterIsNotURL struct {
	indexes []string
}

func (e *errParameterIsNotURL) Error() string {
	return fmt.Sprintf("following indexes are not an url: %v", e.indexes)
}
