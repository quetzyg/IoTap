package httpclient

import "errors"

var (
	errRequestUnauthorised = errors.New("unauthorised HTTP request")
	errRequestUnsuccessful = errors.New("unsuccessful HTTP request")
)
