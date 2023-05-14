package routing

import (
	"encoding/json"
	"net/http"

	"github.com/mayankpahwa/aspire-loan-app/internal/resources/types"
	"github.com/pkg/errors"
)

const (
	// ContentTypeKey is content type key header.
	ContentTypeKey = "Content-Type"
	// ContentTypeJSON is content type for JSON.
	ContentTypeJSON = "application/json"
)

// ResponseError represents error response to UI.
type ResponseError struct {
	Result        string `json:"result"`
	StatusMessage string `json:"status_message"`
	Reason        string `json:"reason"`
}

func WriteJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WriteError writes error response with defined header.
func WriteError(r *http.Request, w http.ResponseWriter, err error) {

	errCode := http.StatusInternalServerError
	causer := errors.Cause(err)

	for e, status := range map[error]int{
		types.ErrMalformedRequest:    http.StatusBadRequest,
		types.ErrNoResultForPage:     http.StatusNotFound,
		types.ErrNoResourceFound:     http.StatusNotFound,
		types.ErrUnprocessableEntity: http.StatusUnprocessableEntity,
		types.ErrUnauthorized:        http.StatusUnauthorized,
		types.ErrForbidden:           http.StatusForbidden,
	} {
		if !errors.Is(causer, e) {
			continue
		}

		errCode = status
		break
	}

	Write(w, errCode, err.Error())
}

// Write response.
func Write(w http.ResponseWriter, code int, reason string) {
	w.Header().Add(ContentTypeKey, ContentTypeJSON)
	w.WriteHeader(code)
	// if code == http.StatusUnauthorized || code == http.StatusForbidden {
	// 	return
	// }

	WriteJSON(w, ResponseError{
		Result:        "error",
		StatusMessage: http.StatusText(code),
		Reason:        reason,
	})
}
