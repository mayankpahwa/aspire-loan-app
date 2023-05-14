package types

import "github.com/pkg/errors"

// Predefined error messages
const (
	MsgMalformedRequest    = "got malformed request"
	MsgUnprocessableEntity = "got unprocessable request"
	MsgNoResultForPage     = "no result found for page"
	MsgNoResourceFound     = "no resource found"
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
)
