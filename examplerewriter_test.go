package handlers_test

import (
	"github.com/gorilla/handlers"
)

type (
	APIEndpointError struct {
		Type APIErrorType
		Msg  string
	}
	APIEndpointErrorType string
)

const (
	ErrBadRequest    APIEndpointErrorType = "client_bad_request"
	ErrInternalError APIEndpointErrorType = "api_internal_error"
)

// JSONRewriter rewrites responses which are 4xx or 5xx
// into a specific error type a JSON REST API could use for its endpoints.
type JSONRewriter struct{}

func (jr JSONRewriter) RewriteIf(header http.Header, status int) bool {
	return status >= 400 /* 4xx and 5xx */ && !isContentType(header, "application/json")
}

func (jr JSONRewriter) RewriteHeader(header http.Header, status int) {
	header.Set("Content-Type", "application/json; charset=utf-8")
}

func (jr JSONRewriter) Rewrite(w io.Writer, b []byte, status int) error {
	err := APIEndpointError{Msg: string(b)}
	if status >= 400 && status < 500 {
		err.Type = ErrBadRequest
	} else if status >= 500 {
		err.Type = ErrInternalError
	}
	return json.NewEncoder(w).Encode(err)
}

func JSONResponseRewriteHandler(Fn func([]byte, int) interface{}, h http.Handler) http.Handler {
	return ResponseRewriteHandler(
		JSONRewriter{Fn},
		h,
	)
}
