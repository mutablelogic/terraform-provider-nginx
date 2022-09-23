package util

import (
	"encoding/json"
	"net/http"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ErrorResponse struct {
	Code   uint   `json:"code"`
	Reason string `json:"reason,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ContentTypeKey   = "Content-Type"
	ContentLengthKey = "Content-Length"
	ContentTypeJSON  = "application/json"
	ContentTypeText  = "text/plain"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ServeJSON is a utility function to serve an arbitary object as JSON
func ServeJSON(w http.ResponseWriter, v interface{}, code, indent uint) error {
	w.Header().Add(ContentTypeKey, ContentTypeJSON)
	w.WriteHeader(int(code))
	enc := json.NewEncoder(w)
	if indent > 0 {
		enc.SetIndent("", strings.Repeat(" ", int(indent)))
	}
	return enc.Encode(v)
}

// ServeText is a utility function to serve plaintext
func ServeText(w http.ResponseWriter, v string, code uint) {
	w.Header().Add(ContentTypeKey, ContentTypeText)
	w.WriteHeader(int(code))
	w.Write([]byte(v + "\n"))
}

// ServeEmpty is a utility function to serve an empty response
func ServeEmpty(w http.ResponseWriter, code uint) {
	w.Header().Add(ContentLengthKey, "0")
	w.WriteHeader(int(code))
}

// ServeError is a utility function to serve a JSON error notice
func ServeError(w http.ResponseWriter, code uint, reason ...string) error {
	err := ErrorResponse{code, strings.Join(reason, " ")}
	if len(reason) == 0 {
		err.Reason = http.StatusText(int(code))
	}
	return ServeJSON(w, err, code, 0)
}
