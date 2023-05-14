package types

import "github.com/pkg/errors"

// Predefined error messages
const (
	MsgMalformedRequest    = "got malformed request"
	MsgUnprocessableEntity = "got unprocessable request"
	MsgNoResultForPage     = "no result found for page"
	MsgNoResourceFound     = "no resource found"
	MsgUnauthorized        = "not authorized to perform this request"
	MsgForbidden           = "forbidden"
)

var (
	// ErrNoResultForPage is the error returned when there is no result for passed page.
	ErrNoResultForPage = errors.New(MsgNoResultForPage)

	// ErrNoResourceFound is the error returned when there is no resource found
	ErrNoResourceFound = errors.New(MsgNoResourceFound)

	// ErrMalformedRequest is the error returned when client send malformed request.
	ErrMalformedRequest = errors.New(MsgMalformedRequest)

	// ErrUnprocessableEntity is the error returned when client send unprocessable request.
	ErrUnprocessableEntity = errors.New(MsgUnprocessableEntity)

	// ErrUnauthorized is the error returned when client send unprocessable request.
	ErrUnauthorized = errors.New(MsgUnauthorized)

	// ErrForbidden is the error returned when client send unprocessable request.
	ErrForbidden = errors.New(MsgForbidden)
)
